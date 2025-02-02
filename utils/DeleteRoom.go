package utils

import "Emotion_chat/models"

func deleteRoom(sender *models.User, reader *models.User) {
	sender.Room.Mutex.Lock()
	defer sender.Room.Mutex.Unlock()
	if sender.Conn != nil {
		sender.Conn.Close()
	}

	if reader.Conn != nil {
		reader.Conn.Close()
	}
	delete(models.Rooms, sender.Room.RoomId)

}
