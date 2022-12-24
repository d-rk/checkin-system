package models

import (
	"context"
	"errors"
	"time"

	"github.com/d-rk/checkin-system/pkg/services/database"
	"github.com/jmoiron/sqlx"
	"gopkg.in/guregu/null.v4"
)

type User struct {
	ID        int64       `db:"id" json:"id" csv:"-"`
	CreatedAt time.Time   `db:"created_at" json:"created_at" csv:"-"`
	UpdatedAt null.Time   `db:"updated_at" json:"updated_at" csv:"-"`
	Name      string      `db:"name" json:"name"  csv:"name"`
	MemberID  null.String `db:"member_id" json:"member_id"  csv:"member_id"`
	RFIDuid   string      `db:"rfid_uid" json:"rfid_uid" csv:"rfid_uid"`
}

func ListUsers(db *sqlx.DB) ([]User, error) {

	users := []User{}

	if err := db.Select(&users, "SELECT * FROM users"); err != nil {
		return nil, errors.New("no users found")
	}

	return users, nil
}

func GetUserByID(db *sqlx.DB, uid int64) (User, error) {

	user := User{}

	if err := db.Get(&user, "SELECT * FROM users WHERE id = $1", uid); err != nil {
		return user, errors.New("user not found")
	}

	return user, nil
}

func GetUserByName(db *sqlx.DB, name string, excludeID int64) (User, error) {

	user := User{}

	if err := db.Get(&user, "SELECT * FROM users WHERE name = $1 and id != $2", name, excludeID); err != nil {
		return user, errors.New("user not found")
	}

	return user, nil
}

func GetUserByRfidUid(db *sqlx.DB, rfidUID string, excludeID int64) (User, error) {

	user := User{}

	if err := db.Get(&user, "SELECT * FROM users WHERE rfid_uid = $1 and id != $2", rfidUID, excludeID); err != nil {
		return user, errors.New("user not found")
	}

	return user, nil
}

func DeleteUser(db *sqlx.DB, ctx context.Context, id int64) error {

	return database.WithTransaction(db, func(tx database.Tx) error {

		deleteCheckinsStatement, err := db.PreparexContext(ctx, `DELETE FROM checkins WHERE user_id = $1`)

		if err != nil {
			return err
		}

		_, err = deleteCheckinsStatement.Exec(id)

		if err != nil {
			return err
		}

		deleteUserStatement, err := db.PreparexContext(ctx, `DELETE FROM users WHERE id = $1`)

		if err != nil {
			return err
		}

		_, err = deleteUserStatement.Exec(id)
		return err
	})
}

func DeleteAllUsers(db *sqlx.DB, ctx context.Context) error {

	return database.WithTransaction(db, func(tx database.Tx) error {

		deleteCheckinsStatement, err := db.PreparexContext(ctx, `DELETE FROM checkins`)

		if err != nil {
			return err
		}

		_, err = deleteCheckinsStatement.Exec()

		if err != nil {
			return err
		}

		deleteUserStatement, err := db.PreparexContext(ctx, `DELETE FROM users`)

		if err != nil {
			return err
		}

		_, err = deleteUserStatement.Exec()
		return err
	})
}

func (user *User) Insert(db *sqlx.DB, ctx context.Context) (*User, error) {

	user.CreatedAt = time.Now()

	insertStatement, err := db.PrepareNamedContext(ctx, `INSERT INTO users
    		(created_at, name, rfid_uid, member_id) VALUES
            (:created_at, :name,:rfid_uid, :member_id) RETURNING id`)

	if err != nil {
		return nil, err
	}

	row := insertStatement.QueryRow(user)

	if row.Err() != nil {
		return nil, row.Err()
	}

	if err := row.Scan(&user.ID); err != nil {
		return nil, err
	}

	return user, nil
}

func (user *User) Update(db *sqlx.DB, ctx context.Context) (*User, error) {

	user.UpdatedAt = null.TimeFrom(time.Now())

	updateStatement, err := db.PrepareNamedContext(ctx, `UPDATE users SET
    		(updated_at, name, rfid_uid, member_id) =
            (:updated_at, :name,:rfid_uid, :member_id) WHERE id = :id`)

	if err != nil {
		return nil, err
	}

	_ = updateStatement.MustExecContext(ctx, user)

	return user, nil
}
