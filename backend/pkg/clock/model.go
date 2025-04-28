package clock

import (
	"time"
)

type Clock struct {
	Timestamp time.Time `json:"timestamp"`
}
