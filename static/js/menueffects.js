import { cap, stringify } from "./utils.js";
import { log, dlog } from "./clientutils.js";
import {
  getElement,
  getElements,
  addAttribute,
  removeAttribute,
  addHandler,
} from "./computils.js";

function observeElement() {
  if ("IntersectionObserver" in window) {
    if (getElement("menu")) {
      const menu = getElement("menu");
      const target = getElements("observe");

      if (document.querySelector(".observe")) {
        // Instancing a new IntersectionObserver
        const observer = new IntersectionObserver((entries) => {}, {
          root: document.querySelector(".menu"),
          threshold: [0.1, 1],
          rootMargin: "1px",
        });

        // Adding a target to be observed
        observer.observe(document.querySelector(".observe"));
      }
    }
  }
}

observeElement();
