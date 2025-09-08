package checkin

import (
	"time"

	"github.com/d-rk/checkin-system/pkg/user"
)

type CheckIn struct {
	ID        int64     `db:"id"        json:"id"        csv:"id"`
	Date      time.Time `db:"date"      json:"date"      csv:"date"`
	Timestamp time.Time `db:"timestamp" json:"timestamp" csv:"timestamp"`
	UserID    int64     `db:"user_id"   json:"user_id"   csv:"-"`
}

type WithUser struct {
	CheckIn

	User user.User `db:"user" json:"user" csv:"user"`
}

type Date struct {
	Date time.Time `db:"date" json:"date"`
}

type WebsocketMessage struct {
	RFIDuid string   `db:"rfid_uid" json:"rfid_uid"`
	CheckIn *CheckIn `              json:"check_in"`
}
