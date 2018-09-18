package main

import (
	"github.com/gorilla/websocket"
	"sync"
)

type Chatroom struct {
	Name        string `json:"name"`
	Password    string `json:"-"`
	Description string `json:"description"`

	Users map[string]*User `json:"-"`

	History        []*MessageData `json:"-"`
	RecentHistory  []*MessageData `json:"-"`
	HistoryMutex   sync.Mutex     `json:"-"`
	HistoryUpdated bool           `json:"-"`

	FailedPassAttempts map[string]int `json:"-"`
}

type User struct {
	Name   string          `json:"name"`
	Color  string          `json:"-"`
	Online bool            `json:"isOnline"`
	Socket *websocket.Conn `json:"-"`
}

type MessageData struct {
	Timestamp int64  `json:"timestamp"`
	Username  string `json:"username"`
	Message   string `json:"message"`
	Color     string `json:"color"`
}
