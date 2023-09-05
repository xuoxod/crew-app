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

	newUserData := models.Member{
		FirstName: r.Form.Get("fname"),
		LastName:  r.Form.Get("lname"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
		Password:  r.Form.Get("pwd1"),
	}

	createdUser, err := m.DB.CreateUser(newUserData)

	if err != nil {
		fmt.Printf("\n\t\tError creating new user:\t%s\n\n", err.Error())
		m.App.Session.Put(r.Context(), "error", "Error registering user")
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	/* 	if err != nil {
		fmt.Printf("\n\tAnother Error Occurred\n\t%v\n\n", err.Error())
		helpers.ServerError(w, err)
		return
	} */

	id, _ := strconv.Atoi(createdUser["ID"])

	fmt.Println("\tUser Registered")
	fmt.Println("        ID:\t", id)
	fmt.Println("First Name:\t", createdUser["firstName"])
	fmt.Println(" Last Name:\t", createdUser["lastName"])
	fmt.Println("     Email:\t", createdUser["email"])
	fmt.Println("     Phone:\t", createdUser["phone"])

	m.App.Session.Put(r.Context(), "registration", registration)

	http.Redirect(w, r, "/registrationsummary", http.StatusSeeOther)
}

// @desc        Registration summary
// @route       GET /registrationsummary
// @access      Public
func (m *Repository) RegistrationSummary(w http.ResponseWriter, r *http.Request) {
	registration, ok := m.App.Session.Get(r.Context(), "registration").(models.Registration)

	if !ok {
		log.Println("Cannot get registration data from session")
		m.App.ErrorLog.Println("Can't get registration data from the session")
		m.App.Session.Put(r.Context(), "error", "Can't get registration data from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Remove(r.Context(), "registration")
	data := make(map[string]interface{})
	data["registration"] = registration

	_ = render.Template(w, r, "registrationsummary.page.tmpl", &models.TemplateData{Data: data})
}

// @desc        Login user
// @route       POST /login
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
	userName := results["userName"]
	imgUrl := results["imgUrl"]
	craft := results["craft"]
	address := results["address"]
	city := results["city"]
	state := results["state"]
	displayName := results["displayName"]
	accessLevel, _ := strconv.Atoi(results["accessLevel"])
	createdAt := strings.Split(results["createdAt"], " ")[0]
	updatedAt := strings.Split(results["updatedAt"], " ")[0]
	yearsService := results["yos"]
	showProfile := results["showProfile"]
	showOnlineStatus := results["showOnlineStatus"]
	showAddress := results["showAddress"]
	showCity := results["showCity"]
	showState := results["showState"]
	showDisplayName := results["showDisplayName"]
	showContactInfo := results["showContactInfo"]
	showPhone := results["showPhone"]
	showEmail := results["showEmail"]
	showCraft := results["showCraft"]
	showRun := results["showRun"]
	showNotifications := results["showNotifications"]
	problem := results["err"]

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

	if problem != "" {
		fmt.Println("Authentication error:\t", problem)
		m.App.Session.Put(r.Context(), "error", "Authentication Failed")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	fmt.Println("")
	fmt.Println("User Authenticated")
	fmt.Println("ID:\t", userID)
	fmt.Println("First Name:\t", firstName)
	fmt.Println("Last Name:\t", lastName)
	fmt.Println("Email Address:\t", email)
	fmt.Println("Phone:\t", phone)
	fmt.Println("User Name:\t", userName)
	fmt.Println("Image URL:\t", imgUrl)
	fmt.Println("Craft:\t", craft)
	fmt.Println("Address:\t", address)
	fmt.Println("City:\t", city)
	fmt.Println("State:\t", state)
	fmt.Println("Display Name:\t", displayName)
	fmt.Println("Access Level:\t", accessLevel)
	fmt.Println("Created At:\t", createdAt)
	fmt.Println("Last Updated:\t", updatedAt)
	fmt.Println("Years of Service:\t", yearsService)
	fmt.Println("Show Profile:\t", showProfile)
	fmt.Println("Show Online Status:\t", showOnlineStatus)
	fmt.Println("Show City:\t", showCity)
	fmt.Println("Show State:\t", showState)
	fmt.Println("Show Display Name:\t", showDisplayName)
	fmt.Println("Show Contact Info:\t", showContactInfo)
	fmt.Println("Show Phone:\t", showPhone)
	fmt.Println("Show Email:\t", showEmail)
	fmt.Println("Show Craft:\t", showCraft)
	fmt.Println("Show Run:\t", showRun)
	fmt.Println("Show Notifications:\t", showNotifications)
	fmt.Println("")

	yos, _ := strconv.Atoi(yearsService)
	profileShow, _ := strconv.ParseBool(showProfile)
	onlineStatusShow, _ := strconv.ParseBool(showOnlineStatus)
	addressShow, _ := strconv.ParseBool(showAddress)
	cityShow, _ := strconv.ParseBool(showCity)
	stateShow, _ := strconv.ParseBool(showState)
	displayNameShow, _ := strconv.ParseBool(showDisplayName)
	contactInfoShow, _ := strconv.ParseBool(showContactInfo)
	phoneShow, _ := strconv.ParseBool(showPhone)
	emailShow, _ := strconv.ParseBool(showEmail)
	craftShow, _ := strconv.ParseBool(showCraft)
	runShow, _ := strconv.ParseBool(showRun)
	notificationsShow, _ := strconv.ParseBool(showNotifications)

	member := models.Member{
		ID:           userID,
		FirstName:    firstName,
		LastName:     lastName,
		Email:        emailAddress,
		Phone:        phone,
		AccessLevel:  accessLevel,
		UpdatedAt:    updatedLast,
		Updated:      updatedAt,
		CreatedAt:    creationDate,
		CreationDate: createdAt,
	}

	userProfile := models.Profile{
		MemberID:    userID,
		UserName:    userName,
		ImageURL:    imgUrl,
		Craft:       craft,
		Address:     address,
		City:        city,
		State:       state,
		DisplayName: displayName,
		YOS:         yos,
	}

	userSettings := models.UserSettings{
		MemberID:          userID,
		ShowProfile:       profileShow,
		ShowOnlineStatus:  onlineStatusShow,
		ShowAddress:       addressShow,
		ShowCity:          cityShow,
		ShowState:         stateShow,
		ShowDisplayName:   displayNameShow,
		ShowContactInfo:   contactInfoShow,
		ShowPhone:         phoneShow,
		ShowEmail:         emailShow,
		ShowCraft:         craftShow,
		ShowRun:           runShow,
		ShowNotifications: notificationsShow,
	}

	users := m.DB.AllUsers()

	strErr := users["err"][0]
	if strErr != "" {
		fmt.Println("Post login DB error getting all users:\t", strErr)
		return
	}

	delete(users, "err")

	allUsers := models.Users{
		AllUsers: users,
	}

	m.App.Session.Put(r.Context(), "user_id", member)
	m.App.Session.Put(r.Context(), "loggedin", member)
	m.App.Session.Put(r.Context(), "user_profile", userProfile)
	m.App.Session.Put(r.Context(), "user_settings", userSettings)
	m.App.Session.Put(r.Context(), "allusers", allUsers)

	if results["accessLevel"] != "1" {
		http.Redirect(w, r, "/user/dashboard", http.StatusSeeOther)
	}

	http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)

}

// @desc        Signout user
// @route       GET /signout
// @access      Private
func (m *Repository) SignOut(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// @desc        Logout user
// @route       GET /logout
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

// @desc        Dashboard Page
// @route       GET /user/dashboard
// @access      Private
func (m *Repository) Dashboard(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})

	fmt.Println("Get Dashboard Page")

	loggedin, loggedInOk := m.App.Session.Get(r.Context(), "loggedin").(models.Member)
	profile, profileOk := m.App.Session.Get(r.Context(), "user_profile").(models.Profile)

	if !loggedInOk {
		log.Println("Cannot get loggedin session")
		m.App.ErrorLog.Println("Can't get loggedin from the session")
		m.App.Session.Put(r.Context(), "error", "Can't get loggedin from session")
		http.Redirect(w, r, "/user/dashboard", http.StatusTemporaryRedirect)
		return
	}

	if !profileOk {
		log.Println("Cannot get profile session")
		m.App.ErrorLog.Println("Can't get profile from the session")
		m.App.Session.Put(r.Context(), "error", "Can't get profile from session")
		http.Redirect(w, r, "/user/dashboard", http.StatusTemporaryRedirect)
		return
	}

	data["loggedin"] = loggedin
	data["profile"] = profile

	if loggedin.AccessLevel == 1 {
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
	} else {
		_ = render.Template(w, r, "dashboard.page.tmpl", &models.TemplateData{Data: data})
	}
}

// @desc        Settings Page
// @route       GET /user/settings
// @access      Private
func (m *Repository) SettingsPage(w http.ResponseWriter, r *http.Request) {
	var emptyUserSettingsForm models.UserSettings

	fmt.Println("Get Settings Page")

	loggedin, loggedinOk := m.App.Session.Get(r.Context(), "loggedin").(models.Member)
	settings, settingsOk := m.App.Session.Get(r.Context(), "user_settings").(models.UserSettings)

	if !loggedinOk {
		log.Println("Cannot get loggedin from session")
		m.App.ErrorLog.Println("Can't get loggedin from the session")
		m.App.Session.Put(r.Context(), "error", "Can't get loggedin from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	if !settingsOk {
		log.Println("Cannot get user_settings from session")
		m.App.ErrorLog.Println("Can't get user_settings from the session")
		m.App.Session.Put(r.Context(), "error", "Can't get user_settings from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	data := make(map[string]interface{})
	data["settingsform"] = emptyUserSettingsForm
	data["loggedin"] = loggedin
	data["settings"] = settings

	_ = render.Template(w, r, "settings.page.tmpl", &models.TemplateData{Data: data, Form: forms.New(nil)})
}

// @desc        Update user settings
// @route       POST /user/settings
// @access      Private
func (m *Repository) PostSettingsPage(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	id, _ := strconv.Atoi(r.Form.Get("member_id"))

	settings := models.UserSettings{
		MemberID: id,
	}

	var key string

	for key = range r.Form {
		if key == "show_profile" {
			settings.ShowProfile = true
		}

		if key == "show_online_status" {
			settings.ShowOnlineStatus = true
		}

		if key == "show_address" {
			settings.ShowAddress = true
		}

		if key == "show_city" {
			settings.ShowCity = true
		}

		if key == "show_state" {
			settings.ShowState = true
		}

		if key == "show_display_name" {
			settings.ShowDisplayName = true
		}

		if key == "show_contact_info" {
			settings.ShowContactInfo = true
		}

		if key == "show_phone" {
			settings.ShowPhone = true
		}

		if key == "show_email" {
			settings.ShowEmail = true
		}

		if key == "show_craft" {
			settings.ShowCraft = true
		}

		if key == "show_run" {
			settings.ShowRun = true
		}

		if key == "show_notif" {
			settings.ShowNotifications = true
		}
	}

	fmt.Printf("\n\n")

	results := m.DB.UpdateUserSettings(settings)
	resultsErr := results["err"]

	if resultsErr != "" {
		fmt.Println("Update Settings Query Error:\t", resultsErr)
		return
	}

	// Remove the old user data from the session
	m.App.Session.Remove(r.Context(), "user_settings")

	memberId, _ := strconv.Atoi(results["memberId"])
	showProfile, _ := strconv.ParseBool(results["showProfile"])
	showOnlineStatus, _ := strconv.ParseBool(results["showOnlineStatus"])
	showAddress, _ := strconv.ParseBool(results["showAddress"])
	showCity, _ := strconv.ParseBool(results["showCity"])
	showState, _ := strconv.ParseBool(results["showState"])
	showDisplayName, _ := strconv.ParseBool(results["showDisplayName"])
	showContactInfo, _ := strconv.ParseBool(results["showContactInfo"])
	showPhone, _ := strconv.ParseBool(results["showPhone"])
	showEmail, _ := strconv.ParseBool(results["showEmail"])
	showCraft, _ := strconv.ParseBool(results["showCraft"])
	showRun, _ := strconv.ParseBool(results["showRun"])
	showNotifications, _ := strconv.ParseBool(results["showNotifications"])

	postSettings := models.UserSettings{
		MemberID:          memberId,
		ShowProfile:       showProfile,
		ShowOnlineStatus:  showOnlineStatus,
		ShowAddress:       showAddress,
		ShowCity:          showCity,
		ShowState:         showState,
		ShowDisplayName:   showDisplayName,
		ShowContactInfo:   showContactInfo,
		ShowPhone:         showPhone,
		ShowEmail:         showEmail,
		ShowCraft:         showCraft,
		ShowRun:           showRun,
		ShowNotifications: showNotifications,
	}

	m.App.Session.Put(r.Context(), "user_settings", postSettings)

	http.Redirect(w, r, "/user/settings", http.StatusSeeOther)
}

// @desc        User Profile Page
// @route       GET /user/update
// @access      Private
func (m *Repository) ProfilePage(w http.ResponseWriter, r *http.Request) {
	var emptyUserForm models.Member
	data := make(map[string]interface{})
	data["profile"] = emptyUserForm

	fmt.Println("Get Profile Page")

	loggedin, loggedInOk := m.App.Session.Get(r.Context(), "loggedin").(models.Member)
	profile, profileOk := m.App.Session.Get(r.Context(), "user_profile").(models.Profile)

	if !loggedInOk {
		log.Println("Cannot get loggedin session")
		m.App.ErrorLog.Println("Can't get loggedin from the session")
		m.App.Session.Put(r.Context(), "error", "Can't get loggedin from session")
		http.Redirect(w, r, "/user/dashboard", http.StatusTemporaryRedirect)
		return
	}

	if !profileOk {
		log.Println("Cannot get profile session")
		m.App.ErrorLog.Println("Can't get profile from the session")
		m.App.Session.Put(r.Context(), "error", "Can't get profile from session")
		http.Redirect(w, r, "/user/dashboard", http.StatusTemporaryRedirect)
		return
	}

	data["loggedin"] = loggedin
	data["profile"] = profile

	_ = render.Template(w, r, "profile.page.tmpl", &models.TemplateData{
		Data: data,
		Form: forms.New(nil),
	})
}

// @desc        Update user profile
// @route       POST /signin
// @access      private
func (m *Repository) PostProfilePage(w http.ResponseWriter, r *http.Request) {
	// Remove the old user data from the session
	m.App.Session.Remove(r.Context(), "loggedin")
	m.App.Session.Remove(r.Context(), "user_profile")

	fmt.Println("Post to the UpdateTheUserProfile Route")

	err := r.ParseForm()

	if err != nil {
		fmt.Printf("\n\tError parsing user profile form")
		helpers.ServerError(w, err)
		return
	}

	yos, _ := strconv.Atoi(r.Form.Get("yos"))

	updatedMemberInfo := models.Member{
		FirstName: r.Form.Get("fname"),
		LastName:  r.Form.Get("lname"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}

	updatedProfileInfo := models.Profile{
		ImageURL: r.Form.Get("imgurl"),
		UserName: r.Form.Get("uname"),
		Craft:    r.Form.Get("craft"),
		YOS:      yos,
		Address:  r.Form.Get("address"),
		City:     r.Form.Get("city"),
		State:    r.Form.Get("state"),
	}

	// form validation

	form := forms.New(r.PostForm)
	form.MinLength("fname", 2, r)
	form.MinLength("lname", 2, r)
	form.IsEmail("email")
	form.Required("email", "fname", "lname", "uname", "phone", "imgurl", "craft", "address", "city", "state")

	if !form.Valid() {
		fmt.Printf("\n\tForm Error:\t%v\n\n", form.Errors)
		loggedin, ok := m.App.Session.Get(r.Context(), "loggedin").(models.Member)

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
	results := m.DB.UpdateUserProfile(updatedMemberInfo, updatedProfileInfo)

	// fmt.Printf("\n\n\t\tUpdated User Profile Data\n\t\t\t%v\n\n", results)

	if results["memberQueryErr"] != "" {
		return
	}

	if results["memberRowsScanErr"] != "" {
		return
	}

	if results["memberRerr"] != "" {
		return
	}

	if results["memberRowsErr"] != "" {
		return
	}

	if results["profileRowsScanError"] != "" {
		return
	}

	if results["profileRerr"] != "" {
		return
	}

	if results["profileRowsErr"] != "" {
		return
	}

	memberID := results["userID"]
	userId, _ := strconv.Atoi(memberID)
	firstName := results["firstName"]
	lastName := results["lastName"]
	email := results["email"]
	phone := results["phone"]
	userName := results["userName"]
	displayName := results["displayName"]
	imageUrl := results["imageUrl"]
	craft := results["craft"]
	address := results["address"]
	city := results["city"]
	state := results["state"]
	createdAt := results["createdAt"]
	profileUpdatedAt := results["profileUpdatedAt"]
	yearsService := results["yearsService"]

	const layout = "2006-01-02"
	_createdAt, _ := time.Parse(layout, results["createdAt"])
	updatedDate, _ := time.Parse(layout, profileUpdatedAt)

	serviceYears, _ := strconv.Atoi(yearsService)

	loggedIn := models.Member{
		FirstName:    firstName,
		LastName:     lastName,
		ID:           userId,
		Email:        email,
		Phone:        phone,
		CreationDate: createdAt,
		CreatedAt:    _createdAt,
		UpdatedAt:    updatedDate,
		Updated:      profileUpdatedAt,
	}

	profile := models.Profile{
		UserName:    userName,
		DisplayName: displayName,
		ImageURL:    imageUrl,
		Craft:       craft,
		Address:     address,
		City:        city,
		State:       state,
		UpdatedAt:   updatedDate,
		YOS:         serviceYears,
	}

	fmt.Println("\tUpdated User Profile Data")
	fmt.Println("User ID: ", userId)
	fmt.Println("First Name: ", firstName)
	fmt.Println("Last Name: ", lastName)
	fmt.Println("User Name: ", userName)
	fmt.Println("Display Name: ", displayName)
	fmt.Println("Email: ", email)
	fmt.Println("Phone: ", phone)
	fmt.Println("Craft: ", craft)
	fmt.Println("Created: ", createdAt)
	fmt.Println("Updated: ", profileUpdatedAt)
	fmt.Println("Image URL: ", imageUrl)
	fmt.Println("Address: ", address)
	fmt.Println("City: ", city)
	fmt.Println("State: ", state)
	fmt.Println("Service Years: ", serviceYears)
	fmt.Printf("\n\n")

	// Add the update user data to the session
	m.App.Session.Put(r.Context(), "loggedin", loggedIn)
	m.App.Session.Put(r.Context(), "user_profile", profile)

	http.Redirect(w, r, "/user/profile", http.StatusSeeOther)
}

// @desc        AdminPage Dashboard Page
// @route       GET /admin
// @access      Private
func (m *Repository) AdminPage(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})

	fmt.Println("Get Admin Page")

	loggedin, loggedInOk := m.App.Session.Get(r.Context(), "loggedin").(models.Member)
	profile, profileOk := m.App.Session.Get(r.Context(), "user_profile").(models.Profile)
	usersettings, usersettingsOk := m.App.Session.Get(r.Context(), "user_settings").(models.UserSettings)
	allUsers, allUsersOk := m.App.Session.Get(r.Context(), "allusers").(models.Users)

	if !loggedInOk {
		log.Println("Cannot get loggedin session")
		m.App.ErrorLog.Println("Can't get loggedin from the session")
		m.App.Session.Put(r.Context(), "error", "Can't get loggedin from session")
		http.Redirect(w, r, "/admin/dashboard", http.StatusTemporaryRedirect)
		return
	}

	if !profileOk {
		log.Println("Cannot get profile session")
		m.App.ErrorLog.Println("Can't get profile from the session")
		m.App.Session.Put(r.Context(), "error", "Can't get profile from session")
		http.Redirect(w, r, "/admin/dashboard", http.StatusTemporaryRedirect)
		return
	}

	if !usersettingsOk {
		log.Println("Cannot get usersettings session")
		m.App.ErrorLog.Println("Can't get usersettings from the session")
		m.App.Session.Put(r.Context(), "error", "Can't get usersettings from session")
		http.Redirect(w, r, "/admin/dashboard", http.StatusTemporaryRedirect)
		return
	}

	if !allUsersOk {
		log.Println("Cannot get alluser session")
		m.App.ErrorLog.Println("Can't get alluser from the session")
		m.App.Session.Put(r.Context(), "error", "Can't get alluser from session")
		http.Redirect(w, r, "/admin/dashboard", http.StatusTemporaryRedirect)
		return
	}

	if loggedin.AccessLevel != 1 {
		http.Redirect(w, r, "/user/dashboard", http.StatusSeeOther)
	}

	fmt.Println("All users:\t", allUsers.AllUsers)

	data["loggedin"] = loggedin
	data["profile"] = profile
	data["settings"] = usersettings
	data["users"] = allUsers.AllUsers
	_ = render.Template(w, r, "adminuser.page.tmpl", &models.TemplateData{Data: data})

}

// @desc        AdminPage Dashboard Page
// @route       GET /admin/users
// @access      Private
func (m *Repository) UsersPage(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})

	fmt.Println("Get Admin Page")

	loggedin, loggedInOk := m.App.Session.Get(r.Context(), "loggedin").(models.Member)
	profile, profileOk := m.App.Session.Get(r.Context(), "user_profile").(models.Profile)
	usersettings, usersettingsOk := m.App.Session.Get(r.Context(), "user_settings").(models.UserSettings)
	allUsers, allUsersOk := m.App.Session.Get(r.Context(), "allusers").(models.Users)

	if !loggedInOk {
		log.Println("Cannot get loggedin session")
		m.App.ErrorLog.Println("Can't get loggedin from the session")
		m.App.Session.Put(r.Context(), "error", "Can't get loggedin from session")
		http.Redirect(w, r, "/admin/dashboard", http.StatusTemporaryRedirect)
		return
	}

	if !profileOk {
		log.Println("Cannot get profile session")
		m.App.ErrorLog.Println("Can't get profile from the session")
		m.App.Session.Put(r.Context(), "error", "Can't get profile from session")
		http.Redirect(w, r, "/admin/dashboard", http.StatusTemporaryRedirect)
		return
	}

	if !usersettingsOk {
		log.Println("Cannot get usersettings from session")
		m.App.ErrorLog.Println("Can't get usersettings from the session")
		m.App.Session.Put(r.Context(), "error", "Can't get usersettings from session")
		http.Redirect(w, r, "/admin/dashboard", http.StatusTemporaryRedirect)
		return
	}

	if !allUsersOk {
		log.Println("Cannot get alluser session")
		m.App.ErrorLog.Println("Can't get alluser from the session")
		m.App.Session.Put(r.Context(), "error", "Can't get alluser from session")
		http.Redirect(w, r, "/admin/dashboard", http.StatusTemporaryRedirect)
		return
	}

	if loggedin.AccessLevel != 1 {
		http.Redirect(w, r, "/user/dashboard", http.StatusSeeOther)
	} else {
		data["loggedin"] = loggedin
		data["profile"] = profile
		data["settings"] = usersettings
		data["users"] = allUsers.AllUsers
		m.App.Session.Put(r.Context(), "allusers", allUsers)
		_ = render.Template(w, r, "users.page.tmpl", &models.TemplateData{Data: data})
	}
}

// @desc        AdminPage Dashboard Page
// @route       GET /admin/user
// @access      Private
func (m *Repository) UserPage(w http.ResponseWriter, r *http.Request) {
	var emptyMemberForm models.Registration
	data := make(map[string]interface{})
	data["memberform"] = emptyMemberForm

	loggedin, loggedInOk := m.App.Session.Get(r.Context(), "loggedin").(models.Member)
	profile, profileOk := m.App.Session.Get(r.Context(), "user_profile").(models.Profile)
	usersettings, usersettingsOk := m.App.Session.Get(r.Context(), "user_settings").(models.UserSettings)
	allUsers, allUsersOk := m.App.Session.Get(r.Context(), "allusers").(models.Users)

	if !loggedInOk {
		log.Println("Cannot get loggedin session")
		m.App.ErrorLog.Println("Can't get loggedin from the session")
		m.App.Session.Put(r.Context(), "error", "Can't get loggedin from session")
		http.Redirect(w, r, "/admin/dashboard", http.StatusTemporaryRedirect)
		return
	}

	if !profileOk {
		log.Println("Cannot get profile session")
		m.App.ErrorLog.Println("Can't get profile from the session")
		m.App.Session.Put(r.Context(), "error", "Can't get profile from session")
		http.Redirect(w, r, "/admin/dashboard", http.StatusTemporaryRedirect)
		return
	}

	if !usersettingsOk {
		log.Println("Cannot get usersettings from session")
		m.App.ErrorLog.Println("Can't get usersettings from the session")
		m.App.Session.Put(r.Context(), "error", "Can't get usersettings from session")
		http.Redirect(w, r, "/admin/dashboard", http.StatusTemporaryRedirect)
		return
	}

	if !allUsersOk {
		log.Println("Cannot get alluser session")
		m.App.ErrorLog.Println("Can't get alluser from the session")
		m.App.Session.Put(r.Context(), "error", "Can't get alluser from session")
		http.Redirect(w, r, "/admin/dashboard", http.StatusTemporaryRedirect)
		return
	}

	fmt.Println("Get UserPage Page")

	paramUserId := r.URL.Query().Get("code")
	num, _ := strconv.ParseInt(paramUserId, 0, 32)
	var userId int = int(num)

	member, err := m.DB.GetUserByID(userId)

	if err != nil {
		fmt.Println("UserPage handler DB error:\t", err.Error())
		return
	}

	data["member"] = member

	if loggedin.AccessLevel != 1 {
		http.Redirect(w, r, "/user/dashboard", http.StatusSeeOther)
	}

	data["loggedin"] = loggedin
	data["profile"] = profile
	data["settings"] = usersettings
	data["users"] = allUsers.AllUsers
	data["member"] = member

	_ = render.Template(w, r, "user.page.tmpl", &models.TemplateData{
		Data: data,
		Form: forms.New(nil)})
}

// @desc        Update user
// @route       POST /admin/user/update
// @access      private
func (m *Repository) PostUserPage(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	fmt.Println("Posted to the /admin/user/update Route")

	err := r.ParseForm()

	if err != nil {
		fmt.Printf("\n\tError parsing member form")
		helpers.ServerError(w, err)
		return
	}

	memberId, _ := strconv.Atoi(r.Form.Get("id"))
	accessLevel, _ := strconv.Atoi(r.Form.Get("accesslevel"))

	memberForm := models.Member{
		ID:          memberId,
		FirstName:   r.Form.Get("fname"),
		LastName:    r.Form.Get("lname"),
		Email:       r.Form.Get("email"),
		Phone:       r.Form.Get("phone"),
		Password:    r.Form.Get("password"),
		AccessLevel: accessLevel,
	}

	// form validation

	form := forms.New(r.PostForm)
	form.MinLength("fname", 2, r)
	form.MinLength("lname", 2, r)
	form.IsEmail("email")
	form.Required("email", "fname", "lname", "accesslevel", "phone")

	if !form.Valid() {
		fmt.Printf("\n\tForm Error:\t%v\n\n", form.Errors)

		data["member"] = memberForm

		_ = render.Template(w, r, "user.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	updatedMember, err := m.DB.UpdateUser(memberForm)

	if err != nil {
		fmt.Println("Error updating member:\t", err.Error())
		http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
	}

	fmt.Println("Member updated:\t", updatedMember)

	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}

// @desc        Remove user
// @route       GET /admin/user/remove
// @access      private
func (m *Repository) RemoveUser(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}
