package main

import (
	"context"
	"github.com/d-rk/checkin-system/internal/checkin"
	"github.com/d-rk/checkin-system/internal/database"
	"github.com/d-rk/checkin-system/internal/user"
	"github.com/d-rk/checkin-system/internal/websocket"
	"time"

	"github.com/d-rk/checkin-system/internal/middleware"
	"github.com/gin-gonic/gin"
)

func main() {

	db := database.Connect()

	ws := &websocket.Server{}

	checkInRepo := checkin.NewRepo(db)
	userRepo := user.NewRepo(db)

	userHandler := user.CreateHandler(userRepo, ws)
	checkInHandler := checkin.CreateHandler(checkInRepo, userRepo, ws)
	websocketHandler := websocket.CreateHandler(ws)

	checkInHandler.DeleteOldCheckIns(context.Background())

	r := gin.Default()
	r.Use(middleware.Cors())

	api := r.Group("/api/v1")
	api.Use(middleware.Timeout(5 * time.Second))

	api.GET("/users", userHandler.ListUsers)
	api.POST("/users", userHandler.AddUser)
	api.GET("/users/:id", userHandler.GetUserByID)
	api.PUT("/users/:id", userHandler.UpdateUser)
	api.DELETE("/users/all", userHandler.DeleteAllUsers)
	api.DELETE("/users/:id", userHandler.DeleteUser)
	api.GET("/users/:id/checkins", checkInHandler.ListUserCheckIns)
	api.DELETE("/users/:id/checkins", checkInHandler.DeleteUserCheckIns)

	api.GET("/checkins", checkInHandler.ListCheckIns)
	api.GET("/checkins/per-day", checkInHandler.ListCheckInsPerDay)
	api.GET("/checkins/all", checkInHandler.ListAllCheckIns)
	api.POST("/checkins", checkInHandler.AddCheckIn)
	api.DELETE("/checkins/:id", checkInHandler.DeleteCheckIn)

	api.GET("/checkins/dates", checkInHandler.ListCheckInDates)

	r.GET("/websocket", websocketHandler)

	r.Run(":8080")
}
