package repository

import "github.com/xuoxod/crew-app/internal/models"

type DatabaseRepo interface {
	AllUsers() bool
	CreateUser(res models.User) (int, error)
	GetUserByID(id int) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	UpdateUser(u models.User) error
	Authenticate(email, testPassword string) (int, string, error)
	AuthenticateUser(email, testPassword string) map[string]string
	InsertCraft(c models.Craft) (int, error)
}
