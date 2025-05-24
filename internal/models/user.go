package models

import "golang.org/x/net/websocket"

type User struct {
	Username string
	RoomId   string
	Conn     *websocket.Conn
}
