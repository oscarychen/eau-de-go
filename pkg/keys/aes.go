package keys

import "crypto/rand"

var aesKey *[]byte

func GetInMemoryAesKey() *[]byte {
	if aesKey == nil {
		aesKey = makeKey()
	}
	return aesKey
}

func makeKey() *[]byte {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return nil
	}
	return &key
}
