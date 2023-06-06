package services

import (
	"errors"
	"github.com/jackc/pgx/v4"
	"github.com/leavemeal0ne/Golang-2023/internal/models"
)

var ErrNotFound = errors.New("order is not found")

type UserRepository interface {
	InsertUser(user models.Users) error
	GetUserByEmail(email string) (models.Users, error)
	IsContainsUserByEmail(email string) bool
	GetUserByEmailAndPassword(email, password string) (models.Users, error)
}

type UserService struct {
	userRepository UserRepository
}

func NewUserService(userRepository UserRepository) *UserService {
	return &UserService{userRepository}
}

func (s *UserService) InsertUser(user models.Users) error {
	err := s.userRepository.InsertUser(user)
	return err
}

func (s *UserService) GetUserByEmail(email string) (models.Users, error) {
	user, err := s.GetUserByEmail(email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return models.Users{}, ErrNotFound
		}
		return models.Users{}, err
	}
	return user, nil
}

func (s *UserService) IsContainsUserByEmail(email string) bool {
	return s.userRepository.IsContainsUserByEmail(email)
}
func (s *UserService) GetUserByEmailAndPassword(email, password string) (models.Users, error) {
	user, err := s.GetUserByEmailAndPassword(email, password)
	if err != nil {
		if err == pgx.ErrNoRows {
			return models.Users{}, ErrNotFound
		}
		return models.Users{}, err
	}
	return user, nil
}
