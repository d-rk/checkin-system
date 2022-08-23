package models

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/d-rk/checkin-system/pkg/services/database"
	"github.com/jmoiron/sqlx"
)

type CheckIn struct {
	ID        int64     `db:"id" json:"id"`
	Timestamp time.Time `db:"timestamp" json:"timestamp"`
	UserID    int64     `db:"user_id" json:"user_id"`
}

type CheckInWithUser struct {
	CheckIn
	User User `db:"user" json:"user"`
}

func ListCheckIns(db *sqlx.DB) ([]CheckIn, error) {

	checkIns := []CheckIn{}

	if err := db.Select(&checkIns, "SELECT * FROM checkins"); err != nil {
		return nil, errors.New("no checkins found")
	}

	return checkIns, nil
}

func ListCheckInsPerDay(db *sqlx.DB, day time.Time) ([]CheckInWithUser, error) {

	checkIns := []CheckInWithUser{}

	if err := db.Select(&checkIns, `SELECT
			checkins.*,
			users.id "user.id",
			users.name "user.name",
			users.created_at "user.created_at",
			users.updated_at "user.updated_at",
			users.rfid_uid "user.rfid_uid"
			FROM checkins JOIN users ON checkins.user_id = users.id
			WHERE checkins.timestamp >= $1 and checkins.timestamp < $1 + interval '1 day'
			ORDER BY checkins.timestamp ASC`, day); err != nil {
		return nil, fmt.Errorf("unable to query checkins: %s", err.Error())
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
