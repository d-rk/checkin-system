package middlewares

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		requestProcessed := make(chan bool)
		go func() {
			c.Next()
			requestProcessed <- true
		}()

		select {
		case <-ctx.Done():
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": fmt.Sprintf("request aborted after %v", timeout)})
			return
		case <-requestProcessed:
		}
	}
}
