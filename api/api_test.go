package api

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
)

func TestNewAPI(t *testing.T) {
	ctx := context.Background()

	db, err := pgx.Connect(ctx, "postgres://postgres:123@localhost:6433/test_db")
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	defer db.Close(ctx)

	api, err := NewAPI(db)
	assert.NoError(t, err)
	assert.NotNil(t, api)
	assert.Equal(t, db, api.DB)
}

func TestServeHTTP(t *testing.T) {
	ctx := context.Background()

	db, err := pgx.Connect(ctx, "postgres://postgres:123@localhost:6433/test_db")
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	defer db.Close(ctx)

	api, err := NewAPI(db)
	assert.NoError(t, err)

	req := httptest.NewRequest("GET", "/graphql", nil)
	w := httptest.NewRecorder()

	api.ServeHTTP(w, req)

	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "POST, GET, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(t, "Content-Type, Authorization", w.Header().Get("Access-Control-Allow-Headers"))
}

func TestPlaygroundHandler(t *testing.T) {
	ctx := context.Background()

	db, err := pgx.Connect(ctx, "postgres://postgres:123@localhost:6433/test_db")
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	defer db.Close(ctx)

	api, err := NewAPI(db)
	assert.NoError(t, err)

	handler := api.PlaygroundHandler()
	assert.NotNil(t, handler)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.NotEqual(t, 0, w.Body.Len())
}
