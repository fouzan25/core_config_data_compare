package main

import (
	"log"

	"github.com/joho/godotenv"
)

func loadENV() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}
