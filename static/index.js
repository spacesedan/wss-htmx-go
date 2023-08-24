const wsStatus = document.querySelector("#wsStatus");
const messageBox = document.querySelector("#messageBox");
const message = document.querySelector("#message");
const messages = document.querySelector("#messages");
const currentUser = document.querySelector("#current_user");

console.log(htmx);

document.body.addEventListener("htmx:wsConnecting", function(_) {
  wsStatus.innerHTML = "Connecting...";
});

document.body.addEventListener("htmx:wsOpen", function(e) {
  e.preventDefault();
  wsStatus.innerHTML = "Connected";
  e.detail.socketWrapper.send("Work plz");
});

document.body.addEventListener("htmx:wsAfterMessage", function(e) {
  const sw = e.detail.socketWrapper;
  messages.scrollIntoView({
    block: "end",
    behavior: "smooth",
    inline: "nearest",
  });

  sw.send("Poop");
});

document.body.addEventListener("htmx:wsAfterSend", function(_) {
  messageBox.value = "";
});
