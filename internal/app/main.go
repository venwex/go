package app

import (
	"context"
	"example/test/internal/config"
	"example/test/internal/handlers"
	m "example/test/internal/middleware"
	"example/test/internal/repository"
	"example/test/internal/repository/postgres"
	"example/test/internal/service"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

func Run() {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbConfig := InitPostgresConfig()
	postgre := postgres.NewDialect(dbConfig)

	h := buildHandler(postgre) // DI
	mux := setUpRoutes(h)      // handlers

	handler := m.AuthMiddleware(m.LoggingMiddleware("message")(mux))
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func buildHandler(db *postgres.Dialect) *handlers.Handlers {
	repo := repository.NewRepositories(db)
	svc := service.NewServices(repo)
	h := handlers.NewHandlers(svc)

	return h
}

func setUpRoutes(h *handlers.Handlers) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /tasks", h.Task.HandleGetTasks)
	mux.HandleFunc("POST /tasks", h.Task.HandlePostTask)
	mux.HandleFunc("PATCH /tasks", h.Task.HandlePatchTask)
	mux.HandleFunc("DELETE /tasks", h.Task.HandleDeleteTask)

	mux.HandleFunc("GET /users/{id}", h.User.HandleGetUserById)
	mux.HandleFunc("GET /users", h.User.HandleGetUsers) // get all users
	mux.HandleFunc("POST /users", h.User.HandleCreateUser)
	mux.HandleFunc("PATCH /users/{id}", h.User.HandleUpdateUser)  // update specific user
	mux.HandleFunc("DELETE /users/{id}", h.User.HandleDeleteUser) // delete specific user

	return mux
}

func InitPostgresConfig() *config.PostgresConfig {
	return &config.PostgresConfig{
		Host:        "localhost",
		Port:        "5433",
		Username:    "postgres",
		Password:    "secret",
		DBName:      "gopgtest",
		SSLMode:     "disable",
		ExecTimeout: 5 * time.Second,
	}
}
