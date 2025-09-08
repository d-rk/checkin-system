package user

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type User struct {
	ID             int64       `db:"id"              json:"id"         csv:"-"`
	CreatedAt      time.Time   `db:"created_at"      json:"created_at" csv:"-"`
	UpdatedAt      null.Time   `db:"updated_at"      json:"updated_at" csv:"-"`
	Name           string      `db:"name"            json:"name"       csv:"name"`
	Group          null.String `db:"group_name"      json:"group"      csv:"group"`
	Role           string      `db:"role"            json:"role"       csv:"-"`
	PasswordDigest null.String `db:"password_digest" json:"-"          csv:"-"`
	MemberID       null.String `db:"member_id"       json:"member_id"  csv:"member_id"`
	RFIDuid        null.String `db:"rfid_uid"        json:"rfid_uid"   csv:"rfid_uid"`
}
