// const wsStatus = document.querySelector("#wsStatus");
// const messageBox = document.querySelector("#messageBox");
// const message = document.querySelector("#message");
// const messages = document.querySelector("#messages");
// const currentUser = document.querySelector("#current_user");
//
// document.body.addEventListener("htmx:wsConnecting", function(_) {
//   wsStatus.innerHTML = "Connecting...";
// });
//
document.body.addEventListener("htmx:wsOpen", function(e) {
  e.preventDefault();
  const msg = JSON.stringify({ HEADERS: { action: "Poop" } });
  e.detail.socketWrapper.send(msg);
  // wsStatus.innerHTML = "Connected";
});
//
// document.body.addEventListener("htmx:wsAfterMessage", function(e) {
//   e.preventDefault();
//   console.log(e.detail);
//   messages.scrollIntoView({
//     block: "end",
//     behavior: "smooth",
//     inline: "nearest",
//   });
// });
//
// document.body.addEventListener("htmx:wsAfterSend", function(_) {
//   messageBox.value = "";
// });
//
// document.body.addEventListener("htmx:wsClose", function() {
//   htmx.trigger("#left");
//   // e.detail.send("UO").then((res) => console.log(res));
// });
