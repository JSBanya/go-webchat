package main

import (
	"fmt"
	"log"
	"net/http"
)

var PORT int = 8080
var FILE_PATH string = "www/"

func main() {
	// Define paths and handler functions
	http.HandleFunc("/", serveContent)
	http.HandleFunc("/connect", wsConnect)
	http.HandleFunc("/rooms", getRooms)
	http.HandleFunc("/auth", auth)
	http.HandleFunc("/checkauth", isAuthenticated)
	http.HandleFunc("/chat", serveChatroom)
	http.HandleFunc("/users", getUsers)
	http.HandleFunc("/create", createChannelRequest)

	// Create chat map
	chats = make(map[string]*Chatroom)

	// Start server
	log.Printf("Webserver serving files from %s started on port %d.", FILE_PATH, PORT)
	log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%d", PORT), "server.pem", "server.key", nil))
}
