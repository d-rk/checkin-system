package user

import (
	"context"
	"fmt"
	"github.com/d-rk/checkin-system/internal/app"
	"github.com/d-rk/checkin-system/internal/websocket"
	"os"
)

type Service interface {
	ListUsers(ctx context.Context) ([]User, error)
	GetUserByNameAndPassword(ctx context.Context, name, password string) (*User, error)
	GetUserByRfidUid(ctx context.Context, rfidUid string, excludeID int64) (*User, error)
	GetUserByID(ctx context.Context, id int64) (*User, error)
	CreateUser(ctx context.Context, user *User) (*User, error)
	UpdateUser(ctx context.Context, user *User) (*User, error)
	UpdateUserPassword(ctx context.Context, id int64, password string) error
	DeleteUser(ctx context.Context, id int64) error
	DeleteAllUsers(ctx context.Context) error
}

type service struct {
	repo      Repository
	websocket *websocket.Server
}

func NewService(repo Repository, websocket *websocket.Server) Service {

	adminPassword := os.Getenv("ADMIN_PASSWORD")

	if adminPassword != "" {
		if err := repo.updateAdminPassword(context.Background(), adminPassword); err != nil {
			panic(err)
		}
	}
	return &service{repo, websocket}
}

func (s *service) ListUsers(ctx context.Context) ([]User, error) {
	return s.repo.ListUsers(ctx)
}

func (s *service) GetUserByID(ctx context.Context, id int64) (*User, error) {
	return s.repo.GetUserByID(ctx, id)
}

func (s *service) GetUserByRfidUid(ctx context.Context, rfidUid string, excludeID int64) (*User, error) {
	return s.repo.GetUserByRfidUid(ctx, rfidUid, excludeID)
}

func (s *service) GetUserByNameAndPassword(ctx context.Context, name, password string) (*User, error) {
	return s.repo.GetUserByNameAndPassword(ctx, name, password)
}

func (s *service) UpdateUser(ctx context.Context, user *User) (*User, error) {

	user, err := s.repo.GetUserByID(ctx, user.ID)
	if err != nil {
		return nil, app.NotFoundErr
	}

	_, err = s.repo.GetUserByName(ctx, user.Name, &user.ID)

	if err == nil {
		return nil, fmt.Errorf("user with name already exists: %w", app.ConflictErr)
	}

	if user.RFIDuid.Ptr() != nil {
		_, err = s.repo.GetUserByRfidUid(ctx, user.RFIDuid.ValueOrZero(), user.ID)

		if err == nil {
			return nil, fmt.Errorf("user with rfid_uid already exists: %w", app.ConflictErr)
		}
	}

	return s.repo.UpdateUser(ctx, user)
}

func (s *service) DeleteUser(ctx context.Context, id int64) error {
	return s.repo.DeleteUser(ctx, id)
}

func (s *service) DeleteAllUsers(ctx context.Context) error {
	return s.repo.DeleteAllUsers(ctx)
}

func (s *service) CreateUser(ctx context.Context, user *User) (*User, error) {

	_, err := s.repo.GetUserByName(ctx, user.Name, nil)

	if err == nil {
		return nil, fmt.Errorf("user already exists: %w", app.ConflictErr)
	}

	if user.RFIDuid.Ptr() != nil {
		_, err = s.repo.GetUserByRfidUid(ctx, user.RFIDuid.ValueOrZero(), -1)

		if err == nil {
			return nil, fmt.Errorf("user with rfid_uid already exists: %w", app.ConflictErr)
		}
	}

	return s.repo.SaveUser(ctx, user)
}

func (s *service) UpdateUserPassword(ctx context.Context, id int64, password string) error {
	return s.repo.UpdateUserPassword(ctx, id, password)
}
