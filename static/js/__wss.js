import { parse, stringify } from "./utils.js";
import { log, dlog, tlog } from "./clientutils.js";
import { getElement, addAttribute, addClickHandler } from "./computils.js";
import {
  updateUsersList,
  showMessage,
  showCallAlert,
  showCallResponse,
  showCallRequest,
} from "./ui.js";

let socketIO = null,
  userDetails = {};

export const registerSocketEvents = (socket) => {
  socketIO = socket;

  socket.on("connect", () => {
    dlog(`\n\tSuccessfully connected to socket.io server\n`);
    userDetails = {};
    userDetails.uid = document.querySelector("#rmtid-input").value;
    socket.emit("registerme", userDetails);
  });

  socket.on("registered",()=>{
    dlog(`I'm regiestered`,`wss script: line 26`)
  });

  socket.on("updateonlineuserlist", (data) => {
    const { users, msg } = data;
    const pUsers = parse(users);
    const currentUser = document.querySelector("#rmtid-input").value;
    const currentUserBlockedUsers = pUsers[currentUser].blockedUsers;
    const arrUsers = [];

    if (msg) {
      alert(msg);
    }

    dlog(
      `${pUsers[currentUser].fname} blocked users ${stringify(
        currentUserBlockedUsers
      )}\nID: ${currentUser}`,
      `wss script: line 32`
    );

    for (const u in pUsers) {
      const user = pUsers[u];
      const uid = user._id;

      if (uid != currentUser) {
        arrUsers.push({
          ...user,
          isBlocked: currentUserBlockedUsers.findIndex((x) => x == uid) != -1,
        });
      }
    }

    const listItemClickHandler = (e) => {
      dlog(`List item ${e.target.id} was clicked`, `${document.title}`);
      userDetails = {};
      userDetails.receiver = e.target.id.trim().split("-")[1];
      userDetails.sender = document.querySelector("#rmtid-input").value;
      userDetails.conntype = e.target.dataset.connectiontype;

      dlog(
        `Sending ${userDetails.conntype} connection request\t Sender: ${userDetails.sender} Receiver: ${userDetails.receiver}`,
        "wss script"
      );
      socket.emit("userclicked", userDetails);
    };

    updateUsersList(
      arrUsers,
      listItemClickHandler,
      detectWebcam,
      blockUser,
      blockedBy
    );
  });

  socket.on("registered", (data) => {
    const { uid } = data;
    dlog(`I am registered`);

    try {
      const xmlHttp = new XMLHttpRequest();

      xmlHttp.onload = () => {
        // location.href = `/chat/profile/create/${uid}`;

        const { status, doc, hasDoc } = parse(xmlHttp.responseText);

        if (status) {
          dlog(`chat profile created`);
          userDetails = {};
          userDetails.uid = document.querySelector("#rmtid-input").value;
          userDetails.doc = doc;
          userDetails.hasDoc = hasDoc;
          socket.emit("updateuser", userDetails);
        }
      };

      xmlHttp.open("GET", `/chat/profile/create/${uid}`);

      xmlHttp.send(true);
    } catch (err) {
      dlog(err);
      return;
    }
  });

  socket.on("clickeduser", (data) => {
    const { strUser } = data;
    const user = parse(strUser);
    userDetails = {};
    userDetails.userInfo = user;
    userDetails.alertType = `alert-info`;

    showCallAlert(userDetails);
  });

  socket.on("connectionrequested", (data) => {
    const { strUserDetails } = data;

    userDetails = parse(strUserDetails);
    userDetails.callee = document.querySelector("#rmtid-input").value;

    /* TODO:
        display user's connection request in a mini window */

    showCallRequest(userDetails, acceptCall);
  });

  socket.on("connectionrequestresponse", (data) => {
    const { responseData } = data;
    const userReponseData = parse(responseData);
    const { receiver, response, roomName, connType, sender } = userReponseData;

    // dlog(`Made it to connectionrequestresponse`);

    if (response == "accepted") {
      dlog(`User ${receiver.fname} ${response} your request`);

      userReponseData.alertType = "alert-success";
      showCallResponse(userReponseData);
      createRoom(roomName);
      userDetails = {};
      userDetails.sender = sender;
      userDetails.receiver = receiver;
      userDetails.roomName = roomName;
      userDetails.connectionType = connType;
      userDetails.from = `sender`;
      socket.emit("enterroom", userDetails);
      getRoomTokenAndEnterRoom(roomName, connType, sender, receiver._id);
    } else if (response == "rejected") {
      userReponseData.alertType = "alert-warning";
      showCallResponse(userReponseData);
    } else {
      dlog(`No response`);
    }
  });

  socket.on("enterroom", (data) => {
    userDetails = parse(data);
    const { roomName, connectionType, sender, receiver } = userDetails;
    getRoomTokenAndEnterRoom(roomName, connectionType, sender, receiver._id);
  });

  socket.on("updateyourblockedbylist", (data) => {
    const { list } = data;

    const listItemClickHandler = (e) => {
      dlog(`List item ${e.target.id} was clicked`, `${document.title}`);
      userDetails = {};
      userDetails.receiver = e.target.id.trim().split("-")[1];
      userDetails.sender = document.querySelector("#rmtid-input").value;
      userDetails.conntype = e.target.dataset.connectiontype;

      dlog(
        `Sending ${userDetails.conntype} connection request\t Sender: ${userDetails.sender} Receiver: ${userDetails.receiver}`,
        "wss script"
      );
      socket.emit("userclicked", userDetails);
    };

    updateUsersList(
      list,
      listItemClickHandler,
      detectWebcam,
      blockUser,
      blockedBy
    );
  });
};

function detectWebcam(callback) {
  let md = navigator.mediaDevices;
  if (!md || !md.enumerateDevices) return callback(false);
  md.enumerateDevices().then((devices) => {
    callback(devices.some((device) => "videoinput" === device.kind));
  });
}

function acceptCall(senderUid, receiverUid, connType) {
  dlog(`You accepted ${senderUid}'s connection request`);
  userDetails = {};
  userDetails.sender = senderUid;
  userDetails.receiver = receiverUid;
  userDetails.connType = connType;
  socketIO.emit("callaccepted", userDetails);
}

function rejectCall(senderUid, receiverUid) {
  dlog(`You rejected ${senderUid}'s connection request`);
  userDetails = {};
  userDetails.sender = senderUid;
  userDetails.receiver = receiverUid;
  socketIO.emit("callrejected", userDetails);
}

function noResponseToCall(senderUid, receiverUid) {}

function blockUser(blockerUid, blockeeUid) {
  let xmlHttp;
  try {
    xmlHttp = new XMLHttpRequest();

    xmlHttp.open("POST", "/chat/block", true);

    xmlHttp.setRequestHeader(
      "Content-type",
      "application/x-www-form-urlencoded"
    );

    xmlHttp.onload = () => {
      const responseText = xmlHttp.responseText;

      if (responseText) {
        log(`\n\tResponse Text: ${stringify(responseText)}\n`);
        const responseJson = parse(responseText);
        const status = responseJson.status;

        if (status) {
          const { blockerdoc, blockeedoc } = responseJson;

          userDetails = {};
          userDetails.blocker = { blockerUid, blockerDoc: blockerdoc };
          userDetails.blockee = { blockeeUid, blockeeDoc: blockeedoc };
          socketIO.emit("iblockedauser", userDetails);
        } else {
          dlog(`Something went wrong blocking user`);
        }
        return;
      }
    };

    xmlHttp.send(`blocker=${blockerUid}&blockee=${blockeeUid}`, true);
  } catch (err) {
    tlog(err);
    return;
  }
}

function unblockUser(blockerUid, blockeeUid) {
  let xmlHttp;

  try {
    xmlHttp = new XMLHttpRequest();

    xmlHttp.open("POST", "/chat/unblock", true);

    xmlHttp.setRequestHeader(
      "Content-type",
      "application/x-www-form-urlencoded"
    );

    xmlHttp.onload = () => {
      const responseText = xmlHttp.responseText;

      if (responseText) {
        // log(`\n\tResponse Text: ${stringify(responseText)}\n`);
        const responseJson = parse(responseText);
        const status = responseJson.status;

        if (status) {
          dlog(`${blockerUid} unblocked ${blockeeUid}`);
          userDetails = {};
          userDetails.blocker = blockerUid;
          userDetails.blockee = blockeeUid;
          socketIO.emit("iunblockedauser", userDetails);
        } else {
          dlog(`Something went wrong unblocking user`);
        }
        return;
      }
    };

    xmlHttp.send(`blocker=${blockerUid}&blockee=${blockeeUid}`, true);
  } catch (err) {
    tlog(err);
    return;
  }
}

function createRoom(roomName) {
  let xmlHttp;

  try {
    xmlHttp = new XMLHttpRequest();

    xmlHttp.open("POST", "/chat/room/create", true);

    xmlHttp.setRequestHeader(
      "Content-type",
      "application/x-www-form-urlencoded"
    );

    xmlHttp.onload = () => {
      const responseText = xmlHttp.responseText;

      if (responseText) {
        // log(`\n\tResponse Text: ${stringify(responseText)}\n`);
        const responseJson = parse(responseText);
        const status = responseJson.status;

        if (status) {
          dlog(
            `room:\t${roomName} successfully created.`,
            "createRoom xhr post"
          );
        } else {
          const cause = responseJson.cause;
          const detail = responseJson.detail;
          const err = responseJson.err;

          dlog(
            `Cause:\t${cause}\nDetails:\t${detail}\nError:\n\t${err}\n\n`,
            `createRoom xhr post`
          );
        }
        return;
      }
    };

    xmlHttp.send(`roomName=${roomName}`, true);
  } catch (err) {
    tlog(err);
    return;
  }
}

function getRoomTokenAndEnterRoom(
  roomName,
  connectionType,
  senderId,
  receiverId
) {
  let xmlHttp;

  try {
    xmlHttp = new XMLHttpRequest();

    xmlHttp.open("POST", `/chat/room/gettoken`);

    xmlHttp.setRequestHeader(
      "Content-type",
      "application/x-www-form-urlencoded"
    );

    xmlHttp.onload = () => {
      const responseText = xmlHttp.responseText;

      if (responseText) {
        // log(`\n\tResponse Text: ${stringify(responseText)}\n`);
        const responseJson = parse(responseText);
        const status = responseJson.status;
        const _roomName = responseJson.roomName;
        const _connectionType = responseJson.connectionType;
        const _senderId = responseJson.senderId;
        const _receiverId = responseJson.receiverId;
        const _token = responseJson.token;

        if (status) {
          dlog(
            `roomName:\t${_roomName}\nConn Type:\t${_connectionType}\nSender:\t${_senderId}\nReceiver:\t${_receiverId}\nToken:\t${_token}\n`,
            "sender xhr request"
          );

          location.href = `/chat/room/connect?roomName=${roomName}&connectionType=${connectionType}&senderId=${senderId}&receiverId=${receiverId}&token=${_token}`;
        }
      }
    };

    xmlHttp.send(
      `roomName=${roomName}&connectionType=${connectionType}&senderId=${senderId}&receiverId=${receiverId}`,
      true
    );
  } catch (err) {
    tlog(err);
    return;
  }
}

export const cloakMe = () => {
  dlog(`Going invisible`, `wss script: line 398`);
  const uid = document.querySelector("#rmtid-input").value;
  let xmlHttp;

  try {
    xmlHttp = new XMLHttpRequest();

    xmlHttp.open("POST", "/chat/user/hide", true);

    xmlHttp.setRequestHeader(
      "Content-type",
      "application/x-www-form-urlencoded"
    );

    xmlHttp.onload = () => {
      const responseText = xmlHttp.responseText;

      if (responseText) {
        const responseJson = parse(responseText);
        const status = responseJson.status;

        if (status) {
          const userDoc = responseJson.doc;
          userDetails = {};
          userDetails.uid = uid;
          userDetails.doc = userDoc;

          const hidemeLink = getElement("hideme");
          hidemeLink.innerText = "Show";
          hidemeLink.removeEventListener("click", cloakMe);
          addClickHandler(hidemeLink, uncloakMe);

          socketIO.emit("makemeinvisible", userDetails);
        } else {
          const reason = responseJson.reason;
          dlog(`${reason}`);
        }
        return;
      }
    };

    xmlHttp.send(`userId=${uid}`, true);
  } catch (err) {
    tlog(err);
    return;
  }
};

export const uncloakMe = () => {
  dlog(`Going visible`, `wss script: line 447`);
  const uid = document.querySelector("#rmtid-input").value;
  let xmlHttp;

  try {
    xmlHttp = new XMLHttpRequest();

    xmlHttp.open("POST", "/chat/user/unhide", true);

    xmlHttp.setRequestHeader(
      "Content-type",
      "application/x-www-form-urlencoded"
    );

    xmlHttp.onload = () => {
      const responseText = xmlHttp.responseText;

      if (responseText) {
        const responseJson = parse(responseText);
        const status = responseJson.status;

        if (status) {
          const userDoc = responseJson.doc;
          userDetails = {};
          userDetails.uid = uid;
          userDetails.doc = userDoc;

          const hidemeLink = getElement("hideme");
          hidemeLink.innerText = "Hide";
          hidemeLink.removeEventListener("click", uncloakMe);
          addClickHandler(hidemeLink, cloakMe);

          socketIO.emit("makemevisible", userDetails);
        } else {
          const reason = responseJson.reason;
          dlog(`${reason}`);
        }
        return;
      }
    };

    xmlHttp.send(`userId=${uid}`, true);
  } catch (err) {
    tlog(err);
    return;
  }
};

function blockedBy(blockedByList, blocker) {
  const blockerIndex = blockedByList.findIndex((x) => x == blocker);
  return blockerIndex != -1;
}
