package main

import (
	"Emotion_chat/Handler"
	"Emotion_chat/model"
	"Emotion_chat/util"
	"github.com/gin-gonic/gin"
)

var rooms map[string]*model.Room

func main() {
	rooms = make(map[string]*model.Room)

	util.RoomsInit(rooms)

	router := gin.Default()
	router.Static("/static", "./static")
	router.GET("/ws", func(c *gin.Context) {
		Handler.MatchingHandler(c, rooms)
	})
	router.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	router.Run(":8080")
}
