const wsStatus = document.querySelector("#wsStatus");
const messageBox = document.querySelector("#messageBox");
const message = document.querySelector("#message");
const messages = document.querySelector("#messages");

document.body.addEventListener("htmx:wsConnecting", function(_) {
  wsStatus.innerHTML = "Connecting...";
});

document.body.addEventListener("htmx:wsOpen", function(evt) {
  evt.detail.socketWrapper.send();
  wsStatus.innerHTML = "Connected";
});

document.body.addEventListener("htmx:wsBeforeSend", function(_) {
  if (messageBox.value == "") {
    console.log("sending empty");
    return;
  }
});

document.body.addEventListener("htmx:wsAfterMessage", function(_) {
  messages.scrollIntoView({
    block: "end",
    behavior: "smooth",
    inline: "nearest",
  });
});

document.body.addEventListener("htmx:wsAfterSend", function(_) {
  messageBox.value = "";
  console.log(messages.scrollHeight);
});
