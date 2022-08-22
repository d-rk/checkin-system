package database

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
	"log"
	"os"
)

func Connect() *sqlx.DB {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	dbDriver := "postgres"
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dbSslMode := os.Getenv("DB_SSL_MODE")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC", dbHost, dbUser, dbPassword, dbName, dbPort, dbSslMode)
	db, err := sqlx.Connect(dbDriver, dsn)

	if err != nil {
		log.Println(fmt.Sprintf("cannot connect to database %s", dbName))
		log.Fatal("connection error:", err)
	} else {
		log.Println(fmt.Sprintf("connected to database %s", dbName))
	}

	err = runMigration(db.DB)

	if err != nil {
		log.Fatal("migration failed:", err)
	}

	return db
}

func runMigration(db *sql.DB) error {
	migrations := &migrate.FileMigrationSource{
		Dir: "db/migrations",
	}

	migrate.SetTable("migrations")

	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		return err
	}
	fmt.Printf("applied %d migrations\n", n)
	return nil
}
