package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

var upgrader = websocket.Upgrader{}

type Room struct {
	RoomId  string
	Users   [2]*User
	Message chan string
	mutex   sync.Mutex
}

type User struct {
	ID      string `json:"id"`
	Emotion string `json:"emotion"`
	Conn    *websocket.Conn
	Room    *Room
}

var (
	happyQueue []User
	sadQueue   []User
	angryQueue []User
	mutex      sync.Mutex
	rooms      map[string]*Room
)

func matchingHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatalln("소켓 연결 실패", err.Error())
		return
	}

	userID := generateId()
	user := &User{
		ID:   userID,
		Conn: conn,
	}

	_, emotion, err := conn.ReadMessage()
	if err != nil {
		log.Fatalln("감정 읽기 실패", err.Error())
		return
	}
	fmt.Println(string(emotion))
	user.Emotion = string(emotion)

	Match(*user)

}

func Match(user User) {
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
		user1 := &happyQueue[0]
		user2 := &happyQueue[1]

		roomId := generateId()
		room := &Room{
			RoomId: roomId,
		}
		rooms[roomId] = room

		user1.Room = room
		user2.Room = room

		room.Users[0] = user1
		room.Users[1] = user2

		user1.Conn.WriteMessage(websocket.TextMessage, []byte("matched"))
		user2.Conn.WriteMessage(websocket.TextMessage, []byte("matched"))

		go startChat(user1, user2)
		go startChat(user2, user1)

		happyQueue = happyQueue[2:] //사용자 두명 지워버리
	}

	if len(sadQueue) >= 2 {
		user1 := &sadQueue[0]
		user2 := &sadQueue[1]

		roomId := generateId()
		room := &Room{
			RoomId: roomId,
		}
		rooms[roomId] = room

		user1.Room = room
		user2.Room = room

		room.Users[0] = user1
		room.Users[1] = user2
		user1.Conn.WriteMessage(websocket.TextMessage, []byte("matched"))
		user2.Conn.WriteMessage(websocket.TextMessage, []byte("matched"))

		go startChat(user1, user2)
		go startChat(user2, user1)

		sadQueue = sadQueue[2:] //사용자 두명 지워버리
	}

	if len(angryQueue) >= 2 {
		user1 := &angryQueue[0]
		user2 := &angryQueue[1]

		roomId := generateId()
		room := &Room{
			RoomId: roomId,
		}
		rooms[roomId] = room

		user1.Room = room
		user2.Room = room

		room.Users[0] = user1
		room.Users[1] = user2

		user1.Conn.WriteMessage(websocket.TextMessage, []byte("matched"))
		user2.Conn.WriteMessage(websocket.TextMessage, []byte("matched"))

		go startChat(user1, user2)
		go startChat(user2, user1)

		angryQueue = angryQueue[2:] //사용자 두명 지워버리
	}
}

func startChat(sender *User, reader *User) {
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

func deleteRoom(sender *User, reader *User) {
	sender.Room.mutex.Lock()
	defer sender.Room.mutex.Unlock()
	if sender.Conn != nil {
		sender.Conn.Close()
	}

	if reader.Conn != nil {
		reader.Conn.Close()
	}
	delete(rooms, sender.Room.RoomId)

}

func generateId() string {
	return uuid.New().String()
}

func main() {
	rooms = make(map[string]*Room)
	router := gin.Default()
	router.Static("/static", "./static")
	router.GET("/ws", matchingHandler)
	router.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	router.Run(":8080")
}
