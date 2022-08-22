package websocket

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type Server struct {
	Clients []Client
}

// each client consists of auto-generated ID & connection
type Client struct {
	ID         string
	Connection *websocket.Conn
}

// type for a valid message.
type Message struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func (s *Server) send(client *Client, message []byte) {
	client.Connection.WriteMessage(1, message)
}

func (s *Server) RemoveClient(client Client) {
	// Read all client
	for i := 0; i < len(s.Clients); i++ {
		if client.ID == (s.Clients)[i].ID {
			// If found, remove client
			if i == len(s.Clients)-1 {
				// if it's stored as the last element, crop the array length
				s.Clients = (s.Clients)[:len(s.Clients)-1]
			} else {
				// if it's stored in between elements, overwrite the element and reduce iterator to prevent out-of-bound
				s.Clients = append((s.Clients)[:i], (s.Clients)[i+1:]...)
				i--
			}
		}
	}
}

func (s *Server) ProcessMessage(client Client, messageType int, payload []byte) {

	log.Printf("received message %s", payload)
	s.PublishClient(&client, Message{Message: "cannot handle client message", Data: payload})
}

func (s *Server) Publish(message any) error {

	rawMessage, err := json.Marshal(message)

	if err != nil {
		return err
	}

	// send to clients
	for _, client := range s.Clients {
		s.send(&client, rawMessage)
	}

	return nil
}

func (s *Server) PublishClient(client *Client, message any) error {

	rawMessage, err := json.Marshal(message)

	if err != nil {
		return err
	}

	s.send(client, rawMessage)

	return nil
}
