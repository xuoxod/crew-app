package repository

import "github.com/xuoxod/crew-app/internal/models"

type DatabaseRepo interface {
	AllUsers() bool
	CreateUser(res models.Member) (map[string]string, error)
	GetUserByID(id int) (models.Member, error)
	GetUserByEmail(email string) (models.Member, error)
	UpdateUserProfile(u models.Member, p models.Profile) map[string]string
	AuthenticateUser(email, testPassword string) map[string]string
	UpdateUserSettings(u models.UserSettings) map[string]string
}
