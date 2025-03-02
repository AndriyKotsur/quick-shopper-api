package database

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Seed struct {
	DB *sql.DB
}

type User struct {
	ID       string
	Name     string
	Email    string
	Password string
	Role     string
}

func Seeder(db *sql.DB) *Seed {
	return &Seed{
		DB: db,
	}
}

func (s *Seed) SeedUsers() {
	users := []User{
		{
			Name:     "Drew",
			Email:    "admin@booking.com",
			Password: "12345678",
			Role:     "ADMIN",
		},
		{
			Name:     "Peter",
			Email:    "manager@booking.com",
			Password: "12345678",
			Role:     "MANAGER",
		},
		{
			Name:     "Monica",
			Email:    "customer@booking.com",
			Password: "12345678",
			Role:     "CUSTOMER",
		},
	}

	for _, u := range users {
		password, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalln(err)
		}

		_, err = s.DB.ExecContext(
			context.Background(),
			`INSERT INTO users (id, name, email, role, password, created_at, updated_at) 
				VALUES ($1, $2, $3, $4, $5, $6, $7);`,
			uuid.New(),
			u.Name,
			u.Email,
			u.Role,
			password,
			time.Now(),
			time.Now(),
		)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
