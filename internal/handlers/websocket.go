package handlers

import (
	"fmt"
	"io"

	"golang.org/x/net/websocket"
)

type server struct {
	clients map[*websocket.Conn]bool
}

func NewServer() *server {
	return &server{
		clients: make(map[*websocket.Conn]bool),
	}
}

func (s *server) HandleWS(ws *websocket.Conn) {
	s.clients[ws] = true

	fmt.Println("New client connected:", ws.RemoteAddr())

	for {
		msg := make([]byte, 512)
		if _, err := ws.Read(msg); err != nil {
			if err == io.EOF {
				fmt.Println("Client disconnected")
				delete(s.clients, ws)
				ws.Close()
				break
			}
			fmt.Println("Error reading message", err)
			continue
		}

		go s.broadcast(string(msg))
	}

}

func (s *server) broadcast(msg string) {
	for client := range s.clients {
		if _, err := client.Write([]byte(msg)); err != nil {
			fmt.Println("Error broadcasting message", err)
			continue
		}
	}
}
