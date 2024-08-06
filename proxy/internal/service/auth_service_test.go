package service_test

import (
	"context"
	"geo-controller/proxy/internal/models"
	"geo-controller/proxy/internal/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

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

func TestAuthService_RegisterUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := service.NewAuthService(mockRepo)

	mockRepo.On("GetByUsername", mock.Anything, "testuser").Return(models.User{}, nil)
	mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	err := svc.RegisterUser("testuser", "testpassword")
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_AuthenticateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := service.NewAuthService(mockRepo)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testpassword"), bcrypt.DefaultCost)
	mockUser := models.User{
		Username: "testuser",
		Password: string(hashedPassword),
	}

	mockRepo.On("GetByUsername", mock.Anything, "testuser").Return(mockUser, nil)

	authenticated := svc.AuthenticateUser("testuser", "testpassword")
	assert.True(t, authenticated)
	mockRepo.AssertExpectations(t)
}
func TestAuthService_GetByID(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := service.NewAuthService(mockRepo)

	mockUser := models.User{
		ID:       1,
		Username: "testuser",
		Password: "hashedpassword",
	}

	mockRepo.On("GetByID", mock.Anything, uint32(1)).Return(mockUser, nil)

	user, err := svc.GetByID(1)
	assert.NoError(t, err)
	assert.Equal(t, mockUser, user)
	mockRepo.AssertExpectations(t)
}
func TestAuthService_ListUsers(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := service.NewAuthService(mockRepo)

	mockUsers := []models.User{
		{ID: 1, Username: "user1", Password: "hashedpassword1"},
		{ID: 2, Username: "user2", Password: "hashedpassword2"},
		{ID: 3, Username: "user3", Password: "hashedpassword3"},
	}

	mockRepo.On("List", mock.Anything).Return(mockUsers, nil)

	users, err := svc.ListUsers()
	assert.NoError(t, err)
	assert.Equal(t, &mockUsers, users)
	assert.Len(t, *users, 3)
	mockRepo.AssertExpectations(t)
}
func TestAuthService_UpdateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := service.NewAuthService(mockRepo)

	mockUser := models.User{
		ID:       1,
		Username: "testuser",
		Password: "newhashedpassword",
	}

	mockRepo.On("Update", mock.Anything, mockUser).Return(nil)

	err := svc.UpdateUser(mockUser)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
func TestAuthService_DeleteByID(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := service.NewAuthService(mockRepo)

	mockRepo.On("Delete", mock.Anything, uint32(1)).Return(nil)

	err := svc.DeleteByID(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
