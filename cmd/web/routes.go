package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/xuoxod/crew-app/internal/config"
	"github.com/xuoxod/crew-app/internal/handlers"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(WriteToConsole)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)
	mux.Get("/dummy", handlers.Repo.DummyHandler)
	mux.Get("/admin", handlers.Repo.AdminPage)

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
		mux.Post("/login", handlers.Repo.LoginPage)
	})

	mux.Route("/user", func(mux chi.Router) {
		mux.Use(Auth)
		mux.Get("/dashboard", handlers.Repo.Dashboard)
		mux.Get("/settings", handlers.Repo.SettingsPage)
		mux.Post("/settings/update", handlers.Repo.PostSettingsPage)
		mux.Get("/profile", handlers.Repo.ProfilePage)
		mux.Post("/profile/update", handlers.Repo.PostProfilePage)
		mux.Get("/signout", handlers.Repo.SignOut)
		mux.Get("/logout", handlers.Repo.Logout)
	})

	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(Admin)
		mux.Get("/dashboard", handlers.Repo.AdminPage)
		mux.Get("/users", handlers.Repo.UsersPage)
		mux.Get("/user", handlers.Repo.UserPage)
		mux.Post("/user/update", handlers.Repo.PostUserPage)
		mux.Get("/user/remove", handlers.Repo.RemoveUser)
	})

	return mux
}
