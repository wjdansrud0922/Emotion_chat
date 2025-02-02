package models

import "sync"

var (
	HappyQueue []User
	SadQueue   []User
	AngryQueue []User
	Mutex      sync.Mutex
)
