package model

import (
	"database/sql"
	"time"
)

type Post struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	Author        string    `json:"author"`
	AllowComments bool      `json:"allow_comments"`
	CreatedAt     time.Time `json:"created_at"`
}

type NewPost struct {
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	Author        string    `json:"author"`
	AllowComments bool      `json:"allow_comments"`
	CreatedAt     time.Time `json:"created_at"`
}

type Comment struct {
	ID        string         `json:"id"`
	PostID    string         `json:"post_id"`
	ParentID  sql.NullString `json:"parent_id"`
	Author    string         `json:"author"`
	Content   string         `json:"content"`
	CreatedAt time.Time      `json:"created_at"`
}

type NewComment struct {
	PostID    string         `json:"post_id"`
	ParentID  sql.NullString `json:"parent_id"`
	Author    string         `json:"author"`
	Content   string         `json:"content"`
	CreatedAt time.Time      `json:"created_at"`
}

type CommentEdge struct {
	Node   *Comment `json:"node"`
	Cursor string   `json:"cursor"`
}

type PageInfo struct {
	HasNextPage     bool    `json:"has_next_page"`
	HasPreviousPage bool    `json:"has_previous_page"`
	StartCursor     *string `json:"start_cursor"`
	EndCursor       *string `json:"end_cursor"`
}

type CommentConnection struct {
	Edges    []*CommentEdge `json:"edges"`
	PageInfo *PageInfo      `json:"page_info"`
}
