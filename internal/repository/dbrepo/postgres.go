package dbrepo

import (
	"context"
	"time"

	"github.com/xuoxod/crew-app/internal/models"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

func (m *postgresDBRepo) CreateUser(res models.User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	var newID int

	stmt := `insert into krxbyhhs.public.crews(fname, lname, email, phone, password,created_at,updated_at)
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

	stmt := `insert into crafts (fname, lname, userID)
	values($1,$2,$3,$4)`

	_, err := m.DB.ExecContext(ctx, stmt,
		r.Fname,
		r.Lname,
		r.UserID,
	)

	if err != nil {
		return 0, err
	}

	return 0, nil
}
