package service


import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	users map[string]string //map[username]passwordHash
}

func NewAuthService() *AuthService {
	return &AuthService{
		users: make(map[string]string),
	}
}
func (s *AuthService) RegisterUser(username, password string) error {
	if _, ok := s.users[username]; ok {
		return errors.New("username already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	s.users[username] = string(hashedPassword)
	return nil
}

func (s *AuthService) AuthenticateUser(username, password string) bool {
	hashedPassword, ok := s.users[username]
	if !ok {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

