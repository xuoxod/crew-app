import { dlog, log } from "./clientutils.js";
import { stringify, parse, stripTags } from "./utils.js";
import {
  getElement,
  getElements,
  addOnChangeHandler,
  addClickHandler,
  addKeyupHandler,
} from "./computils.js";

const submitButton = getElement("submit-button");
const unameError = getElement("uname-error");
const uNameInput = getElement("uName");
const displayNameInput = getElement("displayName");
const isVisibleInput = getElement("isVisible");
const publicInput = getElement("public");
const photo = getElement("primary-photo");
const unblockUserIcons = getElements(".unblock");
const originalUsername = uNameInput.value.trim();
const usernames = getElements(".username");
const unames = [];
const unblocked = getElement("unblocked-user");
const goback = getElement("goback");

usernames.forEach((u) => {
  // const username = u.innerHTML.split("-")[1];
  const username = u.innerHTML;
  unames.push(username);
});

submitButton.disabled = true;

addKeyupHandler(uNameInput, (e) => {
  const target = e.target;
  target.value = target.value.trim();

  log(`Key Up event fired`);
  log(`Element ${target.id} value: ${target.value}\n\n`);

  checkUsername();
});

addOnChangeHandler(displayNameInput, (e) => {
  const target = e.target;

  log(`Change event fired`);
  log(`Element ${target.id} value: ${target.value}\n\n`);

  checkUsername();
});

addOnChangeHandler(isVisibleInput, (e) => {
  const target = e.target;

  log(`Change event fired`);
  log(`Element ${target.id} value: ${target.value}\n\n`);

  checkUsername();
});

addOnChangeHandler(publicInput, (e) => {
  const target = e.target;

  log(`Change event fired`);
  log(`Element ${target.id} value: ${target.value}\n\n`);

  checkUsername();
});

addClickHandler(photo, (e) => {
  let xmlHttp;

  try {
    xmlHttp = new XMLHttpRequest();

    xmlHttp.open("POST", "/chat/profile/update/photo", true);

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
          dlog(`Photo upload successful`);
          //  location.href = `/chat/profile/view/${blockerUid}`;
        } else {
          dlog(`Photo upload failed`);
        }
        return;
      }
    };

    xmlHttp.send(`photoPath=path/to/photo`, true);
  } catch (err) {
    tlog(err);
    return;
  }
});

addClickHandler(goback, (e) => {
  if (unblocked.value) {
    location.href = `/chat/room/enter?unblockeduserid=${unblocked.value}`;
  } else {
    location.href = `/chat/room/enter?unblockeduserid=${false}`;
  }
});

unblockUserIcons.forEach((icon) => {
  addClickHandler(icon, (e) => {
    const target = e.target.id.split("-")[1];
    const blocker = getElement("rmtId").value;

    dlog(`${blocker} is unblocking ${target}`);
    unblockUser(blocker, target);
  });
});

function nameTaken(str) {
  const index = unames.findIndex((x) => x == str);
  return index != -1;
}

function checkUsername() {
  if (
    uNameInput.value != originalUsername &&
    nameTaken(uNameInput.value.trim())
  ) {
    submitButton.disabled = true;
    unameError.classList.remove("d-none");
    unameError.innerHTML = `<strong>${uNameInput.value}</strong> is taken`;
  } else {
    submitButton.disabled = false;

    if (!unameError.classList.contains("d-none")) {
      unameError.classList.add("d-none");
    }
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
          const blockerdoc = responseJson.blockerdoc;
          const blockeedoc = responseJson.blockeedoc;

          unblocked.value = `${blockeedoc}`;

          dlog(
            `Blockee Doc:\t${stringify(blockeedoc)}`,
            `chatprofile: unblockUser`
          );

          /* TODO: Send back the unblocked user's ID */

          location.href = `/chat/profile/view/${blockerUid}?unblockeduser=${blockeedoc.user}`;
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
