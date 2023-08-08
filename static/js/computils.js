export const addHandler = (theElement, whichEvent, method) => {
  if (null != theElement && null != whichEvent && typeof method == "function") {
    theElement.addEventListener(whichEvent, method);
  }
};

export const addClickHandler = (theElement, handler) => {
  if (null != theElement && typeof handler == "function") {
    addHandler(theElement, "click", handler);
  }
};

export const addKeyupHandler = (theElement, handler) => {
  if (null != theElement && typeof handler == "function") {
    addHandler(theElement, "keyup", handler);
  }
};

export const addKeydownHandler = (theElement, handler) => {
  if (null != theElement && typeof handler == "function") {
    addHandler(theElement, "keydown", handler);
  }
};

export const addOnFocusHandler = (theElement, handler) => {
  if (null != theElement && typeof handler == "function") {
    addHandler(theElement, "focus", handler);
  }
};

export const addOnChangeHandler = (theElement, handler) => {
  if (null != theElement && typeof handler == "function") {
    addHandler(theElement, "change", handler);
  }
};

export const addOffFocusHandler = (theElement, handler) => {
  if (null != theElement && typeof handler == "function") {
    addHandler(theElement, "focusout", handler);
  }
};

export const addAttribute = (theElement, whichAttribute, attributeValue) => {
  if (null != theElement) {
    theElement.setAttribute(whichAttribute, attributeValue);
  }
};

export const setAttribute = (theElement, whichAttribute, attributeValue) => {
  if (null != theElement) {
    theElement.setAttribute(whichAttribute, attributeValue);
  }
};

export const getAttribute = (theElement, whichAttribute) => {
  if (null != theElement && null != whichAttribute) {
    return theElement.getAttribute(`${whichAttribute}`) || null;
  }
  return "Element is null";
};

export const removeAttribute = (theElement, whichAttribute) => {
  if (null != theElement) {
    if (theElement.hasAttribute(whichAttribute)) {
      theElement.removeAttribute(whichAttribute);
    }
  }
};

export const getElement = (nameIdClass) => {
  let element = null;
  if (null != (element = document.querySelector(`${nameIdClass}`))) {
    return element;
  }
  if (null != (element = document.querySelector(`#${nameIdClass}`))) {
    return element;
  }
  if (null != (element = document.querySelector(`.${nameIdClass}`))) {
    return element;
  }
  return null;
};

export const getElements = (nameIdClass) => {
  let elements = null;
  if (null != (elements = document.querySelectorAll(`${nameIdClass}`))) {
    return elements;
  }
  if (null != (elements = document.querySelectorAll(`#${nameIdClass}`))) {
    return elements;
  }
  if (null != (elements = document.querySelectorAll(`.${nameIdClass}`))) {
    return elements;
  }
  return null;
};

export const appendChild = (parent, child) => {
  if (null != parent && null != child) {
    parent.appendChild(child);
  }
};

export const appendBeforeLastChild = (parent, child) => {
  if (null != parent && null != child) {
    const lastChildIndex = parent.children.length - 1;
    const lastChild = parent.children[lastChildIndex];
    parent.insertBefore(child, lastChild);
  }
};

export const append = (parent, child) => {
  parent.append(child);
};

export const removeChildren = (parent) => {
  parent.querySelectorAll("*").forEach((dialog) => {
    dialog.remove();
  });
};

export const countChildren = (parent) => {
  if (null != parent) {
    return parent.children.length;
  }
  return null;
};

export const getLastChild = (parent) => {
  if (null != parent) {
    return parent.lastElementChild;
  }
  return null;
};

export const removeChild = (parent, child) => {
  parent.removeChild(child);
};

export const removeById = (elementID) => {
  const element = document.querySelector(`#${elementID}`) || null;
  if (element != null) {
    element.remove();
  }
};

export const getFirstChild = (parent) => {
  if (null != parent) {
    return parent.firstElementChild;
  }
  return null;
};

export const newElement = (type) => {
  if (null != type && typeof type == "string") {
    return document.createElement(type);
  }
  return null;
};

export const generatePhoneInputBlock = (
  placeHolder = "Phone",
  name = "phonex"
) => {
  const col = newElement("div");
  const group = newElement("div");
  const label = newElement("span");
  const icon = newElement("i");
  const deleteIcon = newElement("i");

  addAttribute(col, "class", "col-12 py-3");
  addAttribute(group, "class", "input-group");
  addAttribute(label, "class", "input-group-text");
  addAttribute(icon, "class", "bi bi-phone");
  addAttribute(deleteIcon, "style", "margin-left:5px;margin-top:10px; ");
  addAttribute(deleteIcon, "class", "bi bi-trash");

  appendChild(col, group);
  appendChild(group, label);
  appendChild(label, icon);
  appendChild(group, generatePhoneInput(placeHolder, name));
  appendChild(group, deleteIcon);

  addClickHandler(deleteIcon, () => {
    col.remove();
  });

  return col;
};

export const generateEmailInputBlock = (
  placeHolder = "example@email.net",
  name = "emailx"
) => {
  const col = newElement("div");
  const group = newElement("div");
  const label = newElement("span");
  const icon = newElement("i");
  const deleteIcon = newElement("i");

  addAttribute(col, "class", "col-12 py-3");
  addAttribute(group, "class", "input-group");
  addAttribute(icon, "class", "bi bi-envelope");
  addAttribute(label, "class", "input-group-text");
  addAttribute(deleteIcon, "style", "margin-left:5px;margin-top:10px; ");
  addAttribute(deleteIcon, "class", "bi bi-trash");

  appendChild(col, group);
  appendChild(group, label);
  appendChild(label, icon);
  appendChild(group, generateEmailInput(placeHolder, name));
  appendChild(group, deleteIcon);

  addClickHandler(deleteIcon, () => {
    col.remove();
  });

  return col;
};

export const generatePhoneInput = (placeHolder = "Phone", name = `phonex`) => {
  const textInput = newElement("input");

  addAttribute(textInput, "type", "tel");
  addAttribute(textInput, "name", `${name}`);
  addAttribute(textInput, "class", `form-control`);
  addAttribute(textInput, "placeholder", `${placeHolder}`);
  addAttribute(textInput, "required", "");

  return textInput;
};

export const generateEmailInput = (placeHolder = "Email", name = `emailx`) => {
  const textInput = newElement("input");

  addAttribute(textInput, "type", "email");
  addAttribute(textInput, "name", `${name}`);
  addAttribute(textInput, "class", `form-control`);
  addAttribute(textInput, "placeholder", `${placeHolder}`);
  addAttribute(textInput, "required", "");

  return textInput;
};
