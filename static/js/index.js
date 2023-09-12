// username is coming from a http secure cookie
const username = document.querySelector("#username");
const messageInput = document.querySelector("#chat_message_input");
const chatMessages = document.querySelector("#chat_messages");
const leftChatMessage = document.querySelector("#leavers");

function uniqueID() {
  return Math.floor(Math.random() * Date.now());
}

document.body.addEventListener("htmx:wsOpen", function (e) {
  console.log(username.dataset.username);
  const msg = {
    action: "entered",
    message: `${username.dataset.username} has entered the chat`,
    user: username.dataset.username,
    id: `${uniqueID()}`,
  };
  e.detail.socketWrapper.send(JSON.stringify(msg), e.detail.elt);
});

document.body.addEventListener("htmx:wsClose", function (e) {
  console.log(e);
  const msg = {
    action: "left",
    message: `${username.dataset.username} has left the chat`,
    user: username.dataset.username,
    id: `${uniqueID()}`,
  };
  e.detail.socketWrapper.send(JSON.stringify(msg), e.detail.elt);
});

document.body.addEventListener("htmx:wsConfigSend", function (e) {
  switch (e.detail.headers["HX-Trigger"]) {
    case "chat_message_form":
      console.log(e);
      e.detail.parameters = {
        message: messageInput.value,
        action: "message",
        user: username.dataset.username,
        id: `${uniqueID()}`,
      };
    default:
      break;
  }
});

document.body.addEventListener("htmx:wsAfterSend", function () {
  messageInput.value = "";
});

document.body.addEventListener("htmx:wsAfterMessage", function () {
  chatMessages.scrollIntoView({
    behavior: "smooth",
    block: "end",
    inline: "nearest",
  });
});
