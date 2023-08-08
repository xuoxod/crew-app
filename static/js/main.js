import { log, dlog } from "./clientutils.js";
import {
  addHandler,
  newElement,
  appendChild,
  addAttribute,
  removeChildren,
} from "./computils.js";

const clientsUl = document.querySelector(".clients");
let source;

if (typeof EventSource !== "undefined") {
  source = new EventSource("/landing");

  log(source);

  addHandler(source, "open", (e) => {
    log(`Event open: ${JSON.stringify(e)}`);
  });

  addHandler(source, "error", (e) => {
    console.log(`Event error: ${JSON.stringify(e)}`);
  });

  addHandler(source, "message", (e) => {
    console.log(`Event message: ${JSON.stringify(e)}`);
    const clients = JSON.parse(e.data);

    removeChildren(clientsUl);

    const liHeader = newElement("li");
    appendChild(clientsUl, liHeader);
    liHeader.innerHTML = `<strong><h3>Online Users: ${clients.length}</h3></strong>`;

    /* clients.forEach((client) => {
      // Create elements
      const li = newElement("li");
      const inputGroup = newElement("div");
      const inputGroupLabel = newElement("span");
      const inputElement = newElement("input");

      // Add attributes
      addAttribute(li, "class", "w3-padding-small");
      addAttribute(inputGroup, "class", "input-group");
      addAttribute(inputGroupLabel, "class", "input-group-text");
      addAttribute(inputElement, "type", "text");
      addAttribute(inputElement, "readonly");
      addAttribute(inputElement, "class", "form-control");
      addAttribute(inputElement, "value", `${client.stamp}`);

      // Add data
      inputGroupLabel.innerHTML = `<strong>${client.platform}</strong>`;

      // Append elements
      appendChild(li, inputGroup);
      appendChild(inputGroup, inputGroupLabel);
      appendChild(inputGroup, inputElement);
      appendChild(clientsUl, li);
    }); */
  });

  addEventListener("beforeunload", (event) => {
    log(`\n\tBefore unload\n`);
    try {
      source.close();
    } catch (err) {
      log(err);
    } finally {
      source = null;
    }
  });
}
