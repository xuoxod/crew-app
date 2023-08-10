package repository

import "github.com/xuoxod/crew-app/internal/models"

type DatabaseRepo interface {
	AllUsers() bool
	CreateUser(res models.User) (int, error)
	InsertCraft(c models.Craft) (int, error)
}
