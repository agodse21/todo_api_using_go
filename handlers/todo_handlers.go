package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-mongo-todos/services"
	"github.com/go-playground/validator"
)

var validate = validator.New()

var todo services.Todo

func healthCheck(w http.ResponseWriter, r *http.Request) {
	res := Response{
		Message: "OK",
		Status:  200,
	}
	jsonResponse, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)

}

func createTodo(w http.ResponseWriter, r *http.Request) {

	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate struct fields
	if err := validate.Struct(todo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	todo.CreatedAt = time.Now()
	todo.UpdatedAt = time.Now()

	err = todo.InsertTodo(todo)

	if err != nil {
		errMsg := Response{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
		json.NewEncoder(w).Encode(errMsg)
		return
	}

	res := Response{
		Message: "Successfully created todo",
		Status:  http.StatusOK,
	}

	jsonResponse, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)

}

func getAllTodos(w http.ResponseWriter, r *http.Request) {

	todos, err := todo.GetAllTodos()
	if err != nil {
		slog.Error("Error getting todos")
		errMsg := Response{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
		json.NewEncoder(w).Encode(errMsg)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(todos)

}

func getTodoByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	todo, err := todo.GetTodoByID(id)
	if err != nil {
		slog.Error("Error getting todo")
		errMsg := Response{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
		json.NewEncoder(w).Encode(errMsg)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(todo)
}

func deleteTodoByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := todo.DeleteTodoByID(id)
	if err != nil {
		slog.Error("Error deleting todo")
		errMsg := Response{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
		json.NewEncoder(w).Encode(errMsg)
		return
	}

	res := Response{
		Message: "Successfully deleted todo",
		Status:  http.StatusOK,
	}

	jsonResponse, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func updateTodoByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = todo.UpdateTodoByID(id, todo)

	if err != nil {
		slog.Error("Error updating todo")
		errMsg := Response{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
		json.NewEncoder(w).Encode(errMsg)
		return
	}

	res := Response{
		Message: "Successfully updated todo",
		Status:  http.StatusOK,
	}

	jsonResponse, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)

}
