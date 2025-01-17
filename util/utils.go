package util

import (
	"Emotion_chat/model"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

var (
	happyQueue []model.User
	sadQueue   []model.User
	angryQueue []model.User
	mutex      sync.Mutex
)

func GenerateId() string {
	return uuid.New().String()
}

func RoomsInit(rooms map[string]*model.Room) {
	rooms = make(map[string]*model.Room)
}

func Match(user model.User, rooms map[string]*model.Room) {
	mutex.Lock()
	defer mutex.Unlock()

	userP := &user

	switch userP.Emotion {
	case "happy":
		happyQueue = append(happyQueue, *userP)
	case "sad":
		sadQueue = append(sadQueue, *userP)
	case "angry":
		angryQueue = append(angryQueue, *userP)
	}

	//각 매칭 큐마다 2명 이상이면 짝찌
	if len(happyQueue) >= 2 {
		matchingUsers(happyQueue, rooms)
	}

	if len(sadQueue) >= 2 {
		matchingUsers(sadQueue, rooms)
	}

	if len(angryQueue) >= 2 {
		matchingUsers(angryQueue, rooms)
	}
}

func matchingUsers(queue []model.User, rooms map[string]*model.Room) {
	user1 := &queue[0]
	user2 := &queue[1]

	roomId := GenerateId()
	room := &model.Room{
		RoomId: roomId,
	}
	rooms[roomId] = room

	user1.Room = room
	user2.Room = room

	room.Users[0] = user1
	room.Users[1] = user2

	user1.Conn.WriteMessage(websocket.TextMessage, []byte("matched"))
	user2.Conn.WriteMessage(websocket.TextMessage, []byte("matched"))

	go startChat(user1, user2, rooms)
	go startChat(user2, user1, rooms)

	queue = queue[2:] //사용자 두명 지워버리
}

func startChat(sender *model.User, reader *model.User, rooms map[string]*model.Room) {
	defer func() {
		sender.Conn.Close()
		reader.Conn.Close()
		deleteRoom(sender, reader, rooms)

	}()

	for {
		_, msg, err := sender.Conn.ReadMessage()
		if err != nil {
			log.Println("Read에러", err.Error())
			return
		}
		if string(msg) != "matched" {
			//if err := sender.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			//	log.Println("sender Write에러")
			//}

			if err := reader.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Println("reader Write에러")
			}
		}

	}
}

func deleteRoom(sender *model.User, reader *model.User, rooms map[string]*model.Room) {
	sender.Room.Mutex.Lock()
	defer sender.Room.Mutex.Unlock()

	if sender.Conn == nil || reader.Conn == nil {
		delete(rooms, sender.Room.RoomId)
	}

}
