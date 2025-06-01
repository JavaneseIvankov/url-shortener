package db

import (
	"context"
	"fmt"
	"javaneseivankov/url-shortener/pkg/logger"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func Init() error {
	pgUser := os.Getenv("POSTGRES_USER")
	pgPassword := os.Getenv("POSTGRES_PASSWORD")
	pgDB := os.Getenv("POSTGRES_DB")

	if pgUser == "" || pgPassword == "" || pgDB == "" {
		log.Fatalln("POSTGRES_USER, POSTGRES_PASSWORD, or POSTGRES_DB not set")
		
	}

	dsn := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s", pgUser, pgPassword, pgDB)
	logger.Debug("postgresql.Init: initializing DB...", "dsn", dsn)
	DB, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		logger.Error("postgresql.Init: Failed to initialize DB", "error", err)
		panic("Failed to initialze DB: " + err.Error())
	}

	err = DB.Ping(context.Background())
	if err != nil {
		logger.Error("postgresql.Init: Failed to ping DB", "error", err)
		panic("Failed to ping DB: " + err.Error())
	}
	logger.Info("postgresql.Init: DB successfully initialized")
	return nil
}

