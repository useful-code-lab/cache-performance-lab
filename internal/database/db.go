package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func ConnectDB() (*pgx.Conn, error) {
	// Параметры из docker-compose
	connStr := "postgres://user:password@localhost:5432/inventory"
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}
	return conn, nil
}
