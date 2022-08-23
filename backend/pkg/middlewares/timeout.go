package middlewares

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
)

func timeoutResponse(timeoutDuration time.Duration) func(*gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusRequestTimeout, gin.H{"error": fmt.Sprintf("request aborted after %v", timeoutDuration)})
	}
}

func TimeoutMiddleware(timeoutDuration time.Duration) gin.HandlerFunc {

	return timeout.New(
		timeout.WithTimeout(timeoutDuration),
		timeout.WithHandler(func(c *gin.Context) {
			c.Next()
		}),
		timeout.WithResponse(timeoutResponse(timeoutDuration)),
	)
}
