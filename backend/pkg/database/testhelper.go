//go:build integration

package database

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/joho/godotenv"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// SetupTestDB creates a new test database, runs migrations, and returns the DB, dbName, and a cleanup function.
func SetupTestDB(t *testing.T) *sqlx.DB {
	t.Helper()

	if err := setWorkingDir(); err != nil {
		t.Fatalf("failed to set working dir: %v", err)
	}

	_ = godotenv.Load(".env")
	adminDB := Connect()

	baseDB := os.Getenv("DB_NAME")
	dbName := fmt.Sprintf("%s_it_%s", baseDB, uuid.New().String()[:8])

	_, err := adminDB.Exec(fmt.Sprintf("CREATE DATABASE \"%s\"", dbName))
	if err != nil {
		_ = adminDB.Close()
		t.Fatalf("failed to create test db: %v", err)
	}
	_ = adminDB.Close()

	testDB := connect(dbName)
	RunMigration(testDB)

	t.Cleanup(func() {
		_ = testDB.Close()
		adminDB := Connect()
		_, _ = adminDB.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS \"%s\" WITH (FORCE)", dbName))
		_ = adminDB.Close()
	})

	return testDB
}

func setWorkingDir() error {

	path, err := os.Getwd()
	if err != nil {
		return errors.New("unable to get working dir")
	}

	_, err = os.Stat(filepath.Join(path, ".env"))

	for os.IsNotExist(err) {
		if strings.HasSuffix(path, "backend") {
			return errors.New("unable to find .env in backend folder")
		}
		oldPath := path

		if path = filepath.Dir(oldPath); path == oldPath {
			return errors.New("unable to find .env file")
		}
		_, err = os.Stat(filepath.Join(path, ".env"))
	}

	return os.Chdir(path)
}
