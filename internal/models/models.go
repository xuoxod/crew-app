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

// User data
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

// Craft data
type Craft struct {
	ID      int
	Title   string
	CraftID int
}

// Run data
type Run struct {
	RunNumber     int
	StartLocation string
	StartTime     time.Time
	EndTime       time.Time
	LeaveTime     time.Time
}

// Crew member data
type Crew struct {
	RunNumber int
	FirstName string
	LastName  string
	Phone     string
	ID        int
}
