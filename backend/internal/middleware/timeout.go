package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	timeout "github.com/vearne/gin-timeout"
)

func Timeout(timeoutDuration time.Duration) gin.HandlerFunc {

	defaultMsg := fmt.Sprintf(`{"error": "request aborted after %v"}`, timeoutDuration)

	return timeout.Timeout(
		timeout.WithTimeout(timeoutDuration),
		timeout.WithErrorHttpCode(http.StatusRequestTimeout),
		timeout.WithDefaultMsg(defaultMsg))
}
