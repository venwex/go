package users

import (
	"database/sql"
	"errors"
	e "example/test/internal/errors"
	m "example/test/internal/models"
	"example/test/internal/repository/postgres"
	"time"
)

type UserRepository struct {
	db      *postgres.Dialect
	timeout time.Duration
}

func NewRepository(db *postgres.Dialect) *UserRepository {
	return &UserRepository{
		db:      db,
		timeout: 5 * time.Second,
	}
}

func (r *UserRepository) GetUsers() ([]m.User, error) {
	var users []m.User

	err := r.db.DB.Select(&users, `
		SELECT id, name, email, created_at, updated_at
		FROM users`)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) GetUserByID(id int) (m.User, error) {
	if id <= 0 {
		return m.User{}, e.ErrInvalidID
	}

	var user m.User
	err := r.db.DB.Get(&user, `
		SELECT id, name, email, created_at, updated_at
		FROM users
		WHERE id = $1`, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return m.User{}, e.ErrUserNotFound
		}
		return m.User{}, err
	}

	return user, nil
}

func (r *UserRepository) CreateUser(user m.User) (m.User, error) {
	err := r.db.DB.Get(&user, `
		INSERT INTO users (name, email)
		VALUES ($1, $2)
		RETURNING id, name, email, created_at, updated_at`,
		user.Name, user.Email,
	)
	if err != nil {
		return m.User{}, err
	}

	return user, nil
}

func (r *UserRepository) UpdateUser(id int, name, email string) (m.User, error) {
	var user m.User

	err := r.db.DB.Get(&user, `
		UPDATE users
		SET name = $1,
		    email = $2,
		    updated_at = NOW()
		WHERE id = $3
		RETURNING id, name, email, created_at, updated_at`,
		name, email, id,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return m.User{}, e.ErrUserNotFound
		}
		return m.User{}, err
	}

	return user, nil
}

func (r *UserRepository) DeleteUser(id int) (m.User, error) {
	var user m.User

	err := r.db.DB.Get(&user, `
		DELETE FROM users
		WHERE id = $1
		RETURNING id, name, email, created_at, updated_at`,
		id,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return m.User{}, e.ErrUserNotFound
		}
		return m.User{}, err
	}

	return user, nil
}
