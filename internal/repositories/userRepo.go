package repositories

import (
	"context"
	"github.com/leavemeal0ne/Golang-2023/internal/driver"
	"github.com/leavemeal0ne/Golang-2023/internal/models"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

type UserRepository struct {
	DB *driver.DB
}

func NewUserRepository(DB *driver.DB) *UserRepository {
	return &UserRepository{DB}
}

func (m *UserRepository) InsertUser(user models.Users) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password_hash), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	stmt := `insert into users(email, password_hash,created_at, updated_at) values ($1, $2, $3, $4)`

	_, err = m.DB.SQL.ExecContext(ctx, stmt,
		user.Email,
		hash,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return err
	}
	return nil
}

func (m *UserRepository) GetUserByEmail(email string) (models.Users, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user models.Users
	user.Email = email

	query := `
				select 
					id, password_hash
				from 
				    users
				where 
					email = $1`

	row := m.DB.SQL.QueryRowContext(ctx, query, email)
	err := row.Scan(
		&user.ID,
		&user.Password_hash,
	)
	if err != nil {
		log.Println(err)
		return models.Users{}, err
	}

	return user, nil
}

func (m *UserRepository) IsContainsUserByEmail(email string) bool {
	_, err := m.GetUserByEmail(email)
	if err != nil {
		return false
	}
	return true
}

func (m *UserRepository) GetUserByEmailAndPassword(email, password string) (models.Users, error) {
	user, err := m.GetUserByEmail(email)
	if err != nil {
		return models.Users{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password_hash), []byte(password))
	if err != nil {
		return models.Users{}, err
	}
	return user, nil
}
