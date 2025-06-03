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
	dbHost := os.Getenv("DB_HOST");
	dbPort := os.Getenv("DB_PORT");

	if pgUser == "" || pgPassword == "" || pgDB == "" || dbHost == "" || dbPort == "" {
		log.Fatalln("Env var(s) is not complete for db initialization");
		
	}

	var err error;
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", pgUser, pgPassword, dbHost, dbPort, pgDB)
	logger.Debug("postgresql.Init: initializing DB...", "dsn", dsn)
	DB, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		logger.Error("postgresql.Init: Failed to initialize DB", "error", err)
		// panic("Failed to initialze DB: " + err.Error())
		return err;
	}

	err = DB.Ping(context.Background())
	if err != nil {
		logger.Error("postgresql.Init: Failed to ping DB", "error", err)
		// panic("Failed to ping DB: " + err.Error())
		return err;
	}
	logger.Info("postgresql.Init: DB successfully initialized")
	return nil
}

