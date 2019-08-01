(function($) {

var socket;

$("#js-join").click(function(){
  StartWs($("#js-token").val());
  $("#js-token").addClass("hide");
  $("#js-join").addClass("hide");
  $("#js-room").removeClass("hide");
  $("#js-join-room").removeClass("hide");
});

$("#js-join-room").click(function(){
  JoinRoom($("#js-room").val());
});

$("#js-logs").on("click", ".js-log-view", function(){
  $(this).toggleClass("big-log");
  document.querySelectorAll('.big-log pre').forEach((block) => {
    hljs.highlightBlock(block);
  });
});

function StartWs(token) {
  var full = 'ws://'+location.hostname+(location.port ? ':'+location.port: '');
  // Create WebSocket connection.
  socket = new WebSocket(`${full}/ws/${token}/join`);

  // Listen for messages
  socket.addEventListener('message', function (event) {
      const log = JSON.parse(event.data);
      AddLogs(log);
  });
}

function AddLogs(log) {
  var msgHTML = "";
  if(isJson(log.m)){
    const message = JSON.stringify(JSON.parse(log.m),null,2);
    msgHTML = `<pre><code class="json">${message}</code></pre>`
  } else {
    msgHTML = `<div class="other-message">${log.m}</div>`
  }
  const html = "<div class=\"js-log-view log\">"+
    `<div>${log.c}</div>` +
    msgHTML +
    "</div>"
  $("#js-logs").append(html);
}

function isJson(str) {
    try {
        JSON.parse(str);
    } catch (e) {
        return false;
    }
    return true;
}

function JoinRoom(room) {
  socket.send(room);
  $("#js-join-room").html("更换房价");
}

})(jQuery);
