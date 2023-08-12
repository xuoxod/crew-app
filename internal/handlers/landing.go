package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/xuoxod/crew-app/internal/config"
	"github.com/xuoxod/crew-app/internal/driver"
	"github.com/xuoxod/crew-app/internal/forms"
	"github.com/xuoxod/crew-app/internal/helpers"
	"github.com/xuoxod/crew-app/internal/models"
	"github.com/xuoxod/crew-app/internal/render"
	"github.com/xuoxod/crew-app/internal/repository"
	"github.com/xuoxod/crew-app/internal/repository/dbrepo"
)

// The repository used by the handlers
var Repo *Repository

// Repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// Creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// Sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// @desc        Home page
// @route       GET /
// @access      Public
func (m *Repository) HomePage(w http.ResponseWriter, r *http.Request) {
	var emptySigninForm models.Signin
	data := make(map[string]interface{})
	data["signin"] = emptySigninForm

	fmt.Printf("\n\tHome page\n")

	render.Template(w, r, "home.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// @desc        About page
// @route       GET /about
// @access      Public
func (m *Repository) AboutPage(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["subheading"] = "Who we are is totally and completely irrelevant"

	render.Template(w, r, "about.page.tmpl", &models.TemplateData{StringMap: stringMap})
}

// @desc        Registration page
// @route       GET /register
// @access      Public
func (m *Repository) RegisterPage(w http.ResponseWriter, r *http.Request) {
	var emptyRegistrationForm models.Registration
	data := make(map[string]interface{})
	data["registration"] = emptyRegistrationForm

	render.Template(w, r, "register.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// @desc        Register user
// @route       POST /auth/register
// @access      Public
func (m *Repository) PostRegisterPage(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	registration := models.Registration{
		FirstName:       r.Form.Get("fname"),
		LastName:        r.Form.Get("lname"),
		Email:           r.Form.Get("email"),
		Phone:           r.Form.Get("phone"),
		PasswordCreate:  r.Form.Get("pwd1"),
		PasswordConfirm: r.Form.Get("pwd2"),
	}

	// form validation
	fmt.Println("registration posted")

	form := forms.New(r.PostForm)
	form.MinLength("fname", 2, r)
	form.MinLength("lname", 2, r)
	form.IsEmail("email")
	form.PasswordsMatch("pwd1", "pwd2", r)
	form.Required("fname", "lname", "email", "phone", "pwd1", "pwd2")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["registration"] = registration

		render.Template(w, r, "register.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})

		return
	}

	newUser := models.User{
		FirstName: r.Form.Get("fname"),
		LastName:  r.Form.Get("lname"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
		Password:  r.Form.Get("pwd1"),
	}

	newUserID, err := m.DB.CreateUser(newUser)

	if err != nil {
		fmt.Printf("\n\tError creating new user\n\t%v\n\n", err.Error())

		m.App.Session.Put(r.Context(), "error", "User already registered")
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	fmt.Printf("\tNew user ID: %v\n", newUserID)

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "registration", registration)

	http.Redirect(w, r, "/registrationsummary", http.StatusSeeOther)
}

// @desc        Signin user
// @route       POST /signin
// @access      Public
func (m *Repository) SigninPage(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	signin := models.Signin{
		Email:    r.Form.Get("email"),
		Password: r.Form.Get("password"),
	}

	// form validation
	fmt.Println("signin posted")

	form := forms.New(r.PostForm)
	form.IsEmail("email")
	form.Required("email", "password")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["signin"] = signin

		render.Template(w, r, "home.page.tmpl", &models.TemplateData{Form: form, Data: data})
		return
	}

	var email string
	var password string

	email = r.Form.Get("email")
	password = r.Form.Get("password")

	// Authenticate user
	id, _, err := m.DB.Authenticate(email, password)

	if err != nil {
		log.Println("Authentication failed")
		m.App.Session.Put(r.Context(), "error", "Invalid signin credentials")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "flash", "Authenticated succesfully")
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

func (m *Repository) DummyHandler(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		OK:      true,
		Message: "Good to go",
	}

	out, err := json.MarshalIndent(resp, "", "   ")

	if err != nil {
		helpers.ServerError(w, err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// @desc        Registration summary
// @route       GET /registrationsummary
// @access      Public
func (m *Repository) RegistrationSummary(w http.ResponseWriter, r *http.Request) {
	registration, ok := m.App.Session.Get(r.Context(), "registration").(models.Registration)

	if !ok {
		log.Println("Cannot get item from session")
		m.App.ErrorLog.Println("Can't get error from the session")
		m.App.Session.Put(r.Context(), "error", "Can't get registration from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Remove(r.Context(), "registration")
	data := make(map[string]interface{})
	data["registration"] = registration

	render.Template(w, r, "registrationsummary.page.tmpl", &models.TemplateData{Data: data})
}

// @desc        Dashboard
// @route       GET /dashboard
// @access      Public
func (m *Repository) Dashboard(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["subheading"] = "Dashboard"
	stringMap["body"] = "This is the house that Jack built."
	stringMap["footer"] = "This is the malt that lay in the house that Jack built."

	render.Template(w, r, "dashboard.page.tmpl", &models.TemplateData{StringMap: stringMap})
}
