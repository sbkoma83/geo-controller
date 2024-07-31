package service

import (
	"testing"
)

func TestAuthService_RegisterAndAuthenticate(t *testing.T) {
	authService := NewAuthService()

	// Test user registration
	username := "testuser"
	password := "testpassword"
	err := authService.RegisterUser(username, password)
	if err != nil {
		t.Errorf("Failed to register user: %v", err)
	}

	// Test successful authentication
	if !authService.AuthenticateUser(username, password) {
		t.Error("Authentication failed for valid credentials")
	}

	// Test failed authentication
	if authService.AuthenticateUser(username, "wrongpassword") {
		t.Error("Authentication succeeded for invalid credentials")
	}

	// Test registering duplicate user
	err = authService.RegisterUser(username, "anotherpassword")
	if err == nil {
		t.Error("Expected error when registering duplicate user, but got nil")
	}
}

func TestAuthService_AuthenticateUser_EmptyFields(t *testing.T) {
	authService := NewAuthService()

	// Регистрация пользователя для теста
	username := "testuser"
	password := "testpassword"
	err := authService.RegisterUser(username, password)
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	// Попытка аутентификации с пустым именем пользователя
	if authService.AuthenticateUser("", password) {
		t.Error("Expected authentication to fail with empty username, but it succeeded")
	}

	// Попытка аутентификации с пустым паролем
	if authService.AuthenticateUser(username, "") {
		t.Error("Expected authentication to fail with empty password, but it succeeded")
	}
}
