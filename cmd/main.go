package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
)

var (
	port        int
	filePath    string
	useHTTPS    bool
	certificate string
	serverKey   string
)

func main() {
	// Define args
	portPtr := flag.Int("port", 443, "Port number.")
	filePathPtr := flag.String("path", "www", "Path to front-end files.")
	useHTTPSPtr := flag.Bool("https", true, "Use HTTPs.")
	certificatePtr := flag.String("cert", "server.pem", "HTTPs server certificate.")
	serverKeyPtr := flag.String("key", "server.key", "HTTPs server key.")

	help := flag.Bool("h", false, "Show help information.")
	help2 := flag.Bool("help", false, "Show help information.")

	flag.Parse()

	if *help || *help2 {
		flag.PrintDefaults()
		return
	}

	// Read args
	port = *portPtr
	filePath = strings.TrimSuffix(*filePathPtr, "/") + "/"
	useHTTPS = *useHTTPSPtr
	certificate = *certificatePtr
	serverKey = *serverKeyPtr

	// Define paths and handler functions
	http.HandleFunc("/", serveContent)
	http.HandleFunc("/connect", wsConnect)
	http.HandleFunc("/rooms", getRooms)
	http.HandleFunc("/auth", auth)
	http.HandleFunc("/checkauth", isAuthenticated)
	http.HandleFunc("/chat", serveChatroom)
	http.HandleFunc("/users", getUsers)
	http.HandleFunc("/create", createChannelRequest)

	// Redirect HTTP to HTTPS
	if useHTTPS && port == 443 {
		go http.ListenAndServe(":80", http.HandlerFunc(redirect))
	}

	// Create chat map
	chats = make(map[string]*Chatroom)

	// Start server
	log.Printf("Webserver serving files from %s started on port %d. (HTTPS: %v)", filePath, port, useHTTPS)
	if useHTTPS {
		log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%d", port), certificate, serverKey, nil))
	} else {
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
	}
}
