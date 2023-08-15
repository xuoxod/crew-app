package dbrepo

import (
	"github.com/xuoxod/crew-app/internal/models"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

func (m *testDBRepo) CreateUser(res models.User) (int, error) {
	return 0, nil
}

func (m *testDBRepo) InsertCraft(r models.Craft) (int, error) {
	return 0, nil
}

// User stuff

func (m *testDBRepo) GetUserByID(id int) (models.User, error) {
	var u models.User
	return u, nil
}

func (m *testDBRepo) GetUserByEmail(email string) (models.User, error) {
	var u models.User

	return u, nil
}

func (m *testDBRepo) UpdateUser(u models.User) error {

	return nil

}

func (m *testDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	var id int
	var hashedPassword string
	return id, hashedPassword, nil
}

func (m *testDBRepo) AuthenticateUser(email, testPassword string) map[string]string {
	results := make(map[string]string)

	return results
}
