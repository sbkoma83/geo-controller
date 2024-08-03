package repositories

import (
	"context"
	"database/sql"
	"errors"
	"geo-controller/proxy/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user models.User) error
	GetByUsername(ctx context.Context, username string) (models.User, error)
	GetByID(ctx context.Context, id int) (models.User, error)
	Update(ctx context.Context, user models.User) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context) ([]models.User, error)
}
type userRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, user models.User) error {
	query := `INSERT INTO users (id, name, password) VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, query, user.ID, user.Username, user.Password)
	return err
}

func (r *userRepo) GetByUsername(ctx context.Context, username string) (models.User, error) {
	query := `SELECT * FROM users WHERE name = $1`
	var user models.User
	err := r.db.QueryRowContext(ctx, query, username).Scan(&user.ID, &user.Username, &user.Password)
	if err == sql.ErrNoRows {
		return models.User{}, errors.New("user not found")
	}
	return user, err
}
func (r *userRepo) GetByID(ctx context.Context, id int) (models.User, error) {
	query := `SELECT * FROM users WHERE id = $1`
	var user models.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Username, &user.Password)
	if err == sql.ErrNoRows {
		return models.User{}, errors.New("user not found")
	}
	return user, err
}

func (r *userRepo) Update(ctx context.Context, user models.User) error {
	query := `UPDATE users SET name = $2, password = $3 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, user.ID, user.Username, user.Password)
	return err
}

func (r *userRepo) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *userRepo) List(ctx context.Context) ([]models.User, error) {
	query := `SELECT * FROM users`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Password); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
