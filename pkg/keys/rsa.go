package keys

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"eau-de-go/settings"
	"encoding/pem"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	log "github.com/sirupsen/logrus"
)

type RsaKeyStore interface {
	CreateKeyPair() (*rsa.PrivateKey, *rsa.PublicKey, error)
	GetVerificationKey() (*rsa.PublicKey, error)
	GetSigningKey() (*rsa.PrivateKey, error)
}

// In-memory RSA key store, for monolithic deployment and development.
// RSA key pair is generated on first access and kept only in memory.
type inMemoryRsaKeyStore struct {
	signingKey      *rsa.PrivateKey
	verificationKey *rsa.PublicKey
}

var inMemoryRsaKeyStoreInstance *inMemoryRsaKeyStore

func GetInMemoryRsaKeyStore() RsaKeyStore {
	if inMemoryRsaKeyStoreInstance == nil {
		inMemoryRsaKeyStoreInstance = &inMemoryRsaKeyStore{}
	}
	return inMemoryRsaKeyStoreInstance
}

func (keyStore *inMemoryRsaKeyStore) CreateKeyPair() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	signingKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to create private key: %s", err))
		return nil, nil, err
	}
	verificationKey := &signingKey.PublicKey

	keyStore.signingKey = signingKey
	keyStore.verificationKey = verificationKey
	fmt.Println("Created new key pair")
	return signingKey, verificationKey, nil
}

func (keyStore *inMemoryRsaKeyStore) GetVerificationKey() (*rsa.PublicKey, error) {
	if keyStore.verificationKey == nil {
		_, _, err := keyStore.CreateKeyPair()
		if err != nil {
			fmt.Println(fmt.Sprintf("Failed to create key pair: %s", err))
			return nil, nil
		}
	}
	return keyStore.verificationKey, nil
}

func (keyStore *inMemoryRsaKeyStore) GetSigningKey() (*rsa.PrivateKey, error) {
	if keyStore.signingKey == nil {
		_, _, err := keyStore.CreateKeyPair()
		if err != nil {
			fmt.Println(fmt.Sprintf("Failed to create key pair: %s", err))
			return nil, err
		}
	}
	return keyStore.signingKey, nil
}

// AWS S3 RSA key store, for distributed deployment.
// RSA key pair is fetched from AWS S3 on first access and kept in memory
type awsS3RsaKeyStore struct {
	signingKey      *rsa.PrivateKey
	verificationKey *rsa.PublicKey
	Session         *session.Session
	Downloader      *s3manager.Downloader
}

var awsS3RsaKeyStoreInstance *awsS3RsaKeyStore

func GetAwsS3RsaKeyStore() RsaKeyStore {
	if awsS3RsaKeyStoreInstance == nil {
		session, err := session.NewSession(
			&aws.Config{Region: aws.String(settings.AwsS3KeyStoreRegion)},
		)
		if err != nil {
			log.Errorf("Failed to create AWS session: %s", err)
		}
		downloader := s3manager.NewDownloader(session)
		awsS3RsaKeyStoreInstance = &awsS3RsaKeyStore{Session: session, Downloader: downloader}
	}
	return awsS3RsaKeyStoreInstance
}

func (keyStore *awsS3RsaKeyStore) CreateKeyPair() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	signingKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to create private key: %s", err))
		return nil, nil, err
	}
	verificationKey := &signingKey.PublicKey
	err = keyStore.pushToS3(signingKey, verificationKey)
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to push key pair to S3: %s", err))
		return nil, nil, err
	}
	keyStore.signingKey = signingKey
	keyStore.verificationKey = verificationKey
	fmt.Println("Created new key pair")
	return signingKey, verificationKey, nil
}

func (keyStore *awsS3RsaKeyStore) pushToS3(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) error {
	uploader := s3manager.NewUploader(keyStore.Session)
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPem := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: privateKeyBytes})
	_, err := uploader.Upload(
		&s3manager.UploadInput{
			Bucket: aws.String(settings.AwsS3KeyStoreBucket),
			Key:    aws.String(settings.JwtSigningKeyPath),
			Body:   bytes.NewReader(privateKeyPem),
		})
	if err != nil {
		log.Errorf("Failed to upload private key to S3: %s", err)
		return err
	}
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		log.Errorf("Failed to marshal public key: %s", err)
		return err
	}
	publicKeyPem := pem.EncodeToMemory(&pem.Block{Type: "RSA PUBLIC KEY", Bytes: publicKeyBytes})
	_, err = uploader.Upload(
		&s3manager.UploadInput{
			Bucket: aws.String(settings.AwsS3KeyStoreBucket),
			Key:    aws.String(settings.JwtVerificationKeyPath),
			Body:   bytes.NewReader(publicKeyPem),
		})
	if err != nil {
		log.Errorf("Failed to upload public key to S3: %s", err)
		return err
	}
	return nil
}

func (keyStore *awsS3RsaKeyStore) GetVerificationKey() (*rsa.PublicKey, error) {
	if keyStore.verificationKey == nil {
		signingKey, verificationKey, err := keyStore.fetchFromS3()
		if err != nil {
			return nil, err
		}
		keyStore.signingKey = signingKey
		keyStore.verificationKey = verificationKey
	}
	return keyStore.verificationKey, nil
}

func (keyStore *awsS3RsaKeyStore) GetSigningKey() (*rsa.PrivateKey, error) {
	if keyStore.signingKey == nil {
		signingKey, verificationKey, err := keyStore.fetchFromS3()
		if err != nil {
			return nil, err
		}
		keyStore.signingKey = signingKey
		keyStore.verificationKey = verificationKey
	}
	return keyStore.signingKey, nil
}

func (keyStore *awsS3RsaKeyStore) fetchFromS3() (*rsa.PrivateKey, *rsa.PublicKey, error) {

	privateKeyBuf := new(aws.WriteAtBuffer)
	_, err := keyStore.Downloader.Download(
		privateKeyBuf,
		&s3.GetObjectInput{
			Bucket: aws.String(settings.AwsS3KeyStoreBucket),
			Key:    aws.String(settings.JwtSigningKeyPath),
		})
	if err != nil {
		log.Errorf("Failed to download private key from S3: %s", err)
		return nil, nil, err
	}
	privateBlock, _ := pem.Decode(privateKeyBuf.Bytes())
	if privateBlock == nil {
		log.Error("Failed to decode private key")
		return nil, nil, err
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(privateBlock.Bytes)
	if err != nil {
		log.Errorf("Failed to parse private key: %s", err)
		return nil, nil, err
	}
	publicKeyBuf := new(aws.WriteAtBuffer)
	_, err = keyStore.Downloader.Download(
		publicKeyBuf,
		&s3.GetObjectInput{
			Bucket: aws.String(settings.AwsS3KeyStoreBucket),
			Key:    aws.String(settings.JwtVerificationKeyPath),
		})
	if err != nil {
		log.Errorf("Failed to download public key from S3: %s", err)
		return nil, nil, err
	}
	publicBlock, _ := pem.Decode(publicKeyBuf.Bytes())
	if publicBlock == nil {
		log.Error("Failed to decode public key")
		return nil, nil, err
	}
	publicInterface, err := x509.ParsePKIXPublicKey(publicBlock.Bytes)
	if err != nil {
		log.Errorf("Failed to parse public key: %s", err)
		return nil, nil, err
	}
	publicKey, ok := publicInterface.(*rsa.PublicKey)
	if !ok {
		log.Error("Failed to cast public key to RSA public key")
		return nil, nil, err
	}
	return privateKey, publicKey, nil
}
