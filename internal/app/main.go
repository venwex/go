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
	"os"
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
	log.Println("Starting server...")
	log.Fatal(http.ListenAndServe(":8000", handler))
}

func buildHandler(db *postgres.Dialect) *handlers.Handlers {
	repo := repository.NewRepositories(db)
	svc := service.NewServices(repo)
	h := handlers.NewHandlers(svc)

	return h
}

func setUpRoutes(h *handlers.Handlers) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", h.Task.HandleHealth)

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
		Host:        getEnv("DB_HOST", "localhost"),
		Port:        getEnv("DB_PORT", "5432"),
		Username:    getEnv("DB_USER", "postgres"),
		Password:    getEnv("DB_PASSWORD", "secret"),
		DBName:      getEnv("DB_NAME", "postgres"),
		SSLMode:     getEnv("DB_SSLMODE", "disable"),
		ExecTimeout: 5 * time.Second,
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
