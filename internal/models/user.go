package models

import "time"

type User struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Gender    *string   `json:"gender" db:"gender"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"password" db:"password"`
	Role      string    `json:"role" db:"role"`
	BirthDate time.Time `json:"birth_date" db:"birth_date"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type UserQuery struct {
	Page     int // offset
	PageSize int // limit

	Filters UserFilters
	Sorting UserSorting
}

type UserFilters struct {
	ID        *int
	Name      string
	Email     string
	Gender    string
	BirthDate string
}

type UserSorting struct {
	OrderBy  string
	OrderDir string
}

type PaginatedResponse struct {
	Data       []User `json:"data"`
	TotalCount int    `json:"totalCount"`
	Page       int    `json:"page"`
	PageSize   int    `json:"pageSize"`
}

type LoginUserDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
