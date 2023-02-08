package main

import (
	"fmt"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	setupMasscan()
	connectRedis()

	startRoutines()
}

func startRoutines() {
	go scan()
	go test()
}
