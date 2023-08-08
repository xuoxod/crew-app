package datastore

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/xuoxod/crew-app/internal/config"
)

var app *config.AppConfig

func DbConnect() error {

	fmt.Printf("\n\t\tConnection String : %s\n\n", app.DBConnection)

	conn, err := pgx.Connect(context.Background(), app.DBConnection)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		// os.Exit(1)
		return err
	}
	defer conn.Close(context.Background())

	var greeting string
	err = conn.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		// os.Exit(1)

		return err
	}

	fmt.Println(greeting)

	return nil

}
