package api

import (
	"OzonTask/db"
	"OzonTask/graph/generated"
	"OzonTask/graph/model"
	"context"
	"log"

	"github.com/jackc/pgx/v4"
)

type API struct {
	Conn *pgx.Conn
}

func NewAPI(conn *pgx.Conn) *API {
	return &API{Conn: conn}
}

type Resolver struct {
	DB *pgx.Conn
}

func (api *API) GetResolver() *Resolver {
	return &Resolver{DB: api.Conn}
}

func (r *Resolver) Mutation() generated.MutationResolver {
	return &mutationResolver{r}
}

func (r *Resolver) Query() generated.QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

func (r *mutationResolver) CreatePost(ctx context.Context, input model.NewPost) (*model.Post, error) {
	post, err := db.CreatePost(r.DB, input.Title, input.Content, input.Author, input.AllowComments)
	if err != nil {
		log.Printf("Error creating post: %v", err)
		return nil, err
	}

	return &model.Post{
		ID:            post.ID,
		Title:         post.Title,
		Content:       post.Content,
		Author:        post.Author,
		AllowComments: post.AllowComments,
	}, nil
}

func (r *queryResolver) Posts(ctx context.Context) ([]*model.Post, error) {
	dbPosts, err := db.GetPosts(r.DB)
	if err != nil {
		log.Printf("Error getting posts: %v", err)
		return nil, err
	}

	var posts []*model.Post
	for _, p := range dbPosts {
		posts = append(posts, &model.Post{
			ID:            p.ID,
			Title:         p.Title,
			Content:       p.Content,
			Author:        p.Author,
			AllowComments: p.AllowComments,
		})
	}

	return posts, nil
}
