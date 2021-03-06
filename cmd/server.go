package main

import (
	"math/rand"
	"strings"
)

var chats map[string]*Chatroom

// Available username colors (randomly selected on user creation)
var colors []string = []string{"#0000ff", "#ff0000", "#009933", "#ee8800", "#cc00cc", "#ee55ee", "#6666ff", "#009999", "#800000", "#e6b800"}

func broadcast(channel string, message *MessageData) {
	for _, user := range chats[channel].Users {
		if !user.Online {
			continue
		}

		user.Socket.WriteJSON(message)
	}

	chats[channel].HistoryMutex.Lock()
	chats[channel].History = append(chats[channel].History, message)
	chats[channel].HistoryMutex.Unlock()
}

// Creates a new chatroom
func createChatroom(name string, password string, description string) {
	c := &Chatroom{Name: name, Password: password, Description: description}
	c.Users = make(map[string]*User)
	c.FailedPassAttempts = make(map[string]int)

	id := chanIdEncode(name)
	chats[id] = c
}

func getRandomColor() string {
	return colors[rand.Intn(len(colors))]
}

// Encode the channel name to a format used internally
func chanIdEncode(name string) string {
	name = strings.Replace(name, " ", "_", -1)
	name = strings.Replace(name, "'", "", -1)
	name = strings.Replace(name, "\"", "", -1)
	name = strings.Replace(name, "<", "", -1)
	name = strings.Replace(name, ">", "", -1)
	name = strings.Replace(name, "&", "", -1)
	name = strings.Replace(name, "%%", "", -1)
	return name
}
