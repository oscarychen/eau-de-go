package settings

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"
)

var (
	DbHost           string
	DbPort           string
	SslMode          string
	DbUsername       string
	DbName           string
	DbPassword       string
	RefreshTokenLife time.Duration
	AccessTokenLife  time.Duration
)

func init() {
	log.Printf("Initializing settings...")
	err := godotenv.Load()
	if err != nil {
		log.Error("Error loading .env file")
	}

	DbHost = getEnv("DB_HOST", "localhost")
	DbPort = getEnv("DB_PORT", "5432")
	SslMode = getEnv("SSL_MODE", "disable")
	DbUsername = getEnv("DB_USERNAME", "postgres")
	DbName = getEnv("DB_NAME", "eau-de-go")
	DbPassword = getEnv("DB_PASSWORD", "")

	if refreshTokenLifeMinutes, err := strconv.Atoi(getEnv("REFRESH_TOKEN_LIFE_MINUTES", "")); err == nil {
		RefreshTokenLife = time.Minute * time.Duration(refreshTokenLifeMinutes)
	} else {
		defaultRefreshTokenLifeMinutes := 10080
		RefreshTokenLife = time.Minute * time.Duration(defaultRefreshTokenLifeMinutes)
	}

	if accessTokenLifeMinutes, err := strconv.Atoi(getEnv("ACCESS_TOKEN_LIFE_MINUTES", "")); err == nil {
		AccessTokenLife = time.Minute * time.Duration(accessTokenLifeMinutes)
	} else {
		defaultAccessTokenLifeMinutes := 15
		AccessTokenLife = time.Minute * time.Duration(defaultAccessTokenLifeMinutes)
	}
}

func getEnv(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Printf("Warning: %s is not set, using default value '%s'", key, defaultValue)
		return defaultValue
	}
	return value
}