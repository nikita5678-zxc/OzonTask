package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"OzonTask/model/dbmodel"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

func ConnectDB() (*pgx.Conn, error) {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is not set")
	}

	if !strings.Contains(connStr, "sslmode=") {
		if strings.Contains(connStr, "?") {
			connStr += "&sslmode=disable"
		} else {
			connStr += "?sslmode=disable"
		}
	}

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	if err := conn.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	log.Println("Successfully connected to database")
	return conn, nil
}

func CreateSchema(conn *pgx.Conn) error {
	_, err := conn.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS posts (
			id UUID PRIMARY KEY,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			author TEXT NOT NULL,
			allow_comments BOOLEAN NOT NULL DEFAULT true,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS comments (
			id UUID PRIMARY KEY,
			post_id UUID NOT NULL REFERENCES posts(id),
			parent_id UUID REFERENCES comments(id),
			author TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
	`)
	return err
}

func GetPosts(conn *pgx.Conn) ([]*dbmodel.Post, error) {
	rows, err := conn.Query(context.Background(), "SELECT * FROM posts ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*dbmodel.Post
	for rows.Next() {
		post := &dbmodel.Post{}
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Author, &post.AllowComments, &post.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func GetPost(conn *pgx.Conn, id uuid.UUID) (*dbmodel.Post, error) {
	post := &dbmodel.Post{}
	err := conn.QueryRow(context.Background(), "SELECT * FROM posts WHERE id = $1", id).
		Scan(&post.ID, &post.Title, &post.Content, &post.Author, &post.AllowComments, &post.CreatedAt)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func CreatePost(conn *pgx.Conn, post *dbmodel.Post) error {
	_, err := conn.Exec(context.Background(), `
		INSERT INTO posts (id, title, content, author, allow_comments, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, post.ID, post.Title, post.Content, post.Author, post.AllowComments, post.CreatedAt)
	return err
}

func GetComments(conn *pgx.Conn, postID uuid.UUID, limit int, after string) ([]*dbmodel.Comment, error) {
	query := `
			SELECT id, post_id, parent_id, author, content, created_at
			FROM comments
			WHERE post_id = $1
		`
	args := []interface{}{postID}

	if after != "" {
		query += " AND created_at > (SELECT created_at FROM comments WHERE id = $2)"
		args = append(args, after)
	}

	query += " ORDER BY created_at ASC LIMIT $3"
	args = append(args, limit)

	rows, err := conn.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*dbmodel.Comment
	for rows.Next() {
		comment := &dbmodel.Comment{}
		err := rows.Scan(&comment.ID, &comment.PostID, &comment.ParentID, &comment.Author, &comment.Content, &comment.CreatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func CreateComment(conn *pgx.Conn, comment *dbmodel.Comment) error {
	_, err := conn.Exec(context.Background(), `
		INSERT INTO comments (id, post_id, parent_id, author, content, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, comment.ID, comment.PostID, comment.ParentID, comment.Author, comment.Content, comment.CreatedAt)
	return err
}

func GetCommentReplies(db *pgx.Conn, parentID uuid.UUID, limit int, cursor string) ([]*dbmodel.Comment, error) {
	query := `
		SELECT id, post_id, parent_id, author, content, created_at
		FROM comments
		WHERE parent_id = $1 AND id > $2
		ORDER BY id
		LIMIT $3
	`
	rows, err := db.Query(context.Background(), query, parentID, cursor, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*dbmodel.Comment
	for rows.Next() {
		var c dbmodel.Comment
		err := rows.Scan(&c.ID, &c.PostID, &c.ParentID, &c.Author, &c.Content, &c.CreatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &c)
	}
	return comments, nil
}
