package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"geo-controller/proxy/internal/models"
	"geo-controller/proxy/internal/service"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByUsername(ctx context.Context, username string) (models.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uint32) (models.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uint32) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) List(ctx context.Context) ([]models.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.User), args.Error(1)
}

func newMockAuthController() (*AuthController, *MockUserRepository) {
	mockRepo := new(MockUserRepository)
	authService := service.NewAuthService(mockRepo)
	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)
	authController := NewAuthController(authService, tokenAuth)
	return authController, mockRepo
}

// HashPassword hashes a password using bcrypt
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func TestAuthController_RegisterHandler(t *testing.T) {
	authController, mockRepo := newMockAuthController()

	// Mock repository behavior
	mockRepo.On("GetByUsername", mock.Anything, "testuser").Return(models.User{}, nil)
	mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	// Create request with JSON body
	user := models.User{Username: "testuser", Password: "testpassword"}
	reqBody, _ := json.Marshal(user)
	req, err := http.NewRequest("POST", "/api/register", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(authController.RegisterHandler)

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	mockRepo.AssertExpectations(t)
}

func TestAuthController_LoginHandler(t *testing.T) {
	authController, mockRepo := newMockAuthController()

	// Mock repository behavior
	hashedPassword, _ := hashPassword("testpassword")
	mockRepo.On("GetByUsername", mock.Anything, "testuser").Return(models.User{Username: "testuser", Password: hashedPassword}, nil)

	// Create request with JSON body for login
	user := models.User{Username: "testuser", Password: "testpassword"}
	reqBody, _ := json.Marshal(user)
	req, err := http.NewRequest("POST", "/api/login", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(authController.LoginHandler)

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	var response map[string]string
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}

	if _, exists := response["token"]; !exists {
		t.Error("Expected token in response, got none")
	}

	mockRepo.AssertExpectations(t)
}

func TestAuthController_LoginHandler_InvalidUser(t *testing.T) {
	authController, mockRepo := newMockAuthController()

	// Mock repository behavior
	mockRepo.On("GetByUsername", mock.Anything, "nonexistentuser").Return(models.User{}, errors.New("user not found"))

	// Create request with JSON body for login of non-existent user
	user := models.User{Username: "nonexistentuser", Password: "wrongpassword"}
	reqBody, _ := json.Marshal(user)
	req, err := http.NewRequest("POST", "/api/login", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(authController.LoginHandler)

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

	// Check the response body
	expectedError := "Error encoding token"
	actualError := strings.TrimSpace(rr.Body.String())
	if actualError != expectedError {
		t.Errorf("handler returned unexpected body: got %v want %v", actualError, expectedError)
	}

	mockRepo.AssertExpectations(t)
}

func TestAuthController_RegisterHandler_DuplicateUser(t *testing.T) {
	authController, mockRepo := newMockAuthController()

	// Mock repository behavior
	mockRepo.On("GetByUsername", mock.Anything, "testuser").Return(models.User{Username: "testuser"}, nil)

	// Create request with JSON body for registering the same user
	user := models.User{Username: "testuser", Password: "testpassword"}
	reqBody, _ := json.Marshal(user)
	req, err := http.NewRequest("POST", "/api/register", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(authController.RegisterHandler)

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	// Check the response body
	expectedError := "username exists"
	actualError := strings.TrimSpace(rr.Body.String())
	if actualError != expectedError {
		t.Errorf("handler returned unexpected body: got %v want %v", actualError, expectedError)
	}

	mockRepo.AssertExpectations(t)
}
func TestAuthController_ListHandler(t *testing.T) {
	authController, mockRepo := newMockAuthController()

	// Mock repository behavior
	mockUsers := []models.User{
		{ID: 1, Username: "user1"},
		{ID: 2, Username: "user2"},
	}
	mockRepo.On("List", mock.Anything).Return(mockUsers, nil)

	// Create request
	req, err := http.NewRequest("GET", "/api/users/list", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(authController.ListHandler)

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	var response []models.User
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}

	if len(response) != 2 {
		t.Errorf("handler returned unexpected number of users: got %v want %v", len(response), 2)
	}

	mockRepo.AssertExpectations(t)
}
func TestAuthController_GetByID(t *testing.T) {
	authController, mockRepo := newMockAuthController()

	// Mock repository behavior
	mockUser := models.User{ID: 1, Username: "testuser"}
	mockRepo.On("GetByID", mock.Anything, uint32(1)).Return(mockUser, nil)

	// Create request
	req, err := http.NewRequest("GET", "/api/users/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Create a new router and register the GetByID handler
	r := chi.NewRouter()
	r.Get("/api/users/{id}", authController.GetByID)

	// Serve the request
	r.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	var response models.User
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}

	if response.ID != mockUser.ID || response.Username != mockUser.Username {
		t.Errorf("handler returned unexpected body: got %v want %v", response, mockUser)
	}

	mockRepo.AssertExpectations(t)
}
func TestAuthController_UpdateByID(t *testing.T) {
	authController, mockRepo := newMockAuthController()

	// Mock repository behavior
	updatedUser := models.User{ID: 1, Username: "updateduser", Password: "newpassword"}
	mockRepo.On("Update", mock.Anything, mock.AnythingOfType("models.User")).Return(nil)

	// Create request with JSON body
	reqBody, _ := json.Marshal(updatedUser)
	req, err := http.NewRequest("POST", "/api/users/update/1", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Create a new router and register the UpdateByID handler
	r := chi.NewRouter()
	r.Post("/api/users/update/{id}", authController.UpdateByID)

	// Serve the request
	r.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	expectedBody := "User updated"
	if rr.Body.String() != expectedBody {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expectedBody)
	}

	mockRepo.AssertExpectations(t)
}
