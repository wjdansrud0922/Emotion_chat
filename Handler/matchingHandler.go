package Handler

import (
	"Emotion_chat/model"
	"Emotion_chat/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// 모든 오리진 허용 (필요에 따라 도메인 제한 가능)
		return true
	},
}

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
