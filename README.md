# go-webchat
A chat program written in Go with a web-based frontend.

# Installation
As the backend is written in Go, the Go programming language must be installed. Information about installing Go for a particular system can be readily found online and will not be shown here.

To build the project, first fetch all dependencies by running
```
make fetch
```

If all dependencies were successfully installed, the command should terminate with no errors and a folder named 'vendor' will be created in the project directory. 

The project can be built by running
```
make build
```

Use the -h or -help flag to see available command line arguments. When the server is running, the front-end can be accessed from within a browser (either desktop or mobile) by simply navigating to the server's IP.

# Adding emotes
Emotes are phrases typed into a chat that are replaced by images (png, jpg, or gif). To add new emotes, the file www/emotes.js must be modified by adding a new entry in the 'emotes' array. Simply add a new image file to the www/images folder and add the emotes to the 'emotes' array as follows:
```javascript
{ Name: "NAME", Image: "/images/NAME.png" },
```
where the 'Name' field is the phrase to be replaced and 'Image' is the location of the image. Note that the Name should ideally be a phrase that is not typically used in normal speech to avoid emotes from appearing unintentionally. It may be necessary to clear the browser cache for a new emote to appear.
