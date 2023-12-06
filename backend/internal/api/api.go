package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/d-rk/checkin-system/internal/app"
	"github.com/d-rk/checkin-system/internal/auth"
	"github.com/d-rk/checkin-system/internal/checkin"
	"github.com/d-rk/checkin-system/internal/user"
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

func (h *apiHandler) ListCheckins(w http.ResponseWriter, r *http.Request) {
	checkins, err := h.checkinService.ListCheckins(r.Context())
	if err != nil {
		handlerError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusOK, toAPICheckIns(checkins))
}

func (h *apiHandler) CreateCheckin(w http.ResponseWriter, r *http.Request, userId UserIdPathParam, params CreateCheckinParams) {

	c, err := h.checkinService.CreateCheckinForUser(r.Context(), userId, params.Timestamp)
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

func (h *apiHandler) CreateRfidCheckin(w http.ResponseWriter, r *http.Request, params CreateRfidCheckinParams) {

	c, err := h.checkinService.CreateCheckinForRFID(r.Context(), params.Rfid, params.Timestamp)
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

func (h *apiHandler) ListAllCheckins(w http.ResponseWriter, r *http.Request) {
	checkIns, err := h.checkinService.ListAllCheckins(r.Context())
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

func (h *apiHandler) ListCheckinDates(w http.ResponseWriter, r *http.Request) {

	dates, err := h.checkinService.ListCheckinDates(r.Context())
	if err != nil {
		handlerError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusOK, toAPICheckInsDates(dates))
}

func (h *apiHandler) ListCheckinsPerDay(w http.ResponseWriter, r *http.Request, params ListCheckinsPerDayParams) {
	checkIns, err := h.checkinService.ListCheckinsPerDay(r.Context(), params.Day.Time)
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

func (h *apiHandler) DeleteCheckin(w http.ResponseWriter, r *http.Request, checkinId CheckinIdPathParam) {

	if err := h.checkinService.DeleteCheckinByID(r.Context(), checkinId); err != nil {
		handlerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *apiHandler) DeleteUserCheckins(w http.ResponseWriter, r *http.Request, userId UserIdPathParam) {

	if err := h.checkinService.DeleteCheckInsByUserID(r.Context(), userId); err != nil {
		handlerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *apiHandler) GetUserCheckins(w http.ResponseWriter, r *http.Request, userId UserIdPathParam) {

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

func writeJSON(w http.ResponseWriter, r *http.Request, status int, response any) {

	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		handlerError(w, r, err)
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

	err = gocsv.Marshal(response, w)
	if err != nil {
		handlerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, saneFilename))
	w.Header().Add("X-Filename", saneFilename)
}
