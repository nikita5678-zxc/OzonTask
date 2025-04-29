package main

import (
	"OzonTask/api"
	"OzonTask/db"
	"OzonTask/graph"
	"OzonTask/graph/generated"
	"context"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

func main() {
	conn, err := db.ConnectDB()
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer conn.Close(context.Background())

	apiHandler := api.NewAPI(conn)

	port := os.Getenv("DATABASE_URL")
	if port == "" {
		log.Fatal("Missing DATABASE_URL")
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: &graph.Resolver{
			Resolver: apiHandler.GetResolver(),
		},
	}))

	http.Handle("/", playground.Handler("GraphQL Playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("Server started on :7070")
}
