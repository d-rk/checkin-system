package checkin

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/d-rk/checkin-system/pkg/app"
	"github.com/d-rk/checkin-system/pkg/database"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	ListCheckIns(ctx context.Context) ([]CheckIn, error)
	ListCheckInsPerDay(ctx context.Context, date time.Time) ([]WithUser, error)
	ListAllCheckIns(ctx context.Context) ([]WithUser, error)
	ListUserCheckIns(ctx context.Context, userID int64) ([]CheckIn, error)
	GetLatestCheckinDate(ctx context.Context) (*time.Time, error)
	DeleteCheckInByID(ctx context.Context, id int64) error
	DeleteCheckInsByUserID(ctx context.Context, userID int64) error
	DeleteCheckInsOlderThan(ctx context.Context, thresholdDays int64) error
	SaveCheckIn(ctx context.Context, checkIn *CheckIn) (*CheckIn, error)
	ListCheckInDates(ctx context.Context) ([]Date, error)
}

type repository struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) Repository {
	return &repository{db}
}

func (r *repository) ListCheckIns(ctx context.Context) ([]CheckIn, error) {

	var checkIns []CheckIn

	if err := r.db.SelectContext(ctx, &checkIns, "SELECT * FROM checkins"); err != nil {
		return nil, errors.New("no checkins found")
	}

	return checkIns, nil
}

func (r *repository) ListCheckInsPerDay(ctx context.Context, date time.Time) ([]WithUser, error) {

	var checkIns []WithUser

	if err := r.db.SelectContext(ctx, &checkIns, `SELECT
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

func (r *repository) ListAllCheckIns(ctx context.Context) ([]WithUser, error) {

	var checkIns []WithUser

	if err := r.db.SelectContext(ctx, &checkIns, `SELECT
			checkins.*,
			users.id "user.id",
			users.name "user.name",
			users.created_at "user.created_at",
			users.updated_at "user.updated_at",
			users.member_id "user.member_id",
			users.rfid_uid "user.rfid_uid"
			FROM checkins JOIN users ON checkins.user_id = users.id
			ORDER BY checkins.timestamp ASC`); err != nil {
		return nil, fmt.Errorf("unable to query checkins: %s", err.Error())
	}

	return checkIns, nil
}

func (r *repository) ListUserCheckIns(ctx context.Context, userID int64) ([]CheckIn, error) {

	var checkIns []CheckIn

	if err := r.db.SelectContext(ctx, &checkIns, "SELECT * FROM checkins WHERE user_id = $1", userID); err != nil {
		return nil, errors.New("no checkins found")
	}

	return checkIns, nil
}

func (r *repository) GetLatestCheckinDate(ctx context.Context) (*time.Time, error) {

	var timestamp sql.NullTime

	if err := r.db.GetContext(ctx, &timestamp, "SELECT max(timestamp) FROM checkins"); err != nil {
		return nil, err
	}

	if !timestamp.Valid {
		return nil, app.ErrNotFound
	}

	return &timestamp.Time, nil
}

func (r *repository) DeleteCheckInByID(ctx context.Context, id int64) error {

	return database.WithTransaction(r.db, func(_ database.Tx) error {

		deleteCheckinsStatement, err := r.db.PreparexContext(ctx, `DELETE FROM checkins WHERE id = $1`)
		if err != nil {
			return err
		}
		defer deleteCheckinsStatement.Close()

		_, err = deleteCheckinsStatement.ExecContext(ctx, id)
		return err
	})
}

func (r *repository) DeleteCheckInsByUserID(ctx context.Context, userID int64) error {

	return database.WithTransaction(r.db, func(_ database.Tx) error {

		deleteCheckinsStatement, err := r.db.PreparexContext(ctx, `DELETE FROM checkins WHERE user_id = $1`)
		if err != nil {
			return err
		}
		defer deleteCheckinsStatement.Close()

		_, err = deleteCheckinsStatement.ExecContext(ctx, userID)
		return err
	})
}

func (r *repository) DeleteCheckInsOlderThan(ctx context.Context, thresholdDays int64) error {

	return database.WithTransaction(r.db, func(_ database.Tx) error {

		var query string

		switch r.db.DriverName() {
		case "postgres":
			query = `DELETE FROM checkins WHERE DATE_PART('day', now() - date) > $1`
		case "sqlite3":
			query = `DELETE FROM checkins WHERE julianday('now') - julianday(date) > $1`
		default:
			return fmt.Errorf("unknown driver %s", r.db.DriverName())
		}

		deleteCheckinsStatement, err := r.db.PreparexContext(ctx, query)
		if err != nil {
			return err
		}
		defer deleteCheckinsStatement.Close()

		_, err = deleteCheckinsStatement.ExecContext(ctx, thresholdDays)
		return err
	})
}

func (r *repository) SaveCheckIn(ctx context.Context, checkIn *CheckIn) (*CheckIn, error) {

	insertStatement, err := r.db.PrepareNamedContext(ctx, `INSERT INTO checkins
		(date, timestamp, user_id) VALUES
		(:date, :timestamp, :user_id) RETURNING id`)
	if err != nil {
		return nil, err
	}
	defer insertStatement.Close()

	row := insertStatement.QueryRow(checkIn)

	if row.Err() != nil {
		return nil, row.Err()
	}

	if err = row.Scan(&checkIn.ID); err != nil {
		return nil, err
	}

	return checkIn, nil
}

func (r *repository) ListCheckInDates(ctx context.Context) ([]Date, error) {

	var dates []Date

	if err := r.db.SelectContext(ctx, &dates, "SELECT distinct date as date FROM checkins"); err != nil {
		return nil, errors.New("no checkIn dates found")
	}

	return dates, nil
}
