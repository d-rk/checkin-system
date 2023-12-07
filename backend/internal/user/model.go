package user

import (
	"gopkg.in/guregu/null.v4"
	"time"
)

type User struct {
	ID             int64       `db:"id" json:"id" csv:"-"`
	CreatedAt      time.Time   `db:"created_at" json:"created_at" csv:"-"`
	UpdatedAt      null.Time   `db:"updated_at" json:"updated_at" csv:"-"`
	Name           string      `db:"name" json:"name"  csv:"name"`
	PasswordDigest null.String `db:"password_digest" json:"-"  csv:"-"`
	MemberID       null.String `db:"member_id" json:"member_id"  csv:"member_id"`
	RFIDuid        null.String `db:"rfid_uid" json:"rfid_uid" csv:"rfid_uid"`
}
