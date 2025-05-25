package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"regexp"
	"sync"

	"github.com/Dilgo-dev/RoutineChat/internal/models"
	"golang.org/x/net/websocket"
)

type server struct {
	rooms map[string]map[*models.User]bool
	mu    sync.RWMutex
}

func NewServer() *server {
	return &server{
		rooms: make(map[string]map[*models.User]bool),
	}
}

func (s *server) HandleWS(ws *websocket.Conn) {
	roomId := ws.Request().URL.Query().Get("roomId")
	user := &models.User{
		Username: ws.Request().URL.Query().Get("username"),
		RoomId:   roomId,
		Conn:     ws,
	}

	if err := validateInput(user.Username, roomId); err != nil {
		slog.Error("Invalid input", "error", err)
		ws.Close()
		return
	}

	s.joinRoom(user, roomId)

	slog.Info("New client connected to room", "roomId", roomId, "client", ws.RemoteAddr())

	for {
		var msg string
		if err := websocket.Message.Receive(ws, &msg); err != nil {
			if err == io.EOF {
				fmt.Println("Client disconnected")
				s.leaveRoom(user, roomId)
				ws.Close()
				break
			}
			fmt.Println("Error reading message", err)
			continue
		}

		message := models.Message{
			Message:  msg,
			Username: user.Username,
		}

		sendMessageJson, err := json.Marshal(message)
		if err != nil {
			fmt.Println("Error marshaling send message", err)
			continue
		}

		go s.broadcastToRoom(roomId, string(sendMessageJson))
	}

}

func validateInput(username string, roomId string) error {
	if username == "" {
		return errors.New("username is required")
	}
	if roomId == "" {
		return errors.New("roomId is required")
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(roomId) {
		return errors.New("invalid roomId format")
	}
	return nil
}

func (s *server) broadcastToRoom(roomId string, msg string) {
	for client := range s.rooms[roomId] {
		if _, err := client.Conn.Write([]byte(msg)); err != nil {
			fmt.Println("Error broadcasting message", err)
			continue
		}
	}
}

func (s *server) joinRoom(user *models.User, roomId string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.rooms[roomId]; !ok {
		s.rooms[roomId] = make(map[*models.User]bool)
	}
	s.rooms[roomId][user] = true
	s.sendRoomNumber(roomId)
}

func (s *server) leaveRoom(user *models.User, roomId string) {
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
