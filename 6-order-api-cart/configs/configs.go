package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Verify VerifyConfig
	Db     DbCongig
	Secret string
}

type VerifyConfig struct {
	Email    string
	Password string
	Address  string
}

type DbCongig struct {
	Dsn string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using default configuration.")
	}
	return &Config{
		Verify: VerifyConfig{
			Email:    "kirill.zemskoff@gmail.com",
			Password: "123",
			Address:  "smtp.gmail.com",
		},
		Db: DbCongig{
			Dsn: os.Getenv("DSN"),
		},
		Secret: os.Getenv("SECRET"),
	}
}
