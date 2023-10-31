package checkin

import (
	"context"
	"fmt"
	"github.com/d-rk/checkin-system/internal/user"
	"github.com/d-rk/checkin-system/internal/websocket"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/flytam/filenamify"
	"github.com/gin-gonic/gin"
	"github.com/gocarina/gocsv"
)

type Handler interface {
	ListCheckIns(c *gin.Context)
	ListCheckInsPerDay(c *gin.Context)
	ListAllCheckIns(c *gin.Context)
	ListUserCheckIns(c *gin.Context)
	AddCheckIn(c *gin.Context)
	DeleteCheckIn(c *gin.Context)
	DeleteUserCheckIns(c *gin.Context)
	DeleteOldCheckIns(ctx context.Context)
	ListCheckInDates(c *gin.Context)
}

type handler struct {
	repo      Repository
	userRepo  user.Repository
	websocket *websocket.Server
}

type Request struct {
	RFIDuid string `db:"rfid_uid" json:"rfid_uid"`
}

type WebsocketMessage struct {
	Request
	CheckIn *CheckIn `json:"check_in"`
}

func CreateHandler(repo Repository, userRepo user.Repository, websocket *websocket.Server) Handler {
	return &handler{repo, userRepo, websocket}
}

func (h *handler) ListCheckIns(c *gin.Context) {

	checkIns, err := h.repo.ListCheckIns(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("not found: %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, checkIns)
}

func (h *handler) ListCheckInsPerDay(c *gin.Context) {

	dayParam, ok := c.GetQuery("day")

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "query param day is required"})
		return
	}

	day, err := time.Parse("2006-01-02", dayParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("day param is not a date: %s", err.Error())})
		return
	}

	checkIns, err := h.repo.ListCheckInsPerDay(c, day)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("not found: %s", err.Error())})
		return
	}

	switch c.Request.Header.Get("Accept") {
	case "application/csv":
		writeCSV(c, fmt.Sprintf("%s.csv", dayParam), checkIns)
	case "application/json":
		fallthrough
	default:
		c.JSON(http.StatusOK, checkIns)
	}
}

func (h *handler) ListAllCheckIns(c *gin.Context) {

	checkIns, err := h.repo.ListAllCheckIns(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("not found: %s", err.Error())})
		return
	}

	switch c.Request.Header.Get("Accept") {
	case "application/csv":
		writeCSV(c, fmt.Sprintf("%s_all_checkins.csv", time.Now().Format("2006-01-02")), checkIns)
	case "application/json":
		fallthrough
	default:
		c.JSON(http.StatusOK, checkIns)
	}
}

func (h *handler) ListUserCheckIns(c *gin.Context) {

	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("param is not an id: %s", err.Error())})
		return
	}

	checkIns, err := h.repo.ListUserCheckIns(c, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("not found: %s", err.Error())})
		return
	}

	switch c.Request.Header.Get("Accept") {
	case "application/csv":

		u, err := h.userRepo.GetUserByID(c, userID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("user not found: %s", err.Error())})
			return
		}

		writeCSV(c, fmt.Sprintf("%s.csv", u.Name), checkIns)

	case "application/json":
		fallthrough
	default:
		c.JSON(http.StatusOK, checkIns)
	}
}

func (h *handler) AddCheckIn(c *gin.Context) {

	var checkInRequest Request

	if err := c.BindJSON(&checkInRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to extract checkInRequest from request body"})
		return
	}

	u, err := h.userRepo.GetUserByRfidUid(c, checkInRequest.RFIDuid, -1)

	websocketMessage := WebsocketMessage{}
	websocketMessage.RFIDuid = checkInRequest.RFIDuid

	if err != nil {
		h.websocket.Publish(websocketMessage)
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("No user found with rfid_uuid = %s", checkInRequest.RFIDuid)})
		return
	}

	checkIn := CheckIn{
		ID:        -1,
		Date:      truncateToStartOfDay(time.Now()),
		Timestamp: time.Now(),
		UserID:    u.ID,
	}

	savedCheckIn, err := h.repo.SaveCheckIn(c, &checkIn)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Unable to save checkIn: %v", err)})
		return
	}

	websocketMessage.CheckIn = savedCheckIn

	h.websocket.Publish(websocketMessage)
	c.JSON(http.StatusOK, savedCheckIn)
}

func (h *handler) DeleteCheckIn(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("param is not an id: %s", err.Error())})
		return
	}

	err = h.repo.DeleteCheckInByID(c, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot delete checkIn"})
		return
	}

	c.Writer.WriteHeader(http.StatusNoContent)
}

func (h *handler) DeleteUserCheckIns(c *gin.Context) {

	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("param is not an id: %s", err.Error())})
		return
	}

	err = h.repo.DeleteCheckInsByUserID(c, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot delete user checkIns"})
		return
	}

	c.Writer.WriteHeader(http.StatusNoContent)
}

func (h *handler) DeleteOldCheckIns(ctx context.Context) {

	retentionDaysEnv := os.Getenv("CHECKIN_RETENTION_DAYS")
	if retentionDaysEnv == "" {
		retentionDaysEnv = "365"
	}

	retentionDays, err := strconv.ParseInt(retentionDaysEnv, 10, 64)
	if err != nil {
		log.Fatal("parsing CHECKIN_RETENTION_DAYS failed", err)
	}

	err = h.repo.DeleteCheckInsOlderThan(ctx, retentionDays)

	if err != nil {
		log.Fatal("error deleting old checkIns", err)
	}
}

func (h *handler) ListCheckInDates(c *gin.Context) {

	dates, err := h.repo.ListCheckInDates(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("unable to list dates: %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, dates)
}

func truncateToStartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

func writeCSV(c *gin.Context, filename string, data interface{}) {

	saneFilename, err := filenamify.Filenamify(filename, filenamify.Options{
		Replacement: "_",
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("internal error: %s", err.Error())})
	}

	err = gocsv.Marshal(data, c.Writer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("internal error: %s", err.Error())})
	}

	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, saneFilename))
	c.Writer.Header().Add("X-Filename", saneFilename)
}
