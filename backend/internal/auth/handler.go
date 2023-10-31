package auth

import (
	"github.com/d-rk/checkin-system/internal/user"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	CurrentUser(c *gin.Context)
	Login(c *gin.Context)
	Register(c *gin.Context)
}

type handler struct {
	userRepo user.Repository
}

func CreateHandler(userRepo user.Repository) Handler {
	return &handler{userRepo}
}

func (h *handler) CurrentUser(c *gin.Context) {

	u, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to extract user from request context"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": u})
}

func (h *handler) Login(c *gin.Context) {

	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := h.userRepo.VerifyCredentials(c, input.Username, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "username or password is incorrect."})
		return
	}

	token, err := GenerateToken(u.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *handler) Register(c *gin.Context) {
	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//turn password into hash
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	u := user.User{
		CreatedAt:    time.Now(),
		Name:         input.Username,
		PasswordHash: string(hashedPassword),
	}

	_, err = h.userRepo.SaveUser(c, &u)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "registration success"})
}
