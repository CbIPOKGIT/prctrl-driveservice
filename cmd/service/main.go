package main

import (
	_ "embed"
	"log"
	"os"

	"github.com/CbIPOKGIT/prctrl-driveservice/internal/app"
	"github.com/joho/godotenv"
)

//go:embed .env
var envContent []byte

func init() {
	if envData, err := godotenv.UnmarshalBytes(envContent); err == nil {
		for key, value := range envData {
			if err := os.Setenv(key, value); err != nil {
				log.Fatalf("Failed to set environment key %s: %s", key, err)
			}
		}
	} else {
		log.Fatal("Failed to parse .env file: ", err)
	}

	var requiredEnv = []string{"PORT", "DRIVE_SECRET_TOKEN"}
	for _, key := range requiredEnv {
		if os.Getenv(key) == "" {
			log.Fatalf("%s environment variable is required", key)
		}
	}

	if content, err := os.ReadFile("service-token.json"); err == nil {
		os.Setenv("SERVICE_TOKEN", string(content))
	} else {
		log.Fatal("Failed to read service token file: ", err)
	}
}

func main() {
	app.StartApplication()
}
