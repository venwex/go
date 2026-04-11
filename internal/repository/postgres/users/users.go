package users

import (
	"database/sql"
	"errors"
	e "example/test/internal/errors"
	m "example/test/internal/models"
	"example/test/internal/repository/postgres"
	"fmt"
	"strings"
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

func (r *UserRepository) GetUsers(q m.UserQuery) (m.PaginatedResponse, error) {
	offset := (q.Page - 1) * q.PageSize

	baseQuery := "FROM users WHERE 1=1"
	args := []interface{}{}
	argID := 1

	if q.Filters.ID != nil {
		baseQuery += fmt.Sprintf(" AND id = $%d", argID)
		args = append(args, *q.Filters.ID)
		argID++
	}

	if q.Filters.Name != "" {
		baseQuery += fmt.Sprintf(" AND name ILIKE $%d", argID)
		args = append(args, "%"+q.Filters.Name+"%")
		argID++
	}

	if q.Filters.Email != "" {
		baseQuery += fmt.Sprintf(" AND email ILIKE $%d", argID)
		args = append(args, "%"+q.Filters.Email+"%")
		argID++
	}

	if q.Filters.Gender != "" {
		baseQuery += fmt.Sprintf(" AND gender = $%d", argID)
		args = append(args, q.Filters.Gender)
		argID++
	}

	if q.Filters.BirthDate != "" {
		baseQuery += fmt.Sprintf(" AND birth_date = $%d", argID)
		args = append(args, q.Filters.BirthDate)
		argID++
	}

	var totalCount int

	countQuery := "SELECT COUNT(*) " + baseQuery

	err := r.db.DB.Get(&totalCount, countQuery, args...)
	if err != nil {
		return m.PaginatedResponse{}, err
	}

	dataQuery := `
		SELECT id, name, email, gender, birth_date
	` + baseQuery

	allowedSort := map[string]bool{
		"id":         true,
		"name":       true,
		"email":      true,
		"birth_date": true,
		"gender":     true,
	}

	orderBy := "id"
	if allowedSort[q.Sorting.OrderBy] {
		orderBy = q.Sorting.OrderBy
	}

	orderDir := "ASC"
	if strings.ToUpper(q.Sorting.OrderDir) == "DESC" {
		orderDir = "DESC"
	}

	dataQuery += " ORDER BY " + orderBy + " " + orderDir

	dataQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argID, argID+1)

	args = append(args, q.PageSize, offset)

	var users []m.User

	err = r.db.DB.Select(&users, dataQuery, args...)
	if err != nil {
		return m.PaginatedResponse{}, err
	}

	return m.PaginatedResponse{
		Data:       users,
		TotalCount: totalCount,
		Page:       q.Page,
		PageSize:   q.PageSize,
	}, nil
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
		INSERT INTO users (name, email, password, gender, birth_date, role)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, name, gender, email, password, role, birth_date, created_at, updated_at`,
		user.Name,
		user.Email,
		user.Password,
		user.Gender,
		user.BirthDate,
		user.Role,
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

func (r *UserRepository) GetCommonFriends(u1, u2 int) ([]m.User, error) {
	query := `
        SELECT u.*
        FROM user_friends uf1
        JOIN user_friends uf2
            ON uf1.friend_id = uf2.friend_id
        JOIN users u
            ON u.id = uf1.friend_id
        WHERE uf1.user_id = $1
        AND uf2.user_id = $2
    ` // join to solve n + 1 problem

	rows, err := r.db.DB.Query(query, u1, u2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []m.User

	for rows.Next() {
		var u m.User

		err := rows.Scan(
			&u.ID,
			&u.Name,
			&u.Gender,
			&u.Email,
			&u.BirthDate,
			&u.CreatedAt,
			&u.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) GetUserByEmail(email string) (m.User, error) {
	var user m.User

	err := r.db.DB.Get(&user, "SELECT * FROM users WHERE email = $1", email)

	return user, err
}

func (r *UserRepository) PromoteUser(id int) error {
	_, err := r.db.DB.Exec(
		`UPDATE users SET role = 'admin' WHERE id = $1`,
		id,
	)

	return err
}
