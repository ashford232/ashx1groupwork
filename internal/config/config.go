package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type CloudflareConfig struct {
	AccountID       string
	AccessKeyID     string
	SecretAccessKey string
	BucketName      string
	PublicURL       string
}
type Config struct {
	DatabaseURL string
	Port        string
	JWTSecret   string
	Cloudflare  CloudflareConfig
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println(err)
	}
	config := &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Port:        os.Getenv("PORT"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		Cloudflare: CloudflareConfig{
			AccountID:       os.Getenv("R2_ACCOUNT_ID"),
			AccessKeyID:     os.Getenv("R2_ACCESS_KEY_ID"),
			SecretAccessKey: os.Getenv("R2_SECRET_ACCESS_KEY"),
			BucketName:      os.Getenv("R2_BUCKET_NAME"),
			PublicURL:       os.Getenv("R2_PUBLIC_URL"),
		},
	}
	return config, nil
}
