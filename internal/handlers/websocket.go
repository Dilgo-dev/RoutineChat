package handlers

import (
	"fmt"
	"io"

	"golang.org/x/net/websocket"
)

type server struct {
	rooms map[string]map[*websocket.Conn]bool
}

func NewServer() *server {
	return &server{
		rooms: make(map[string]map[*websocket.Conn]bool),
	}
}

func (s *server) HandleWS(ws *websocket.Conn) {
	roomId := ws.Request().URL.Query().Get("roomId")

	if roomId == "" {
		fmt.Println("No roomId provided")
		ws.Close()
		return
	}

	s.joinRoom(ws, roomId)

	fmt.Printf("New client connected to room %s from client %s 🐡\n", roomId, ws.RemoteAddr())

	for {
		msg := make([]byte, 512)
		if _, err := ws.Read(msg); err != nil {
			if err == io.EOF {
				fmt.Println("Client disconnected")
				s.leaveRoom(ws, roomId)
				ws.Close()
				break
			}
			fmt.Println("Error reading message", err)
			continue
		}

		go s.broadcastToRoom(roomId, string(msg))
	}

}

func (s *server) broadcastToRoom(roomId string, msg string) {
	for client := range s.rooms[roomId] {
		if _, err := client.Write([]byte(msg)); err != nil {
			fmt.Println("Error broadcasting message", err)
			continue
		}
	}
}

func (s *server) joinRoom(ws *websocket.Conn, roomId string) {
	if _, ok := s.rooms[roomId]; !ok {
		s.rooms[roomId] = make(map[*websocket.Conn]bool)
	}
	s.rooms[roomId][ws] = true
}

func (s *server) leaveRoom(ws *websocket.Conn, roomId string) {
	delete(s.rooms[roomId], ws)
	if len(s.rooms[roomId]) == 0 {
		delete(s.rooms, roomId)
	}
}
