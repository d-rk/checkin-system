package websocket

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func CreateHandler(server *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// trust all origin to avoid CORS
		upgrader.CheckOrigin = func(r *http.Request) bool {
			return true
		}

		// upgrades connection to websocket
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer conn.Close()

		// create new client & add to client list
		client := Client{
			ID:         uuid.Must(uuid.NewRandom()).String(),
			Connection: conn,
		}

		server.Clients = append(server.Clients, client)

		// greet the new client
		greeting := fmt.Sprintf("Server: Welcome! Your ID is %s", client.ID)
		server.PublishClient(&client, Message{Message: greeting, Data: client.ID})

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
