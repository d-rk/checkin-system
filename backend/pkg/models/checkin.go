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
	ID        int64     `db:"id" json:"id" csv:"id"`
	Date      time.Time `db:"date" json:"date" csv:"date"`
	Timestamp time.Time `db:"timestamp" json:"timestamp" csv:"timestamp"`
	UserID    int64     `db:"user_id" json:"user_id" csv:"-"`
}

type CheckInWithUser struct {
	CheckIn
	User User `db:"user" json:"user" csv:"-"`
}

type CheckInDate struct {
	Date time.Time `db:"date" json:"date"`
}

func ListCheckIns(db *sqlx.DB) ([]CheckIn, error) {

	checkIns := []CheckIn{}

	if err := db.Select(&checkIns, "SELECT * FROM checkins"); err != nil {
		return nil, errors.New("no checkins found")
	}

	return checkIns, nil
}

func ListCheckInsPerDay(db *sqlx.DB, date time.Time) ([]CheckInWithUser, error) {

	checkIns := []CheckInWithUser{}

	if err := db.Select(&checkIns, `SELECT
			checkins.*,
			users.id "user.id",
			users.name "user.name",
			users.created_at "user.created_at",
			users.updated_at "user.updated_at",
			users.member_id "user.member_id",
			users.rfid_uid "user.rfid_uid"
			FROM checkins JOIN users ON checkins.user_id = users.id
			WHERE checkins.date = $1
			ORDER BY checkins.timestamp ASC`, date); err != nil {
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
		(date, timestamp, user_id) VALUES
		(:date, :timestamp, :user_id) RETURNING id`)

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

func ListCheckInDates(db *sqlx.DB) ([]CheckInDate, error) {

	dates := []CheckInDate{}

	if err := db.Select(&dates, "SELECT distinct date as date FROM checkins"); err != nil {
		return nil, errors.New("no checkIn dates found")
	}

	return dates, nil
}
