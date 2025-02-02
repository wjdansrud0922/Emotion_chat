package handlers

import (
	"Emotion_chat/models"
	"Emotion_chat/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

func MatchingHandler(c *gin.Context) {
	conn, err := models.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatalln("소켓 연결 실패", err.Error())
		return
	}

	userID := utils.GenerateId()
	user := &models.User{
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

	utils.Match(*user)

}
