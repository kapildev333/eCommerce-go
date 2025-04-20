package db

import (
	"database/sql"
	config "eCommerce-go/utils"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"os"
)

var DB *sql.DB

func InitDB() {
	log := config.GetLogger()
	log.With("component", "db")
	var err error
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://user:password@localhost/dbname?sslmode=disable"
	}
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Error("Failed to open database connection", "error", err)
	}

	if err = DB.Ping(); err != nil {
		log.Error("Failed to ping database", "error", err)
	}

	_, err = NewRedisStore()
	if err != nil {
		log.Error("Failed to connect to Redis", "error", err)
	} else {
		log.Info("Connected to Redis successfully")
	}
	migration, err := migrate.New("file://db/migration", connStr)
	if err != nil {
		log.Error("Failed to create migration instance", "error", err)
		return
	}
	if err = migration.Up(); err != nil {
		if err != migrate.ErrNoChange {
			log.Error("Failed to migrate database", "error", err)
		} else {
			log.Info("No new migrations to apply")
		}
	} else {
		log.Info("Database migrated successfully")
	}

}
