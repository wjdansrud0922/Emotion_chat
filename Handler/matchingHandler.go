package Handler

import (
	"Emotion_chat/model"
	"Emotion_chat/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
)

var upgrader = websocket.Upgrader{}

func MatchingHandler(c *gin.Context, rooms map[string]*model.Room) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {

		log.Fatalln("소켓 연결 실패", err.Error())
		return
	}

	userID := util.GenerateId()
	user := &model.User{
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

	util.Match(*user, rooms)

}
