package repositories

import (
	"context"
	"geo-controller/proxy/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user models.User) error
	GetByUsername(ctx context.Context, username string) (models.User, error)
	GetByID(ctx context.Context, id uint32) (models.User, error)
	Update(ctx context.Context, user models.User) error
	Delete(ctx context.Context, id uint32) error
	List(ctx context.Context) ([]models.User, error)
}
type UserRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, user models.User) error {
	result := r.db.Create(&user)
	return result.Error
}

func (r *UserRepo) GetByUsername(ctx context.Context, username string) (models.User, error) {
	var user models.User
	result := r.db.Where("username = ?", username).First(&user)
	return user, result.Error
}

func (r *UserRepo) GetByID(ctx context.Context, id uint32) (models.User, error) {
	var user models.User
	result := r.db.Where("id = ?", id).First(&user)
	return user, result.Error
}

func (r *UserRepo) Update(ctx context.Context, user models.User) error {
	result := r.db.Save(&user)
	return result.Error
}

func (r *UserRepo) Delete(ctx context.Context, id uint32) error {
	result := r.db.Delete(&models.User{}, id)
	return result.Error
}

func (r *UserRepo) List(ctx context.Context) ([]models.User, error) {
	var users []models.User
	result := r.db.Find(&users)
	return users, result.Error
}
