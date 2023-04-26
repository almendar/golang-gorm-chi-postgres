package main

import (
	"fmt"
	"os"
	"time"

	"github.com/almendar/golang-gorm-chi-postgres/dogowners"
	"github.com/almendar/golang-gorm-chi-postgres/shared"
)

var MIGRATE = os.Getenv("MIGRATE") == "true"

func main() {
	// Migrate the schema if set so
	db, err := shared.DefaultGormHandle()
	if err != nil {
		fmt.Printf("failed to connect to database: %s\n", err)
		return
	}

	dbStorage := dogowners.NewDatabase(db)

	if !MIGRATE {
		fmt.Println("MIGRATE is not set to true, skipping migration")
		return
	}

	if err := dogowners.RunMigration(db); err != nil {
		fmt.Printf("failed to migrate database: %s\n", err)
	}

	email := "j.doe@email.com"
	owner := dogowners.OwnerDBModel{
		Name:     "John Doe",
		Email:    &email,
		Age:      42,
		Birthday: time.Now().Add(-42 * 365 * 24 * time.Hour),
		Dogs: []dogowners.DogDBModel{
			{
				Name:     "Fido",
				Birthday: time.Now().Add(-10 * 365 * 24 * time.Hour),
			},
			{
				Name:     "Rex",
				Birthday: time.Now().Add(-5 * 365 * 24 * time.Hour),
			},
		},
	}
	if err := dbStorage.SaveOwner(&owner); err != nil {
		fmt.Printf("failed to create user: %s\n", err)
		return
	}
}
