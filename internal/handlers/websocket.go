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

		fmt.Println("Message received:", string(msg))
		ws.Write([]byte("Message received from server 🎣"))
	}

}
