var socket = new WebSocket("wss://"+location.hostname+(location.port ? ':'+location.port: '')+"/connect?channel="+getChannelName());

$(document).ready(function(){
  $("#button-submit").click(sendMessage);

  // Callback for enter key to send message
  $("#box-chat").on('keyup', function (e) {
      let code = (e.keyCode ? e.keyCode : e.which);
      if (code == 13) {
          sendMessage();
      }
  });

  $("#box-chat").on('keydown', function (e) {
      let code = (e.keyCode ? e.keyCode : e.which);
      if (code == 13) {
          return false;
      }
  });

  // Callback for clicking the home (back) button
  $("#home-button").click(function() {
    window.location.replace("/index.html");
  });

  socket.onmessage = updateHistory;
  userListWorker();
});

// Send the message contained in the chat box
function sendMessage() {
  let msg = $("#box-chat").val();
  if(msg == "") {
    return;
  }

  let data = {timestamp: Date.now(), message: msg}
  socket.send(JSON.stringify(data));
  $("#box-chat").val('');
}

// Update chat history with incoming data
function updateHistory(event) {
    let data = JSON.parse(event.data);
    let message = parseEmotes(escapeHTML(data.message));

    let chatHTML = `<span class="chat-entry"><span class="chat-username" style="color:`+data.color+`">`+escapeHTML(data.username)+`:</span><span class="chat-text">`+message+`</span></span>`;
    $('#box-history').append(chatHTML);
    $('#box-history').animate({scrollTop : document.getElementById("box-history").scrollHeight }, 0);
}

// Returns the current channel name
function getChannelName() {
  let href = window.location.href;
  let regex = /channel=[^&]+/;
  let match = regex.exec(href);
  if (match === null) {
    return "";
  }
  return match[0].split("=")[1];
}

function escapeHTML(unsafe_str) {
    return unsafe_str.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/\"/g, '&quot;').replace(/\'/g, '&#39;');
}

// Periodically updates the list of online/offline users
function userListWorker() {
  $.ajax({
    url: '/users?channel='+getChannelName(), 
    success: updateUserList,
    complete: function() {
      setTimeout(userListWorker, 5000);
    }
  });
}

function updateUserList(event) {
  let data = JSON.parse(event);
  data.sort(function(a, b) { 
    if(a.name < b.name && a.isOnline && b.isOnline) {
      return -1;
    } else if(a.name > b.name && a.isOnline && b.isOnline) {
      return 1;
    } else if(a.isOnline && !b.isOnline) {
      return -1;
    } else if(!a.isOnline && b.isOnline) {
      return 1;
    } else if(!a.isOnline && !b.isOnline && a.name < b.name) {
      return -1;
    } else {
      return 1;
    }
  });

  $("#chat-user-list").empty()
  for(let i = 0; i < data.length; i++) {
    let userHTML = `<span class="user-list-item"`+(data[i].isOnline ? "" : `style="color:#888888"`)
      +`><span class="user-list-status" style="color:`+(data[i].isOnline ? "green" : "#888888")+`">&bull;</span>`+data[i].name+`</span>`;
    $('#chat-user-list').append(userHTML);
  }
}