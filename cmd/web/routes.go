package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/xuoxod/crew-app/internal/config"
	"github.com/xuoxod/crew-app/internal/handlers"
)

func routes(app *config.AppConfig) http.Handler {
	// mux := pat.New()

	// mux.Get("/", http.HandlerFunc(handlers.Repo.Index))
	// mux.Get("/about", http.HandlerFunc(handlers.Repo.About))

	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(WriteToConsole)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)
	// mux.Get("/about", handlers.Repo.AboutPage)
	// mux.Get("/register", handlers.Repo.RegisterPage)
	// mux.Post("/register", handlers.Repo.PostRegisterPage)
	// mux.Get("/registrationsummary", handlers.Repo.RegistrationSummary)
	// mux.Post("/signin", handlers.Repo.SigninPage)
	// mux.Get("/signout", handlers.Repo.SignOut)
	mux.Get("/dummy", handlers.Repo.DummyHandler)

	// Static files
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	mux.Route("/", func(mux chi.Router) {
		mux.Use(Unauth)
		mux.Get("/", handlers.Repo.HomePage)
		mux.Get("/about", handlers.Repo.AboutPage)
		mux.Get("/register", handlers.Repo.RegisterPage)
		mux.Post("/register", handlers.Repo.PostRegisterPage)
		mux.Get("/registrationsummary", handlers.Repo.RegistrationSummary)
		mux.Post("/signin", handlers.Repo.SigninPage)
	})

	mux.Route("/user", func(mux chi.Router) {
		mux.Use(Auth)
		mux.Get("/dashboard", handlers.Repo.Dashboard)
		mux.Get("/signout", handlers.Repo.SignOut)
	})

	return mux
}
