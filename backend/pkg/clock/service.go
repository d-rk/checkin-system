package clock

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/d-rk/checkin-system/pkg/cmd"
)

type Service interface {
	GetClock(ctx context.Context) (Clock, error)
	SetClock(ctx context.Context, timestamp time.Time) error
}

type service struct {
	executor cmd.Executor
}

func NewService() Service {
	return &service{executor: cmd.NewExecutor()}
}

func (s *service) GetClock(_ context.Context) (Clock, error) {
	return Clock{
		Timestamp: time.Now(),
	}, nil
}

func (s *service) SetClock(ctx context.Context, timestamp time.Time) error {
	// Format the time for the 'date' command (YYYY-MM-DD HH:MM:SS)
	loc, _ := time.LoadLocation("Local")
	dateStr := timestamp.In(loc).Format("2006-01-02 15:04:05")

	slog.Info("setting system clock", "timestamp", timestamp, "dateStr", dateStr)

	if err := s.executor.Call(ctx, "date", "-s", dateStr); err != nil {
		return fmt.Errorf("failed to call date command: %w", err)
	}

	if err := s.executor.Call(ctx, "hwclock", "--systohc"); err != nil {
		return fmt.Errorf("failed to call hwclock command: %w", err)
	}

	return nil
}
