package checkin

import (
	"context"
	"errors"
	"github.com/d-rk/checkin-system/internal/app"
	"github.com/d-rk/checkin-system/internal/user"
	"github.com/d-rk/checkin-system/internal/websocket"
	"os"
	"strconv"
	"time"
)

type Service interface {
	ListCheckins(ctx context.Context) ([]CheckIn, error)
	ListAllCheckins(ctx context.Context) ([]CheckInWithUser, error)
	DeleteCheckinByID(ctx context.Context, checkinID int64) error
	DeleteCheckInsByUserID(ctx context.Context, userID int64) error
	DeleteOldCheckIns(ctx context.Context) error
	CreateCheckinForUser(ctx context.Context, userID int64, timestamp *time.Time) (*CheckIn, error)
	CreateCheckinForRFID(ctx context.Context, rfidUid string, timestamp *time.Time) (*CheckIn, error)
	ListCheckinsPerDay(ctx context.Context, day time.Time) ([]CheckInWithUser, error)
	ListCheckinDates(ctx context.Context) ([]CheckInDate, error)
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

func (s *service) ListCheckins(ctx context.Context) ([]CheckIn, error) {
	return s.repo.ListCheckIns(ctx)
}

func (s *service) ListAllCheckins(ctx context.Context) ([]CheckInWithUser, error) {
	return s.repo.ListAllCheckIns(ctx)
}

func (s *service) ListCheckinsPerDay(ctx context.Context, day time.Time) ([]CheckInWithUser, error) {
	return s.repo.ListCheckInsPerDay(ctx, day)
}

func (s *service) ListCheckinDates(ctx context.Context) ([]CheckInDate, error) {
	return s.repo.ListCheckInDates(ctx)
}

func (s *service) ListUserCheckIns(ctx context.Context, userID int64) ([]CheckIn, error) {
	return s.repo.ListUserCheckIns(ctx, userID)
}

func (s *service) CreateCheckinForUser(ctx context.Context, userID int64, timestamp *time.Time) (*CheckIn, error) {

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

func (s *service) CreateCheckinForRFID(ctx context.Context, rfidUid string, timestamp *time.Time) (*CheckIn, error) {

	websocketMessage := WebsocketMessage{}
	websocketMessage.RFIDuid = rfidUid

	checkinTimestamp := time.Now()
	if timestamp != nil {
		checkinTimestamp = *timestamp
	}

	u, err := s.userService.GetUserByRfidUid(ctx, rfidUid, -1)

	if err != nil && errors.Is(err, app.NotFoundErr) {
		s.websocket.Publish(websocketMessage)
		return nil, err
	} else if err != nil {
		return nil, err
	}

	checkin, err := s.createCheckinForUser(ctx, u, checkinTimestamp)
	if err != nil {
		return nil, err
	}

	websocketMessage.CheckIn = checkin
	s.websocket.Publish(websocketMessage)

	return checkin, nil
}

func (s *service) createCheckinForUser(ctx context.Context, user *user.User, timestamp time.Time) (*CheckIn, error) {

	checkIn := CheckIn{
		ID:        -1,
		Date:      truncateToStartOfDay(timestamp),
		Timestamp: timestamp,
		UserID:    user.ID,
	}

	return s.repo.SaveCheckIn(ctx, &checkIn)
}

func (s *service) DeleteCheckinByID(ctx context.Context, checkinID int64) error {
	return s.repo.DeleteCheckInByID(ctx, checkinID)
}

func (s *service) DeleteCheckInsByUserID(ctx context.Context, userID int64) error {
	return s.repo.DeleteCheckInsByUserID(ctx, userID)
}

func (s *service) DeleteOldCheckIns(ctx context.Context) error {

	retentionDaysEnv := os.Getenv("CHECKIN_RETENTION_DAYS")
	if retentionDaysEnv == "" {
		retentionDaysEnv = "365"
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
