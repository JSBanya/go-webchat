// Add new emotes here
// Adding the name and location of the image here is sufficient to have it display in chat when a user types the name
var emotes = [
	{ Name: "HackerMan", Image: "/images/hackerman.png" },
	{ Name: "Kappa", Image: "/images/kappa.png" },
	{ Name: "ditto", Image: "/images/ditto.gif" },
];

// Replaces all instances of emote[x].Name with the appropriate image for the given string
function parseEmotes(message) {
	let parsedMessage = message;
	for(let i = 0; i < emotes.length; i++) {
		let regex = new RegExp("\\b"+emotes[i].Name+"\\b", "g");
		parsedMessage = parsedMessage.replace(regex, `<img src="`+emotes[i].Image+`" alt="`+emotes[i].Name+`" title="`+emotes[i].Name+`" class="emote" width="25" height="25">`);
	}

	return parsedMessage;
}