package database

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func ConnectDB() (*pgxpool.Pool, error) {

	dbUser := "admin"
	dbPass := "admin"
	dbHost := "localhost"
	dbPort := "5455"
	dbName := "products_db"

	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPass, dbHost, dbPort, dbName,
	)

	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("não foi possível criar o pool de conexões: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, fmt.Errorf("não foi possível conectar ao banco de dados: %w", err)
	}

	log.Println("Conexão com o PostgreSQL estabelecida com sucesso!")
	DB = pool
	return DB, nil
}
