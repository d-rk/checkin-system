package handlers

import (
	"fmt"
	"net/http"

	ws "github.com/d-rk/checkin-system/pkg/services/websocket"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func CreateWebsocketHandler(server *ws.Server) func(*gin.Context) {
	return func (ctx *gin.Context) {
		// trust all origin to avoid CORS
		upgrader.CheckOrigin = func(r *http.Request) bool {
			return true
		}

		// upgrades connection to websocket
		conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			ctx.Status(http.StatusInternalServerError)
			return
		}
		defer conn.Close()

		// create new client & add to client list
		client := ws.Client{
			ID:         uuid.Must(uuid.NewRandom()).String(),
			Connection: conn,
		}

		server.Clients = append(server.Clients, client)

		// greet the new client
		greeting := fmt.Sprintf("Server: Welcome! Your ID is %s", client.ID)
		server.PublishClient(&client, ws.Message{Message: greeting, Data: client.ID})


		// message handling
		for {
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				server.RemoveClient(client)
				return
			}
			server.ProcessMessage(client, messageType, p)
		}
	}
}



