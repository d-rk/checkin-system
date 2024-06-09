package user

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/d-rk/checkin-system/internal/app"
	"github.com/d-rk/checkin-system/internal/websocket"
	"golang.org/x/crypto/bcrypt"
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
	ListUserGroups(ctx context.Context) ([]string, error)
}

type service struct {
	repo      Repository
	websocket *websocket.Server
}

func NewService(repo Repository, websocket *websocket.Server) Service {

	adminPassword := os.Getenv("ADMIN_PASSWORD")

	service := &service{repo, websocket}
	if err := service.updateAdminPassword(context.Background(), adminPassword); err != nil {
		panic(err)
	}

	return service
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

	user, err := s.repo.GetUserByName(ctx, name, -1)
	if err != nil {
		return nil, err
	}

	if s.passwordEquals(user, password) {
		return user, nil
	} else {
		return nil, app.NotFoundErr
	}
}

func (s *service) UpdateUser(ctx context.Context, user *User) (*User, error) {

	_, err := s.repo.GetUserByID(ctx, user.ID)
	if err != nil {
		return nil, app.NotFoundErr
	}

	_, err = s.repo.GetUserByName(ctx, user.Name, user.ID)

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

	_, err := s.repo.GetUserByName(ctx, user.Name, -1)

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

	if password == "" {
		return fmt.Errorf("empty password provided: %w", app.ConflictErr)
	}

	user, err := s.repo.GetUserByID(context.Background(), id)
	if err != nil {
		return err
	}

	return s.updateUserPassword(ctx, user, password)
}

func (s *service) ListUserGroups(ctx context.Context) ([]string, error) {
	return s.repo.ListUserGroups(ctx)
}

func (s *service) updateAdminPassword(ctx context.Context, password string) error {

	if password == "" {
		return nil
	}

	admin, err := s.repo.GetUserByName(ctx, "admin", -1)
	if err != nil && errors.Is(err, app.NotFoundErr) {
		log.Printf("not updating admin password. user not found")
		return nil
	} else if err != nil {
		return err
	}

	return s.updateUserPassword(ctx, admin, password)
}

func (s *service) updateUserPassword(ctx context.Context, user *User, password string) error {

	if s.passwordEquals(user, password) {
		return nil
	}

	digest, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.repo.UpdateUserPasswordDigest(ctx, user.ID, string(digest))
}

func (s *service) passwordEquals(user *User, password string) bool {

	if user.PasswordDigest.IsZero() {
		return true
	}
	return bcrypt.CompareHashAndPassword([]byte(user.PasswordDigest.String), []byte(password)) != nil
}
