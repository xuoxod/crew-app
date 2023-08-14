package handlers

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/justinas/nosurf"
	"github.com/xuoxod/crew-app/internal/config"
	"github.com/xuoxod/crew-app/internal/helpers"
	"github.com/xuoxod/crew-app/internal/models"
	"github.com/xuoxod/crew-app/internal/render"
)

var app config.AppConfig
var session *scs.SessionManager

// var infoLog *log.Logger
// var errorLog *log.Logger
var pathToTemplates = "./../../templates"
var functions = template.FuncMap{}

func getRoutes() http.Handler {
	// Store a value in session
	gob.Register(models.Registration{})

	// Get the template cache from appConfg

	// Application mode
	app.InProduction = false

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// Session middleware
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	// Set the app level session
	app.Session = session

	tc, err := CreateTestTemplateCache()

	if err != nil {
		log.Fatal("Cannot create template cache")
		// return err
	}

	app.TemplateCache = tc
	app.UseCache = true

	repo := NewTestRepo(&app)
	NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	// Router
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(WriteToConsole)
	// mux.Use(NoSurf)
	mux.Use(SessionLoad)
	mux.Get("/", Repo.HomePage)
	mux.Get("/about", Repo.AboutPage)
	mux.Get("/register", Repo.RegisterPage)
	mux.Post("/register", Repo.PostRegisterPage)
	mux.Get("/registrationsummary", Repo.RegistrationSummary)
	mux.Post("/signin", Repo.SigninPage)
	mux.Get("/dummy", Repo.DummyHandler)
	mux.Get("/dashboard", Repo.Dashboard)

	// Static files
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux
}

func WriteToConsole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		remoteAddr := r.RemoteAddr
		host := r.Host
		path := r.URL.Path
		method := r.Method
		protocol := r.Proto
		protocolMajor := r.ProtoMajor
		protocolMinor := r.ProtoMinor

		fmt.Printf("\nPage Hit\n\tHost: %v\n\taddress: %v\n\tPath: %v\n\tMethod: %v\n\tProtocol: %v\n\t\tMajor: %v\n\t\tMinor: %v\n", host, remoteAddr, path, method, protocol, protocolMajor, protocolMinor)
		next.ServeHTTP(w, r)
	})
}

func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

func CreateTestTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))

	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		fmt.Printf("\n\tCurrent page is: %v\n", page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)

		if err != nil {
			return myCache, err
		}

		// Find base layout and parse
		matches, err := filepath.Glob(fmt.Sprintf(".%s/*.layout.tmpl", pathToTemplates))

		if err != nil {
			fmt.Println("Layout match? ", matches)
			return myCache, err
		}

		if len(matches) > 0 {
			fmt.Println("Matched Layout: ", matches)
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))

			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil

}
