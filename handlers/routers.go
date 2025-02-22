package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
)

type Response struct {
	Message string
	Status  int
}

func CreateRouter() *chi.Mux {

	router := chi.NewRouter()

	// Initialize CORS middleware
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any major browser
	})

	// Use CORS middleware with Chi
	router.Use(corsMiddleware.Handler)

	router.Route("/todo", func(route chi.Router) {
		route.Get("/health-check", healthCheck)
		route.Post("/create", createTodo)
		route.Get("/", getAllTodos)
		route.Get("/{id}", getTodoByID)
		route.Delete("/{id}", deleteTodoByID)
		route.Patch("/{id}", updateTodoByID)
	})

	return router

}
