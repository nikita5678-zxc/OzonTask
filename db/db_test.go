package db

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
)

func TestConnectDB(t *testing.T) {
	ctx := context.Background()

	db, err := pgx.Connect(ctx, "postgres://postgres:123@localhost:6433/test_db")
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	defer db.Close(ctx)

	err = db.Ping(ctx)
	assert.NoError(t, err)
}

func TestCreateSchema(t *testing.T) {
	ctx := context.Background()

	db, err := pgx.Connect(ctx, "postgres://postgres:123@localhost:6433/test_db")
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	defer db.Close(ctx)

	err = CreateSchema(db)
	assert.NoError(t, err)

	var count int
	err = db.QueryRow(ctx, `
		SELECT COUNT(*) 
		FROM information_schema.tables 
		WHERE table_schema = 'public'
	`).Scan(&count)
	assert.NoError(t, err)
	assert.Greater(t, count, 0)
}

func TestTransaction(t *testing.T) {
	ctx := context.Background()

	db, err := pgx.Connect(ctx, "postgres://postgres:123@localhost:6433/test_db")
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	defer db.Close(ctx)

	tx, err := db.Begin(ctx)
	assert.NoError(t, err)

	assert.NotNil(t, tx)

	err = tx.Rollback(ctx)
	assert.NoError(t, err)
}
