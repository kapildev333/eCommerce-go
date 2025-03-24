package db

import (
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var DB *sql.DB

func InitDB() {
	var err error
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://user:password@localhost/dbname?sslmode=disable"
	}
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal(err)
	}

	migration, err := migrate.New("file://db/migration", connStr)
	if err != nil {
		log.Fatal(err)
	}
	if err = migration.Up(); err != nil {
		if err.Error() == "no change" {
			log.Println("No migration to run")
		} else {
			log.Fatal(err)
		}
	}

}
