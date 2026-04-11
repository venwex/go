package users

import (
	"example/test/internal/service/users"
	u "example/test/internal/utils"
	"log"
	"net/http"
	"strconv"
)

type H map[string]any
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

func (h *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	user, err := u.DecodeUser(r)
	if err != nil {
		u.RenderError(w, http.StatusBadRequest, err.Error())
		return
	}

	hashedPassword, err := u.HashPassword(user.Password)
	if err != nil {
		u.RenderError(w, http.StatusInternalServerError, err.Error())
		return
	}

	user.Password = hashedPassword
	user.Role = "user"

	createdUser, sessionId, err := h.Users.RegisterUser(user)
	if err != nil {
		u.RenderError(w, http.StatusInternalServerError, err.Error())
		return
	}

	u.RenderJSON(w, http.StatusOK, H{
		"message":   "user created successfully",
		"sessionId": sessionId,
		"user":      createdUser,
	})
}

func (h *UserHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	login, err := u.DecodeLogin(r)
	if err != nil {
		u.RenderError(w, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.Users.SignIn(login)
	if err != nil {
		u.RenderError(w, http.StatusInternalServerError, err.Error())
		return
	}

	u.RenderJSON(w, http.StatusOK, H{
		"token": token,
	})
}

func (h *UserHandler) ProtectedHello(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID")
	role := r.Context().Value("role")

	u.RenderJSON(w, http.StatusOK, H{
		"message": "OK",
		"user_id": userID,
		"role":    role,
	})
}

var jwtSecret = []byte("super-secret-key")

func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)

	user, err := h.Users.ServiceGetUser(userID)
	if err != nil {
		u.RenderError(w, http.StatusInternalServerError, err.Error())
		return
	}

	u.RenderJSON(w, http.StatusOK, H{
		"user": user,
	})
}

func (h *UserHandler) PromoteUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, _ := strconv.Atoi(idStr)

	err := h.Users.PromoteUser(id)
	if err != nil {
		log.Printf("error promoting user: %v, path: %v", err, r.URL.Path)
		u.RenderError(w, http.StatusInternalServerError, err.Error())
		return
	}

	u.RenderJSON(w, http.StatusOK, H{
		"message": "user promoted",
	})
}

func (h *UserHandler) Admin(w http.ResponseWriter, r *http.Request) {
	u.RenderJSON(w, http.StatusOK, H{
		"message": "Access for admin endpoint is successful",
	})
}
