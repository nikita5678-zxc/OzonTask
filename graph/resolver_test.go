package graph

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
)

func TestNewResolver(t *testing.T) {
	ctx := context.Background()

	db, err := pgx.Connect(ctx, "postgres://postgres:123@localhost:6433/test_db")
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	defer db.Close(ctx)

	resolver := &Resolver{DB: db}
	assert.NotNil(t, resolver)
	assert.Equal(t, db, resolver.DB)
}

func TestResolverContext(t *testing.T) {
	ctx := context.Background()

	db, err := pgx.Connect(ctx, "postgres://postgres:123@localhost:6433/test_db")
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	defer db.Close(ctx)

	resolver := &Resolver{DB: db}

	testCtx := context.WithValue(ctx, "resolver", resolver)

	assert.NotNil(t, testCtx)

	assert.NotNil(t, testCtx.Value("resolver"))
}
