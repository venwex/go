package users

import (
	"example/test/internal/service/users"
	u "example/test/internal/utils"
	"net/http"
	"strconv"
)

type UserHandler struct {
	Users *users.UserService
}

func NewUserHandler(users *users.UserService) *UserHandler {
	return &UserHandler{Users: users}
}

func (handler *UserHandler) HandleGetUserById(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		u.RenderError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := handler.Users.ServiceGetUser(id)
	if err != nil {
		u.RenderError(w, http.StatusNotFound, err.Error())
		return
	}

	u.RenderJSON(w, http.StatusOK, user)
}

/*
pagination: /users?page=1&page_size=10
sorting: 	/users?order_by=birth_date (по дефолту id)
filtering:  /users?gender=female&name=Anna
*/

func (handler *UserHandler) HandleGetUsers(w http.ResponseWriter, r *http.Request) {
	query := u.ParseUserQuery(r)

	users, err := handler.Users.ServiceGetUsers(query)
	if err != nil {
		u.RenderError(w, http.StatusInternalServerError, err.Error())
		return
	}

	u.RenderJSON(w, http.StatusOK, users)
}

func (handler *UserHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	user, err := u.DecodeUser(r)
	if err != nil {
		u.RenderError(w, http.StatusBadRequest, err.Error())
		return
	}

	if user.Name == "" {
		u.RenderError(w, http.StatusBadRequest, "name is required")
		return
	}

	if user.Email == "" {
		u.RenderError(w, http.StatusBadRequest, "email is required")
		return
	}

	user, err = handler.Users.ServiceCreateUser(user)
	if err != nil {
		u.RenderError(w, http.StatusInternalServerError, err.Error())
		return
	}

	u.RenderJSON(w, http.StatusCreated, user)
}

func (handler *UserHandler) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		u.RenderError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := u.DecodeUser(r)
	if err != nil {
		u.RenderError(w, http.StatusBadRequest, err.Error())
		return
	}

	if user.Name == "" {
		u.RenderError(w, http.StatusBadRequest, "name is required")
		return
	}

	if user.Email == "" {
		u.RenderError(w, http.StatusBadRequest, "email is required")
		return
	}

	user, err = handler.Users.ServiceUpdateUser(id, user.Name, user.Email)
	if err != nil {
		u.RenderError(w, http.StatusInternalServerError, err.Error())
		return
	}

	u.RenderJSON(w, http.StatusOK, user)
}

func (handler *UserHandler) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		u.RenderError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := handler.Users.ServiceDeleteUser(id)
	if err != nil {
		u.RenderError(w, http.StatusInternalServerError, err.Error())
		return
	}

	u.RenderJSON(w, http.StatusOK, user)
}

func (h *UserHandler) CommonFriends(w http.ResponseWriter, r *http.Request) {
	firstStr := r.URL.Query().Get("u1")
	secondStr := r.URL.Query().Get("u2")

	user1, err := strconv.Atoi(firstStr)
	if err != nil {
		u.RenderError(w, http.StatusBadRequest, err.Error())
		return
	}

	user2, err := strconv.Atoi(secondStr)
	if err != nil {
		u.RenderError(w, http.StatusBadRequest, err.Error())
		return
	}

	users, err := h.Users.ServiceGetCommonFriends(user1, user2)
	if err != nil {
		u.RenderError(w, http.StatusInternalServerError, err.Error())
		return
	}

	u.RenderJSON(w, http.StatusOK, users)
}
