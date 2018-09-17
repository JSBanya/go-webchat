var channel = "";

$(document).ready(function(){
  $.ajax({
    url: "/rooms",
    success: populateRooms,
    dataType: "json"
  });

  $(".modal-close").click(function() {
    channel = "";
    $('#chat-modal').hide();
  });

  $("#chat-modal").on('keyup', function (e) {
    let code = (e.keyCode ? e.keyCode : e.which);
    if (code == 13) {
        submit();
    }
  });

  $("#chat-modal").on('keydown', function (e) {
    let code = (e.keyCode ? e.keyCode : e.which);
    if (code == 13) {
        return false;
    }
  });
});

// Display list of rooms
function populateRooms(data) {
  data.sort(function(a,b) {
    if(a.name < b.name) {
      return -1;
    }
    return 1;
  });
  for(let i = 0; i < data.length; i++) {
    let name = data[i].name;
    let roomHTML = `<div class="chat-item" onclick="enterRoom('`+chanIdEncode(name)+`')"><span class="chat-name">`+escapeHTML(name)+`</span><span class="chat-description">`+escapeHTML(data[i].description)+`</span></div>`;
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
    $('#chat-modal').show();
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

function escapeHTML(unsafe_str) {
    return unsafe_str.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/\"/g, '&quot;').replace(/\'/g, '&#39;');
}

function chanIdEncode(str) {
  return str.replace(" ", "_").replace("'", "").replace("\"", "").replace("<", "").replace(">", "").replace("&", "").replace("%", "");
}