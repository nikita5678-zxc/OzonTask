package main

import (
	"OzonTask/api"
	"OzonTask/db"
	"context"
	"log"
	"net/http"
)

func main() {
	conn, err := db.ConnectDB()
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	defer conn.Close(context.Background())

	apiHandler := api.NewAPI(conn)
	http.HandleFunc("/graphql", apiHandler.GraphqlHandler)

	log.Println("Server started at :7070")
	err = http.ListenAndServe(":7070", nil)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
