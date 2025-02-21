package actions

import (
	"context"

	"github.com/nithish-95/to-do/Storage"
)

// TodoService defines the interface for todo operations.
type TodoService interface {
	CreateTodo(ctx context.Context, todo *Storage.Todo) error
	GetTodo(ctx context.Context, id int) (*Storage.Todo, error)
	UpdateTodo(ctx context.Context, todo *Storage.Todo) error
	DeleteTodo(ctx context.Context, id int) error
	ListTodos(ctx context.Context) ([]*Storage.Todo, error)
}

// todoService implements TodoService.
type todoService struct {
	storage *Storage.SQLiteStorage
}

// NewTodoService creates a new TodoService instance.
func NewTodoService(s *Storage.SQLiteStorage) TodoService {
	return &todoService{storage: s}
}

func (s *todoService) CreateTodo(ctx context.Context, todo *Storage.Todo) error {
	return s.storage.CreateTodo(ctx, todo)
}

func (s *todoService) GetTodo(ctx context.Context, id int) (*Storage.Todo, error) {
	return s.storage.GetTodo(ctx, id)
}

func (s *todoService) UpdateTodo(ctx context.Context, todo *Storage.Todo) error {
	return s.storage.UpdateTodo(ctx, todo)
}

func (s *todoService) DeleteTodo(ctx context.Context, id int) error {
	return s.storage.DeleteTodo(ctx, id)
}

func (s *todoService) ListTodos(ctx context.Context) ([]*Storage.Todo, error) {
	return s.storage.ListTodos(ctx)
}
