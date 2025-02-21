package Storage

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteStorage struct {
	db *sql.DB
}

func NewSQLiteStorage(dbPath string) (*SQLiteStorage, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Create todos table if it doesn't exist.
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS todos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			description TEXT,
			completed BOOLEAN DEFAULT FALSE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return nil, err
	}

	return &SQLiteStorage{db: db}, nil
}

// Close closes the underlying database connection.
func (s *SQLiteStorage) Close() error {
	return s.db.Close()
}

// CreateTodo inserts a new todo into the database.
func (s *SQLiteStorage) CreateTodo(ctx context.Context, todo *Todo) error {
	query := `
		INSERT INTO todos (title, description, completed, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`
	now := time.Now()
	result, err := s.db.ExecContext(ctx, query,
		todo.Title,
		todo.Description,
		todo.Completed,
		now,
		now,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	todo.ID = int(id)
	todo.CreatedAt = now
	todo.UpdatedAt = now
	return nil
}

// GetTodo retrieves a todo by ID.
func (s *SQLiteStorage) GetTodo(ctx context.Context, id int) (*Todo, error) {
	todo := &Todo{}
	query := `SELECT id, title, description, completed, created_at, updated_at FROM todos WHERE id = ?`
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Description,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return todo, nil
}

// UpdateTodo updates an existing todo.
func (s *SQLiteStorage) UpdateTodo(ctx context.Context, todo *Todo) error {
	query := `
		UPDATE todos 
		SET title = ?, description = ?, completed = ?, updated_at = ?
		WHERE id = ?
	`
	now := time.Now()
	_, err := s.db.ExecContext(ctx, query,
		todo.Title,
		todo.Description,
		todo.Completed,
		now,
		todo.ID,
	)
	if err != nil {
		return err
	}
	todo.UpdatedAt = now
	return nil
}

// DeleteTodo removes a todo by ID.
func (s *SQLiteStorage) DeleteTodo(ctx context.Context, id int) error {
	query := `DELETE FROM todos WHERE id = ?`
	_, err := s.db.ExecContext(ctx, query, id)
	return err
}

// ListTodos retrieves all todos from the database.
func (s *SQLiteStorage) ListTodos(ctx context.Context) ([]*Todo, error) {
	query := `SELECT id, title, description, completed, created_at, updated_at FROM todos`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []*Todo
	for rows.Next() {
		todo := &Todo{}
		err := rows.Scan(
			&todo.ID,
			&todo.Title,
			&todo.Description,
			&todo.Completed,
			&todo.CreatedAt,
			&todo.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return todos, nil
}
