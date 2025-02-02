package utils

import (
	"Emotion_chat/models"
)

func Match(user models.User) {
	models.Mutex.Lock()
	defer models.Mutex.Unlock()

	userP := &user

	switch userP.Emotion {
	case "happy":
		models.HappyQueue = append(models.HappyQueue, *userP)
	case "sad":
		models.SadQueue = append(models.SadQueue, *userP)
	case "angry":
		models.AngryQueue = append(models.AngryQueue, *userP)
	}

	//각 매칭 큐마다 2명 이상이면 짝찌
	if len(models.HappyQueue) >= 2 {
		Matching(models.HappyQueue)
	}

	if len(models.SadQueue) >= 2 {
		Matching(models.SadQueue)
	}

	if len(models.AngryQueue) >= 2 {
		Matching(models.AngryQueue)
	}
}
