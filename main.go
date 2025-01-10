package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
)

type User struct {
	Feel string
	Id   string
	conn *websocket.Conn
}

type MatchingPool struct {
	mu    sync.Mutex
	queue chan *User
}

type Rooms struct {
	Rooms map[string]*Room
	mutex sync.Mutex
}

type Room struct {
	users map[*User]bool
	mu    sync.Mutex
}

var rooms = Rooms{
	Rooms: make(map[string]*Room),
}

func NewMatchingPool(maxQueueSize int) *MatchingPool {
	return &MatchingPool{
		queue: make(chan *User, maxQueueSize),
	}
}

func GenerateNewRoom(user1 *User, user2 *User) string {
	roomID := generateRoomId()
	rooms.mutex.Lock()
	defer rooms.mutex.Unlock()
	rooms.Rooms[roomID] = &Room{
		users: map[*User]bool{
			user1: true,
			user2: true,
		},
	}
	return roomID
}

func (p *MatchingPool) AddUser(user *User) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.queue <- user
}

func (p *MatchingPool) GetUsers() (*User, *User, bool) {
	user1, ok1 := <-p.queue
	user2, ok2 := <-p.queue
	if ok1 && ok2 {
		return user1, user2, true
	}
	return nil, nil, false
}

func generateUserId() string {
	ud, err := uuid.NewRandom()
	if err != nil {
		panic("generateUserId 함수에서 오류 발생")
	}
	return ud.String()
}

func generateRoomId() string {
	ud, err := uuid.NewRandom()
	if err != nil {
		panic("generateRoomId 함수에서 오류 발생")
	}
	return ud.String()
}

func handleWebSocket(c *gin.Context, pool *MatchingPool) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "WebSocket upgrade failed"})
		return
	}

	user := &User{
		Feel: "기쁨",
		Id:   generateUserId(),
		conn: conn,
	}
	pool.AddUser(user)
	fmt.Println("New user added to matching pool", user.Id)
}

func StartMatching(pool *MatchingPool) {
	go func() {
		for {
			user1, user2, ok := pool.GetUsers()
			if !ok {
				time.Sleep(100 * time.Millisecond)
				continue
			}
			GenerateNewRoom(user1, user2)
		}
	}()
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	r := gin.Default()

	pool := NewMatchingPool(100)

	StartMatching(pool)

	r.GET("/ws", func(c *gin.Context) {
		handleWebSocket(c, pool)
	})

	r.Run(":8080")
}
