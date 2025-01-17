package model

import "sync"

type Room struct {
	RoomId string
	Users  [2]*User
	Mutex  sync.Mutex
}
