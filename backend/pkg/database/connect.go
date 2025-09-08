package database

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"           //revive: postgres driver
	_ "github.com/mattn/go-sqlite3" //revive: sqlite3 driver
	migrate "github.com/rubenv/sql-migrate"
)

func Connect() *sqlx.DB {
	return connect(os.Getenv("DB_NAME"))
}

func connect(dbName string) *sqlx.DB {

	dbDriver := os.Getenv("DB_DRIVER")

	var dsn string

	switch dbDriver {
	case "postgres":
		dbHost := os.Getenv("DB_HOST")
		dbUser := os.Getenv("DB_USER")
		dbPassword := os.Getenv("DB_PASSWORD")
		dbPort := os.Getenv("DB_PORT")
		dbSslMode := os.Getenv("DB_SSL_MODE")

		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
			dbHost,
			dbUser,
			dbPassword,
			dbName,
			dbPort,
			dbSslMode,
		)
	case "sqlite3":
		dsn = fmt.Sprintf("file:%s?_loc=UTC", dbName)
	}

	db, err := sqlx.Connect(dbDriver, dsn)

	if err != nil {
		slog.Error("cannot connect to database", "dsn", dsn, "error", err)
		os.Exit(1)
	}

	slog.Info("connected to database", "dsn", dsn, "driver", dbDriver)
	return db
}

func RunMigration(db *sqlx.DB) {
	migrations := &migrate.FileMigrationSource{
		Dir: "db/migrations/" + db.DriverName(),
	}

	migrate.SetTable("migrations")

	n, err := migrate.Exec(db.DB, db.DriverName(), migrations, migrate.Up)
	if err != nil {
		slog.Error("migration failed", "error", err)
		os.Exit(1)
	}
	slog.Info("migration applied", "count", n)
}
