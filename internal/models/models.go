package models

import "time"

// User registration data
type Registration struct {
	FirstName       string
	LastName        string
	Email           string
	Phone           string
	PasswordCreate  string
	PasswordConfirm string
}

// User signin data
type Signin struct {
	Email    string
	Password string
}

// Member data
type Member struct {
	ID           int
	FirstName    string
	LastName     string
	Email        string
	Phone        string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Password     string
	CreationDate string
	Updated      string
}

// User profile
type Profile struct {
	MemberID    int
	UserName    string
	ImageURL    string
	Craft       string
	Address     string
	City        string
	State       string
	DisplayName string
	YOS         int
	UpdatedAt   time.Time
}

// Craft data
type UserSettings struct {
	MemberID          int
	ShowProfile       bool
	ShowOnlineStatus  bool
	ShowAddress       bool
	ShowCity          bool
	ShowState         bool
	ShowDisplayName   bool
	ShowContactInfo   bool
	ShowPhone         bool
	ShowEmail         bool
	ShowCraft         bool
	ShowRun           bool
	ShowNotifications bool
}

// User ID
type UserID struct {
	UserID int
}

// Crew member data
type Crew struct {
	RunNumber int
	FirstName string
	LastName  string
	Phone     string
	ID        int
}
