const wsStatus = document.querySelector("#wsStatus");
const messageBox = document.querySelector("#messageBox");
const message = document.querySelector("#message");
const messages = document.querySelector("#messages");

document.body.addEventListener("htmx:wsConnecting", function(_) {
  wsStatus.innerHTML = "Connecting...";
});

document.body.addEventListener("htmx:wsOpen", function(e) {
  wsStatus.innerHTML = "Connected";
  e.detail.socketWrapper.send("Hello");
});

document.body.addEventListener("htmx:wsAfterMessage", function() {
  messages.scrollIntoView({
    block: "end",
    behavior: "smooth",
    inline: "nearest",
  });
});

document.body.addEventListener("htmx:wsAfterSend", function(_) {
  messageBox.value = "";
});
