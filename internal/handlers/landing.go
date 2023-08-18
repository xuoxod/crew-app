package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/xuoxod/crew-app/internal/config"
	"github.com/xuoxod/crew-app/internal/driver"
	"github.com/xuoxod/crew-app/internal/forms"
	"github.com/xuoxod/crew-app/internal/helpers"
	"github.com/xuoxod/crew-app/internal/libs"
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

func NewTestRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewTestingRepo(a),
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

	_ = render.Template(w, r, "home.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// @desc        About page
// @route       GET /about
// @access      Public
func (m *Repository) AboutPage(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["appname"] = "CrewMate"
	stringMap["appver"] = "Version 1.0"
	stringMap["appdate"] = fmt.Sprint("Date: ", time.Now())

	_ = render.Template(w, r, "about.page.tmpl", &models.TemplateData{StringMap: stringMap})
}

// @desc        Registration page
// @route       GET /register
// @access      Public
func (m *Repository) RegisterPage(w http.ResponseWriter, r *http.Request) {
	var emptyRegistrationForm models.Registration
	data := make(map[string]interface{})
	data["registration"] = emptyRegistrationForm

	_ = render.Template(w, r, "register.page.tmpl", &models.TemplateData{
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

		_ = render.Template(w, r, "register.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})

		return
	}

	// $2a$12$wPXDZd9LKDchbmN3cUn9K.Jqr73IfVqpx9XabDcGq5H0/8dMwqnYW

	newUser := models.User{
		FirstName: r.Form.Get("fname"),
		LastName:  r.Form.Get("lname"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
		Password:  r.Form.Get("pwd1"),
	}

	newUserID, err := m.DB.CreateUser(newUser)

	if err != nil {
		errMsg := strings.TrimSpace(strings.Split(err.Error(), ":")[1])
		fmt.Printf("\n\t\tError creating new user:\t%s\n\n", errMsg)
		m.App.Session.Put(r.Context(), "error", errMsg)
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	fmt.Printf("\tNew user ID: %v\n", newUserID)

	if err != nil {
		fmt.Printf("\n\tAnother Error Occurred\n\t%v\n\n", err.Error())
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

		_ = render.Template(w, r, "home.page.tmpl", &models.TemplateData{Form: form, Data: data})
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
	http.Redirect(w, r, "/user/dashboard", http.StatusSeeOther)
}

// @desc        Signin user
// @route       POST /signin
// @access      Public
func (m *Repository) LoginPage(w http.ResponseWriter, r *http.Request) {
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

		_ = render.Template(w, r, "home.page.tmpl", &models.TemplateData{Form: form, Data: data})
		return
	}

	var email string
	var password string

	email = r.Form.Get("email")
	password = r.Form.Get("password")

	// Authenticate user
	results := m.DB.AuthenticateUser(email, password)

	if results["err"] != "" {
		log.Println("unable to authenticate user")
		err = errors.New("unable to authenticate user")
	}

	if err != nil {
		log.Println("Authentication failed")
		m.App.Session.Put(r.Context(), "error", "Invalid signin credentials")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	userID, _ := strconv.Atoi(results["userID"])
	firstName := libs.Cap(results["firstName"])
	lastName := libs.Cap(results["lastName"])
	emailAddress := results["email"]
	phone := results["phone"]
	accessLevel, _ := strconv.Atoi(results["accessLevel"])
	craftID, _ := strconv.Atoi(results["craftID"])
	updatedAt := strings.Split(results["updatedAt"], " ")[0]
	createdAt := strings.Split(results["createdAt"], " ")[0]
	userName := results["userName"]
	imgUrl := results["imgUrl"]
	profileStatus := results["profileStatus"]

	const layout = "2006-01-02"
	creationDate, creationErr := time.Parse(layout, createdAt)

	if creationErr != nil {
		helpers.ServerError(w, creationErr)
		return
	}
	updatedLast, updateErr := time.Parse(layout, updatedAt)

	if updateErr != nil {
		helpers.ServerError(w, updateErr)
		return
	}

	switch profileStatus {
	case "noprofile":
		fmt.Printf("User is authenticated without a profile\n\tID:\t%d\n\tFirst Name:\t%s\n\tLast Name:\t%s\n\tEmail:\t%s\n\tPhone:\t%s\n\tCreated At:\t%v\n\tUpdated At:\t%v\n\n", userID, firstName, lastName, emailAddress, phone, createdAt, updatedAt)

		loggedIn := models.User{
			ID:           userID,
			FirstName:    firstName,
			LastName:     lastName,
			Email:        emailAddress,
			Phone:        phone,
			AccessLevel:  accessLevel,
			CraftID:      craftID,
			CreatedAt:    creationDate,
			UpdatedAt:    updatedLast,
			CreationDate: createdAt,
			Updated:      updatedAt,
			ImageURL:     "no",
			HasID:        "yes",
		}

		data := make(map[string]interface{})
		data["loggedin"] = loggedIn

		m.App.Session.Put(r.Context(), "user_id", userID)
		m.App.Session.Put(r.Context(), "loggedin", loggedIn)

		http.Redirect(w, r, "/user/dashboard", http.StatusSeeOther)
	case "hasprofile":
		fmt.Printf("User is authenticated with a profile\n\tFirst Name:\t%s\n\tLast Name:\t%s\n\tEmail:\t%s\n\tPhone:\t%s\n\tCreated At:\t%v\n\tUpdated At:\t%v\n\tUsername:\t%v\n\tImage URL:\t%v\n\n", firstName, lastName, emailAddress, phone, createdAt, updatedAt, userName, imgUrl)

		loggedIn := models.User{
			FirstName:    firstName,
			LastName:     lastName,
			Email:        emailAddress,
			Phone:        phone,
			AccessLevel:  accessLevel,
			CreatedAt:    creationDate,
			UpdatedAt:    updatedLast,
			CreationDate: createdAt,
			Updated:      updatedAt,
			ImageURL:     imgUrl,
			Username:     userName,
		}

		data := make(map[string]interface{})
		data["loggedin"] = loggedIn

		m.App.Session.Put(r.Context(), "user_id", emailAddress)
		m.App.Session.Put(r.Context(), "loggedin", loggedIn)

		http.Redirect(w, r, "/user/dashboard", http.StatusSeeOther)

	}
}

// @desc        Signout user
// @route       GET /signin
// @access      Private
func (m *Repository) SignOut(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// @desc        Logout user
// @route       GET /signin
// @access      Private
func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/", http.StatusSeeOther)
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
	_, _ = w.Write(out)
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

	_ = render.Template(w, r, "registrationsummary.page.tmpl", &models.TemplateData{Data: data})
}

// @desc        Dashboard Page
// @route       GET /user/dashboard
// @access      Private
func (m *Repository) Dashboard(w http.ResponseWriter, r *http.Request) {
	loggedin, ok := m.App.Session.Get(r.Context(), "loggedin").(models.User)

	if !ok {
		log.Println("Cannot get item from session")
		m.App.ErrorLog.Println("Can't get error from the session")
		m.App.Session.Put(r.Context(), "error", "Can't get loggedin from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	data := make(map[string]interface{})
	data["loggedin"] = loggedin

	_ = render.Template(w, r, "dashboard.page.tmpl", &models.TemplateData{Data: data})
}

// @desc        Settings Page
// @route       GET /user/settings
// @access      Private
func (m *Repository) SettingsPage(w http.ResponseWriter, r *http.Request) {
	var emptyUserForm models.User
	data := make(map[string]interface{})
	data["setting"] = emptyUserForm

	loggedin, ok := m.App.Session.Get(r.Context(), "loggedin").(models.User)

	if !ok {
		log.Println("Cannot get item from session")
		m.App.ErrorLog.Println("Can't get error from the session")
		m.App.Session.Put(r.Context(), "error", "Can't get loggedin from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	data["loggedin"] = loggedin

	_ = render.Template(w, r, "settings.page.tmpl", &models.TemplateData{Data: data, Form: forms.New(nil)})
}

// @desc        User Profile Page
// @route       GET /user/update
// @access      Private
func (m *Repository) ProfilePage(w http.ResponseWriter, r *http.Request) {
	var emptyUserForm models.User
	data := make(map[string]interface{})
	data["profile"] = emptyUserForm

	fmt.Println("Get Profile Page")

	loggedin, ok := m.App.Session.Get(r.Context(), "loggedin").(models.User)

	if !ok {
		log.Println("Cannot get item from session")
		m.App.ErrorLog.Println("Can't get error from the session")
		m.App.Session.Put(r.Context(), "error", "Can't get loggedin from session")
		http.Redirect(w, r, "/user/dashboard", http.StatusTemporaryRedirect)
		return
	}

	// fmt.Printf("\n\n\t\tLoggedin Data:\t\t%v\n\n", loggedin)

	data["loggedin"] = loggedin
	_ = render.Template(w, r, "profile.page.tmpl", &models.TemplateData{
		Data: data,
		Form: forms.New(nil),
	})
}

// @desc        Update user profile
// @route       POST /signin
// @access      private
func (m *Repository) UpdateTheUserProfile(w http.ResponseWriter, r *http.Request) {
	// Remove the old user data from the session
	m.App.Session.Remove(r.Context(), "loggedin")

	err := r.ParseForm()

	if err != nil {
		fmt.Printf("\n\tError parsing user profile form")
		helpers.ServerError(w, err)
		return
	}

	updatedUser := models.User{
		FirstName: r.Form.Get("fname"),
		LastName:  r.Form.Get("lname"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
		Username:  r.Form.Get("uname"),
		ImageURL:  r.Form.Get("imgurl"),
	}

	// fmt.Printf("Parsed Form:\t%v\n", updatedUser)

	// form validation

	form := forms.New(r.PostForm)
	form.MinLength("fname", 2, r)
	form.MinLength("lname", 2, r)
	form.IsEmail("email")
	form.Required("email", "fname", "lname", "uname", "phone", "imgurl")

	if !form.Valid() {
		fmt.Printf("\n\tForm Error:\t%v\n\n", form.Errors)
		loggedin, ok := m.App.Session.Get(r.Context(), "loggedin").(models.User)

		if !ok {
			log.Println("Cannot get item from session")
			m.App.ErrorLog.Println("Can't get error from the session")
			m.App.Session.Put(r.Context(), "error", "Can't get loggedin from session")
			http.Redirect(w, r, "/user/dashboard", http.StatusTemporaryRedirect)
			return
		}

		data := make(map[string]interface{})
		data["loggedin"] = loggedin

		_ = render.Template(w, r, "profile.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	// fmt.Printf("\nUpdate the user's profile\n")
	results := m.DB.UpdateUserProfile(updatedUser)

	// fmt.Printf("\n\n\t\tUpdated User Profile Data\n\t\t\t%v\n\n", results)

	if results["err"] != "" {
		fmt.Printf("\n\tProfile update error: %s\n\n", results["err"])
	}

	userId, _ := strconv.Atoi(results["userID"])
	craftId, _ := strconv.Atoi(results["craftID"])

	/* 	const layout = "2006-01-02"
	   	creationDate, _ := time.Parse(layout, results["createdAt"])
	   	updatedDate, _ := time.Parse(layout, results["updatedAt"]) */

	loggedIn := models.User{
		FirstName:    results["firstName"],
		LastName:     results["lastName"],
		ID:           userId,
		Email:        results["email"],
		Phone:        results["phone"],
		CraftID:      craftId,
		Username:     results["userName"],
		ImageURL:     results["imageUrl"],
		CreationDate: results["createdAt"],
		Updated:      results["updatedAt"],
		HasID:        "yes",
	}

	fmt.Println("\tUpdated User Profile Data")
	fmt.Println("User ID: ", results["userID"])
	fmt.Println("First Name: ", results["firstName"])
	fmt.Println("Last Name: ", results["lastName"])
	fmt.Println("Email: ", results["email"])
	fmt.Println("Phone: ", results["phone"])
	fmt.Println("Craft ID: ", results["craftID"])
	fmt.Println("Created: ", results["createdAt"])
	fmt.Println("Updated: ", results["updatedAt"])
	fmt.Printf("\n\n")

	// Add the update user data to the session
	m.App.Session.Put(r.Context(), "loggedin", loggedIn)

	http.Redirect(w, r, "/user/profile", http.StatusSeeOther)
}

// @desc        Create user profile
// @route       POST /signin
// @access      private
func (m *Repository) CreateUserProfile(w http.ResponseWriter, r *http.Request) {
	// Remove the old user data from the session
	m.App.Session.Remove(r.Context(), "loggedin")

	err := r.ParseForm()

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	user_id, err := strconv.Atoi(r.Form.Get("user_id"))

	if err != nil {
		fmt.Println("user_id conversion failed")
	}

	updatedUser := models.User{
		ID:        user_id,
		FirstName: r.Form.Get("fname"),
		LastName:  r.Form.Get("lname"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
		Username:  r.Form.Get("uname"),
		ImageURL:  r.Form.Get("imgurl"),
	}

	// form validation

	form := forms.New(r.PostForm)
	form.MinLength("fname", 2, r)
	form.MinLength("lname", 2, r)
	form.IsEmail("email")
	// form.Required("email", "fname", "lname", "uname", "phone", "imgurl")

	if !form.Valid() {
		fmt.Printf("\n\tForm Error:\t%v\n\n", form.Errors)
		loggedin, ok := m.App.Session.Get(r.Context(), "loggedin").(models.User)

		if !ok {
			log.Println("Cannot get item from session")
			m.App.ErrorLog.Println("Can't get error from the session")
			m.App.Session.Put(r.Context(), "error", "Can't get loggedin from session")
			return
		}

		data := make(map[string]interface{})
		data["loggedin"] = loggedin

		_ = render.Template(w, r, "profile.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
	}

	results := m.DB.CreateUserProfile(updatedUser)
	userId, _ := strconv.Atoi(results["userID"])
	craftId, _ := strconv.Atoi(results["craftID"])

	/* 	const layout = "2006-01-02"
	   	creationDate, _ := time.Parse(layout, results["createdAt"])
	   	updatedDate, _ := time.Parse(layout, results["updatedAt"]) */

	loggedIn := models.User{
		FirstName:    results["firstName"],
		LastName:     results["lastName"],
		ID:           userId,
		Email:        results["email"],
		Phone:        results["phone"],
		CraftID:      craftId,
		Username:     results["userName"],
		ImageURL:     results["imageUrl"],
		CreationDate: results["createdAt"],
		Updated:      results["updatedAt"],
		HasID:        "yes",
	}

	fmt.Println("\tCreate User Profile Data")
	fmt.Println("User ID: ", results["userID"])
	fmt.Println("First Name: ", results["firstName"])
	fmt.Println("Last Name: ", results["lastName"])
	fmt.Println("Email: ", results["email"])
	fmt.Println("Phone: ", results["phone"])
	fmt.Println("Craft ID: ", results["craftID"])
	fmt.Println("User Name: ", results["userName"])
	fmt.Println("Image URL: ", results["imageUrl"])
	fmt.Println("Created: ", results["createdAt"])
	fmt.Println("Updated: ", results["updatedAt"])
	fmt.Printf("\n\n")

	// Add the update user data to the session
	m.App.Session.Put(r.Context(), "loggedin", loggedIn)

	http.Redirect(w, r, "/user/profile", http.StatusSeeOther)
}
