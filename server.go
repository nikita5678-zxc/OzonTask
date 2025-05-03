package main

import (
	"context"
	"log"
	"net/http"

	"OzonTask/api"
	"OzonTask/db"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	conn, err := db.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer conn.Close(context.Background())

	if err := db.CreateSchema(conn); err != nil {
		log.Fatalf("Failed to create schema: %v", err)
	}

	api, err := api.NewAPI(conn)
	if err != nil {
		log.Fatalf("Failed to create API: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/graphql", api)
	mux.Handle("/", api.PlaygroundHandler())

	log.Println("Server started at http://localhost:3000")
	log.Println("GraphQL Playground available at http://localhost:3000/")
	log.Fatal(http.ListenAndServe(":3000", mux))
}
