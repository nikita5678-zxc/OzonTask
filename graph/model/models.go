package model

import (
	"github.com/google/uuid"
)

type Post struct {
	ID            uuid.UUID `json:"id"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	Author        string    `json:"author"`
	AllowComments bool      `json:"allowComments"`
	CreatedAt     string    `json:"createdAt"`
}

type Comment struct {
	ID        uuid.UUID  `json:"id"`
	PostID    uuid.UUID  `json:"postId"`
	ParentID  *uuid.UUID `json:"parentId,omitempty"`
	Author    string     `json:"author"`
	Content   string     `json:"content"`
	CreatedAt string     `json:"createdAt"`
}

type NewPost struct {
	Title         string `json:"title"`
	Content       string `json:"content"`
	Author        string `json:"author"`
	AllowComments bool   `json:"allowComments"`
}

type NewComment struct {
	PostID   string  `json:"postId"`
	ParentID *string `json:"parentId,omitempty"`
	Author   string  `json:"author"`
	Content  string  `json:"content"`
}

type CommentEdge struct {
	Node   *Comment `json:"node"`
	Cursor string   `json:"cursor"`
}

type CommentsConnection struct {
	Edges    []*CommentEdge `json:"edges"`
	PageInfo *PageInfo      `json:"pageInfo"`
}

type PageInfo struct {
	HasNextPage     bool    `json:"hasNextPage"`
	HasPreviousPage bool    `json:"hasPreviousPage"`
	StartCursor     *string `json:"startCursor"`
	EndCursor       *string `json:"endCursor"`
}
