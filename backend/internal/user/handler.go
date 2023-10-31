package user

import (
	"fmt"
	"github.com/d-rk/checkin-system/internal/websocket"
	"net/http"
	"strconv"

	"github.com/jmoiron/sqlx"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	ListUsers(c *gin.Context)
	GetUserByID(c *gin.Context)
	AddUser(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
	DeleteAllUsers(c *gin.Context)
}

type handler struct {
	db        *sqlx.DB
	websocket *websocket.Server
}

func CreateHandler(db *sqlx.DB, websocket *websocket.Server) Handler {
	return &handler{db, websocket}
}

func (h *handler) ListUsers(c *gin.Context) {

	users, err := ListUsers(h.db)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No users found"})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *handler) GetUserByID(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("param is not an id: %s", err.Error())})
		return
	}

	user, err := GetUserByID(h.db, id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *handler) AddUser(c *gin.Context) {

	var user User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to extract user from request body"})
		return
	}

	_, err := GetUserByName(h.db, user.Name, -1)

	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	_, err = GetUserByRfidUid(h.db, user.RFIDuid, -1)

	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User with rfid_uid already exists"})
		return
	}

	savedUser, err := user.Insert(h.db, c.Request.Context())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Unable to save user: %v", err)})
		return
	}

	c.JSON(http.StatusOK, savedUser)
}

func (h *handler) UpdateUser(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("param is not an id: %s", err.Error())})
		return
	}

	user, err := GetUserByID(h.db, id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to extract user from request body"})
		return
	}

	_, err = GetUserByName(h.db, user.Name, user.ID)

	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User with name already exists"})
		return
	}

	_, err = GetUserByRfidUid(h.db, user.RFIDuid, user.ID)

	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User with rfid_uid already exists"})
		return
	}

	user.ID = id

	savedUser, err := user.Update(h.db, c.Request.Context())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Unable to update user: %v", err)})
		return
	}

	c.JSON(http.StatusOK, savedUser)
}

func (h *handler) DeleteUser(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("param is not an id: %s", err.Error())})
		return
	}

	err = DeleteUser(h.db, c.Request.Context(), id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot delete user"})
		return
	}

	c.Writer.WriteHeader(http.StatusNoContent)
}

func (h *handler) DeleteAllUsers(c *gin.Context) {

	err := DeleteAllUsers(h.db, c.Request.Context())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot delete user"})
		return
	}

	c.Writer.WriteHeader(http.StatusNoContent)
}
