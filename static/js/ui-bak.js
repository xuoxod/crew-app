import { stringify, parse } from "./utils.js";
import { log, hasCam, dlog, cap } from "./clientutils.js";
import {
  getElement,
  newElement,
  removeChildren,
  appendChild,
  addClickHandler,
  addHandler,
  addAttribute,
  countChildren,
  removeChild,
  append,
  getLastChild,
} from "./computils.js";

export const _updateUsersList = async (
  userList,
  listItemClickHandler,
  detectWebcam,
  blockUser,
  sendMessage
) => {
  const usersParent = document.querySelector("#users-parent");
  const currentUser = getElement("rmtid-input").value;

  // if (userList.length > 0) {
  removeChildren(usersParent);
  for (const user in userList) {
    const userObject = userList[user];

    if (userObject.online) {
      const blockedUsersIndex = userObject.blockedUsers.findIndex(
        (x) => x == currentUser
      );

      if (blockedUsersIndex == -1) {
        const blockedByIndex = userObject.blockedBy.findIndex(
          (x) => x == currentUser
        );

        if (blockedByIndex == -1) {
          const friendsIndex = userObject.friends.findIndex(
            (x) => x == currentUser
          );

          // Create comps
          const li = newElement("li");
          const liImg = newElement("img");
          const divName = newElement("div");
          const spanName = newElement("span");
          const divConnect = newElement("div");
          const divIconFriend = newElement("div");
          const iconConnect = newElement("i");
          const divBlock = newElement("div");
          const iconBlock = newElement("i");
          const spanConnect = newElement("span");
          const cardBody = newElement("div");
          const cardBodyLayout = newElement("div");
          const cardTitle = newElement("h5");
          const acceptCallIcon = newElement("i");
          const divConnectIcon = newElement("div");
          const divBlockIcon = newElement("div");
          const divFriendIcon = newElement("div");
          const divSendMessage = newElement("div");
          const iconSendMessage = newElement("i");
          const iconFriend = newElement("i");

          const displayName = userObject.displayName.fname
            ? `${cap(userObject.fname)}`
            : `${cap(userObject.uname)}`;

          liImg.alt = `${displayName}`;

          // Add attributes
          if (friendsIndex !== -1) {
            addAttribute(iconFriend, "class", "bi bi-person-fill-dash");
            addAttribute(iconFriend, "data-toggle", "tooltip");
            addAttribute(iconFriend, "data-placement", "top");
            addAttribute(iconFriend, "data-html", "true");
            addAttribute(
              iconFriend,
              "title",
              `Add ${displayName} to friend list`
            );
          } else {
            addAttribute(iconFriend, "class", "bi bi-person-fill-plus");
            addAttribute(iconFriend, "class", "bi bi-person-dash");
            addAttribute(iconFriend, "data-toggle", "tooltip");
            addAttribute(iconFriend, "data-placement", "top");
            addAttribute(iconFriend, "data-html", "true");
            addAttribute(
              iconFriend,
              "title",
              `Remove ${displayName} from friend list`
            );
          }
          addAttribute(li, "class", "w3-bar w3-round-large w3-mobile");
          addAttribute(li, "id", `li-${userObject._id}`);
          addAttribute(li, "style", "position:relative;");

          addAttribute(liImg, "class", "w3-bar-item w3-circle");
          addAttribute(liImg, "style", "width:85px;margin:0;");
          addAttribute(liImg, "data-toggle", "tooltip");
          addAttribute(liImg, "data-placement", "top");
          addAttribute(liImg, "data-html", "true");
          addAttribute(liImg, "title", `${displayName}`);

          addAttribute(divName, "class", "w3-bar-item");
          addAttribute(divConnect, "class", "w3-bar-item");
          addAttribute(divBlock, "class", "w3-bar-item");
          addAttribute(divSendMessage, "class", "w3-bar-item");
          addAttribute(divFriendIcon, "class", "w3-bar-item");

          addAttribute(iconBlock, "id", `block-${userObject._id}`);
          addAttribute(iconBlock, "class", "bi bi-person-fill-slash");
          addAttribute(iconBlock, "data-toggle", "tooltip");
          addAttribute(iconBlock, "data-placement", "top");
          addAttribute(iconBlock, "data-html", "true");
          addAttribute(iconBlock, "title", `Block ${displayName}`);

          addAttribute(
            iconSendMessage,
            "class",
            "bi bi-chat-square-dots w3-hover-opacity"
          );
          addAttribute(iconSendMessage, "id", `send-${userObject._id}`);
          addAttribute(iconSendMessage, "data-connectiontype", "message");
          addAttribute(iconSendMessage, "data-toggle", "tooltip");
          addAttribute(iconSendMessage, "data-placement", "top");
          addAttribute(iconSendMessage, "data-html", "true");
          addAttribute(
            iconSendMessage,
            "title",
            `Send ${displayName} a message`
          );

          addAttribute(acceptCallIcon, "class", "bi bi-check-lg text-success");

          if (userObject.photoUrl) {
            liImg.src = userObject.photoUrl;
          } else {
            liImg.src = "/assets/graphics/silhouette.png";
          }

          // Append comps
          appendChild(usersParent, li);
          appendChild(li, liImg);
          appendChild(li, divName);
          appendChild(li, divConnect);
          appendChild(li, divSendMessage);
          appendChild(li, divBlock);
          appendChild(li, divFriendIcon);
          appendChild(divSendMessage, iconSendMessage);
          appendChild(divName, spanName);
          appendChild(divConnect, iconConnect);
          appendChild(divBlock, iconBlock);
          appendChild(divFriendIcon, iconFriend);

          spanName.innerHTML = `<strong class="fs">${displayName}</strong>`;

          detectWebcam((results) => {
            if (!results) {
              addAttribute(iconConnect, "id", `connect-${userObject._id}`);
              addAttribute(iconConnect, "data-connectiontype", "audio");
              addAttribute(
                iconConnect,
                "class",
                "bi bi-mic-fill w3-hover-opacity"
              );
              addAttribute(iconConnect, "data-toggle", "tooltip");
              addAttribute(iconConnect, "data-placement", "top");
              addAttribute(iconConnect, "data-html", "true");
              addAttribute(iconConnect, "title", `Connect with ${displayName}`);
            } else {
              addAttribute(iconConnect, "id", `connect-${userObject._id}`);
              addAttribute(iconConnect, "data-connectiontype", "video");
              addAttribute(
                iconConnect,
                "class",
                "bi bi-camera-video-fill w3-hover-opacity"
              );
              addAttribute(iconConnect, "data-toggle", "tooltip");
              addAttribute(iconConnect, "data-placement", "top");
              addAttribute(iconConnect, "data-html", "true");
              addAttribute(iconConnect, "title", `Connect with ${displayName}`);
            }
          });

          // Register click handlers

          addClickHandler(iconConnect, listItemClickHandler);

          addClickHandler(iconBlock, (e) => {
            const blockee = e.target.id.split("-")[1];
            const blocker = currentUser;
            blockUser(blocker, blockee);
          });

          addClickHandler(iconSendMessage, () => {
            if (getElement("messenger")) {
              const messengerElement = getElement("messenger");
              const messengerHeader = getElement("messenger-header");

              messengerElement.style.width = "50%";
              messengerElement.style.display = "block";
              messengerHeader.innerText = `Send ${displayName} a message`;
            }
          });

          if (getElement("messenger-send-button")) {
            addClickHandler(getElement("messenger-send-button"), () => {
              if (getElement("messenger-message").value) {
                const messengerElement = getElement("messenger");
                const messengerHeader = getElement("messenger-header");
                const messengerMessage = getElement("messenger-message");
                const messageDetails = {};

                messageDetails.from = `${currentUser}`;
                messageDetails.to = userObject._id;
                messageDetails.message = messengerMessage.value;

                sendMessage(messageDetails);

                messengerHeader.innerText = "";
                messengerMessage.value = "";
                messengerElement.style.width = "0%";
                messengerElement.style.display = "none";
              }
            });
          }
        }
      }
    }
  }
};

export const updateUsersList = async (
  userList,
  listItemClickHandler,
  detectWebcam,
  blockUser,
  sendMessage,
  befriend,
  unbefriend
) => {
  const usersParent = document.querySelector("#users-parent");
  const currentUser = getElement("rmtid-input").value;

  // if (userList.length > 0) {
  removeChildren(usersParent);
  for (const user in userList) {
    const userObject = userList[user];

    if (userObject.online) {
      const blockedUsersIndex = userObject.blockedUsers.findIndex(
        (x) => x == currentUser
      );

      if (blockedUsersIndex == -1) {
        const blockedByIndex = userObject.blockedBy.findIndex(
          (x) => x == currentUser
        );

        if (blockedByIndex == -1) {
          // Create comps
          const li = newElement("li");

          const liImg = newElement("img");
          const divLiImg = newElement("div");

          const spanName = newElement("span");
          const divName = newElement("div");

          const divTop = newElement("div");
          const divBottom = newElement("div");
          const divContent = newElement("div");

          const divConnect = newElement("div");
          const divSendMessage = newElement("div");
          const divBlock = newElement("div");
          const divFriendIcon = newElement("div");

          const iconConnect = newElement("i");
          const iconSendMessage = newElement("i");
          const iconBlock = newElement("i");
          const iconFriend = newElement("i");

          const acceptCallIcon = newElement("i");

          const displayName = userObject.displayName.fname
            ? `${cap(userObject.fname)}`
            : `${cap(userObject.uname)}`;

          liImg.alt = `${displayName}`;

          const friendsIndex = userObject.befriendedBy
            ? userObject.befriendedBy.findIndex((x) => x == currentUser)
            : -1;

          /* add attributes  */

          detectWebcam((results) => {
            if (!results) {
              addAttribute(iconConnect, "id", `connect-${userObject._id}`);
              addAttribute(iconConnect, "data-connectiontype", "audio");
              addAttribute(
                iconConnect,
                "class",
                "bi bi-mic-fill w3-hover-opacity w3-cell"
              );
              addAttribute(iconConnect, "data-toggle", "tooltip");
              addAttribute(iconConnect, "data-placement", "top");
              addAttribute(iconConnect, "data-html", "true");
              addAttribute(iconConnect, "title", `Connect with ${displayName}`);
            } else {
              addAttribute(iconConnect, "id", `connect-${userObject._id}`);
              addAttribute(iconConnect, "data-connectiontype", "video");
              addAttribute(
                iconConnect,
                "class",
                "bi bi-camera-video-fill w3-hover-opacity w3-cell"
              );
              addAttribute(iconConnect, "data-toggle", "tooltip");
              addAttribute(iconConnect, "data-placement", "top");
              addAttribute(iconConnect, "data-html", "true");
              addAttribute(iconConnect, "title", `Connect with ${displayName}`);
            }
          });

          if (friendsIndex !== -1) {
            addAttribute(iconFriend, "id", `friend-${userObject._id}`);
            addAttribute(iconFriend, "class", "bi bi-person-fill-dash w3-cell");
            addAttribute(iconFriend, "data-toggle", "tooltip");
            addAttribute(iconFriend, "data-placement", "top");
            addAttribute(iconFriend, "data-html", "true");
            addAttribute(
              iconFriend,
              "title",
              `Remove ${displayName} from friend list`
            );
          } else {
            addAttribute(iconFriend, "id", `friend-${userObject._id}`);
            addAttribute(iconFriend, "class", "bi bi-person-fill-add w3-cell");
            addAttribute(iconFriend, "data-toggle", "tooltip");
            addAttribute(iconFriend, "data-placement", "top");
            addAttribute(iconFriend, "data-html", "true");
            addAttribute(
              iconFriend,
              "title",
              `Add ${displayName} to friend list`
            );
          }

          // list item
          addAttribute(li, "class", "w3-bar w3-round-large w3-mobile");
          addAttribute(li, "id", `li-${userObject._id}`);
          addAttribute(li, "style", "position:relative;width:100%;");

          // div containers
          addAttribute(divContent, "class", "w3-container");
          addAttribute(divTop, "class", "w3-cell-row");
          addAttribute(divBottom, "class", "w3-cell-row");

          addAttribute(divLiImg, "class", "w3-cell w3-cell-middle");
          addAttribute(divName, "class", "w3-cell w3-mobile");

          addAttribute(divConnect, "class", "w3-container w3-cell");
          addAttribute(divConnect, "style", "display:inline-block;");
          addAttribute(divSendMessage, "class", "w3-container w3-cell");
          addAttribute(divSendMessage, "style", "display:inline-block;");
          addAttribute(divBlock, "class", "w3-container w3-cell");
          addAttribute(divBlock, "style", "display:inline-block;");
          addAttribute(divFriendIcon, "class", "w3-container w3-cell");
          addAttribute(divFriendIcon, "style", "display:inline-block;");

          /* components */

          // list image
          addAttribute(
            liImg,
            "class",
            "w3-bar-item w3-circle w3-middle w3-margin-bottom w3-cell"
          );
          addAttribute(liImg, "style", "width:85px;margin:0;margin:0;");
          addAttribute(liImg, "data-toggle", "tooltip");
          addAttribute(liImg, "data-placement", "top");
          addAttribute(liImg, "data-html", "true");
          addAttribute(liImg, "title", `${displayName}`);

          // span name
          addAttribute(spanName, "class", "w3-right");

          // block icon
          addAttribute(iconBlock, "id", `block-${userObject._id}`);
          addAttribute(iconBlock, "class", "bi bi-person-fill-slash w3-cell");
          addAttribute(iconBlock, "data-toggle", "tooltip");
          addAttribute(iconBlock, "data-placement", "top");
          addAttribute(iconBlock, "data-html", "true");
          addAttribute(iconBlock, "title", `Block ${displayName}`);

          // send message icon
          addAttribute(
            iconSendMessage,
            "class",
            "bi bi-chat-square-dots w3-hover-opacity w3-cell"
          );
          addAttribute(iconSendMessage, "id", `send-${userObject._id}`);
          addAttribute(iconSendMessage, "data-connectiontype", "message");
          addAttribute(iconSendMessage, "data-toggle", "tooltip");
          addAttribute(iconSendMessage, "data-placement", "top");
          addAttribute(iconSendMessage, "data-html", "true");
          addAttribute(
            iconSendMessage,
            "title",
            `Send ${displayName} a message`
          );

          // accept connection icon
          addAttribute(acceptCallIcon, "class", "bi bi-check-lg text-success");

          // image src
          if (userObject.photoUrl) {
            liImg.src = userObject.photoUrl;
          } else {
            liImg.src = "/assets/graphics/silhouette.png";
          }

          // Append comps
          appendChild(usersParent, li);
          appendChild(li, divContent);
          appendChild(divContent, divTop);
          appendChild(divContent, divBottom);

          // top div
          appendChild(divTop, divLiImg);
          appendChild(divTop, divName);

          // bottom div
          appendChild(divBottom, divConnect);
          appendChild(divBottom, divSendMessage);
          appendChild(divBottom, divBlock);
          appendChild(divBottom, divFriendIcon);

          // connect div
          appendChild(divConnect, iconConnect);

          // send message div
          appendChild(divSendMessage, iconSendMessage);

          // block user div
          appendChild(divBlock, iconBlock);

          // friend div
          appendChild(divFriendIcon, iconFriend);

          // div image
          appendChild(divLiImg, liImg);

          // div name
          appendChild(divName, spanName);

          // div send message
          appendChild(divSendMessage, iconSendMessage);

          // div block user
          appendChild(divBlock, iconBlock);

          // div friend icon
          appendChild(divFriendIcon, iconFriend);

          // inner HTML
          spanName.innerHTML = `<strong class="fs">${displayName}</strong>`;

          /* Register click handlers */

          // icon connect
          addClickHandler(iconConnect, listItemClickHandler);

          // icon block user
          addClickHandler(iconBlock, (e) => {
            const blockee = e.target.id.split("-")[1];
            const blocker = currentUser;
            blockUser(blocker, blockee);
          });

          // icon befriend user
          addClickHandler(iconFriend, (e) => {
            const icon = e.target;
            const target = icon.id.split("-")[1];

            if (icon.classList.contains("bi-person-fill-dash")) {
              unbefriend(currentUser, target);
            } else {
              befriend(currentUser, target);
            }
          });

          // icon send message
          addClickHandler(iconSendMessage, () => {
            if (getElement("messenger")) {
              const messengerElement = getElement("messenger");
              const messengerHeader = getElement("messenger-header");

              messengerElement.style.width = "50%";
              messengerElement.style.display = "block";
              messengerHeader.innerText = `Send ${displayName} a message`;
            }
          });

          // send message button
          if (getElement("messenger-send-button")) {
            addClickHandler(getElement("messenger-send-button"), () => {
              if (getElement("messenger-message").value) {
                const messengerElement = getElement("messenger");
                const messengerHeader = getElement("messenger-header");
                const messengerMessage = getElement("messenger-message");
                const messageDetails = {};

                messageDetails.from = `${currentUser}`;
                messageDetails.to = userObject._id;
                messageDetails.message = messengerMessage.value;

                sendMessage(messageDetails);

                messengerHeader.innerText = "";
                messengerMessage.value = "";
                messengerElement.style.width = "0%";
                messengerElement.style.display = "none";
              }
            });
          }
        }
      }
    }
  }
};

export const updateFriendsList = async (
  userList,
  listItemClickHandler,
  detectWebcam,
  blockUser,
  sendMessage
) => {
  const usersParent = document.querySelector("#users-parent");
  const currentUser = getElement("rmtid-input").value;

  // if (userList.length > 0) {
  removeChildren(usersParent);
  for (const user in userList) {
    const userObject = userList[user];

    if (userObject.online) {
      const blockedUsersIndex = userObject.blockedUsers.findIndex(
        (x) => x == currentUser
      );

      if (blockedUsersIndex == -1) {
        const blockedByIndex = userObject.blockedBy.findIndex(
          (x) => x == currentUser
        );

        if (blockedByIndex == -1) {
          const friendsIndex = userObject.friends.findIndex(
            (x) => x == currentUser
          );

          if (friendsIndex !== -1) {
            // Create comps
            const li = newElement("li");
            const liImg = newElement("img");
            const divName = newElement("div");
            const spanName = newElement("span");
            const divConnect = newElement("div");
            const iconConnect = newElement("i");
            const divBlock = newElement("div");
            const iconBlock = newElement("i");
            const spanConnect = newElement("span");
            const cardBody = newElement("div");
            const cardBodyLayout = newElement("div");
            const cardTitle = newElement("h5");
            const acceptCallIcon = newElement("i");
            const divConnectIcon = newElement("div");
            const divBlockIcon = newElement("div");
            const divFriendIcon = newElement("div");
            const divSendMessage = newElement("div");
            const iconSendMessage = newElement("i");

            const displayName = userObject.displayName.fname
              ? `${cap(userObject.fname)}`
              : `${cap(userObject.uname)}`;

            liImg.alt = `${displayName}`;

            // Add attributes
            addAttribute(li, "class", "w3-bar w3-round-large w3-mobile");
            addAttribute(li, "id", `li-${userObject._id}`);
            addAttribute(li, "style", "position:relative;");

            addAttribute(liImg, "class", "w3-bar-item w3-circle");
            addAttribute(liImg, "style", "width:85px;margin:0;");
            addAttribute(liImg, "data-toggle", "tooltip");
            addAttribute(liImg, "data-placement", "top");
            addAttribute(liImg, "data-html", "true");
            addAttribute(liImg, "title", `${displayName}`);

            addAttribute(divName, "class", "w3-bar-item");
            addAttribute(divConnect, "class", "w3-bar-item");
            addAttribute(divBlock, "class", "w3-bar-item");
            addAttribute(divSendMessage, "class", "w3-bar-item");

            addAttribute(iconBlock, "id", `block-${userObject._id}`);
            addAttribute(
              iconBlock,
              "class",
              "bi bi-eye-slash-fill w3-hover-opacity"
            );
            addAttribute(iconBlock, "data-toggle", "tooltip");
            addAttribute(iconBlock, "data-placement", "top");
            addAttribute(iconBlock, "data-html", "true");
            addAttribute(iconBlock, "title", `Block ${displayName}`);

            addAttribute(
              iconSendMessage,
              "class",
              "bi bi-chat-square-dots w3-hover-opacity"
            );
            addAttribute(iconSendMessage, "id", `send-${userObject._id}`);
            addAttribute(iconSendMessage, "data-connectiontype", "message");
            addAttribute(iconSendMessage, "data-toggle", "tooltip");
            addAttribute(iconSendMessage, "data-placement", "top");
            addAttribute(iconSendMessage, "data-html", "true");
            addAttribute(
              iconSendMessage,
              "title",
              `Send ${displayName} a message`
            );

            addAttribute(
              acceptCallIcon,
              "class",
              "bi bi-check-lg text-success"
            );

            if (userObject.photoUrl) {
              liImg.src = userObject.photoUrl;
            } else {
              liImg.src = "/assets/graphics/silhouette.png";
            }

            // Append comps
            appendChild(usersParent, li);
            appendChild(li, liImg);
            appendChild(li, divName);
            appendChild(li, divConnect);
            appendChild(li, divSendMessage);
            appendChild(li, divBlock);
            appendChild(divSendMessage, iconSendMessage);
            appendChild(divName, spanName);
            appendChild(divConnect, iconConnect);
            appendChild(divBlock, iconBlock);

            spanName.innerHTML = `<strong>${displayName}</strong>`;

            detectWebcam((results) => {
              if (!results) {
                addAttribute(iconConnect, "id", `connect-${userObject._id}`);
                addAttribute(iconConnect, "data-connectiontype", "audio");
                addAttribute(
                  iconConnect,
                  "class",
                  "bi bi-mic-fill w3-hover-opacity"
                );
                addAttribute(iconConnect, "data-toggle", "tooltip");
                addAttribute(iconConnect, "data-placement", "top");
                addAttribute(iconConnect, "data-html", "true");
                addAttribute(
                  iconConnect,
                  "title",
                  `Connect with ${displayName}`
                );
              } else {
                addAttribute(iconConnect, "id", `connect-${userObject._id}`);
                addAttribute(iconConnect, "data-connectiontype", "video");
                addAttribute(
                  iconConnect,
                  "class",
                  "bi bi-camera-video-fill w3-hover-opacity"
                );
                addAttribute(iconConnect, "data-toggle", "tooltip");
                addAttribute(iconConnect, "data-placement", "top");
                addAttribute(iconConnect, "data-html", "true");
                addAttribute(
                  iconConnect,
                  "title",
                  `Connect with ${displayName}`
                );
              }
            });

            // Register click handlers

            addClickHandler(iconConnect, listItemClickHandler);

            addClickHandler(iconBlock, (e) => {
              const blockee = e.target.id.split("-")[1];
              const blocker = currentUser;
              blockUser(blocker, blockee);
            });

            addClickHandler(iconSendMessage, () => {
              if (getElement("messenger")) {
                const messengerElement = getElement("messenger");
                const messengerHeader = getElement("messenger-header");

                messengerElement.style.width = "50%";
                messengerElement.style.display = "block";
                messengerHeader.innerText = `Send ${displayName} a message`;
              }
            });

            if (getElement("messenger-send-button")) {
              addClickHandler(getElement("messenger-send-button"), () => {
                if (getElement("messenger-message").value) {
                  const messengerElement = getElement("messenger");
                  const messengerHeader = getElement("messenger-header");
                  const messengerMessage = getElement("messenger-message");
                  const messageDetails = {};

                  messageDetails.from = `${currentUser}`;
                  messageDetails.to = userObject._id;
                  messageDetails.message = messengerMessage.value;

                  sendMessage(messageDetails);

                  messengerHeader.innerText = "";
                  messengerMessage.value = "";
                  messengerElement.style.width = "0%";
                  messengerElement.style.display = "none";
                }
              });
            }
          }
        }
      }
    }
  }
};

export const showMessage = (userDetails, iconClickHandler) => {
  const { userInfo, hasWebcam, messageBody, alertType } = userDetails;
  const messageParent = document.querySelector("#message-container");
  const alert = newElement("div");
  const alertCloseButton = newElement("button");
  const container = newElement("div");
  const row1 = newElement("div");
  const row2 = newElement("div");
  const row1Col = newElement("div");
  const row2Col = newElement("div");
  const row1P = newElement("p");
  const webcamIcon = newElement("i");
  const messageIcon = newElement("i");

  removeChildren(messageParent);

  /* Set attributes */

  // Alert attributes
  addAttribute(
    alert,
    "class",
    `alert ${alertType} d-flex align-items-center alert-dismissible fade show`
  );
  addAttribute(alert, "role", "alert");
  addAttribute(alert, "style", "display:inline-flex;");
  addAttribute(alert, "id", "alert");

  // Alert close button attributes
  addAttribute(alertCloseButton, "class", "btn-close");
  addAttribute(alertCloseButton, "type", "button");
  addAttribute(alertCloseButton, "data-bs-dismiss", "alert");
  addAttribute(alertCloseButton, "aria-label", "Close");

  // Container, rows and columns attributes
  addAttribute(container, "class", "container-fluid m-0 p-0");
  addAttribute(container, "style", "margin:0;padding:0;");
  addAttribute(row1, "class", "row m-0 p-0");
  addAttribute(row2, "class", "row m-0 p-0");
  addAttribute(row1Col, "class", "col-12 m-0 p-0");
  addAttribute(row2Col, "class", "col-12 m-0 p-0");

  // Icon attributes
  addAttribute(webcamIcon, "class", "bi bi-webcam-fill icon");
  addAttribute(webcamIcon, "id", `wc-${userInfo._id}`);
  addAttribute(messageIcon, "class", "bi bi-chat-left-dots-fill icon");
  addAttribute(messageIcon, "id", `mi-${userInfo._id}`);

  /* Append elements */

  // Append to view
  appendChild(messageParent, alert);

  // Append to alert
  appendChild(alert, container);
  appendChild(alert, alertCloseButton);

  // Append to container, rows and columns
  appendChild(container, row1);
  appendChild(container, row2);
  appendChild(row1, row1Col);
  appendChild(row2, row2Col);

  // Append elements
  appendChild(row1, row1P);

  // Detect user's webcam
  if (hasWebcam) {
    appendChild(row2Col, webcamIcon);
  } else {
    appendChild(row2Col, messageIcon);
  }

  // Set alert's title
  row1P.innerHTML = `${messageBody}`;

  addClickHandler(webcamIcon, iconClickHandler);
  addClickHandler(messageIcon, iconClickHandler);
};

export const showCallAlert = (userDetails) => {
  const { userInfo, alertType } = userDetails;
  const messageParent = document.querySelector("#message-container");
  const alert = newElement("div");
  const alertCloseButton = newElement("button");
  const strong = newElement("strong");
  const p = newElement("p");

  dlog(`Alert Type:\t${alertType}`, `ui.js: showCallAlert`);

  /* Set attributes */

  // Alert attributes
  addAttribute(
    alert,
    "class",
    `alert ${alertType} d-flex align-items-center alert-dismissible fade show`
  );
  addAttribute(alert, "role", "alert");
  addAttribute(alert, "style", "display:inline-flex;");
  addAttribute(alert, "id", "alert");

  // Paragraph element
  addAttribute(p, "style", "display:inline-flex;");

  // Alert close button attributes
  addAttribute(alertCloseButton, "class", "btn-close");
  addAttribute(alertCloseButton, "type", "button");
  addAttribute(alertCloseButton, "data-bs-dismiss", "alert");
  addAttribute(alertCloseButton, "aria-label", "Close");

  // Alert title
  strong.innerHTML = `... Calling ${cap(userInfo.fname)}`;

  /* Append elements */

  // Append to view
  appendChild(messageParent, alert);

  // Append to alert
  appendChild(p, strong);
  appendChild(alert, p);
  appendChild(alert, alertCloseButton);

  setTimeout(() => {
    alert.remove();
  }, [4000]);
};

export const showCallResponse = (userDetails) => {
  const { receiver, response, alertType } = userDetails;
  const messageParent = document.querySelector("#message-container");
  const alert = newElement("div");
  const alertCloseButton = newElement("button");
  const strong = newElement("strong");

  /* Set attributes */

  // Alert attributes
  addAttribute(
    alert,
    "class",
    `alert ${alertType} d-flex align-items-center alert-dismissible fade show`
  );
  addAttribute(alert, "role", "alert");
  addAttribute(alert, "style", "display:inline-flex;");
  addAttribute(alert, "id", "alert");

  // Alert close button attributes
  addAttribute(alertCloseButton, "class", "btn-close");
  addAttribute(alertCloseButton, "type", "button");
  addAttribute(alertCloseButton, "data-bs-dismiss", "alert");
  addAttribute(alertCloseButton, "aria-label", "Close");

  // Alert title
  strong.innerHTML = `${cap(
    receiver.fname
  )} ${response} your connection request`;

  /* Append elements */

  // Append to view
  appendChild(messageParent, alert);

  // Append to alert
  appendChild(alert, strong);
  appendChild(alert, alertCloseButton);

  setTimeout(() => {
    alert.remove();
  }, [7000]);
};

export const showMediaControls = (micControlHandler, videoControlHandler) => {
  const parent = getElement("local-media-controls-parent");
  const container = newElement("div");
  const nav = newElement("nav");
  const containerFluid = newElement("div");
  const ul = newElement("ul");
  const liRecord = newElement("li");
  const liMicrophone = newElement("li");
  const liVideo = newElement("li");
  const liScreenShare = newElement("li");
  const recordIcon = newElement("i");
  const microphoneIcon = newElement("i");
  const videoIcon = newElement("i");
  const screenshareIcon = newElement("i");

  // Add attributes
  addAttribute(container, "class", "container");
  addAttribute(nav, "class", "navbar navbar-expand bg-light");
  addAttribute(containerFluid, "class", "container-fluid");

  // List
  addAttribute(ul, "class", "navbar-nav position-relative controls");

  // List items
  addAttribute(liRecord, "class", "mx-4  d-flex align-items-center");
  addAttribute(liMicrophone, "class", "mx-4  d-flex align-items-center");
  addAttribute(liVideo, "class", "mx-4 a d-flex align-items-center");
  addAttribute(liScreenShare, "class", "mx-4  d-flex align-items-center");

  // Icons
  addAttribute(recordIcon, "class", "mx-4 bi bi-record");
  addAttribute(microphoneIcon, "class", "mx-4 bi bi-mic-fill");
  addAttribute(videoIcon, "class", "mx-4 bi bi-camera-video-fill");
  addAttribute(screenshareIcon, "class", "mx-4 bi bi-window-stack");

  // Append comps
  appendChild(parent, container);
  appendChild(container, nav);
  appendChild(nav, containerFluid);
  appendChild(containerFluid, ul);
  appendChild(ul, recordIcon);
  appendChild(ul, microphoneIcon);
  appendChild(ul, videoIcon);
  appendChild(ul, screenshareIcon);

  /*   appendChild(liRecord, recordIcon);
  appendChild(liMicrophone, microphoneIcon);
  appendChild(liVideo, videoIcon);
  appendChild(liScreenShare, screenshareIcon); */

  // Register click events
  addClickHandler(microphoneIcon, micControlHandler);
  addClickHandler(videoIcon, videoControlHandler);
};

export const showCallRequest = (userDetails, acceptCall) => {
  const { user, conntype, callee } = userDetails;
  const name = user.displayName.fname ? `${user.fname}` : `${user.uname}`;
  const msg = `${name} wants to connect`;
  dlog(`${msg}`, `ui.js: showCallRequest`);

  const li = getElement(`li-${user.uid}`);

  if (countChildren(li) > 5) {
    removeChild(li, getLastChild(li));
  }

  if (li) {
    const panelMsg = newElement("div");
    const panelButtons = newElement("div");
    const panelMain = newElement("div");
    const panelContent = newElement("div");
    const acceptButton = newElement("button");
    const closeButton = newElement("span");
    const para = newElement("p");

    // Panel Attributes
    addAttribute(panelMsg, "class", "w3-mobile w3-round-large w3-cell-row");
    addAttribute(panelButtons, "class", "w3-mobile w3-round-large w3-cell-row");
    addAttribute(panelMain, "class", "w3-mobile w3-round-large w3-container");
    addAttribute(
      panelMain,
      "style",
      "background-color:rgba(5,5,5,0.5);background-size:cover;position:absolute;left:0;top:0;height:100%;width:100%;margin:0;"
    );
    addAttribute(panelContent, "class", "w3-section");

    // Element Attributes
    addAttribute(
      closeButton,
      "class",
      "w3-button w3-display-topright w3-center-align w3-text-white w3-opacity-min w3-circle w3-hover-white w3-text-hover-black w3-margin-right w3-margin-top"
    );
    addAttribute(closeButton, "style", "background-color:rgba(10,10,10,0.5);");
    addAttribute(para, "class", "w3-text-white w3-center-align w3-cell-middle");
    addAttribute(
      acceptButton,
      "class",
      "w3-button w3-text-white w3-border w3-border-white w3-opacity-min w3-round-large w3-hover-white w3-cell-middle"
    );
    addAttribute(acceptButton, "style", "background:transparent;");

    // Inner Text
    para.innerHTML = `<strong>${cap(name)} wants to connect</strong>`;
    closeButton.innerHTML = `&times;`;
    acceptButton.innerHTML = `<strong>Accept</strong>`;

    // Append Elements
    appendChild(li, panelMain);
    appendChild(panelMain, panelContent);
    appendChild(panelContent, panelMsg);
    appendChild(panelContent, closeButton);
    appendChild(panelMain, panelButtons);
    appendChild(panelButtons, acceptButton);
    appendChild(panelMsg, para);

    // Register click handler
    addClickHandler(acceptButton, () => {
      acceptCall(user._id, callee, conntype);
    });

    addClickHandler(closeButton, (e) => {
      const target = e.target;
      const parent = target.parentElement;
      const grandParent = parent.parentElement;
      grandParent.remove();
    });
  }
};

export const showPrivateMessageAlert = (
  userDetails,
  privateMessageReplyHandler
) => {
  const { from, text } = parse(userDetails);
  const userName = from.displayName.fname ? from.fname : from.uname;
  const userId = from._id;
  const li = getElement(`li-${userId}`);
  const currentUser = getElement("rmtid-input").value;

  if (countChildren(li) > 5) {
    removeChild(li, getLastChild(li));
  }

  if (li) {
    const panelMsg = newElement("div");
    const panelButtons = newElement("div");
    const panelMain = newElement("div");
    const panelContent = newElement("div");
    const showMessageButton = newElement("button");
    const closeButton = newElement("span");
    const para = newElement("p");

    // Panel Attributes
    addAttribute(panelMsg, "class", "w3-mobile w3-round-large w3-cell-row");
    addAttribute(panelButtons, "class", "w3-mobile w3-round-large w3-cell-row");
    addAttribute(panelMain, "class", "w3-mobile w3-round-large w3-container");
    addAttribute(
      panelMain,
      "style",
      "background-color:rgba(5,5,5,0.5);background-size:cover;position:absolute;left:0;top:0;height:100%;width:100%;margin:0;"
    );

    // Element Attributes
    addAttribute(
      closeButton,
      "class",
      "w3-button w3-display-topright w3-center-align w3-text-white w3-opacity-min w3-circle w3-hover-white w3-text-hover-black w3-margin-right w3-margin-top"
    );
    addAttribute(closeButton, "style", "background-color:rgba(10,10,10,0.5);");
    addAttribute(para, "class", "w3-text-white w3-center-align w3-cell-middle");
    addAttribute(
      showMessageButton,
      "class",
      "w3-button w3-text-white w3-border w3-border-white w3-opacity-min w3-round-large w3-hover-white w3-middle"
    );
    addAttribute(showMessageButton, "style", "background:transparent;");
    addAttribute(panelContent, "class", "w3-row w3-margin");

    // Inner Text
    para.innerHTML = `<strong>${cap(userName)} sent you a message`;
    closeButton.innerHTML = `&times;`;
    showMessageButton.innerHTML = `<strong>Accept</strong>`;

    // Append Elements
    appendChild(li, panelMain);
    appendChild(panelMain, panelContent);
    appendChild(panelContent, panelMsg);
    appendChild(panelContent, closeButton);
    appendChild(panelMain, panelButtons);
    appendChild(panelButtons, showMessageButton);
    appendChild(panelMsg, para);

    // Register click handler

    addClickHandler(showMessageButton, () => {
      dlog(
        `\n\tMessage from ${userName}\n\tMessage:  ${text}\n`,
        `ui.js: showPrivateMessageAlert`
      );

      panelMain.remove();

      if (!getElement(`messenger-${from._id}`)) {
        const messengerContainer = newElement("div");
        const messageContainer = newElement("div");
        const inputContainer = newElement("div");
        const messagePara = newElement("p");
        const input = newElement("input");
        const replyButton = newElement("button");
        const messengerCloseButton = newElement("span");

        // add attributes

        addAttribute(
          messengerContainer,
          "class",
          "w3-mobile w3-round-large w3-container"
        );
        addAttribute(messengerContainer, "id", `messenger-${from._id}`);
        addAttribute(
          messengerContainer,
          "style",
          "background-color:rgba(5,5,5,0.5);background-size:cover;position:absolute;top:0;left:0;height:100%;width:100%;margin:0;overflow:scroll;"
        );

        addAttribute(inputContainer, "class", "w3-row w3-margin-top");

        addAttribute(messageContainer, "class", "w3-row w3-margin w3-cell");
        addAttribute(messageContainer, "id", `message-${from._id}`);

        addAttribute(
          para,
          "class",
          "w3-card w3-round-xxlarge w3-text-white w3-middle w3-marging-right"
        );
        addAttribute(para, "style", "word-wrap:break-word;width:27%;");

        addAttribute(
          messengerCloseButton,
          "class",
          "w3-button w3-display-topright w3-center-align w3-text-white w3-opacity-min w3-circle w3-hover-white w3-text-hover-black w3-cell"
        );
        addAttribute(messengerCloseButton, "style", "margin:10px;");

        addAttribute(
          input,
          "class",
          "w3-mobile w3-input w3-left w3-small w3-round-xxlarge w3-cell"
        );
        addAttribute(input, "type", "text");
        addAttribute(input, "style", "width:80%; height:35px;");

        addAttribute(
          replyButton,
          "class",
          "w3-button w3-border w3-border-white w3-text-white w3-round-xxlarge w3-small w3-cell w3-hover-white w3-text-hover-black"
        );

        // inner HTML
        messagePara.innerHTML = `<small class="w3-text-white"><strong>${text}</strong></small>`;
        replyButton.innerHTML = `<small class="w3-text-white"><strong>Reply</strong></small>`;

        // inner text
        messengerCloseButton.innerHTML = `&times;`;

        // append elements

        appendChild(li, messengerContainer);
        appendChild(messengerContainer, messengerCloseButton);
        appendChild(messengerContainer, messageContainer);
        appendChild(messageContainer, messagePara);
        appendChild(messengerContainer, inputContainer);
        appendChild(inputContainer, input);
        appendChild(inputContainer, replyButton);

        // register click events
        addClickHandler(messengerCloseButton, () => {
          messengerContainer.remove();
        });

        addClickHandler(replyButton, () => {
          if (input.value) {
            const replyDetails = {};

            replyDetails.replyFrom = currentUser;
            replyDetails.replyTo = from._id;
            replyDetails.replyMessage = input.value;

            privateMessageReplyHandler(replyDetails);
            messengerContainer.remove();
          }
        });
      }
    });

    addClickHandler(closeButton, (e) => {
      const target = e.target;
      const parent = target.parentElement;
      const grandParent = parent.parentElement;
      grandParent.remove();
    });
  }
};
