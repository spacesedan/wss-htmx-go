import { writable } from "svelte/store";

export default { messageStore: writable() };

// const sendMessage = (message: string) => {
//   if (ws.readyState <= 1) {
//     ws.send(message);
//   }
// };
