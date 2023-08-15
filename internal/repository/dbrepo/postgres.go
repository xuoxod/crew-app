package dbrepo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/xuoxod/crew-app/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

func (m *postgresDBRepo) CreateUser(res models.User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	var newID int

	stmt := `insert into krxbyhhs.public.users(first_name, last_name, email, phone, password,created_at,updated_at)
	values($1,$2,$3,$4,$5,$6,$7) returning id`

	err := m.DB.QueryRowContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.Password,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return 0, nil
}

func (m *postgresDBRepo) InsertCraft(r models.Craft) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	stmt := `insert into krxbyhhs.public.crafts (title, craft_id, userID)
	values($1,$2,$3,$4)`

	_, err := m.DB.ExecContext(ctx, stmt,
		r.Title,
		r.CraftID,
	)

	if err != nil {
		return 0, err
	}

	return 0, nil
}

// User stuff

func (m *postgresDBRepo) GetUserByID(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	query := `select id, first_name, last_name, email, phone, password, access_level, created_at, updated_at, craft_id from users where id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var u models.User
	err := row.Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Phone,
		&u.Password,
		&u.AccessLevel,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		return u, err
	}

	return u, nil
}

func (m *postgresDBRepo) GetUserByEmail(email string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	query := `select id, first_name, last_name, email, phone, password, access_level, created_at, updated_at, craft_id from users where email = $1`

	row := m.DB.QueryRowContext(ctx, query, email)

	var u models.User
	err := row.Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Password,
		&u.AccessLevel,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		return u, err
	}

	return u, nil
}

func (m *postgresDBRepo) UpdateUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	query := `
		update users set first_name = $1, last_name =$2, email = $3, access_level = $4, updated_at = $5
	`

	_, err := m.DB.ExecContext(ctx, query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.AccessLevel,
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil

}

func (m *postgresDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string

	row := m.DB.QueryRowContext(ctx, "select id, password from users where email = $1", email)

	err := row.Scan(&id, &hashedPassword)

	if err != nil {
		return id, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))

	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("authentication failed")
	} else if err != nil {
		return 0, "", err
	}

	return id, hashedPassword, nil
}

func (m *postgresDBRepo) AuthenticateUser(email, testPassword string) map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var results = make(map[string]string)
	var userId int
	var firstName, lastName, _email, phone, hashedPassword, userName, imageUrl string
	var createdAt, updatedAt time.Time

	row := m.DB.QueryRowContext(ctx, "select user_id, first_name, last_name, email, phone, created_at, updated_at, password, user_name, image_url from users u inner join profiles p on p.profile_id = u.profile_id where u.profile_id = p.profile_id and email = $1", email)

	err := row.Scan(&userId, &firstName, &lastName, &_email, &phone, &createdAt, &updatedAt, &hashedPassword, &userName, &imageUrl)

	if err != nil {
		log.Printf("\n\tScan error:\n\t%s\n\n", err.Error())
		results["err"] = err.Error()
		return results
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))

	if err == bcrypt.ErrMismatchedHashAndPassword {
		log.Println("Authentication failed")
		results["err"] = ""
		return results
	} else if err != nil {
		log.Printf("\nError:\n\t%s\n\n", err.Error())
		results["err"] = err.Error()
		return results
	}

	results["id"] = fmt.Sprintf("%d", userId)
	results["firstName"] = firstName
	results["lastName"] = lastName
	results["userName"] = userName
	results["email"] = email
	results["phone"] = phone
	results["imageUrl"] = imageUrl
	results["createdAt"] = createdAt.String()
	results["updatedAt"] = updatedAt.String()
	results["err"] = ""

	return results
}
