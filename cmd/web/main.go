package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/xuoxod/crew-app/internal/config"
	"github.com/xuoxod/crew-app/internal/driver"
	"github.com/xuoxod/crew-app/internal/envloader"
	"github.com/xuoxod/crew-app/internal/handlers"
	"github.com/xuoxod/crew-app/internal/helpers"
	"github.com/xuoxod/crew-app/internal/models"
	"github.com/xuoxod/crew-app/internal/render"
)

// Application configuration
var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	app.DBConnection = os.Getenv("DB_URL")

	err := envloader.LoadEnvVars()

	if err != nil {
		log.Fatal("error loading .env file")
	}

	db, err := run()

	if err != nil {
		log.Fatal(err)
	}

	defer db.SQL.Close()

	port := os.Getenv("PORT")

	if port == "" {
		port = ":8080"
	}

	fmt.Printf("\n\tServer listening on port %v\n\n", port)

	srv := &http.Server{
		Addr:    port,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()

	log.Fatal(err)
}

func run() (*driver.DB, error) {

	// Store a value in session
	gob.Register(models.Registration{})
	gob.Register(models.Member{})
	gob.Register(models.Profile{})
	gob.Register(models.UserSettings{})
	gob.Register(models.Signin{})
	gob.Register(models.Users{})

	// Get the template cache from appConfg

	// Application mode
	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// Session middleware
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = app.InProduction
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	// Set the app level session
	app.Session = session

	// Connect to database
	log.Println("Connecting to database ...")

	var host string = os.Getenv("DB_HOST")
	var user string = os.Getenv("DB_USER")
	var password string = os.Getenv("DB_PASSWD")
	var dbname string = os.Getenv("DB_NAME")
	var port int = 5432

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := driver.ConnectSql(psqlInfo)

	if err != nil {
		log.Fatal("Cannot connect to database! Dying ...")
	}

	tc, err := render.CreateTemplateCache()

	if err != nil {
		log.Fatal("Cannot create template cache")
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = true

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
