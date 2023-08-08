import { parse, stringify } from "./utils.js";
import { log, dlog, tlog } from "./clientutils.js";
import { getElement, addAttribute, addClickHandler } from "./computils.js";
import {
  updateUsersList,
  updateFriendsList,
  showMessage,
  showCallAlert,
  showCallResponse,
  showCallRequest,
  showPrivateMessageAlert,
} from "./ui.js";

let socketIO = null,
  userDetails = {};

export const registerSocketEvents = ( socket ) => {
  socketIO = socket;

  socket.on( "connect", () => {
    dlog( `connect event fired`, `wss.js: registerSocketEvents` );

    getChatUserProfile( ( results ) => {
      const { status, doc, hasDoc } = results;

      if ( status ) {
        userDetails = {};
        userDetails.doc = doc;
        userDetails.hasDoc = hasDoc;
        userDetails.uid = document.querySelector( "#rmtid-input" ).value;

        if ( getElement( "unblockeduser" ) ) {
          const unblockedUserId = getElement( "unblockeduser" ).value;
          userDetails.unblockedUser = unblockedUserId != null;
          userDetails.unblockeduserid = getElement( "unblockeduser" ).value;
        }

        socket.emit( "registerme", userDetails );
      }
    } );
  } );

  socket.on( "updateyourself", ( data ) => {
    dlog( `updateyourself event fired`, `wss.js: updateyourself` );
    const { unblockedBy } = data;

    getChatUserProfile( ( results ) => {
      const { status, doc } = results;

      if ( status ) {
        userDetails = {};
        userDetails.unblockedBy = unblockedBy;
        userDetails.doc = doc;

        socketIO.emit( "updateme", userDetails );
      }
    } );
  } );

  socket.on( "clearunblocked", () => {
    dlog( `clearunblocked event fired`, `wss.js: clearunblocked` );
    if ( getElement( "unblockeduser" ) ) {
      getElement( "unblockeduser" ).value = "";
      getElement( "unblockeduser" ).remove();
    }
  } );

  socket.on( "updateonlineuserlist", ( data ) => {
    dlog( `updateonlineuserlist event fired`, `wss.js` );
    const { users } = data;
    const pUsers = parse( users );
    const currentUserInputValue = document.querySelector( "#rmtid-input" ).value;
    const currentUser = pUsers[ currentUserInputValue ];
    const currentUserBlockedUsers = currentUser.blockedUsers;
    const arrUsers = [];

    /* TODO:
        Check if current user visibility to only friends is true */

    for ( const u in pUsers ) {
      const user = pUsers[ u ];
      const uid = user._id;

      if ( uid != currentUserInputValue ) {
        arrUsers.push( {
          ...user,
          isBlocked: currentUserBlockedUsers.findIndex( ( x ) => x == uid ) != -1,
        } );
      }
    }

    const listItemClickHandler = ( e ) => {
      dlog( `List item ${ e.target.id } was clicked`, `${ document.title }` );
      userDetails = {};
      userDetails.receiver = e.target.id.trim().split( "-" )[ 1 ];
      userDetails.sender = document.querySelector( "#rmtid-input" ).value;
      userDetails.conntype = e.target.dataset.connectiontype;

      dlog(
        `Sending ${ userDetails.conntype } connection request\t Sender: ${ userDetails.sender } Receiver: ${ userDetails.receiver }`,
        "wss script"
      );
      socket.emit( "userclicked", userDetails );
    };

    const sendMessage = ( messageDetails ) => {
      const { from, to, message } = messageDetails;

      dlog(
        `\tMessage Details\n\t\tFrom:  ${ from }\n\t\tTo:  ${ to }\n\t\tMessage:  ${ message }\n\n`
      );

      socket.emit( "sendprivatemessage", messageDetails );
    };

    updateUsersList(
      arrUsers,
      listItemClickHandler,
      detectWebcam,
      blockUser,
      sendMessage,
      befriendUser,
      unbefriendUser
    );
  } );

  socket.on( "registered", ( data ) => {
    dlog( `registered event fired`, `wss.js: registerSocketEvents: registered` );
    const { doc } = data;
    const pDoc = parse( doc );
    try {
      dlog(
        `You are now registered as ${ pDoc.uname }`,
        `wss.js: registerSocketEvents: registered`
      );
      getElement( "online" ).value = pDoc.online;
    } catch ( err ) {
      dlog( `${ err }`, `wss.js: registered` );
      return;
    }
  } );

  socket.on( "clickeduser", ( data ) => {
    dlog( `clickeduser event fired`, `wss.js` );
    const { strUser } = data;
    const user = parse( strUser );
    userDetails = {};
    userDetails.userInfo = user;
    userDetails.alertType = `alert-info`;

    showCallAlert( userDetails );
  } );

  socket.on( "connectionrequested", ( data ) => {
    dlog( `connectionrequested event fired`, `wss.js` );
    const { strUserDetails } = data;

    userDetails = parse( strUserDetails );
    userDetails.callee = document.querySelector( "#rmtid-input" ).value;

    /* TODO:
        display user's connection request in a mini window */

    showCallRequest( userDetails, acceptCall );
  } );

  socket.on( "connectionrequestresponse", ( data ) => {
    dlog( `connectionrequestresponse event fired`, `wss.js` );
    const { responseData } = data;
    const userReponseData = parse( responseData );
    const { receiver, response, roomName, connType, sender } = userReponseData;

    // dlog(`Made it to connectionrequestresponse`);

    if ( response == "accepted" ) {
      dlog( `User ${ receiver.fname } ${ response } your request` );

      userReponseData.alertType = "alert-success";
      showCallResponse( userReponseData );
      createRoom( roomName );
      userDetails = {};
      userDetails.sender = sender;
      userDetails.receiver = receiver;
      userDetails.roomName = roomName;
      userDetails.connectionType = connType;
      userDetails.from = `sender`;
      socket.emit( "enterroom", userDetails );
      getRoomTokenAndEnterRoom( roomName, connType, sender, receiver._id );
    } else if ( response == "rejected" ) {
      userReponseData.alertType = "alert-warning";
      showCallResponse( userReponseData );
    } else {
      dlog( `No response` );
    }
  } );

  socket.on( "enterroom", ( data ) => {
    dlog( `enterroom event fired`, `wss.js` );
    userDetails = parse( data );
    const { roomName, connectionType, sender, receiver } = userDetails;
    getRoomTokenAndEnterRoom( roomName, connectionType, sender, receiver._id );
  } );

  socket.on( "privatemessage", ( data ) => {
    dlog( `privatemessage event fired`, `wss.js` );
    showPrivateMessageAlert( data.messageDetails, privateMessageReplyHandler );
  } );
};

addEventListener( "beforeunload", ( event ) => {
  dlog( `\tbeforeunload window event fired`, `wss.js` );
  userDetails = {};
  userDetails.uid = document.querySelector( "#rmtid-input" ).value;
  socketIO.emit( "disconnectme", userDetails );
  return;
} );

function detectWebcam ( callback ) {
  let md = navigator.mediaDevices;
  if ( !md || !md.enumerateDevices ) return callback( false );
  md.enumerateDevices().then( ( devices ) => {
    callback( devices.some( ( device ) => "videoinput" === device.kind ) );
  } );
}

function acceptCall ( senderUid, receiverUid, connType ) {
  dlog( `You accepted ${ senderUid }'s connection request` );
  userDetails = {};
  userDetails.sender = senderUid;
  userDetails.receiver = receiverUid;
  userDetails.connType = connType;
  socketIO.emit( "callaccepted", userDetails );
}

function rejectCall ( senderUid, receiverUid ) {
  dlog( `You rejected ${ senderUid }'s connection request` );
  userDetails = {};
  userDetails.sender = senderUid;
  userDetails.receiver = receiverUid;
  socketIO.emit( "callrejected", userDetails );
}

function noResponseToCall ( senderUid, receiverUid ) { }

function blockUser ( blockerUid, blockeeUid ) {
  let xmlHttp;
  try {
    xmlHttp = new XMLHttpRequest();

    xmlHttp.open( "POST", "/chat/block", true );

    xmlHttp.setRequestHeader(
      "Content-type",
      "application/x-www-form-urlencoded"
    );

    xmlHttp.onload = () => {
      const responseText = xmlHttp.responseText;

      if ( responseText ) {
        const responseJson = parse( responseText );
        const status = responseJson.status;

        if ( status ) {
          const { blockerdoc, blockeedoc } = responseJson;

          userDetails = {};
          userDetails.blocker = { blockerUid, blockerDoc: blockerdoc };
          userDetails.blockee = { blockeeUid, blockeeDoc: blockeedoc };
          socketIO.emit( "iblockedauser", userDetails );
        } else {
          dlog( `Something went wrong blocking user` );
        }
        return;
      }
    };

    xmlHttp.send( `blocker=${ blockerUid }&blockee=${ blockeeUid }`, true );
  } catch ( err ) {
    tlog( err );
    return;
  }
}

function unblockUser ( blockerUid, blockeeUid ) {
  let xmlHttp;

  try {
    xmlHttp = new XMLHttpRequest();

    xmlHttp.open( "POST", "/chat/unblock", true );

    xmlHttp.setRequestHeader(
      "Content-type",
      "application/x-www-form-urlencoded"
    );

    xmlHttp.onload = () => {
      const responseText = xmlHttp.responseText;

      if ( responseText ) {
        // log(`\n\tResponse Text: ${stringify(responseText)}\n`);
        const responseJson = parse( responseText );
        const status = responseJson.status;

        if ( status ) {
          dlog( `${ blockerUid } unblocked ${ blockeeUid }` );
          userDetails = {};
          userDetails.blocker = blockerUid;
          userDetails.blockee = blockeeUid;
          socketIO.emit( "iunblockedauser", userDetails );
        } else {
          dlog( `Something went wrong unblocking user` );
        }
        return;
      }
    };

    xmlHttp.send( `blocker=${ blockerUid }&blockee=${ blockeeUid }`, true );
  } catch ( err ) {
    tlog( err );
    return;
  }
}

function befriendUser ( befrienderUid, befriendeeUid ) {
  let xmlHttp;
  try {
    xmlHttp = new XMLHttpRequest();

    xmlHttp.open( "POST", "/chat/user/befriend", true );

    xmlHttp.setRequestHeader(
      "Content-type",
      "application/x-www-form-urlencoded"
    );

    xmlHttp.onload = () => {
      const responseText = xmlHttp.responseText;

      if ( responseText ) {
        const responseJson = parse( responseText );
        const status = responseJson.status;

        if ( status ) {
          // const { befrienderdoc, befriendeedoc } = responseJson;

          const befrienderdoc = responseJson.befrienderdoc;
          const befriendeedoc = responseJson.befriendeedoc;

          userDetails = {};

          userDetails.befriender = {
            befrienderUid,
            befrienderDoc: befrienderdoc,
          };

          userDetails.befriendee = {
            befriendeeUid,
            befriendeeDoc: befriendeedoc,
          };

          socketIO.emit( "ibefriendedauser", userDetails );
        } else {
          dlog( `Something went wrong befriending user` );
        }
        return;
      }
    };

    xmlHttp.send(
      `befriender=${ befrienderUid }&befriendee=${ befriendeeUid }`,
      true
    );
  } catch ( err ) {
    tlog( err );
    return;
  }
}

function unbefriendUser ( unbefrienderUid, unbefriendeeUid ) {
  let xmlHttp;

  try {
    xmlHttp = new XMLHttpRequest();

    xmlHttp.open( "POST", "/chat/user/unbefriend", true );

    xmlHttp.setRequestHeader(
      "Content-type",
      "application/x-www-form-urlencoded"
    );

    xmlHttp.onload = () => {
      const responseText = xmlHttp.responseText;

      if ( responseText ) {
        // log(`\n\tResponse Text: ${stringify(responseText)}\n`);
        const responseJson = parse( responseText );
        const status = responseJson.status;
        const unbefrienderdoc = responseJson.unbefrienderdoc;
        const unbefriendeedoc = responseJson.unbefriendeedoc;

        if ( status ) {
          userDetails = {};

          userDetails.unbefriender = {};
          userDetails.unbefriendee = {};

          userDetails.unbefriender.unbefrienderUid = unbefrienderUid;
          userDetails.unbefriender.unbefrienderDoc = unbefrienderdoc;

          userDetails.unbefriendee.unbefriendeeUid = unbefriendeeUid;
          userDetails.unbefriendee.unbefriendeeDoc = unbefriendeedoc;

          socketIO.emit( "iunbefriendedauser", userDetails );
        } else {
          dlog( `Something went wrong unbefriending user` );
        }
        return;
      }
    };

    xmlHttp.send(
      `unbefriender=${ unbefrienderUid }&unbefriendee=${ unbefriendeeUid }`,
      true
    );
  } catch ( err ) {
    tlog( err );
    return;
  }
}

function createRoom ( roomName ) {
  let xmlHttp;

  try {
    xmlHttp = new XMLHttpRequest();

    xmlHttp.open( "POST", "/chat/room/create", true );

    xmlHttp.setRequestHeader(
      "Content-type",
      "application/x-www-form-urlencoded"
    );

    xmlHttp.onload = () => {
      const responseText = xmlHttp.responseText;

      if ( responseText ) {
        // log(`\n\tResponse Text: ${stringify(responseText)}\n`);
        const responseJson = parse( responseText );
        const status = responseJson.status;

        if ( status ) {
          dlog(
            `room:\t${ roomName } successfully created.`,
            "wss.js: createRoom"
          );
        } else {
          const cause = responseJson.cause;
          const detail = responseJson.detail;
          const err = responseJson.err;

          dlog(
            `Cause:\t${ cause }\nDetails:\t${ detail }\nError:\n\t${ err }\n\n`,
            `wss.js: createRoom`
          );
        }
        return;
      }
    };

    xmlHttp.send( `roomName=${ roomName }`, true );
  } catch ( err ) {
    tlog( err );
    return;
  }
}

function getRoomTokenAndEnterRoom (
  roomName,
  connectionType,
  senderId,
  receiverId
) {
  let xmlHttp;

  try {
    xmlHttp = new XMLHttpRequest();

    xmlHttp.open( "POST", `/chat/room/gettoken` );

    xmlHttp.setRequestHeader(
      "Content-type",
      "application/x-www-form-urlencoded"
    );

    xmlHttp.onload = () => {
      const responseText = xmlHttp.responseText;

      if ( responseText ) {
        const responseJson = parse( responseText );
        const status = responseJson.status;
        const _roomName = responseJson.roomName;
        const _connectionType = responseJson.connectionType;
        const _senderId = responseJson.senderId;
        const _receiverId = responseJson.receiverId;
        const _token = responseJson.token;

        if ( status ) {
          location.href = `/chat/room/connect?roomName=${ roomName }&connectionType=${ connectionType }&senderId=${ senderId }&receiverId=${ receiverId }&token=${ _token }`;
        }
      }
    };

    xmlHttp.send(
      `roomName=${ roomName }&connectionType=${ connectionType }&senderId=${ senderId }&receiverId=${ receiverId }`,
      true
    );
  } catch ( err ) {
    tlog( err );
    return;
  }
}

function blockedBy ( blockedByList, blocker ) {
  const blockerIndex = blockedByList.findIndex( ( x ) => x == blocker );
  return blockerIndex != -1;
}

function getChatUserProfile ( cb ) {
  let xmlHttp;

  dlog( `Requesting my profile`, `wss.js: getChatUserProfile` );

  try {
    xmlHttp = new XMLHttpRequest();

    xmlHttp.open( "POST", `/chat/profile` );

    xmlHttp.setRequestHeader(
      "Content-type",
      "application/x-www-form-urlencoded"
    );

    xmlHttp.onload = () => {
      const responseText = xmlHttp.responseText;

      if ( responseText ) {
        const responseJson = parse( responseText );
        const status = responseJson.status;
        const doc = responseJson.doc;
        const hasDoc = responseJson.hasDoc;

        dlog(
          `Received a response for my profile request`,
          `wss.js: getChatUserProfile: xmlHttp response`
        );

        if ( status ) {
          return cb( { status: status, doc: doc, hasDoc: hasDoc } );
        } else {
          return cb( { status: false, msg: `Got nothing` } );
        }
      }
    };
    const uid = document.querySelector( "#rmtid-input" ).value;
    xmlHttp.send( `uid=${ uid }`, true );
  } catch ( err ) {
    tlog( err );
    return cb( { status: false, err: err, msg: `Got nothing` } );
  }
}

function privateMessageReplyHandler ( replyDetails ) {
  socketIO.emit( "privatemessageresponse", replyDetails );
}

const cloakMe = () => {
  dlog( `Going invisible`, `wss.js: cloakMe` );
  const uid = document.querySelector( "#rmtid-input" ).value;
  let xmlHttp;

  try {
    xmlHttp = new XMLHttpRequest();

    xmlHttp.open( "POST", "/chat/user/hide", true );

    xmlHttp.setRequestHeader(
      "Content-type",
      "application/x-www-form-urlencoded"
    );

    xmlHttp.onload = () => {
      const responseText = xmlHttp.responseText;

      if ( responseText ) {
        const responseJson = parse( responseText );
        const status = responseJson.status;

        if ( status ) {
          const userDoc = responseJson.doc;
          userDetails = {};
          userDetails.uid = uid;
          userDetails.doc = userDoc;

          socketIO.emit( "makemeinvisible", userDetails );
        } else {
          const reason = responseJson.reason;
          dlog( `${ reason }` );
        }
        return;
      }
    };

    xmlHttp.send( `userId=${ uid }`, true );
  } catch ( err ) {
    tlog( err );
    return;
  }
};

const uncloakMe = () => {
  dlog( `Going visible`, `wss.js: uncloakMe` );
  const uid = document.querySelector( "#rmtid-input" ).value;
  let xmlHttp;

  try {
    xmlHttp = new XMLHttpRequest();

    xmlHttp.open( "POST", "/chat/user/unhide", true );

    xmlHttp.setRequestHeader(
      "Content-type",
      "application/x-www-form-urlencoded"
    );

    xmlHttp.onload = () => {
      const responseText = xmlHttp.responseText;

      if ( responseText ) {
        const responseJson = parse( responseText );
        const status = responseJson.status;

        if ( status ) {
          const userDoc = responseJson.doc;
          userDetails = {};
          userDetails.uid = uid;
          userDetails.doc = userDoc;

          socketIO.emit( "makemevisible", userDetails );
        } else {
          const reason = responseJson.reason;
          dlog( `${ reason }` );
        }
        return;
      }
    };

    xmlHttp.send( `userId=${ uid }`, true );
  } catch ( err ) {
    tlog( err );
    return;
  }
};

const checkOnlineStatus = () => {
  if (
    getElement( "online" ) &&
    getElement( "cloak" ) &&
    getElement( "cloak-label" )
  ) {
    const online = getElement( "online" ).value.trim() == "true";
    const cloakButton = getElement( "cloak" );
    const cloakLabel = getElement( "cloak-label" );

    if ( online ) {
      cloakLabel.innerHTML = "<strong><small>Online</strong></small>";
      cloakButton.checked = true;
    } else {
      cloakLabel.innerHTML = "<strong><small>Offline</small></strong>";
      cloakButton.checked = false;
    }
  }
};

if ( getElement( "cloak" ) && getElement( "cloak-label" ) ) {
  const onlineInput = getElement( "online" );
  const cloakButton = getElement( "cloak" );
  const cloakLabel = getElement( "cloak-label" );

  addClickHandler( cloakButton, ( e ) => {
    const target = e.target;

    cloakLabel.innerHTML = target.checked
      ? "<strong><small>Online</small></strong>"
      : "<strong><small>Offline</small></strong>";

    onlineInput.value = target.checked ? "true" : "false";

    if ( target.checked ) {
      uncloakMe();
    } else {
      cloakMe();
    }
  } );
}

window.onload = setTimeout( () => checkOnlineStatus(), [ 250 ] );
