package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/d-rk/checkin-system/pkg/models"
	"github.com/d-rk/checkin-system/pkg/services/websocket"
	"github.com/flytam/filenamify"
	"github.com/gin-gonic/gin"
	"github.com/gocarina/gocsv"
	"github.com/jmoiron/sqlx"
)

type CheckInHandler struct {
	db        *sqlx.DB
	websocket *websocket.Server
}

type CheckinRequest struct {
	RFIDuid string `db:"rfid_uid" json:"rfid_uid"`
}

type CheckinWebsocketMessage struct {
	CheckinRequest
	CheckIn *models.CheckIn `json:"check_in"`
}

func CreateCheckInHandler(db *sqlx.DB, websocket *websocket.Server) *CheckInHandler {
	return &CheckInHandler{db, websocket}
}

func (h *CheckInHandler) ListCheckIns(c *gin.Context) {

	checkIns, err := models.ListCheckIns(h.db)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("not found: %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, checkIns)
}

func (h *CheckInHandler) ListCheckInsPerDay(c *gin.Context) {

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

	checkIns, err := models.ListCheckInsPerDay(h.db, day)
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

func (h *CheckInHandler) ListAllCheckIns(c *gin.Context) {

	checkIns, err := models.ListAllCheckIns(h.db)
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

func (h *CheckInHandler) ListUserCheckIns(c *gin.Context) {

	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("param is not an id: %s", err.Error())})
		return
	}

	checkIns, err := models.ListUserCheckIns(h.db, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("not found: %s", err.Error())})
		return
	}

	switch c.Request.Header.Get("Accept") {
	case "application/csv":

		user, err := models.GetUserByID(h.db, userID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("user not found: %s", err.Error())})
			return
		}

		writeCSV(c, fmt.Sprintf("%s.csv", user.Name), checkIns)

	case "application/json":
		fallthrough
	default:
		c.JSON(http.StatusOK, checkIns)
	}
}

func (h *CheckInHandler) AddCheckIn(c *gin.Context) {

	var checkInRequest CheckinRequest

	if err := c.BindJSON(&checkInRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to extract checkInRequest from request body"})
		return
	}

	user, err := models.GetUserByRfidUid(h.db, checkInRequest.RFIDuid, -1)

	websocketMessage := CheckinWebsocketMessage{}
	websocketMessage.RFIDuid = checkInRequest.RFIDuid

	if err != nil {
		h.websocket.Publish(websocketMessage)
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("No user found with rfid_uuid = %s", checkInRequest.RFIDuid)})
		return
	}

	checkIn := models.CheckIn{
		ID:        -1,
		Date:      truncateToStartOfDay(time.Now()),
		Timestamp: time.Now(),
		UserID:    user.ID,
	}

	savedCheckIn, err := checkIn.Save(h.db, c.Request.Context())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Unable to save checkIn: %v", err)})
		return
	}

	websocketMessage.CheckIn = savedCheckIn

	h.websocket.Publish(websocketMessage)
	c.JSON(http.StatusOK, savedCheckIn)
}

func (h *CheckInHandler) DeleteCheckIn(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("param is not an id: %s", err.Error())})
		return
	}

	err = models.DeleteCheckInByID(h.db, c.Request.Context(), id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot delete checkIn"})
		return
	}

	c.Writer.WriteHeader(http.StatusNoContent)
}

func (h *CheckInHandler) DeleteUserCheckIns(c *gin.Context) {

	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("param is not an id: %s", err.Error())})
		return
	}

	err = models.DeleteCheckInsByUserID(h.db, c.Request.Context(), userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot delete user checkIns"})
		return
	}

	c.Writer.WriteHeader(http.StatusNoContent)
}

func (h *CheckInHandler) DeleteOldCheckIns() {

	retentionDaysEnv := os.Getenv("CHECKIN_RETENTION_DAYS")
	if retentionDaysEnv == "" {
		retentionDaysEnv = "365"
	}

	retentionDays, err := strconv.ParseInt(retentionDaysEnv, 10, 64)
	if err != nil {
		log.Fatal("parsing CHECKIN_RETENTION_DAYS failed", err)
	}

	err = models.DeleteCheckInsOlderThan(h.db, retentionDays)

	if err != nil {
		log.Fatal("error deleting old checkIns", err)
	}
}

func (h *CheckInHandler) ListCheckInDates(c *gin.Context) {

	dates, err := models.ListCheckInDates(h.db)
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
