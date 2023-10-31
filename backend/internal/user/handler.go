package user

import (
	"fmt"
	"github.com/d-rk/checkin-system/internal/websocket"
	"net/http"
	"strconv"

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
	repo      Repository
	websocket *websocket.Server
}

func CreateHandler(repo Repository, websocket *websocket.Server) Handler {
	return &handler{repo, websocket}
}

func (h *handler) ListUsers(c *gin.Context) {

	users, err := h.repo.ListUsers(c)

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

	user, err := h.repo.GetUserByID(c, id)

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

	_, err := h.repo.GetUserByName(c, user.Name, -1)

	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	_, err = h.repo.GetUserByRfidUid(c, user.RFIDuid, -1)

	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User with rfid_uid already exists"})
		return
	}

	savedUser, err := h.repo.SaveUser(c, &user)

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

	user, err := h.repo.GetUserByID(c, id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to extract user from request body"})
		return
	}

	_, err = h.repo.GetUserByName(c, user.Name, user.ID)

	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User with name already exists"})
		return
	}

	_, err = h.repo.GetUserByRfidUid(c, user.RFIDuid, user.ID)

	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User with rfid_uid already exists"})
		return
	}

	user.ID = id

	savedUser, err := h.repo.UpdateUser(c, &user)

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

	err = h.repo.DeleteUser(c, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot delete user"})
		return
	}

	c.Writer.WriteHeader(http.StatusNoContent)
}

func (h *handler) DeleteAllUsers(c *gin.Context) {

	err := h.repo.DeleteAllUsers(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot delete user"})
		return
	}

	c.Writer.WriteHeader(http.StatusNoContent)
}
