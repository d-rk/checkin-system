package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/d-rk/checkin-system/internal/app"
	"github.com/d-rk/checkin-system/internal/database"
	"time"

	"github.com/jmoiron/sqlx"
	"gopkg.in/guregu/null.v4"
)

type Repository interface {
	ListUsers(ctx context.Context) ([]User, error)
	GetUserByID(ctx context.Context, id int64) (*User, error)
	GetUserByName(ctx context.Context, name string, excludeID *int64) (*User, error)
	GetUserByNameAndPassword(ctx context.Context, name, password string) (*User, error)
	GetUserByRfidUid(ctx context.Context, rfidUID string, excludeID int64) (*User, error)
	DeleteUser(ctx context.Context, id int64) error
	DeleteAllUsers(ctx context.Context) error
	SaveUser(ctx context.Context, user *User) (*User, error)
	UpdateUser(ctx context.Context, user *User) (*User, error)
	UpdateUserPassword(ctx context.Context, id int64, password string) error
	updateAdminPassword(ctx context.Context, password string) error
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
			return nil, app.NotFoundErr
		} else {
			return nil, err
		}
	}

	return &user, nil
}

func (r *repository) GetUserByName(ctx context.Context, name string, excludeID *int64) (*User, error) {

	user := User{}

	if err := r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE name = $1 and id != $2", name, excludeID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, app.NotFoundErr
		} else {
			return nil, err
		}
	}

	return &user, nil
}

func (r *repository) GetUserByNameAndPassword(ctx context.Context, name, password string) (*User, error) {

	user := User{}

	stmt, err := r.db.PrepareNamedContext(ctx, `SELECT * FROM users
		WHERE name = :name and password_digest = crypt(:password, password_digest)`)

	if err != nil {
		return nil, err
	}

	if err := stmt.GetContext(ctx, &user, map[string]interface{}{"name": name, "password": password}); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, app.NotFoundErr
		} else {
			return nil, err
		}
	}

	return &user, nil
}

func (r *repository) GetUserByRfidUid(ctx context.Context, rfidUID string, excludeID int64) (*User, error) {

	user := User{}

	if err := r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE rfid_uid = $1 and id != $2", rfidUID, excludeID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, app.NotFoundErr
		} else {
			return nil, err
		}
	}

	return &user, nil
}

func (r *repository) DeleteUser(ctx context.Context, id int64) error {

	return database.WithTransaction(r.db, func(tx database.Tx) error {

		deleteCheckinsStatement, err := r.db.PreparexContext(ctx, `DELETE FROM checkins WHERE user_id = $1`)

		if err != nil {
			return err
		}

		_, err = deleteCheckinsStatement.Exec(id)

		if err != nil {
			return err
		}

		deleteUserStatement, err := r.db.PreparexContext(ctx, `DELETE FROM users WHERE id = $1`)

		if err != nil {
			return err
		}

		_, err = deleteUserStatement.Exec(id)
		return err
	})
}

func (r *repository) DeleteAllUsers(ctx context.Context) error {

	return database.WithTransaction(r.db, func(tx database.Tx) error {

		deleteCheckinsStatement, err := r.db.PreparexContext(ctx, `DELETE FROM checkins`)

		if err != nil {
			return err
		}

		_, err = deleteCheckinsStatement.Exec()

		if err != nil {
			return err
		}

		deleteUserStatement, err := r.db.PreparexContext(ctx, `DELETE FROM users`)

		if err != nil {
			return err
		}

		_, err = deleteUserStatement.Exec()
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

	row := insertStatement.QueryRow(user)

	if row.Err() != nil {
		return nil, row.Err()
	}

	if err := row.Scan(&user.ID); err != nil {
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

	_ = updateStatement.MustExecContext(ctx, user)

	return user, nil
}

func (r *repository) UpdateUserPassword(ctx context.Context, id int64, password string) error {

	updateStatement, err := r.db.PrepareNamedContext(ctx, `UPDATE users SET
    		(updated_at, password_digest) =
            (current_timestamp, crypt(:password, gen_salt('bf')))
             WHERE id = :id and (password_digest is null or password_digest != crypt(:password, password_digest))`)

	if err != nil {
		return err
	}

	_, err = updateStatement.ExecContext(ctx, map[string]interface{}{"id": id, "password": password})
	return err
}

func (r *repository) updateAdminPassword(ctx context.Context, password string) error {

	updateStatement, err := r.db.PrepareNamedContext(ctx, `UPDATE users SET
    		(updated_at, password_digest) =
            (current_timestamp, crypt(:password, gen_salt('bf')))
             WHERE name = 'admin' and (password_digest is null or password_digest != crypt(:password, password_digest))`)

	if err != nil {
		return err
	}

	_, err = updateStatement.ExecContext(ctx, map[string]interface{}{"password": password})
	return err
}

func (r *repository) ListUserGroups(ctx context.Context) ([]string, error) {

	groups := make([]string, 0)
	rows, err := r.db.QueryContext(ctx, "SELECT distinct group_name FROM users where group_name is not null")

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return groups, nil
		} else {
			return nil, err
		}
	}

	for rows.Next() {
		var group string
		if err = rows.Scan(&group); err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}

	return groups, nil
}
