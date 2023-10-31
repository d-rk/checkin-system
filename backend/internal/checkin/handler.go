package checkin

import (
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
	"github.com/jmoiron/sqlx"
)

type Handler interface {
	ListCheckIns(c *gin.Context)
	ListCheckInsPerDay(c *gin.Context)
	ListAllCheckIns(c *gin.Context)
	ListUserCheckIns(c *gin.Context)
	AddCheckIn(c *gin.Context)
	DeleteCheckIn(c *gin.Context)
	DeleteUserCheckIns(c *gin.Context)
	DeleteOldCheckIns()
	ListCheckInDates(c *gin.Context)
}

type handler struct {
	db        *sqlx.DB
	websocket *websocket.Server
}

type Request struct {
	RFIDuid string `db:"rfid_uid" json:"rfid_uid"`
}

type WebsocketMessage struct {
	Request
	CheckIn *CheckIn `json:"check_in"`
}

func CreateHandler(db *sqlx.DB, websocket *websocket.Server) Handler {
	return &handler{db, websocket}
}

func (h *handler) ListCheckIns(c *gin.Context) {

	checkIns, err := ListCheckIns(h.db)
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

	checkIns, err := ListCheckInsPerDay(h.db, day)
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

	checkIns, err := ListAllCheckIns(h.db)
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

	checkIns, err := ListUserCheckIns(h.db, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("not found: %s", err.Error())})
		return
	}

	switch c.Request.Header.Get("Accept") {
	case "application/csv":

		user, err := user.GetUserByID(h.db, userID)
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

func (h *handler) AddCheckIn(c *gin.Context) {

	var checkInRequest Request

	if err := c.BindJSON(&checkInRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to extract checkInRequest from request body"})
		return
	}

	user, err := user.GetUserByRfidUid(h.db, checkInRequest.RFIDuid, -1)

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

func (h *handler) DeleteCheckIn(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("param is not an id: %s", err.Error())})
		return
	}

	err = DeleteCheckInByID(h.db, c.Request.Context(), id)

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

	err = DeleteCheckInsByUserID(h.db, c.Request.Context(), userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot delete user checkIns"})
		return
	}

	c.Writer.WriteHeader(http.StatusNoContent)
}

func (h *handler) DeleteOldCheckIns() {

	retentionDaysEnv := os.Getenv("CHECKIN_RETENTION_DAYS")
	if retentionDaysEnv == "" {
		retentionDaysEnv = "365"
	}

	retentionDays, err := strconv.ParseInt(retentionDaysEnv, 10, 64)
	if err != nil {
		log.Fatal("parsing CHECKIN_RETENTION_DAYS failed", err)
	}

	err = DeleteCheckInsOlderThan(h.db, retentionDays)

	if err != nil {
		log.Fatal("error deleting old checkIns", err)
	}
}

func (h *handler) ListCheckInDates(c *gin.Context) {

	dates, err := ListCheckInDates(h.db)
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
