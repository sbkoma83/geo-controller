package service

import (
	"context"
	"errors"
	"geo-controller/proxy/internal/models"
	"geo-controller/proxy/internal/repositories"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type AuthService struct {
	usersRepo repositories.UserRepository
}

func NewAuthService(uRepo repositories.UserRepository) *AuthService {
	return &AuthService{
		usersRepo: uRepo,
	}
}
func (s *AuthService) RegisterUser(username, password string) error {
	user, _ := s.usersRepo.GetByUsername(context.Background(), username)
	if user.Username != "" {
		return errors.New("username exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u := models.User{Username: username, Password: string(hashedPassword)}

	err = s.usersRepo.Create(context.Background(), u)
	if err != nil {
		return err
	}
	return nil
}

func (s *AuthService) AuthenticateUser(username, password string) bool {
	user, err := s.usersRepo.GetByUsername(context.Background(), username)
	if err != nil {
		log.Println(err)
		return false
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}
func (s *AuthService) ListUsers() (*[]models.User, error) {
	users, err := s.usersRepo.List(context.Background())
	if err != nil {
		return nil, err
	}
	return &users, nil
}
func (s *AuthService) GetByID(id int) (models.User, error) {
	user, err := s.usersRepo.GetByID(context.Background(), id)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
func (s *AuthService) UpdateUser(user models.User) error {
	err := s.usersRepo.Update(context.Background(), user)
	if err != nil {
		return err
	}
	return nil
}
func (s *AuthService) DeleteByID(id int) error {
	err := s.usersRepo.Delete(context.Background(), id)
	if err != nil {
		return err
	}
	return nil
}
