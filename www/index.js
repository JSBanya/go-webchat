var channel = "";

$(document).ready(function(){
  $.ajax({
    url: "/rooms",
    success: populateRooms,
    dataType: "json"
  });

  $("#modal-close-auth").click(function() {
    channel = "";
    $('#chat-modal-auth').hide();
  });


  $("#modal-close-create").click(function() {
    $('#chat-modal-create').hide();
  });

  // Authentication modal 'enter' key callbacks
  $("#chat-modal-auth").on('keyup', function (e) {
    let code = (e.keyCode ? e.keyCode : e.which);
    if (code == 13) {
        submit();
    }
  });

  $("#chat-modal-auth").on('keydown', function (e) {
    let code = (e.keyCode ? e.keyCode : e.which);
    if (code == 13) {
        return false;
    }
  });

  // Creation modal 'enter' key callbacks
  $("#chat-modal-create").on('keyup', function (e) {
    let code = (e.keyCode ? e.keyCode : e.which);
    if (code == 13) {
        createRoom();
    }
  });

  $("#chat-modal-create").on('keydown', function (e) {
    let code = (e.keyCode ? e.keyCode : e.which);
    if (code == 13) {
        return false;
    }
  });
});

// Display list of rooms
function populateRooms(data) {
  if(data.length == 0) {
    $('#empty-warning').show();
    return;
  }

  data.sort(function(a,b) {
    if(a.name < b.name) {
      return -1;
    }
    return 1;
  });
  for(let i = 0; i < data.length; i++) {
    let name = data[i].name;
    let description = data[i].description;
    let roomHTML = `<div class="chat-item" onclick="enterRoom('`+chanIdEncode(name)+`')"><span class="chat-name">`+escapeHTML(name)+`</span><span class="chat-description">`+escapeHTML(description)+`</span></div>`;
    $('#chat-list').append(roomHTML);
  }
}

// Check if user is authenticated for the given room
// Display login if user is not authenticated
function enterRoom(name) {
  $.post("/checkauth?channel="+name, function() {
    window.location.replace("/chat?channel="+encodeURI(name));
  }).fail(function() {
    channel = name;
    $('#chat-modal-auth').show();
  });
}

// Submit user entered information
function submit() {
  let username = $("#modal-username").val().trim();
  let password = $("#modal-password").val().trim();
  if(username == "" || password == "") {
    return;
  }

  let data = "channel="+channel+"&username="+username+"&password="+password;

  $.ajax({
    type: "POST",
    url: "/auth",
    data: data
  }).done(function() {
    window.location.replace("/chat?channel="+encodeURI(channel));
  }).fail(function(xhr, textStatus, error) {
      alert(xhr.responseText);
  });

  $("#modal-password").val("");
}

function showRoomCreation() {
  $('#chat-modal-create').show();
}

function createRoom() {
  let name = $("#modal-room-name").val().trim();
  let description = $("#modal-room-description").val().trim();
  let password = $("#modal-room-password").val().trim();
  let password2 = $("#modal-room-password-confirm").val().trim();
  if(name == "" || password == "" || password2 == "") {
    return;
  }

  if(password != password2) {
    alert("Passwords do not match.");
    $("#modal-room-password").val("");
    $("#modal-room-password-confirm").val("");
    return;
  }

  let data = "name="+name+"&password="+password;
  if(description != "") {
    data += "&desc="+description.replace(/;/g, '').replace(/&/g, '');
  }

  $.ajax({
    type: "POST",
    url: "/create",
    data: data
  }).done(function() {
    window.location.replace("/index.html");
  }).fail(function(xhr, textStatus, error) {
      alert(xhr.responseText);
  });
}

function escapeHTML(unsafe_str) {
    return unsafe_str.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/\"/g, '&quot;').replace(/\'/g, '&#39;');
}

function chanIdEncode(str) {
  return str.replace(" ", "_").replace("'", "").replace("\"", "").replace("<", "").replace(">", "").replace("&", "").replace("%", "");
}