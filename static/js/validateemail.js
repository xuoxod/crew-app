import { stringify, parse, stripTags } from "./utils.js";
import { log, cls, dlog } from "./clientutils.js";
import { addClickHandler, addKeyupHandler, addOnFocusHandler, getElement } from "./computils.js";

if ( window ) {
  cls();

  log( `Validate email` );

  addEventListener( "beforeunload", ( event ) => {
    log( `\n\tBefore unload\n` );
  } );
}

// Comps
const emailInput = getElement( "email" );
const submitButton = getElement( "reset-submit-button" );
const divEmailError = getElement( "divemailerror" );
const emailErrorInput = getElement( "emailerror" );
const resetForm = getElement( "reset-password" );
// const emailPattern = /(\w)+(\.(\w)+)?@(\w)+\.(\w){2,3}/i;

submitButton.disabled = true;
// emailInput.pattern = emailPattern;

addKeyupHandler( emailInput, ( e ) => {
  const target = e.target;
  const value = target.value.trim();
  const sanitizedValue = stripTags( value );
  target.value = sanitizedValue;
  cls();
  dlog( `${ sanitizedValue }\n` );

  submitButton.disabled = sanitizedValue.length > 0 ? false : true;
} );

addOnFocusHandler( emailInput, ( e ) => {
  if ( !divEmailError.classList.contains( "d-none" ) ) {
    divEmailError.classList.add( "d-none" );
    emailErrorInput.value = "";
  }
} );

resetForm.addEventListener( "submit", ( e ) => {
  e.preventDefault();
  validateEmail();
} );

function validateEmail () {
  let xmlHttp;

  try {
    xmlHttp = new XMLHttpRequest();

    xmlHttp.open( "POST", "/auth/validateuser", true );

    xmlHttp.setRequestHeader(
      "Content-type",
      "application/x-www-form-urlencoded"
    );

    xmlHttp.onload = () => {
      const responseText = xmlHttp.responseText;
      submitButton.disabled = true;

      if ( responseText ) {
        const responseJson = parse( responseText );
        const status = responseJson.status;

        if ( status ) {
          dlog( `User validated` );

          if ( !divEmailError.classList.contains( "d-none" ) ) {
            divEmailError.classList.add( "d-none" );
          }

          emailErrorInput.value = "";
          location.href = `/resetpassword/${ emailInput.value }`;
          //   location.href = `/chat/profile/view/${blockerUid}`;
        } else {
          const cause = responseJson.cause;
          dlog( `${ cause }` );

          if ( divEmailError.classList.contains( "d-none" ) ) {
            divEmailError.classList.remove( "d-none" );
          }

          emailErrorInput.value = `${ cause }`;
        }
        return;
      }
    };

    xmlHttp.send( `email=${ emailInput.value }`, true );
  } catch ( err ) {
    tlog( err );
    return;
  }
}
