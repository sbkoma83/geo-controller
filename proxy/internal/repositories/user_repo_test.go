package repositories

import (
	"context"
	"geo-controller/proxy/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"regexp"
	"testing"
)

func TestUserRepo_Create(t *testing.T) {
	// Создаем новый SQL-мок
	sqlDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer sqlDB.Close()

	// Создаем диалект GORM для Postgres с нашим моком
	dialector := postgres.New(postgres.Config{
		Conn: sqlDB,
	})

	// Создаем экземпляр GORM с нашим моком
	db, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	// Создаем репозиторий с мок-базой данных
	repo := NewUserRepository(db)

	// Подготавливаем ожидаемый запрос
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users" ("username","password") VALUES ($1,$2) RETURNING "id"`)).
		WithArgs("testuser", "testpassword").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	// Выполняем тестируемый метод
	user := models.User{
		Username: "testuser",
		Password: "testpassword",
	}
	err = repo.Create(context.Background(), user)

	// Проверяем результаты
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepo_GetByUsername(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer sqlDB.Close()

	dialector := postgres.New(postgres.Config{
		Conn: sqlDB,
	})

	db, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	repo := NewUserRepository(db)

	expectedUser := models.User{
		ID:       1,
		Username: "testuser",
		Password: "hashedpassword",
	}

	rows := sqlmock.NewRows([]string{"id", "username", "password"}).
		AddRow(expectedUser.ID, expectedUser.Username, expectedUser.Password)

	// Обновляем ожидаемый SQL-запрос
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE username = $1 ORDER BY "users"."id" LIMIT $2`)).
		WithArgs("testuser", 1).
		WillReturnRows(rows)

	user, err := repo.GetByUsername(context.Background(), "testuser")

	assert.NoError(t, err)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Username, user.Username)
	assert.Equal(t, expectedUser.Password, user.Password)
	assert.NoError(t, mock.ExpectationsWereMet())
}
func TestUserRepo_GetByID(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer sqlDB.Close()

	dialector := postgres.New(postgres.Config{
		Conn: sqlDB,
	})

	db, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	repo := NewUserRepository(db)

	expectedUser := models.User{
		ID:       1,
		Username: "testuser",
		Password: "hashedpassword",
	}

	rows := sqlmock.NewRows([]string{"id", "username", "password"}).
		AddRow(expectedUser.ID, expectedUser.Username, expectedUser.Password)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1 ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(uint32(1), 1).
		WillReturnRows(rows)

	user, err := repo.GetByID(context.Background(), uint32(1))

	assert.NoError(t, err)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Username, user.Username)
	assert.Equal(t, expectedUser.Password, user.Password)
	assert.NoError(t, mock.ExpectationsWereMet())
}
func TestUserRepo_Update(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer sqlDB.Close()

	dialector := postgres.New(postgres.Config{
		Conn: sqlDB,
	})

	db, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	repo := NewUserRepository(db)

	user := models.User{
		ID:       1,
		Username: "updateduser",
		Password: "updatedpassword",
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "username"=$1,"password"=$2 WHERE "id" = $3`)).
		WithArgs(user.Username, user.Password, user.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = repo.Update(context.Background(), user)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
func TestUserRepo_Delete(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer sqlDB.Close()

	dialector := postgres.New(postgres.Config{
		Conn: sqlDB,
	})

	db, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	repo := NewUserRepository(db)

	userID := uint32(1)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "users" WHERE "users"."id" = $1`)).
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err = repo.Delete(context.Background(), userID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
