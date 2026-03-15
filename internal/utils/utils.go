package utils

import (
	"encoding/json"
	e "example/test/internal/errors"
	m "example/test/internal/models"
	"log"
	"net/http"
	"strconv"
)

type H map[string]any

func GetIDFromQuery(r *http.Request) (int, error) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		return 0, e.ErrMissingId
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, e.ErrInvalidID
	}

	return id, nil
}

func RenderJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Println("json encode error: ", err)
	}
}

func DecodeTask(r *http.Request) (string, error) {
	title := m.Task{}

	err := json.NewDecoder(r.Body).Decode(&title)
	if err != nil {
		log.Println("error during decoding json")
		return "", err
	}

	return title.Title, nil
}

func DecodeUser(r *http.Request) (m.User, error) {
	var user m.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println("error during decoding json")
		return m.User{}, err
	}

	return user, nil
}

func RenderError(w http.ResponseWriter, status int, text string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(H{"error": text}); err != nil {
		log.Println("error during encoding an error")
	}
}

func ParseUserQuery(r *http.Request) m.UserQuery {
	q := r.URL.Query()

	page := 1
	pageSize := 10 // by default pagination is 10

	if p, err := strconv.Atoi(q.Get("page")); err == nil && p > 0 {
		page = p
	}

	if ps, err := strconv.Atoi(q.Get("page_size")); err == nil && ps > 0 {
		pageSize = ps
	}

	return m.UserQuery{
		Page:     page,
		PageSize: pageSize,

		Filters: m.UserFilters{
			Name:   q.Get("name"),
			Email:  q.Get("email"),
			Gender: q.Get("gender"),
		},

		Sorting: m.UserSorting{
			OrderBy:  q.Get("order_by"),
			OrderDir: q.Get("order_dir"),
		},
	}
}
