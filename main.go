package main

import (
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env only once at the very start
	if err := godotenv.Load(); err != nil {
		log.Println("Notice: No .env file found, falling back to system environment variables")
	}

	store, err := NewPostgresAccountStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	server := NewAPIServer(":3000", store)
	server.Run()
}
