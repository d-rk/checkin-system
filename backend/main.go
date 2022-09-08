package main

import (
	"time"

	"github.com/d-rk/checkin-system/pkg/handlers"
	"github.com/d-rk/checkin-system/pkg/middlewares"
	"github.com/d-rk/checkin-system/pkg/services/database"
	"github.com/d-rk/checkin-system/pkg/services/websocket"
	"github.com/gin-gonic/gin"
)

func main() {

	db := database.Connect()

	websocket := &websocket.Server{}

	userHandler := handlers.CreateUserHandler(db, websocket)
	checkInHandler := handlers.CreateCheckInHandler(db, websocket)
	websocketHandler := handlers.CreateWebsocketHandler(websocket)

	r := gin.Default()
	r.Use(middlewares.CORSMiddleware())

	api := r.Group("/api/v1")
	api.Use(middlewares.TimeoutMiddleware(5 * time.Second))

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
	api.POST("/checkins", checkInHandler.AddCheckIn)
	api.DELETE("/checkins/:id", checkInHandler.DeleteCheckIn)

	api.GET("/checkins/dates", checkInHandler.ListCheckInDates)

	api.GET("/websocket", websocketHandler)

	r.Run(":8080")
}
