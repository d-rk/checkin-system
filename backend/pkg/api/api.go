package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/d-rk/checkin-system/pkg/app"
	"github.com/d-rk/checkin-system/pkg/auth"
	"github.com/d-rk/checkin-system/pkg/checkin"
	"github.com/d-rk/checkin-system/pkg/user"
	"github.com/flytam/filenamify"
	"github.com/gocarina/gocsv"
	"net/http"
	"time"
)

type apiHandler struct {
	userService    user.Service
	checkinService checkin.Service
}

func NewHandler(userService user.Service, checkinService checkin.Service) ServerInterface {
	return &apiHandler{
		userService:    userService,
		checkinService: checkinService,
	}
}

func (h *apiHandler) Login(w http.ResponseWriter, r *http.Request) {

	credentials := &LoginCredentials{}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		handlerError(w, r, BadRequestErr.Wrap(err))
		return
	}

	u, err := h.userService.GetUserByNameAndPassword(r.Context(), credentials.Username, credentials.Password)
	if err != nil && errors.Is(err, app.NotFoundErr) {
		handlerError(w, r, InvalidCredentialsErr.Wrap(err))
		return
	} else if err != nil {
		handlerError(w, r, err)
		return
	}

	token, err := auth.GenerateToken(u.ID)
	if err != nil {
		handlerError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusOK, BearerToken{Token: token})
}

func (h *apiHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.ListUsers(r.Context())
	if err != nil {
		handlerError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusOK, toAPIUsers(users))
}

func (h *apiHandler) CreateUser(w http.ResponseWriter, r *http.Request) {

	apiUser := &NewUser{}

	if err := json.NewDecoder(r.Body).Decode(&apiUser); err != nil {
		handlerError(w, r, BadRequestErr.Wrap(err))
		return
	}

	u, err := h.userService.CreateUser(r.Context(), fromAPINewUser(apiUser))
	if err != nil {
		if errors.Is(err, app.ConflictErr) {
			handlerError(w, r, ConflictErr.Wrap(err))
			return
		} else {
			handlerError(w, r, err)
			return
		}
	}

	writeJSON(w, r, http.StatusCreated, toAPIUser(u))
}

func (h *apiHandler) GetAuthenticatedUser(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value(AuthenticatedUserID).(int64)
	if !ok {
		handlerError(w, r, fmt.Errorf("authenticated user not found on context"))
		return
	}

	h.GetUser(w, r, userID)
}

func (h *apiHandler) GetUser(w http.ResponseWriter, r *http.Request, userId UserIdPathParam) {
	u, err := h.userService.GetUserByID(r.Context(), userId)
	if err != nil && errors.Is(err, app.NotFoundErr) {
		handlerError(w, r, NotFoundErr.Wrap(err))
		return
	} else if err != nil {
		handlerError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusOK, toAPIUser(u))
}

func (h *apiHandler) UpdateUser(w http.ResponseWriter, r *http.Request, userId UserIdPathParam) {

	apiUser := &User{}

	if err := json.NewDecoder(r.Body).Decode(&apiUser); err != nil {
		handlerError(w, r, BadRequestErr.Wrap(err))
		return
	}

	if userId != apiUser.Id {
		handlerError(w, r, BadRequestErr.Wrap(fmt.Errorf("id missmatch: %d != %d", userId, apiUser.Id)))
		return
	}

	u, err := h.userService.UpdateUser(r.Context(), fromAPIUser(apiUser))
	if err != nil && errors.Is(err, app.NotFoundErr) {
		handlerError(w, r, NotFoundErr.Wrap(err))
		return
	} else if err != nil && errors.Is(err, app.ConflictErr) {
		handlerError(w, r, ConflictErr.Wrap(err))
		return
	} else if err != nil {
		handlerError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusOK, toAPIUser(u))
}

func (h *apiHandler) DeleteAllUsers(w http.ResponseWriter, r *http.Request) {

	if err := h.userService.DeleteAllUsers(r.Context()); err != nil {
		handlerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *apiHandler) DeleteUser(w http.ResponseWriter, r *http.Request, userId UserIdPathParam) {

	authenticatedUserID, ok := r.Context().Value(AuthenticatedUserID).(int64)
	if !ok {
		handlerError(w, r, fmt.Errorf("authenticated user not found on context"))
		return
	}

	if authenticatedUserID == userId {
		handlerError(w, r, ConflictErr.Wrap(fmt.Errorf("cannot delete yourself")))
		return
	}

	if err := h.userService.DeleteUser(r.Context(), userId); err != nil {
		handlerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *apiHandler) UpdateUserPassword(w http.ResponseWriter, r *http.Request, userId UserIdPathParam) {

	password := &Password{}

	if err := json.NewDecoder(r.Body).Decode(&password); err != nil {
		handlerError(w, r, BadRequestErr.Wrap(err))
		return
	}

	if err := h.userService.UpdateUserPassword(r.Context(), userId, password.Password); err != nil {
		handlerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *apiHandler) ListCheckIns(w http.ResponseWriter, r *http.Request) {
	checkins, err := h.checkinService.ListCheckIns(r.Context())
	if err != nil {
		handlerError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusOK, toAPICheckIns(checkins))
}

func (h *apiHandler) CreateCheckIn(w http.ResponseWriter, r *http.Request, userId UserIdPathParam, params CreateCheckInParams) {

	c, err := h.checkinService.CreateCheckInForUser(r.Context(), userId, params.Timestamp)
	if err != nil {
		if errors.Is(err, app.ConflictErr) {
			handlerError(w, r, ConflictErr.Wrap(err))
			return
		} else {
			handlerError(w, r, err)
			return
		}
	}

	writeJSON(w, r, http.StatusCreated, toAPICheckIn(c))
}

func (h *apiHandler) CreateRfidCheckIn(w http.ResponseWriter, r *http.Request, params CreateRfidCheckInParams) {

	c, err := h.checkinService.CreateCheckInForRFID(r.Context(), params.Rfid, params.Timestamp)
	if err != nil {
		if errors.Is(err, app.ConflictErr) {
			handlerError(w, r, ConflictErr.Wrap(err))
			return
		} else {
			handlerError(w, r, err)
			return
		}
	}

	writeJSON(w, r, http.StatusCreated, toAPICheckIn(c))
}

func (h *apiHandler) ListAllCheckIns(w http.ResponseWriter, r *http.Request) {
	checkIns, err := h.checkinService.ListAllCheckIns(r.Context())
	if err != nil {
		handlerError(w, r, err)
		return
	}

	switch r.Header.Get("Accept") {
	case "application/csv":
		writeCSV(w, r, fmt.Sprintf("%s_all_checkins.csv", time.Now().Format("2006-01-02")), checkIns)
	case "application/json":
		fallthrough
	default:
		writeJSON(w, r, http.StatusOK, toAPICheckInsWithUser(checkIns))
	}
}

func (h *apiHandler) ListCheckInDates(w http.ResponseWriter, r *http.Request) {

	dates, err := h.checkinService.ListCheckInDates(r.Context())
	if err != nil {
		handlerError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusOK, toAPICheckInsDates(dates))
}

func (h *apiHandler) ListCheckInsPerDay(w http.ResponseWriter, r *http.Request, params ListCheckInsPerDayParams) {
	checkIns, err := h.checkinService.ListCheckInsPerDay(r.Context(), params.Day.Time)
	if err != nil {
		handlerError(w, r, err)
		return
	}

	switch r.Header.Get("Accept") {
	case "application/csv":
		writeCSV(w, r, fmt.Sprintf("%s.csv", params.Day.String()), checkIns)
	case "application/json":
		fallthrough
	default:
		writeJSON(w, r, http.StatusOK, toAPICheckInsWithUser(checkIns))
	}
}

func (h *apiHandler) DeleteCheckIn(w http.ResponseWriter, r *http.Request, checkinId CheckInIdPathParam) {

	if err := h.checkinService.DeleteCheckInByID(r.Context(), checkinId); err != nil {
		handlerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *apiHandler) DeleteUserCheckIns(w http.ResponseWriter, r *http.Request, userId UserIdPathParam) {

	if err := h.checkinService.DeleteCheckInsByUserID(r.Context(), userId); err != nil {
		handlerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *apiHandler) GetUserCheckIns(w http.ResponseWriter, r *http.Request, userId UserIdPathParam) {

	u, err := h.userService.GetUserByID(r.Context(), userId)
	if err != nil {
		handlerError(w, r, err)
		return
	}

	checkIns, err := h.checkinService.ListUserCheckIns(r.Context(), userId)
	if err != nil {
		handlerError(w, r, err)
		return
	}

	switch r.Header.Get("Accept") {
	case "application/csv":
		writeCSV(w, r, fmt.Sprintf("%s.csv", u.Name), checkIns)
	case "application/json":
		fallthrough
	default:
		writeJSON(w, r, http.StatusOK, toAPICheckIns(checkIns))
	}
}

func (h *apiHandler) ListUserGroups(w http.ResponseWriter, r *http.Request) {
	groups, err := h.userService.ListUserGroups(r.Context())
	if err != nil {
		handlerError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusOK, groups)
}

func writeJSON(w http.ResponseWriter, r *http.Request, status int, response any) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		handlerError(w, r, err)
		return
	}
}

func writeCSV(w http.ResponseWriter, r *http.Request, filename string, response any) {

	saneFilename, err := filenamify.Filenamify(filename, filenamify.Options{
		Replacement: "_",
	})
	if err != nil {
		handlerError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, saneFilename))
	w.Header().Set("X-Filename", saneFilename)
	w.WriteHeader(http.StatusOK)

	err = gocsv.Marshal(response, w)
	if err != nil {
		handlerError(w, r, err)
		return
	}
}
