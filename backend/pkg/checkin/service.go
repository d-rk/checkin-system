package checkin

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/d-rk/checkin-system/pkg/app"
	"github.com/d-rk/checkin-system/pkg/user"
	"github.com/d-rk/checkin-system/pkg/websocket"
	"github.com/lib/pq"
)

const hoursInDay = 24
const daysInYear = 356

type Service interface {
	ListCheckIns(ctx context.Context) ([]CheckIn, error)
	ListAllCheckIns(ctx context.Context) ([]WithUser, error)
	DeleteCheckInByID(ctx context.Context, checkinID int64) error
	DeleteCheckInsByUserID(ctx context.Context, userID int64) error
	DeleteOldCheckIns(ctx context.Context) error
	CreateCheckInForUser(ctx context.Context, userID int64, timestamp *time.Time) (*CheckIn, error)
	CreateCheckInForRFID(ctx context.Context, rfidUID string, timestamp *time.Time) (*CheckIn, error)
	ListCheckInsPerDay(ctx context.Context, day time.Time) ([]WithUser, error)
	ListCheckInDates(ctx context.Context) ([]Date, error)
	ListUserCheckIns(ctx context.Context, userID int64) ([]CheckIn, error)
}

type service struct {
	repo        Repository
	userService user.Service
	websocket   *websocket.Server
}

func NewService(repo Repository, userService user.Service, websocket *websocket.Server) Service {

	return &service{repo, userService, websocket}
}

func (s *service) ListCheckIns(ctx context.Context) ([]CheckIn, error) {
	return s.repo.ListCheckIns(ctx)
}

func (s *service) ListAllCheckIns(ctx context.Context) ([]WithUser, error) {
	return s.repo.ListAllCheckIns(ctx)
}

func (s *service) ListCheckInsPerDay(ctx context.Context, day time.Time) ([]WithUser, error) {
	return s.repo.ListCheckInsPerDay(ctx, day)
}

func (s *service) ListCheckInDates(ctx context.Context) ([]Date, error) {
	return s.repo.ListCheckInDates(ctx)
}

func (s *service) ListUserCheckIns(ctx context.Context, userID int64) ([]CheckIn, error) {
	return s.repo.ListUserCheckIns(ctx, userID)
}

func (s *service) CreateCheckInForUser(ctx context.Context, userID int64, timestamp *time.Time) (*CheckIn, error) {

	checkinTimestamp := time.Now()
	if timestamp != nil {
		checkinTimestamp = *timestamp
	}

	u, err := s.userService.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.createCheckinForUser(ctx, u, checkinTimestamp)
}

func (s *service) CreateCheckInForRFID(ctx context.Context, rfidUID string, timestamp *time.Time) (*CheckIn, error) {

	websocketMessage := WebsocketMessage{}
	websocketMessage.RFIDuid = rfidUID

	checkinTimestamp := time.Now()
	if timestamp != nil {
		checkinTimestamp = *timestamp
	}

	u, err := s.userService.GetUserByRfidUID(ctx, rfidUID, -1)

	if err != nil && errors.Is(err, app.ErrNotFound) {
		_ = s.websocket.Publish(websocketMessage)
		return nil, err
	} else if err != nil {
		return nil, err
	}

	checkin, err := s.createCheckinForUser(ctx, u, checkinTimestamp)
	if err != nil {
		return nil, err
	}

	websocketMessage.CheckIn = checkin
	_ = s.websocket.Publish(websocketMessage)

	return checkin, nil
}

func (s *service) createCheckinForUser(ctx context.Context, user *user.User, timestamp time.Time) (*CheckIn, error) {

	checkIn := CheckIn{
		ID:        -1,
		Date:      truncateToStartOfDay(timestamp),
		Timestamp: timestamp,
		UserID:    user.ID,
	}

	savedCheckIn, err := s.repo.SaveCheckIn(ctx, &checkIn)

	var pgErr *pq.Error
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return nil, fmt.Errorf("checkIn for day already exists: %w", app.ErrConflict)
		}
	}

	return savedCheckIn, err
}

func (s *service) DeleteCheckInByID(ctx context.Context, checkinID int64) error {
	return s.repo.DeleteCheckInByID(ctx, checkinID)
}

func (s *service) DeleteCheckInsByUserID(ctx context.Context, userID int64) error {
	return s.repo.DeleteCheckInsByUserID(ctx, userID)
}

func (s *service) DeleteOldCheckIns(ctx context.Context) error {

	latestTimestamp, err := s.repo.GetLatestCheckinDate(ctx)
	if err != nil && !errors.Is(err, app.ErrNotFound) {
		return err
	} else if err == nil {
		now := time.Now()
		diffDays := now.Sub(*latestTimestamp).Hours() / hoursInDay
		if diffDays < -1 {
			return fmt.Errorf("%w: too far in past - now (%v) < latest checkin (%v)", app.ErrInternal, now,
				*latestTimestamp)
		} else if diffDays > daysInYear {
			return fmt.Errorf("%w: too far in future - now (%v) > latest checkin (%v)", app.ErrInternal, now,
				*latestTimestamp)
		}
	}

	retentionDaysEnv := os.Getenv("CHECKIN_RETENTION_DAYS")
	if retentionDaysEnv == "" {
		retentionDaysEnv = strconv.Itoa(daysInYear)
	}

	retentionDays, err := strconv.ParseInt(retentionDaysEnv, 10, 64)
	if err != nil {
		return err
	}

	return s.repo.DeleteCheckInsOlderThan(ctx, retentionDays)
}

func truncateToStartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}
