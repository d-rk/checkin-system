package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/d-rk/checkin-system/pkg/models"
	"github.com/d-rk/checkin-system/pkg/services/websocket"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type CheckInHandler struct {
	db *sqlx.DB
	websocket *websocket.Server
}

type CheckinRequest struct {
	RFIDuid string `db:"rfid_uid" json:"rfid_uid"`
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

	c.JSON(http.StatusOK, checkIns)
}

func (h *CheckInHandler) AddCheckIn(c *gin.Context) {

	var checkInRequest CheckinRequest

	if err := c.BindJSON(&checkInRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to extract checkInRequest from request body"})
		return
	}

	h.websocket.Publish(checkInRequest);

	user, err := models.GetUserByRfidUid(h.db, checkInRequest.RFIDuid, -1)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("No user found with rfid_uuid = %s", checkInRequest.RFIDuid)})
		return
	}

	checkIn := models.CheckIn{
		ID:        -1,
		Timestamp: time.Now(),
		UserID:    user.ID,
	}

	savedCheckIn, err := checkIn.Save(h.db, c.Request.Context())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Unable to save checkIn: %v", err)})
		return
	}

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
