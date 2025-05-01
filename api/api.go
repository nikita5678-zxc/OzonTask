package api

import (
	"log"
	"net/http"

	"OzonTask/graph"
	"OzonTask/graph/generated"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/jackc/pgx/v4"
)

type API struct {
	DB *pgx.Conn
}

func NewAPI(db *pgx.Conn) (*API, error) {
	return &API{
		DB: db,
	}, nil
}

func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Настройка CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Создаем GraphQL сервер
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: &graph.Resolver{DB: a.DB},
	}))

	// Логируем входящие запросы
	log.Printf("Received %s request to %s", r.Method, r.URL.Path)

	// Обрабатываем запрос
	srv.ServeHTTP(w, r)
}

func (a *API) PlaygroundHandler() http.Handler {
	return playground.Handler("GraphQL playground", "/graphql")
}
