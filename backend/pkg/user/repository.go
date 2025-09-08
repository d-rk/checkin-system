package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/d-rk/checkin-system/pkg/app"
	"github.com/d-rk/checkin-system/pkg/database"

	"github.com/jmoiron/sqlx"
	"gopkg.in/guregu/null.v4"
)

type Repository interface {
	ListUsers(ctx context.Context) ([]User, error)
	GetUserByID(ctx context.Context, id int64) (*User, error)
	GetUserByName(ctx context.Context, name string, excludeID int64) (*User, error)
	GetUserByRfidUID(ctx context.Context, rfidUID string, excludeID int64) (*User, error)
	DeleteUser(ctx context.Context, id int64) error
	DeleteAllUsers(ctx context.Context) error
	SaveUser(ctx context.Context, user *User) (*User, error)
	UpdateUser(ctx context.Context, user *User) (*User, error)
	UpdateUserPasswordDigest(ctx context.Context, id int64, passwordDigest string) error
	ListUserGroups(ctx context.Context) ([]string, error)
}

type repository struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) Repository {
	return &repository{db}
}

func (r *repository) ListUsers(ctx context.Context) ([]User, error) {

	var users []User

	if err := r.db.SelectContext(ctx, &users, "SELECT * FROM users"); err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return users, nil
}

func (r *repository) GetUserByID(ctx context.Context, uid int64) (*User, error) {

	user := User{}

	if err := r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE id = $1", uid); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, app.ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *repository) GetUserByName(ctx context.Context, name string, excludeID int64) (*User, error) {

	user := User{}

	if err := r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE name = $1 and id != $2", name, excludeID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, app.ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *repository) GetUserByRfidUID(ctx context.Context, rfidUID string, excludeID int64) (*User, error) {

	user := User{}

	if err := r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE rfid_uid = $1 and id != $2", rfidUID, excludeID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, app.ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *repository) DeleteUser(ctx context.Context, id int64) error {

	return database.WithTransaction(r.db, func(_ database.Tx) error {

		deleteCheckinsStatement, err := r.db.PreparexContext(ctx, `DELETE FROM checkins WHERE user_id = $1`)
		if err != nil {
			return err
		}
		defer deleteCheckinsStatement.Close()

		_, err = deleteCheckinsStatement.ExecContext(ctx, id)

		if err != nil {
			return err
		}

		deleteUserStatement, err := r.db.PreparexContext(ctx, `DELETE FROM users WHERE id = $1`)
		if err != nil {
			return err
		}
		defer deleteUserStatement.Close()

		_, err = deleteUserStatement.ExecContext(ctx, id)
		return err
	})
}

func (r *repository) DeleteAllUsers(ctx context.Context) error {

	return database.WithTransaction(r.db, func(_ database.Tx) error {

		deleteCheckinsStatement, err := r.db.PreparexContext(ctx, `DELETE FROM checkins`)
		if err != nil {
			return err
		}
		defer deleteCheckinsStatement.Close()

		_, err = deleteCheckinsStatement.ExecContext(ctx)

		if err != nil {
			return err
		}

		deleteUserStatement, err := r.db.PreparexContext(ctx, `DELETE FROM users`)
		if err != nil {
			return err
		}
		defer deleteUserStatement.Close()

		_, err = deleteUserStatement.ExecContext(ctx)
		return err
	})
}

func (r *repository) SaveUser(ctx context.Context, user *User) (*User, error) {

	user.CreatedAt = time.Now()

	insertStatement, err := r.db.PrepareNamedContext(ctx, `INSERT INTO users
    		(created_at, name, rfid_uid, member_id, role, group_name) VALUES
            (:created_at, :name,:rfid_uid, :member_id, :role, :group_name) RETURNING id`)
	if err != nil {
		return nil, err
	}
	defer insertStatement.Close()

	row := insertStatement.QueryRowContext(ctx, user)

	if row.Err() != nil {
		return nil, row.Err()
	}

	if err = row.Scan(&user.ID); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *repository) UpdateUser(ctx context.Context, user *User) (*User, error) {

	user.UpdatedAt = null.TimeFrom(time.Now())

	updateStatement, err := r.db.PrepareNamedContext(ctx, `UPDATE users SET
    		(updated_at, name, rfid_uid, member_id, role, group_name) =
            (:updated_at, :name,:rfid_uid, :member_id, :role, :group_name) WHERE id = :id`)
	if err != nil {
		return nil, err
	}
	defer updateStatement.Close()

	_ = updateStatement.MustExecContext(ctx, user)

	return user, nil
}

func (r *repository) UpdateUserPasswordDigest(ctx context.Context, id int64, passwordDigest string) error {

	updateStatement, err := r.db.PrepareNamedContext(ctx, `UPDATE users SET
    		(updated_at, password_digest) = (current_timestamp, :passwordDigest)
             WHERE id = :id`)
	if err != nil {
		return err
	}
	defer updateStatement.Close()

	_, err = updateStatement.ExecContext(ctx, map[string]interface{}{"id": id, "passwordDigest": passwordDigest})
	return err
}

func (r *repository) ListUserGroups(ctx context.Context) ([]string, error) {

	groups := make([]string, 0)
	rows, err := r.db.QueryContext(ctx, "SELECT distinct group_name FROM users where group_name is not null")

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return groups, nil
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var group string
		if err = rows.Scan(&group); err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return groups, nil
}
