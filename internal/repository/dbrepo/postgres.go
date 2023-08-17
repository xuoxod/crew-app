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

	stmt := `insert into krxbyhhs.public.members(first_name, last_name, email, phone, password, created_at, updated_at)
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

	return newID, nil
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

func (m *postgresDBRepo) CreateUserProfile(u models.Profile) map[string]string {

	return nil
}

func (m *postgresDBRepo) UpdateUserProfile(u models.User) map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()
	var results = make(map[string]string)

	query := `
		update users set first_name = $1, last_name =$2, email = $3, phone = $4, updated_at = $5
	`

	result, err := m.DB.ExecContext(ctx, query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.Phone,
		time.Now(),
	)

	if err != nil {
		results["err"] = "Update query failed"
		return results
	}

	results["err"] = ""

	rows, err := result.RowsAffected()

	if err != nil {
		results["userUpdateErr"] = "User update failed"
	}

	var userStatus string

	if rows == 1 {
		userStatus = "success"
	} else {
		userStatus = "failed"
	}

	results["userStatus"] = userStatus

	// query = `
	// update profiles set user_name = $1, image_url = $2 where profile_id = $3`

	return results

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
	// var userId, accessLevel int
	var accessLevel, id int
	var firstName, lastName, emailAddress, phone, hashedPassword, userName, imgUrl string
	var createdAt, updatedAt time.Time

	/* row := m.DB.QueryRowContext(ctx, "select id, first_name, last_name, email, phone, access_level, created_at, updated_at, password from members where email = $1", email)

	err := row.Scan(&userId, &firstName, &lastName, &emailAddress, &phone, &accessLevel, &createdAt, &updatedAt, &hashedPassword) */

	/* row := m.DB.QueryRowContext(ctx, "select first_name, last_name, email, phone, access_level, created_at, updated_at, password, user_name, image_url from members m left join profiles p ON p.member_id = p.member_id  where p.member_id = m.id and email = $1", email) */

	row := m.DB.QueryRowContext(ctx, "select first_name, last_name, email, phone, access_level, created_at, updated_at, password, user_name, image_url from members m left join profiles p on p.member_id = p.member_id where p.member_id = m.id and email = $1", email)

	err := row.Scan(&firstName, &lastName, &emailAddress, &phone, &accessLevel, &createdAt, &updatedAt, &hashedPassword, &userName, &imgUrl)

	if err != nil {
		log.Printf("\n\tScan error:\n\t%s\n", err.Error())
		log.Printf("\n\tRunning query for members only\n\n")

		row = m.DB.QueryRowContext(ctx, "select id, first_name, last_name, email, phone, access_level, created_at, updated_at, password from members where email = $1", email)

		err = row.Scan(&id, &firstName, &lastName, &emailAddress, &phone, &accessLevel, &createdAt, &updatedAt, &hashedPassword)

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

		results["userID"] = fmt.Sprintf("%d", id)
		results["firstName"] = firstName
		results["lastName"] = lastName
		results["email"] = email
		results["phone"] = phone
		results["accessLevel"] = fmt.Sprintf("%d", accessLevel)
		results["createdAt"] = createdAt.String()
		results["updatedAt"] = updatedAt.String()
		results["err"] = ""
		results["profileStatus"] = "noprofile"
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

	// results["userID"] = fmt.Sprintf("%d", userId)
	results["firstName"] = firstName
	results["lastName"] = lastName
	results["email"] = email
	results["phone"] = phone
	results["accessLevel"] = fmt.Sprintf("%d", accessLevel)
	results["createdAt"] = createdAt.String()
	results["updatedAt"] = updatedAt.String()
	results["userName"] = userName
	results["imgUrl"] = imgUrl
	results["profileStatus"] = "hasprofile"
	results["err"] = ""
	return results
}
