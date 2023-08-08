package models

import "time"

// Registration data
type Registration struct {
	FirstName       string
	LastName        string
	Email           string
	Phone           string
	PasswordCreate  string
	PasswordConfirm string
}

// Signin data
type Signin struct {
	Email    string
	Password string
}

type User struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Phone       string
	Password    string
	AccessLevel int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Craft struct {
	ID      int
	Fname   string
	Lname   string
	Title   string
	CraftID int
}
