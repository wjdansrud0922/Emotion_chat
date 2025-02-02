package main

import (
	"Emotion_chat/handlers"
	"Emotion_chat/models"
	"github.com/gin-gonic/gin"
)

func main() {
	models.Rooms = make(map[string]*models.Room)
	router := gin.Default()
	router.Static("/static", "./static")
	router.GET("/ws", handlers.MatchingHandler)
	router.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	router.Run(":8080")
}
