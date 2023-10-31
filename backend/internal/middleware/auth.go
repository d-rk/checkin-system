package middleware

import (
	"github.com/d-rk/checkin-system/internal/auth"
	"github.com/d-rk/checkin-system/internal/user"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Auth(userRepo user.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := auth.ValidateToken(c)
		if err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		u, err := userRepo.GetUserByID(c, claims.UserID)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.Set("user", u)

		c.Next()
	}
}
