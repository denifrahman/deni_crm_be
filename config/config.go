package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

var AppConfig DBConfig

func LoadConfig() {
	rootPath, err := filepath.Abs(".")
	if err != nil {
		log.Fatal("Cannot resolve project root:", err)
	}

	paths := []string{
		filepath.Join(rootPath, ".env"),
		filepath.Join(rootPath, "../.env"),
		filepath.Join(rootPath, "../../.env"),
	}

	loaded := false
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			err = godotenv.Load(p)
			if err != nil {
				log.Println("Error loading .env at", p, ":", err)
			} else {
				log.Println("✅ .env loaded from:", p)
				loaded = true
				break
			}
		}
	}

	if !loaded {
		log.Println("⚠️  .env not found. Using system environment variables.")
	}
	AppConfig = DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	}
	// Opsional: log kalau variabel penting kosong
	if AppConfig.Host == "" || AppConfig.User == "" {
		log.Println("Warning: some environment variables are missing!")
	}
}

func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
