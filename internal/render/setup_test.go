package render

import (
	"encoding/gob"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/xuoxod/crew-app/internal/config"
	"github.com/xuoxod/crew-app/internal/models"
)

var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M) {
	// Store a value in session
	gob.Register(models.Registration{})

	// Get the template cache from appConfg

	// Application mode
	testApp.InProduction = false

	// Session middleware
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false

	// Set the app level session
	testApp.Session = session

	app = &testApp

	os.Exit(m.Run())
}

type myWriter struct{}

func (tw *myWriter) Header() http.Header {
	var h http.Header

	return h
}

func (tw *myWriter) WriteHeader(i int) {

}

func (tx *myWriter) Write(b []byte) (int, error) {
	length := len(b)
	return length, nil
}
