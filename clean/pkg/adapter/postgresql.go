package adapter

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

func New() *pgxpool.Pool {
	databaseUrl := "postgres://postgres:mypassword@localhost:5432/postgres"

	// this returns connection pool
	conn, err := pgxpool.Connect(context.Background(), databaseUrl)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	// to close DB pool
	defer conn.Close()

	return conn
}
