package dbrepo

import (
	"context"
	"time"

	"github.com/xuoxod/crew-app/internal/models"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

func (m *postgresDBRepo) CreateUser(res models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	stmt := `insert into krxbyhhs.public.crews(fname, lname, email, phone, password)
	values($1,$2,$3,$4,$5)`

	_, err := m.DB.ExecContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.Password,
	)

	if err != nil {
		return err
	}

	return nil
}
