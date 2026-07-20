package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT          string
	DB_HOST       string
	DB_PORT       string
	DB_NAME       string
	DB_USER       string
	DB_PASSWORD   string
	DB_SSLMODE    string
	RESULTS_FILE  string
	CHECK_TIMEOUT time.Duration
	TZ            string
}

func LoadDotEnv() *Config {
	if err := godotenv.Load("./.env"); err != nil {
		fmt.Printf("aviso: erro ao ler a .env: %v\n", err)
	}

	checkTimeout, _ := time.ParseDuration(os.Getenv("CHECK_TIMEOUT"))
	if checkTimeout == 0 {
		checkTimeout = 10 * time.Second
	}

	resultsFile := os.Getenv("RESULTS_FILE")
	if resultsFile == "" {
		resultsFile = "results.txt"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	tz := os.Getenv("TZ")
	if tz == "" {
		tz = "America/Sao_Paulo"
	}

	return &Config{
		PORT:          port,
		DB_HOST:       os.Getenv("DB_HOST"),
		DB_PORT:       os.Getenv("DB_PORT"),
		DB_NAME:       os.Getenv("DB_NAME"),
		DB_USER:       os.Getenv("DB_USER"),
		DB_PASSWORD:   os.Getenv("DB_PASSWORD"),
		DB_SSLMODE:    os.Getenv("DB_SSLMODE"),
		RESULTS_FILE:  resultsFile,
		CHECK_TIMEOUT: checkTimeout,
		TZ:            tz,
	}
}
