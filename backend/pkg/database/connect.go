package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	migrate "github.com/rubenv/sql-migrate"
	"log"
	"os"
)

func Connect() *sqlx.DB {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	dbDriver := os.Getenv("DB_DRIVER")

	var dsn string

	if dbDriver == "postgres" {
		dbHost := os.Getenv("DB_HOST")
		dbUser := os.Getenv("DB_USER")
		dbPassword := os.Getenv("DB_PASSWORD")
		dbName := os.Getenv("DB_NAME")
		dbPort := os.Getenv("DB_PORT")
		dbSslMode := os.Getenv("DB_SSL_MODE")

		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC", dbHost, dbUser, dbPassword, dbName, dbPort, dbSslMode)

	} else if dbDriver == "sqlite3" {
		dbName := os.Getenv("DB_NAME")
		dsn = fmt.Sprintf("file:%s?_loc=UTC", dbName)
	}

	db, err := sqlx.Connect(dbDriver, dsn)

	if err != nil {
		log.Println(fmt.Sprintf("cannot connect to database: %s", dsn))
		log.Fatal("connection error:", err)
	} else {
		log.Println(fmt.Sprintf("connected to database %s", dsn))
	}

	err = runMigration(db)

	if err != nil {
		log.Fatal("migration failed:", err)
	}

	return db
}

func runMigration(db *sqlx.DB) error {
	migrations := &migrate.FileMigrationSource{
		Dir: "db/migrations/" + db.DriverName(),
	}

	migrate.SetTable("migrations")

	n, err := migrate.Exec(db.DB, db.DriverName(), migrations, migrate.Up)
	if err != nil {
		return err
	}
	fmt.Printf("applied %d migrations\n", n)
	return nil
}
