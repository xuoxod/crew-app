import { log } from "./clientutils.js";

const signedIn = () => {
  log(`\n\My profile\n`);
};

signedIn();

addEventListener("beforeunload", (event) => {
  log(`\n\tsignedin script Before unload\n`);
});
