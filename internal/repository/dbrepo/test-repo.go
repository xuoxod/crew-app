package dbrepo

import "github.com/xuoxod/crew-app/internal/models"

func (m *testDBRepo) AllUsers() map[string][]string {
	var results = make(map[string][]string)

	return results
}

func (m *testDBRepo) CreateUser(res models.Member) (map[string]string, error) {
	results := make(map[string]string)
	return results, nil
}

// User stuff

func (m *testDBRepo) GetUserByID(id int) (models.Member, error) {
	var u models.Member
	return u, nil
}

func (m *testDBRepo) GetUserByEmail(email string) (models.Member, error) {
	var u models.Member

	return u, nil
}

func (m *testDBRepo) AuthenticateUser(email, testPassword string) map[string]string {
	results := make(map[string]string)

	return results
}

func (m *testDBRepo) UpdateUserProfile(u models.Member, p models.Profile) map[string]string {

	results := make(map[string]string)

	return results
}

func (m *testDBRepo) UpdateUserSettings(u models.UserSettings) map[string]string {

	results := make(map[string]string)

	return results
}
