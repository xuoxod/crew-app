import {
  getElement,
  newElement,
  addAttribute,
  removeChildren,
  appendChild,
} from "./computils.js";

import { dlog, log } from "./clientutils.js";

import { stringify, parse } from "./utils.js";

import { showMediaControls } from "./ui.js";

const printObject = (arg = {}, label = `participantmanager`) =>
  dlog(`${stringify(arg)}`, label);

function printNetworkQualityStats(networkQualityLevel, networkQualityStats) {
  // Print in console the networkQualityLevel using bars
  console.log(
    {
      1: "▃",
      2: "▃▄",
      3: "▃▄▅",
      4: "▃▄▅▆",
      5: "▃▄▅▆▇",
    }[networkQualityLevel] || ""
  );

  if (networkQualityStats) {
    // Print in console the networkQualityStats, which is non-null only if Network Quality
    // verbosity is 2 (moderate) or greater
    console.log("Network Quality statistics:", networkQualityStats);
  }
}

const handleTrackPublication = (trackPublication, participant) => {
  function displayTrack(track) {
    // append this track to the participant's div and render it on the page
    const participantDiv = document.getElementById(participant.identity);
    // track.attach creates an HTMLVideoElement or HTMLAudioElement
    // (depending on the type of track) and adds the video or audio stream
    participantDiv.append(track.attach());
  }

  if (trackPublication.track) {
    displayTrack(trackPublication.track);
  }

  // listen for any new subscriptions to this track publication
  trackPublication.on("subscribed", displayTrack);
};

export const handleDisconnection = (participant) => {
  dlog(
    `${participant.identity} disconnected`,
    `participantmanager: on disconnected`
  );

  try {
    document.querySelector(`#${participant.identity}`).remove();
  } catch (err) {
    return;
  }
  // participantDisconnected({ rmtUser: rmtIdInput.value });
};

export const localParticipantHandler = (room) => {
  // Access parent element and create the child elements
  const participant = room.localParticipant;

  const parent = getElement("conn-parent");
  const localPart = newElement("div");
  const rmtIdInput = getElement("rmtid-input");

  // Add element attributes

  addAttribute(localPart, "id", participant.identity);
  addAttribute(localPart, "class", "local-video-container");

  // Append elements

  appendChild(parent, localPart);

  // Handle the participant track publication event

  participant.tracks.forEach((trackPublicaton) => {
    if (trackPublicaton.track) {
      localPart.append(trackPublicaton.track.attach());
    }
  });

  participant.on("trackPublished", (trackPub) => {
    function displayTrack(track) {
      localPart.append(track.attach());
    }
  });

  participant.on("trackDisabled", (track) => {
    const objData = {
      activityType: `Track Disabled`,
      trackKind: track.kind,
      trackId: track.id,
      rmtUserId: rmtIdInput.value,
      participantSid: participant.sid,
      participantIdentity: participant.identity,
    };
    // userActivity(objData);
    printObject(objData, `participantmanager: on trackDisabled`);
  });

  participant.on("trackEnabled", (track) => {
    const objData = {
      activityType: `Track Enabled`,
      trackKind: track.kind,
      trackId: track.id,
      rmtUserId: rmtIdInput.value,
      participantSid: participant.sid,
      participantIdentity: participant.identity,
    };

    // userActivity(objData);
    printObject(objData, `participantmanager: on trackEnabled`);
  });

  participant.on("trackStarted", (track) => {
    const objData = {
      activityType: `Track Started`,
      trackKind: track.kind,
      trackId: track.id,
      rmtUserId: rmtIdInput.value,
      participantSid: participant.sid,
      participantIdentity: participant.identity,
    };
    // userActivity(objData);
    printObject(objData, `participantmanager: on trackStarted`);
  });

  participant.on("disconnected", handleDisconnection);

  const micControlHandler = (e) => {
    const icon = e.target;
    const cl = icon.classList;
    let enabled = cl.contains("bi-mic-fill");

    if (enabled) {
      cl.remove("bi-mic-fill");
      cl.add("bi-mic-mute-fill");
      room.localParticipant.audioTracks.forEach((publication) => {
        publication.track.disable();
      });
    } else {
      cl.remove("bi-mic-mute-fill");
      cl.add("bi-mic-fill");
      room.localParticipant.audioTracks.forEach((publication) => {
        publication.track.enable();
      });
    }

    dlog(`Microphone icon clicked`);
  };

  const vidControlHandler = (e) => {
    const icon = e.target;
    const cl = icon.classList;
    let enabled = cl.contains("bi-camera-video-fill");

    if (enabled) {
      cl.remove("bi-camera-video-fill");
      cl.add("bi-camera-video-off-fill");
      room.localParticipant.videoTracks.forEach((publication) => {
        publication.track.disable();
      });
    } else {
      cl.remove("bi-camera-video-off-fill");
      cl.add("bi-camera-video-fill");
      room.localParticipant.videoTracks.forEach((publication) => {
        publication.track.enable();
      });
    }
    dlog(`Video icon clicked`);
  };

  showMediaControls(micControlHandler, vidControlHandler);
};

export const remoteParticipantHandler = (participant) => {
  // Access parent element and create the child elements

  const parent = getElement("conn-parent");
  const localPart = newElement("div");
  const rmtIdInput = document.querySelector("#rmtid-input");

  // Add element attributes

  addAttribute(localPart, "id", participant.identity);
  addAttribute(localPart, "class", "remote-video-container");

  // Append elements

  appendChild(parent, localPart);

  // Handle the participant track publication event

  participant.tracks.forEach((trackPub) => {
    function displayTrack(track) {
      const remoteDiv = document.querySelector(`#${participant.identity}`);
      remoteDiv.append(track.attach());
    }

    if (trackPub.track) {
      displayTrack(trackPub.track);
    }

    // listen for any new subscriptions to this track publication
    trackPub.on("subscribed", displayTrack);
  });

  participant.on("trackPublished", (trackPublication) => {
    function displayTrack(track) {
      const remoteDiv = document.querySelector(`#${participant.identity}`);
      remoteDiv.append(track.attach());
    }

    if (trackPublication.track) {
      displayTrack(trackPublication.track);
    }

    // listen for any new subscriptions to this track publication
    trackPublication.on("subscribed", displayTrack);
  });

  participant.on("trackDisabled", (track) => {
    const objData = {
      activityType: `Track Disabled`,
      trackKind: track.kind,
      trackId: track.id,
      rmtUserId: rmtIdInput.value,
      participantSid: participant.sid,
      participantIdentity: participant.identity,
    };
    // userActivity(objData);
    printObject(objData, `participantmanager: on trackDisabled`);
  });

  participant.on("trackEnabled", (track) => {
    const objData = {
      activityType: `Track Enabled`,
      trackKind: track.kind,
      trackId: track.id,
      rmtUserId: rmtIdInput.value,
      participantSid: participant.sid,
      participantIdentity: participant.identity,
    };

    // userActivity(objData);
    printObject(objData, `participantmanager: on trackEnabled`);
  });

  participant.on("trackStarted", (track) => {
    const objData = {
      activityType: `Track Started`,
      trackKind: track.kind,
      trackId: track.id,
      rmtUserId: rmtIdInput.value,
      participantSid: participant.sid,
      participantIdentity: participant.identity,
    };
    // userActivity(objData);
    printObject(objData, `participantmanager: on trackStarted`);
  });

  participant.on("disconnected", handleDisconnection);
};

export const logConnectedParticipants = (participant) => {
  dlog(`Participant ${participant.identity}`);
  remoteParticipantHandler(participant);
};
