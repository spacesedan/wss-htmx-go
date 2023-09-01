// username is coming from a http secure cookie
const username = document.querySelector("#username");

function uniqueID() {
  return Math.floor(Math.random() * Date.now());
}

document.body.addEventListener("htmx:wsOpen", function(e) {
  console.log(e);
  const msg = {
    action: "entered",
    message: `${username.innerHTML} has entered the chat`,
    user: username.innerHTML,
    id: `${uniqueID()}`,
  };
  e.detail.socketWrapper.send(JSON.stringify(msg), e.detail.elt);
});

document.body.addEventListener("htmx:wsClose", function(e) {
  console.log(e);
  const msg = {
    action: "left",
    message: `${username.innerHTML} has left the chat`,
    user: username.innerHTML,
    id: `${uniqueID()}`,
  };
  e.detail.socketWrapper.send(JSON.stringify(msg), e.detail.elt);
});

document.body.addEventListener("htmx:wsConfigSend", function(e) {
  console.log(e);
  switch (e.detail.headers["HX-Trigger"]) {
    case "messageForm":
      console.log("chat message");
      e.detail.parameters = {
        ...e.detail.parameters,
        action: "message",
        user: username.innerHTML,
        id: `${uniqueID()}`,
      };
    default:
      break;
  }
});
