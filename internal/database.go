package internal

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDB(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	conn, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}

	log.Println("INFO: Successfully connected to the database")
	return conn, nil
}
