package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	Actions "github.com/nithish-95/to-do/Actions"
	"github.com/nithish-95/to-do/Storage"
)

var tmpl *template.Template

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Initialize SQLite storage.
	store, err := Storage.NewSQLiteStorage("todos.db")
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}
	defer store.Close()

	// Create todo service.
	todoService := Actions.NewTodoService(store)

	// Parse HTML templates.
	tmpl = template.Must(template.ParseGlob("templates/*.html"))

	// Home page route.
	r.Get("/", handleHome(todoService))

	// JSON API routes.
	r.Route("/todos", func(r chi.Router) {
		r.Post("/", handleCreateTodo(todoService))
		r.Get("/", handleListTodos(todoService))
		r.Get("/{id}", handleGetTodo(todoService))
		r.Put("/{id}", handleUpdateTodo(todoService))
		r.Delete("/{id}", handleDeleteTodo(todoService))
	})

	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", r)
}

// handleHome renders the home page with the list of todos.
func handleHome(service Actions.TodoService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		todos, err := service.ListTodos(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := tmpl.ExecuteTemplate(w, "home.html", todos); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// handleCreateTodo creates a new todo from JSON data.
func handleCreateTodo(service Actions.TodoService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var todo Storage.Todo
		if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := service.CreateTodo(r.Context(), &todo); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(todo)
	}
}

// handleListTodos returns all todos as JSON.
func handleListTodos(service Actions.TodoService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		todos, err := service.ListTodos(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(todos)
	}
}

// handleGetTodo returns a single todo as JSON.
func handleGetTodo(service Actions.TodoService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
		todo, err := service.GetTodo(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if todo == nil {
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode(todo)
	}
}

// handleUpdateTodo processes an AJAX PUT request to update a todo.
func handleUpdateTodo(service Actions.TodoService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
		var todo Storage.Todo
		if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		todo.ID = id // enforce URL ID
		if err := service.UpdateTodo(r.Context(), &todo); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(todo)
	}
}

// handleDeleteTodo processes an AJAX DELETE request to remove a todo.
func handleDeleteTodo(service Actions.TodoService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
		if err := service.DeleteTodo(r.Context(), id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
