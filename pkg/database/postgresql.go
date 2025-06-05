package database

import (
	"context"
	"errors"
	"fmt"
	"javaneseivankov/url-shortener/pkg/logger"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

type PgDB struct {
	pgUser string
	pgPassword string
	pgDB string
	dbHost string
	dbPort string
	db *pgxpool.Pool
}

func NewPgDB(user string, password string, dbName string, dbHost string, dbPort string) *PgDB {
	return &PgDB{
		pgUser: user,
		pgPassword: password,
		pgDB: dbName,
		dbHost: dbHost,
		dbPort: dbPort,
	}
}

func createDSN(user string, password string, dbName string, dbHost string, dbPort string) string {
	if user == "" || password == "" || dbName == "" || dbHost == "" || dbPort == "" {
		log.Fatalln("Env var(s) is not complete for db initialization");
		
	}
		return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, dbHost, dbPort, dbName)
}

func (pg  *PgDB) Init() (*pgxpool.Pool,error) {
	var err error
	dsn := createDSN(pg.pgUser, pg.pgPassword, pg.pgDB, pg.dbHost, pg.dbPort)

	var db *pgxpool.Pool;
	logger.Debug("postgresql.Init: initializing DB...", "dsn", dsn)
	db, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		logger.Error("postgresql.Init: Failed to initialize DB", "error", err)
		return nil,err;
	}

	err = db.Ping(context.Background())
	if err != nil {
		logger.Error("postgresql.Init: Failed to ping DB", "error", err)
		return nil,err;
	}

	pg.db = db;

	logger.Info("postgresql.Init: DB successfully initialized")
	return db, nil;
}

func (pg *PgDB) MigrateDB() error {
	logger.Info("postgresql.MigrateDB: Starting migration...")
	dsn := createDSN(pg.pgUser, pg.pgPassword, pg.pgDB, pg.dbHost, pg.dbPort)
	m, err := migrate.New("file:db/migrations/", dsn)
	if err != nil {
		logger.Error("postgresql.MigrateDB: Create new migration failed", "error", err)
		return err
	}

	defer func() {
		if srcErr, dbErr := m.Close(); srcErr != nil || dbErr != nil {
			logger.Error("postgresql.MigrateDB: Failed to close migration", "sourceError", srcErr, "dbError", dbErr)
		}
	}()

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.Error("postgresql.MigrateDB: Up migration failed", "error", err)
		return err
	}
	logger.Info("postgresql.MigrateDB: Migration completed successfuly")

	return nil
}
