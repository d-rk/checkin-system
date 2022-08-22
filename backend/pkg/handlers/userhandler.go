package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/d-rk/checkin-system/pkg/models"
	"github.com/d-rk/checkin-system/pkg/services/websocket"
	"github.com/jmoiron/sqlx"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	db        *sqlx.DB
	websocket *websocket.Server
}

func CreateUserHandler(db *sqlx.DB, websocket *websocket.Server) *UserHandler {
	return &UserHandler{db, websocket}
}

func (h *UserHandler) ListUsers(c *gin.Context) {

	users, err := models.ListUsers(h.db)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No users found"})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("param is not an id: %s", err.Error())})
		return
	}

	users, err := models.GetUserByID(h.db, id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No users found"})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) AddUser(c *gin.Context) {

	var user models.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to extract user from request body"})
		return
	}

	_, err := models.GetUserByName(h.db, user.Name)

	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	savedUser, err := user.Save(h.db, c.Request.Context())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Unable to save user: %v", err)})
		return
	}

	c.JSON(http.StatusOK, savedUser)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("param is not an id: %s", err.Error())})
		return
	}

	err = models.DeleteUser(h.db, c.Request.Context(), id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot delete user"})
		return
	}

	c.Writer.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) DeleteAllUsers(c *gin.Context) {

	err := models.DeleteAllUsers(h.db, c.Request.Context())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot delete user"})
		return
	}

	c.Writer.WriteHeader(http.StatusNoContent)
}
