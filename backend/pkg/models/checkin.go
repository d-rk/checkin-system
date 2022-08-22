package models

import (
	"context"
	"errors"
	"time"

	"github.com/d-rk/checkin-system/pkg/services/database"
	"github.com/jmoiron/sqlx"
)

type CheckIn struct {
	ID        int64     `db:"id" json:"id"`
	Timestamp time.Time `db:"timestamp" json:"timestamp"`
	UserID    int64     `db:"user_id" json:"user_id"`
}

func ListCheckIns(db *sqlx.DB) ([]CheckIn, error) {

	checkIns := []CheckIn{}

	if err := db.Select(&checkIns, "SELECT * FROM checkins"); err != nil {
		return nil, errors.New("no checkins found")
	}

	return checkIns, nil
}

func ListUserCheckIns(db *sqlx.DB, userID int64) ([]CheckIn, error) {

	checkIns := []CheckIn{}

	if err := db.Select(&checkIns, "SELECT * FROM checkins WHERE user_id = $1", userID); err != nil {
		return nil, errors.New("no checkins found")
	}

	return checkIns, nil
}

func DeleteCheckInByID(db *sqlx.DB, ctx context.Context, id int64) error {

	return database.WithTransaction(db, func(tx database.Tx) error {

		deleteCheckinsStatement, err := db.PreparexContext(ctx, `DELETE FROM checkins WHERE user_id = $1`)

		if err != nil {
			return err
		}

		_, err = deleteCheckinsStatement.Exec(id)
		return err
	})
}

func DeleteCheckInsByUserID(db *sqlx.DB, ctx context.Context, userID int64) error {

	return database.WithTransaction(db, func(tx database.Tx) error {

		deleteCheckinsStatement, err := db.PreparexContext(ctx, `DELETE FROM checkins WHERE user_id = $1`)

		if err != nil {
			return err
		}

		_, err = deleteCheckinsStatement.Exec(userID)
		return err
	})
}

func (checkIn *CheckIn) Save(db *sqlx.DB, ctx context.Context) (*CheckIn, error) {

	insertStatement, err := db.PrepareNamedContext(ctx, `INSERT INTO checkins
		(timestamp, user_id) VALUES
		(:timestamp, :user_id) RETURNING id`)

	if err != nil {
		return nil, err
	}

	row := insertStatement.QueryRow(checkIn)

	if row.Err() != nil {
		return nil, row.Err()
	}

	if err := row.Scan(&checkIn.ID); err != nil {
		return nil, err
	}

	return checkIn, nil
}
