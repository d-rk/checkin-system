package checkin

import (
	"github.com/d-rk/checkin-system/internal/user"
	"time"
)

type CheckIn struct {
	ID        int64     `db:"id" json:"id" csv:"id"`
	Date      time.Time `db:"date" json:"date" csv:"date"`
	Timestamp time.Time `db:"timestamp" json:"timestamp" csv:"timestamp"`
	UserID    int64     `db:"user_id" json:"user_id" csv:"-"`
}

type CheckInWithUser struct {
	CheckIn
	User user.User `db:"user" json:"user" csv:"-"`
}

type CheckInDate struct {
	Date time.Time `db:"date" json:"date"`
}
