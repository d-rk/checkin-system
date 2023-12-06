package api

import (
	"github.com/d-rk/checkin-system/internal/checkin"
	"github.com/d-rk/checkin-system/internal/user"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"gopkg.in/guregu/null.v4"
	"time"
)

func toAPIUser(u *user.User) *User {
	return &User{
		Id:        u.ID,
		Name:      u.Name,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt.Ptr(),
	}
}

func toAPIUsers(users []user.User) []User {

	result := make([]User, len(users))

	for i, u := range users {
		uu := u
		result[i] = *toAPIUser(&uu)
	}

	return result
}

func fromAPIUser(u *User) *user.User {
	return &user.User{
		ID:        u.Id,
		CreatedAt: u.CreatedAt,
		UpdatedAt: null.TimeFromPtr(u.UpdatedAt),
		Name:      u.Name,
	}
}

func fromAPINewUser(u *NewUser) *user.User {
	return &user.User{
		Name: u.Name,
	}
}

func toAPICheckIn(c *checkin.CheckIn) *Checkin {
	return &Checkin{
		Id:        c.ID,
		Date:      openapi_types.Date{Time: c.Date},
		Timestamp: c.Timestamp.Format(time.RFC3339),
		UserId:    c.UserID,
	}
}

func toAPICheckIns(checkIns []checkin.CheckIn) []Checkin {

	result := make([]Checkin, len(checkIns))

	for i, c := range checkIns {
		cc := c
		result[i] = *toAPICheckIn(&cc)
	}

	return result
}

func toAPICheckInWithUser(c *checkin.CheckInWithUser) *CheckinWithUser {
	return &CheckinWithUser{
		Id:        c.ID,
		Date:      openapi_types.Date{Time: c.Date},
		Timestamp: c.Timestamp.Format(time.RFC3339),
		UserId:    c.UserID,
		Name:      c.User.Name,
	}
}

func toAPICheckInsWithUser(checkins []checkin.CheckInWithUser) []CheckinWithUser {

	result := make([]CheckinWithUser, len(checkins))

	for i, c := range checkins {
		cc := c
		result[i] = *toAPICheckInWithUser(&cc)
	}

	return result
}

func toAPICheckInsDates(dates []checkin.CheckInDate) []CheckinDate {

	result := make([]CheckinDate, len(dates))

	for i, d := range dates {
		result[i] = CheckinDate{
			Date: openapi_types.Date{d.Date},
		}
	}

	return result
}
