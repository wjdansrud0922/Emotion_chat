package utils

import (
	"Emotion_chat/models"
	"github.com/gorilla/websocket"
)

func Matching(queue []models.User) {
	user1 := &queue[0]
	user2 := &queue[1]

	roomId := GenerateId()
	room := &models.Room{
		RoomId: roomId,
	}
	models.Rooms[roomId] = room

	user1.Room = room
	user2.Room = room

	room.Users[0] = user1
	room.Users[1] = user2

	user1.Conn.WriteMessage(websocket.TextMessage, []byte("matched"))
	user2.Conn.WriteMessage(websocket.TextMessage, []byte("matched"))

	go startChat(user1, user2)
	go startChat(user2, user1)

	queue = queue[2:] //사용자 두명 지워버리
}
