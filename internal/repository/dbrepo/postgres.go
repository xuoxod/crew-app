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

func (m *postgresDBRepo) CreateUserProfile(u models.User) map[string]string {
	var results = make(map[string]string)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var fname, lname, email, phone, uname, iurl, updatedAt, createdAt string
	var memberId, craftId int

	query := `
		update members m set first_name = $1, last_name = $2, email = $3, phone = $4, updated_at = $5 where m.id = $6 returning m.id, m.first_name, m.last_name, m.email, m.phone, m.craft_id, updated_at, created_at
	`

	rows, err := m.DB.QueryContext(ctx, query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.Phone,
		time.Now(),
		u.ID,
	)

	if err != nil {
		fmt.Printf("\tQuery Update Error: %s\n", err.Error())
		results["err"] = err.Error()
		return results
	}

	for rows.Next() {
		if err := rows.Scan(&memberId, &fname, &lname, &email, &phone, &craftId, &updatedAt, &createdAt); err != nil {
			fmt.Printf("\tRow Scan Error: %s\n", err.Error())
			results["err"] = err.Error()
			return results
		}
	}

	rerr := rows.Close()

	if rerr != nil {
		fmt.Printf("rerr error:\t%s\n", rerr.Error())
		results["err"] = rerr.Error()
		return results
	}

	if err := rows.Err(); err != nil {
		fmt.Printf("Row Error:\t%s\n", err.Error())
		results["err"] = err.Error()
		return results
	}

	results["userID"] = fmt.Sprintf("%d", memberId)
	results["firstName"] = fname
	results["lastName"] = lname
	results["email"] = email
	results["phone"] = phone
	results["craftID"] = fmt.Sprintf("%d", craftId)
	results["createdAt"] = fmt.Sprintf("%v", createdAt)
	results["updatedAt"] = fmt.Sprintf("%v", updatedAt)

	// ------------------------------------------------------------------------

	query = `
		update profiles p set user_name = $1, image_url = $2 where p.member_id = $3 returning p.user_name, p.image_url
	`

	rows, err = m.DB.QueryContext(ctx, query,
		u.Username,
		u.ImageURL,
		memberId,
	)

	if err != nil {
		results["err"] = err.Error()
		return results
	}

	for rows.Next() {
		if err := rows.Scan(&uname, &iurl); err != nil {
			fmt.Printf("\t2nd Row Scan Error: %s\n", err.Error())
			results["err"] = err.Error()
			return results
		}
	}

	rerr = rows.Close()

	if rerr != nil {
		fmt.Printf("2nd rerr error:\t%s\n", rerr.Error())
		results["err"] = rerr.Error()
		return results
	}

	if err := rows.Err(); err != nil {
		fmt.Printf("2nd Row Error:\t%s\n", err.Error())
		results["err"] = err.Error()
		return results
	}

	results["userName"] = uname
	results["imageUrl"] = iurl
	results["err"] = ""

	return results
}

func (m *postgresDBRepo) UpdateUserProfile(u models.User) map[string]string {
	var results = make(map[string]string)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var fname, lname, email, phone, uname, iurl, updatedAt, createdAt string
	var memberId, craftId int

	query := `
		update members m set first_name = $1, last_name = $2, email = $3, phone = $4, updated_at = $5 where m.email = $6 returning m.id, m.first_name, m.last_name, m.email, m.phone, m.craft_id, updated_at, created_at
	`

	rows, err := m.DB.QueryContext(ctx, query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.Phone,
		time.Now(),
		u.Email,
	)

	if err != nil {
		fmt.Printf("\tQuery Error: %s\n", err.Error())
		results["err"] = err.Error()
		return results
	}

	for rows.Next() {
		if err := rows.Scan(&memberId, &fname, &lname, &email, &phone, &craftId, &updatedAt, &createdAt); err != nil {
			fmt.Printf("\tRow Scan Error: %s\n", err.Error())
			results["err"] = err.Error()
			return results
		}
	}

	rerr := rows.Close()

	if rerr != nil {
		fmt.Printf("rerr error:\t%s\n", rerr.Error())
		results["err"] = rerr.Error()
		return results
	}

	if err := rows.Err(); err != nil {
		fmt.Printf("Row Error:\t%s\n", err.Error())
		results["err"] = err.Error()
		return results
	}

	results["userID"] = fmt.Sprintf("%d", memberId)
	results["firstName"] = fname
	results["lastName"] = lname
	results["email"] = email
	results["phone"] = phone
	results["craftID"] = fmt.Sprintf("%d", craftId)
	results["createdAt"] = fmt.Sprintf("%v", createdAt)
	results["updatedAt"] = fmt.Sprintf("%v", updatedAt)

	// ------------------------------------------------------------------------

	query = `
		update profiles p set user_name = $1, image_url = $2 where p.member_id = $3 returning p.user_name, p.image_url
	`

	rows, err = m.DB.QueryContext(ctx, query,
		u.Username,
		u.ImageURL,
		memberId,
	)

	if err != nil {
		results["err"] = err.Error()
		return results
	}

	for rows.Next() {
		if err := rows.Scan(&uname, &iurl); err != nil {
			fmt.Printf("\t2nd Row Scan Error: %s\n", err.Error())
			results["err"] = err.Error()
			return results
		}
	}

	rerr = rows.Close()

	if rerr != nil {
		fmt.Printf("2nd rerr error:\t%s\n", rerr.Error())
		results["err"] = rerr.Error()
		return results
	}

	if err := rows.Err(); err != nil {
		fmt.Printf("2nd Row Error:\t%s\n", err.Error())
		results["err"] = err.Error()
		return results
	}

	results["userName"] = uname
	results["imageUrl"] = iurl
	results["err"] = ""
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
	var accessLevel, id, craftId int
	var firstName, lastName, emailAddress, phone, hashedPassword, userName, imgUrl string
	var createdAt, updatedAt time.Time

	row := m.DB.QueryRowContext(ctx, "select first_name, last_name, email, phone, access_level, craft_id, created_at, updated_at, password, user_name, image_url from members m left join profiles p on p.member_id = p.member_id where p.member_id = m.id and email = $1", email)

	err := row.Scan(&firstName, &lastName, &emailAddress, &phone, &accessLevel, &craftId, &createdAt, &updatedAt, &hashedPassword, &userName, &imgUrl)

	if err != nil {
		log.Printf("\n\tScan error:\n\t%s\n", err.Error())
		log.Printf("\n\tRunning query for members only\n\n")

		row = m.DB.QueryRowContext(ctx, "select id, first_name, last_name, email, phone, access_level,craft_id, created_at, updated_at, password from members where email = $1", email)

		err = row.Scan(&id, &firstName, &lastName, &emailAddress, &phone, &accessLevel, &craftId, &createdAt, &updatedAt, &hashedPassword)

		if err != nil {
			log.Printf("\n\tScan error:\n\t%s\n\n", err.Error())
			results["err"] = err.Error()
			return results
		}

		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))

		if err == bcrypt.ErrMismatchedHashAndPassword {
			log.Println("Authentication failed")
			results["err"] = err.Error()
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
		results["craftID"] = fmt.Sprintf("%d", craftId)
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
	results["craftID"] = fmt.Sprintf("%d", craftId)
	results["createdAt"] = createdAt.String()
	results["updatedAt"] = updatedAt.String()
	results["userName"] = userName
	results["imgUrl"] = imgUrl
	results["profileStatus"] = "hasprofile"
	results["err"] = ""
	return results
}
