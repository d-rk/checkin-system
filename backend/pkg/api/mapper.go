package api

import (
	"os"
	"strconv"
	"time"

	"github.com/d-rk/checkin-system/pkg/auth"
	"github.com/d-rk/checkin-system/pkg/checkin"
	"github.com/d-rk/checkin-system/pkg/clock"
	"github.com/d-rk/checkin-system/pkg/user"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"gopkg.in/guregu/null.v4"
)

func toAPIUser(u *user.User) *User {
	return &User{
		Id:       u.ID,
		Name:     u.Name,
		Group:    u.Group.Ptr(),
		Role:     u.Role,
		MemberId: u.MemberID.Ptr(),
		RfidUid:  u.RFIDuid.Ptr(),
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
		ID:       u.Id,
		Name:     u.Name,
		Group:    null.StringFromPtr(u.Group),
		Role:     u.Role,
		MemberID: null.StringFromPtr(u.MemberId),
		RFIDuid:  null.StringFromPtr(u.RfidUid),
	}
}

func fromAPINewUser(u *NewUser) *user.User {
	return &user.User{
		Name:     u.Name,
		Group:    null.StringFromPtr(u.Group),
		Role:     u.Role,
		MemberID: null.StringFromPtr(u.MemberId),
		RFIDuid:  null.StringFromPtr(u.RfidUid),
	}
}

func toAPICheckIn(c *checkin.CheckIn) *CheckIn {
	return &CheckIn{
		Id:        c.ID,
		Date:      openapi_types.Date{Time: c.Date},
		Timestamp: c.Timestamp.Format(time.RFC3339),
		UserId:    c.UserID,
	}
}

func toAPICheckIns(checkIns []checkin.CheckIn) []CheckIn {

	result := make([]CheckIn, len(checkIns))

	for i, c := range checkIns {
		cc := c
		result[i] = *toAPICheckIn(&cc)
	}

	return result
}

func toAPICheckInWithUser(c *checkin.WithUser) *CheckInWithUser {
	return &CheckInWithUser{
		Id:        c.ID,
		Date:      openapi_types.Date{Time: c.Date},
		Timestamp: c.Timestamp.Format(time.RFC3339),
		UserId:    c.UserID,
		User: User{
			Id:       c.User.ID,
			Name:     c.User.Name,
			Group:    c.User.Group.Ptr(),
			MemberId: c.User.MemberID.Ptr(),
			RfidUid:  c.User.RFIDuid.Ptr(),
			Role:     c.User.Role,
		},
	}
}

func toAPICheckInsWithUser(checkins []checkin.WithUser) []CheckInWithUser {

	result := make([]CheckInWithUser, len(checkins))

	for i, c := range checkins {
		cc := c
		result[i] = *toAPICheckInWithUser(&cc)
	}

	return result
}

func toAPICheckInsDates(dates []checkin.Date) []CheckInDate {

	result := make([]CheckInDate, len(dates))

	for i, d := range dates {
		result[i] = CheckInDate{
			Date: openapi_types.Date{Time: d.Date},
		}
	}

	return result
}

func toAPIClock(refTimestamp string, c *clock.Clock) *Clock {
	return &Clock{
		RefTimestamp: refTimestamp,
		Timestamp:    c.Timestamp.Format(time.RFC3339),
	}
}

func toAPIWifiNetworks(networks []string) []WifiNetwork {
	result := make([]WifiNetwork, len(networks))

	for i, ssid := range networks {
		result[i] = WifiNetwork{
			Ssid: ssid,
		}
	}
	return result
}

func fromAPIRefTimestamp(timestamp string) (time.Time, error) {

	loc, _ := time.LoadLocation("Local")

	val, e := time.ParseInLocation(time.RFC3339, timestamp, loc)
	if e != nil {
		return time.Time{}, e
	}
	return val, nil
}

func generateBearerToken(userID int64) (BearerToken, error) {

	var (
		bearerToken BearerToken
		err         error
	)

	bearerToken.Token, err = auth.GenerateToken(userID)
	if err != nil {
		return BearerToken{}, err
	}

	bearerToken.RefreshToken, err = auth.GenerateRefreshToken(userID)
	if err != nil {
		return BearerToken{}, err
	}

	tokenExpiryMinutesStr := os.Getenv("TOKEN_EXPIRY_MINUTES")
	if tokenExpiryMinutesStr != "" {
		if tokenExpiryMinutes, parseErr := strconv.Atoi(tokenExpiryMinutesStr); parseErr == nil {
			expirySeconds := int((time.Duration(tokenExpiryMinutes) * time.Minute).Seconds())
			bearerToken.ExpiresIn = &expirySeconds
		}
	}

	return bearerToken, nil
}
