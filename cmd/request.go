package main

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

const (
	MAX_PASSWORD_ATTEMPTS = 10
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Serve all directories not otherwise specifically handled
func serveContent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	logRequest(r)

	path := filepath.Clean(r.URL.Path)
	if strings.Trim(path, "/") == "images" || strings.Trim(path, "/") == "audio" {
		http.Error(w, "403 Forbidden", http.StatusForbidden)
		return
	}

	http.ServeFile(w, r, FILE_PATH+path)
}

// Serve main chat directory
func serveChatroom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	logRequest(r)

	channel := chanIdEncode(strings.Trim(r.URL.Query().Get("channel"), " "))
	if channel == "" || chats[channel] == nil {
		http.Error(w, "403 Forbidden (Channel not found)", http.StatusForbidden)
		return
	}

	cookie, _ := r.Cookie(channel)
	if cookie == nil || chats[channel].Users[cookie.Value] == nil {
		http.Error(w, "403 Forbidden (Not authenticated)", http.StatusForbidden)
		return
	}

	http.ServeFile(w, r, FILE_PATH+"chat.html")
}

// Handle websocket connection for chatrooms
func wsConnect(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")

	channel := chanIdEncode(strings.Trim(r.URL.Query().Get("channel"), " "))
	if channel == "" || chats[channel] == nil {
		http.Error(w, "403 Forbidden (Channel not found)", http.StatusForbidden)
		return
	}

	cookie, _ := r.Cookie(channel)
	if cookie == nil || chats[channel].Users[cookie.Value] == nil {
		http.Error(w, "403 Forbidden (Not authenticated)", http.StatusForbidden)
		return
	}
	sid := cookie.Value

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Unable to upgrade websocket connection:", err)
		http.Error(w, fmt.Sprintf("Unable to handle websocket: %s", err), http.StatusInternalServerError)
		return
	}
	defer c.Close()

	chats[channel].Users[sid].Online = true
	chats[channel].Users[sid].Socket = c

	// Send the entire history to the user on first connection
	chats[channel].HistoryMutex.Lock()
	for _, data := range chats[channel].History {
		c.WriteJSON(data)
	}
	chats[channel].HistoryMutex.Unlock()

	// Handle incoming messages
	for {
		var data MessageData

		err := c.ReadJSON(&data)
		if err != nil {
			log.Println("Unable to read client data:", err)
			c.Close()
			chats[channel].Users[sid].Socket = nil
			chats[channel].Users[sid].Online = false
			break
		}

		if data.Message != "" && len(data.Message) <= 500 {
			chats[channel].HistoryMutex.Lock()

			data.Username = chats[channel].Users[sid].Name
			data.Color = chats[channel].Users[sid].Color
			data.Timestamp = (time.Now().UnixNano() / int64(time.Millisecond))

			chats[channel].RecentHistory = append(chats[channel].RecentHistory, &data)
			chats[channel].HistoryUpdated = true

			chats[channel].HistoryMutex.Unlock()
		}
	}
}

// Authenticate a user's connection for a chatroom via entered username and password
// All usernames are valid if they are not already in use and do not violate certain requirements (e.g. A-Za-Z0-9)
// Only one password is valid for a given chatroom
// If the user is authenticated, this function will complete by setting the user's session ID as a cookie
func auth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

	defer r.Body.Close()

	channel := chanIdEncode(r.FormValue("channel"))
	username := strings.Trim(r.FormValue("username"), " ")
	password := r.FormValue("password")

	if len(strings.Split(r.RemoteAddr, ":")) != 2 {
		http.Error(w, "Bad Remote Address", http.StatusInternalServerError)
		return
	}
	ip := strings.Split(r.RemoteAddr, ":")[0]

	// Validate username
	valid := true
	if username == "" || len(username) > 20 {
		valid = false
	}

	if valid {
		for _, r := range username {
			if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') {
				valid = false
			}
		}
	}

	if valid {
		for _, user := range chats[channel].Users {
			if user.Name == username {
				// User already exists
				http.Error(w, "That username already exists.", http.StatusBadRequest)
				return
			}
		}
	}

	if !valid {
		http.Error(w, "Invalid Username.", http.StatusBadRequest)
		return
	}

	// Check if channel exists
	if chats[channel] == nil {
		http.Error(w, "That channel does not exist! How did you get here?", http.StatusTeapot)
		return
	}

	if chats[channel].FailedPassAttempts[ip] >= MAX_PASSWORD_ATTEMPTS {
		log.Printf("\u001b[31mIP %s blocked from entered %s (Max attempts reached)\u001b[0m", ip, channel)
		http.Error(w, "Max number of password attempts reached. Please contact an administrator.", http.StatusBadRequest)
		return
	}

	// Check password for channel
	if chats[channel].Password != password {
		chats[channel].FailedPassAttempts[ip]++
		http.Error(w, "Invalid Password.", http.StatusBadRequest)
		return
	}

	// Entry valid
	// Generate session ID
	sid := make([]byte, 16)
	_, err := rand.Read(sid)
	if err != nil {
		http.Error(w, "Unable to create session", http.StatusInternalServerError)
		return
	}
	sidEncoded := hex.EncodeToString(sid)

	user := User{Name: username, Color: getRandomColor(), Online: false}
	chats[channel].Users[sidEncoded] = &user

	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := http.Cookie{Name: channel, Value: sidEncoded, Expires: expiration, HttpOnly: true}
	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
}

// Check if a user is already authenticated by checking their cookie session ID
func isAuthenticated(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")

	channel := chanIdEncode(strings.Trim(r.URL.Query().Get("channel"), " "))
	if channel == "" || chats[channel] == nil {
		http.Error(w, "", http.StatusForbidden)
		return
	}

	cookie, err := r.Cookie(channel)
	if err != nil {
		http.Error(w, "", http.StatusForbidden)
		return
	}

	if chats[channel] == nil || chats[channel].Users[cookie.Value] == nil {
		http.Error(w, "", http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Serves a list of available chatrooms
func getRooms(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	logRequest(r)

	list := make([]*Chatroom, len(chats))
	i := 0
	for _, c := range chats {
		list[i] = c
		i++
	}

	content, err := json.Marshal(list)
	if err != nil {
		log.Printf("Error marshalling chat room list: %s", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	io.Copy(w, bytes.NewBuffer(content))
}

// Serves a list of users (online or offline) for the given chatroom
// Ensures that the connecting client is authenticated before serving information
func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")

	channel := chanIdEncode(strings.Trim(r.URL.Query().Get("channel"), " "))
	if channel == "" || chats[channel] == nil {
		http.Error(w, "", http.StatusForbidden)
		return
	}

	cookie, err := r.Cookie(channel)
	if err != nil {
		http.Error(w, "", http.StatusForbidden)
		return
	}

	if chats[channel] == nil || chats[channel].Users[cookie.Value] == nil {
		http.Error(w, "", http.StatusForbidden)
		return
	}

	users := make([]*User, len(chats[channel].Users))
	i := 0
	for _, user := range chats[channel].Users {
		users[i] = user
		i++
	}

	content, err := json.Marshal(users)
	if err != nil {
		log.Printf("Error marshalling chat room user list: %s", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	io.Copy(w, bytes.NewBuffer(content))
}

// Create a channel based on the received request
func createChannelRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

	name := strings.Trim(r.FormValue("name"), " ")
	password := strings.Trim(r.FormValue("password"), " ")
	description := strings.Trim(r.FormValue("desc"), " ")

	// Validate name
	valid := true
	if name == "" || len(name) > 30 {
		valid = false
	}

	if valid {
		for _, r := range name {
			if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') && r != '\'' && r != ' ' {
				valid = false
			}
		}
	}

	if !valid {
		http.Error(w, "Invalid channel name.", http.StatusBadRequest)
		return
	}

	if chats[name] != nil {
		http.Error(w, "A room with that name already exists.", http.StatusBadRequest)
		return
	}

	// Validate description
	if len(description) > 100 {
		http.Error(w, "Invalid description.", http.StatusBadRequest)
		return
	}

	// Validate password
	if len(password) > 30 {
		http.Error(w, "Invalid password.", http.StatusBadRequest)
		return
	}

	// Create room
	createChatroom(name, password, description)
	log.Printf("\u001b[32mNew channel \"%s\" created by %s\u001b[0m", name, r.RemoteAddr)

	w.WriteHeader(http.StatusOK)
}

// Logs a connection to the server
// Should be called by major handler functions to maintain detailed logs.
func logRequest(r *http.Request) {
	log.Printf("%s %s request on %s from %s.", r.Proto, r.Method, r.URL.Path, r.RemoteAddr)
}
