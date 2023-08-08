import { registerSocketEvents, cloakMe } from "./wss.js";
import { addClickHandler, getElement } from "./computils.js";

// init socket connection
const socket = io("/");

registerSocketEvents(socket);

addClickHandler(
  getElement("hideme", (e) => {
    cloakMe();
  })
);
