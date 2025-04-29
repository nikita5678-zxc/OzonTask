package api

import (
	"OzonTask/db"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v4"
)

type API struct {
	Conn *pgx.Conn
}

func NewAPI(conn *pgx.Conn) *API {
	return &API{Conn: conn}
}

func (api *API) GraphqlHandler(w http.ResponseWriter, r *http.Request) {
	var query string

	if r.Method == http.MethodGet {
		query = r.URL.Query().Get("query")
	} else if r.Method == http.MethodPost {
		var req map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		q, ok := req["query"].(string)
		if !ok {
			http.Error(w, "Query parameter must be a string", http.StatusBadRequest)
			return
		}
		query = q
	} else {
		http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
		return
	}

	if query == "" {
		http.Error(w, "Query parameter must be a string", http.StatusBadRequest)
		return
	}

	if strings.Contains(query, "posts") {
		api.GetPostHandler(w)
		return
	}

	if strings.Contains(query, "createPost") {
		api.CreatePostHandler(w, query)
		return
	}

	http.Error(w, "Unsupported query", http.StatusBadRequest)
}

func (api *API) GetPostHandler(w http.ResponseWriter) {
	posts, err := db.GetPosts(api.Conn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{
		"data": map[string]interface{}{
			"posts": posts,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (api *API) CreatePostHandler(w http.ResponseWriter, query string) {
	title := ExtractField(query, "title")
	content := ExtractField(query, "content")
	author := ExtractField(query, "author")
	allowComments := strings.Contains(query, "allowComments: true")

	post, err := db.CreatePost(api.Conn, title, content, author, allowComments)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{
		"data": map[string]interface{}{
			"createPost": post,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func ExtractField(query, field string) string {
	start := strings.Index(query, field+":")
	if start == -1 {
		return ""
	}
	start += len(field) + 2
	end := start
	for end < len(query) && query[end] != '"' && query[end] != ',' && query[end] != '}' {
		end++
	}
	return strings.Trim(query[start:end], "\" ")
}
