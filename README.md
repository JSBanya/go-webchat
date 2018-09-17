# go-webchat
A chat program written in Go with a web-based frontend.

# Installation
As the backend is written in Go, the Go programming language must be installed. Information about installing Go for a particular system can be readily found online and will not be shown here.

To build the project, first fetch all dependencies by running
```
make fetch
```

If all dependencies were successfully installed, the command should terminate with no errors and a folder named 'vendor' will be created in the project directory. 

To modify the program's port number, change the value in cmd/main.go (note that if the port is set to <1024 then the server must be run as root):
```go
var PORT int = 8080
```

If the server's frontend files (located in www/) are moved to a different location, then FILE_PATH value in cmd/main.go must be changed to reflect the new location.
```go
var FILE_PATH string = "www/"
```

By default, the server runs over TLS/HTTPS and thus requires a certificate. It is possible to self sign a certificate, which is not shown here. To run the server without HTTPS, replace the line
```go
log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%d", PORT), "server.pem", "server.key", nil))
```
in cmd/main.go with
```go
log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil))
```

Once all dependencies are fetched and all configuration changes are made, the project can be built by running
```
make build
```

When the server is running, the front-end can be accessed from within a browser (either desktop or mobile) by simply navigating to the server's IP:PORT.

# Adding emotes
Emotes are phrases typed into a chat that are replaced by images (png, jpg, or gif). To add new emotes, the file www/emotes.js must be modified by adding a new entry in the 'emotes' array. Simply add a new image file to the www/images folder and add the emotes to the 'emotes' array as follows:
```javascript
{ Name: "NAME", Image: "/images/NAME.png" },
```
where the 'Name' field is the phrase to be replaced and 'Image' is the location of the image. Note that the Name should ideally be a phrase that is not typically used in normal speech to avoid emotes from appearing unintentionally. It may be necessary to clear the browser cache for a new emote to appear.
