import { registerSocketEvents } from "./wss.js";
import { getElement, addClickHandler } from "./computils.js";

// init socket connection
const socket = io("/");

// register socket events
registerSocketEvents(socket);

addClickHandler(getElement("messenger-close-button"), () => {
  const messengerSidebar = getElement("messenger");

  messengerSidebar.style.width = "0%";
  messengerSidebar.style.display = "none";
});
