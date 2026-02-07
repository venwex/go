package main

import (
	"example/test/internal/handlers"
	m "example/test/internal/middleware"
	"example/test/internal/service"
	"example/test/internal/store"
	"log"
	"net/http"
)

func main() {
	store := store.NewStore()
	service := service.NewService(store)
	handlers := handlers.NewHandler(service)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks", handlers.HandleGetTasks)
	mux.HandleFunc("POST /tasks", handlers.HandlePostTask)
	mux.HandleFunc("PATCH /tasks", handlers.HandlePatchTask)
	mux.HandleFunc("DELETE /tasks", handlers.HandleDeleteTask)

	handler := m.AuthMiddleware(m.LoggingMiddleware("message")(mux))
	log.Fatal(http.ListenAndServe(":8080", handler))
}
