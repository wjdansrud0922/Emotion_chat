package utils

import (
	"Emotion_chat/models"
	"github.com/gorilla/websocket"
	"log"
)

func startChat(sender *models.User, reader *models.User) {
	defer deleteRoom(sender, reader)

	for {
		_, msg, err := sender.Conn.ReadMessage()

		if err != nil {
			log.Println("연결 종료", err.Error())
			return
		}

		if err := reader.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Println("연결 종료", err.Error())
		}

	}
}
