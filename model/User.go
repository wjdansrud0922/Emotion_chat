package model

import "github.com/gorilla/websocket"

type User struct {
	ID      string `json:"id"`
	Emotion string `json:"emotion"`
	Conn    *websocket.Conn
	Room    *Room
}
