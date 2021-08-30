package main

import (
	"log"
	"simple-go-auth/pkg/server"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	server.Run()
}
