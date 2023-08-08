import { log } from "./clientutils.js";

const signedIn = () => {
  log(`\n\tSigned In\n`);
};

signedIn();

addEventListener("beforeunload", (event) => {
  log(`\n\tsignedin script Before unload\n`);
});
