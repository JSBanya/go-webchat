package main

import (
	"math/rand"
	"strings"
	"time"
)

var chats map[string]*Chatroom

// Available username colors (randomly selected on user creation)
var colors []string = []string{"#0000ff", "#ff0000", "#009933", "#ee8800", "#cc00cc", "#ee55ee", "#6666ff", "#009999", "#800000", "#e6b800"}

// Periodically update all connected users with data send from other users
// Should be called as a go routine for all chatrooms
func handleOutgoing(channel string) {
	for {
		time.Sleep(100 * time.Millisecond)
		if chats[channel].HistoryUpdated {
			chats[channel].HistoryMutex.Lock()
			for _, user := range chats[channel].Users {
				if !user.Online {
					continue
				}

				for _, data := range chats[channel].RecentHistory {
					user.Socket.WriteJSON(data)
				}
			}
			chats[channel].History = append(chats[channel].History, chats[channel].RecentHistory[:]...)
			chats[channel].RecentHistory = []*MessageData{}
			chats[channel].HistoryMutex.Unlock()
		}
	}
}

// Creates a new chatroom
// TODO: Check if chatroom already exists; Allow creation of chatrooms from frontend
func createChatroom(name string, password string, description string) {
	c := &Chatroom{Name: name, Password: password, Description: description}
	c.Users = make(map[string]*User)
	c.FailedPassAttempts = make(map[string]int)

	id := chanIdEncode(name)
	chats[id] = c
	go handleOutgoing(id)
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
