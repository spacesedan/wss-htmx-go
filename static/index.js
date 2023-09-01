const closeConnectionBtn = document.querySelector("#closeConnection");

closeConnectionBtn.addEventListener("click", function () {
  htmx.on("htmx:wsClose");
});

document.body.addEventListener("htmx:wsConfigSend", function (e) {
  console.log(e);
  switch (e.detail.headers["HX-Trigger"]) {
    case "messageForm":
      console.log("chat message");
      e.detail.parameters = {
        action: "broadcast",
        // message: e.detail.parameters.message,
        ...e.detail.parameters,
      };
    default:
      break;
  }
  e.detail.socketWrapper.send();
});
