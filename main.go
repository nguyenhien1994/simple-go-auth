package main

import (
	"fmt"
	"log"
	"os"
	"simple-go-auth/controller"

	"github.com/joho/godotenv"
)

func main() {
	Run()
}

func Run() {
	server := controller.Server{}
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, not comming through %v", err)
	} else {
		fmt.Println("We are getting the env values")
	}

	appAddr := ":" + os.Getenv("PORT")

	server.Initialize()

	server.Run(appAddr)
}
