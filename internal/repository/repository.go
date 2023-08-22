package repository

import "github.com/xuoxod/crew-app/internal/models"

type DatabaseRepo interface {
	AllUsers() bool
	CreateUser(res models.Member) (int, error)
	GetUserByID(id int) (models.Member, error)
	GetUserByEmail(email string) (models.Member, error)
	UpdateUserProfile(u models.Member, p models.Profile) map[string]string
	AuthenticateUser(email, testPassword string) map[string]string
}
