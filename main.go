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
	RoomId string
	Users  [2]*User
	mutex  sync.Mutex
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
		matchingUsers(happyQueue)
	}

	if len(sadQueue) >= 2 {
		matchingUsers(sadQueue)
	}

	if len(angryQueue) >= 2 {
		matchingUsers(angryQueue)
	}
}

func matchingUsers(queue []User) {
	user1 := &queue[0]
	user2 := &queue[1]

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

	queue = queue[2:] //사용자 두명 지워버리
}

func startChat(sender *User, reader *User) {
	defer func() {
		sender.Conn.Close()
		reader.Conn.Close()
		deleteRoom(sender, reader)

	}()

	for {
		_, msg, err := sender.Conn.ReadMessage()
		if err != nil {
			log.Println("Read에러")
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

func deleteRoom(sender *User, reader *User) {
	sender.Room.mutex.Lock()
	defer sender.Room.mutex.Unlock()

	if sender.Conn == nil || reader.Conn == nil {
		delete(rooms, sender.Room.RoomId)
	}

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
		c.File("./index.html")
	})

	router.Run(":8080")
}
