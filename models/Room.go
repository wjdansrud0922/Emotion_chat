package models

import "sync"

type Room struct {
	RoomId  string
	Users   [2]*User
	Message chan string
	Mutex   sync.Mutex
}
