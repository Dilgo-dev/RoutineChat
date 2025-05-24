package handlers

import (
	"encoding/json"
	"fmt"
	"io"

	"golang.org/x/net/websocket"
)

type User struct {
	Username string
	RoomId   string
	Conn     *websocket.Conn
}

type server struct {
	rooms map[string]map[*User]bool
}

func NewServer() *server {
	return &server{
		rooms: make(map[string]map[*User]bool),
	}
}

type sendMessage struct {
	Message  string `json:"message"`
	Username string `json:"username"`
}

func (s *server) HandleWS(ws *websocket.Conn) {
	roomId := ws.Request().URL.Query().Get("roomId")
	user := &User{
		Username: ws.Request().URL.Query().Get("username"),
		RoomId:   roomId,
		Conn:     ws,
	}

	if user.Username == "" {
		fmt.Println("No username provided")
		ws.Close()
		return
	}

	if roomId == "" {
		fmt.Println("No roomId provided")
		ws.Close()
		return
	}

	s.joinRoom(user, roomId)

	fmt.Printf("New client connected to room %s from client %s 🐡\n", roomId, ws.RemoteAddr())

	for {
		msg := make([]byte, 512)
		if _, err := ws.Read(msg); err != nil {
			if err == io.EOF {
				fmt.Println("Client disconnected")
				s.leaveRoom(user, roomId)
				ws.Close()
				break
			}
			fmt.Println("Error reading message", err)
			continue
		}

		sendMessage := sendMessage{
			Message:  string(msg),
			Username: user.Username,
		}

		sendMessageJson, err := json.Marshal(sendMessage)
		if err != nil {
			fmt.Println("Error marshaling send message", err)
			continue
		}

		go s.broadcastToRoom(roomId, string(sendMessageJson))
	}

}

func (s *server) broadcastToRoom(roomId string, msg string) {
	for client := range s.rooms[roomId] {
		if _, err := client.Conn.Write([]byte(msg)); err != nil {
			fmt.Println("Error broadcasting message", err)
			continue
		}
	}
}

func (s *server) joinRoom(user *User, roomId string) {
	if _, ok := s.rooms[roomId]; !ok {
		s.rooms[roomId] = make(map[*User]bool)
	}
	s.rooms[roomId][user] = true
	s.sendRoomNumber(roomId)
}

func (s *server) leaveRoom(user *User, roomId string) {
	delete(s.rooms[roomId], user)
	fmt.Printf("Client %s left room %s 🐡\n", user.Conn.RemoteAddr(), roomId)
	s.sendRoomNumber(roomId)
	if len(s.rooms[roomId]) == 0 {
		delete(s.rooms, roomId)
		fmt.Printf("Room %s is empty, removing it 🐡\n", roomId)
	}
}

type roomInfo struct {
	Number int    `json:"number"`
	RoomId string `json:"roomId"`
}

func (s *server) sendRoomNumber(roomId string) {
	info := roomInfo{
		Number: len(s.rooms[roomId]) - 1,
		RoomId: roomId,
	}
	jsonData, err := json.Marshal(info)
	if err != nil {
		fmt.Println("Error marshaling room info:", err)
		return
	}
	s.broadcastToRoom(roomId, string(jsonData))
}
