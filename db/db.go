package db

import (
	"OzonTask/graph/model"
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func ConnectDB() (*pgx.Conn, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}

	conn, err := pgx.Connect(context.Background(), databaseUrl)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return conn, nil
}

func GetPosts(conn *pgx.Conn) ([]model.Post, error) {
	rows, err := conn.Query(context.Background(), "SELECT id, title, content, author, allow_comments FROM posts")
	if err != nil {
		log.Println("Error querying posts:", err)
		return nil, err
	}
	defer rows.Close()

	var posts []model.Post
	for rows.Next() {
		var post model.Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Author, &post.AllowComments)
		if err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}
		posts = append(posts, post)
	}
	if rows.Err() != nil {
		log.Println("Error reading rows:", rows.Err())
		return nil, rows.Err()
	}
	return posts, nil
}

func CreatePost(conn *pgx.Conn, title, content, author string, allowComments bool) (model.Post, error) {
	var post model.Post
	err := conn.QueryRow(context.Background(),
		"INSERT INTO posts (title, content, author, allow_comments) VALUES ($1, $2, $3, $4) RETURNING id",
		title, content, author, allowComments).Scan(&post.ID)
	if err != nil {
		log.Println("Error inserting post:", err)
		return post, err
	}
	post.Title = title
	post.Content = content
	post.Author = author
	post.AllowComments = allowComments
	return post, nil
}
