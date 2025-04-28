package clock

import (
	"context"
	"time"
)

type Service interface {
	GetClock(ctx context.Context) (Clock, error)
	SetClock(ctx context.Context) error
}

type service struct {
}

func NewService() Service {
	return &service{}
}

func (s *service) GetClock(ctx context.Context) (Clock, error) {
	return Clock{
		Timestamp: time.Now(),
	}, nil
}

func (s *service) SetClock(ctx context.Context) error {
	return nil
}
