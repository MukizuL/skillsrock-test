package models

import (
	"context"
	"github.com/jackc/pgx/v5"
	"time"
)

type Task struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
type TaskModel struct {
	DB *pgx.Conn
}

func (m *TaskModel) Insert(title, description string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := m.DB.Exec(ctx, "INSERT INTO tasks (title, description) VALUES ($1, $2)", title, description)
	if err != nil {
		return err
	}

	return nil
}

func (m *TaskModel) SelectAll() ([]Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rows, err := m.DB.Query(ctx, "SELECT id, title, description, status, created_at, updated_at FROM tasks")
	if err != nil {
		return nil, err
	}

	tasks, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (Task, error) {
		var t Task
		err := row.Scan(
			&t.Id,
			&t.Title,
			&t.Description,
			&t.Status,
			&t.CreatedAt,
			&t.UpdatedAt)
		return t, err
	})
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (m *TaskModel) Update(id int, title, description, status string) (Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := m.DB.Exec(ctx, "UPDATE tasks SET title = $1, description = $2, status = $3, updated_at = now() WHERE id = $4",
		title, description, status, id)
	if err != nil {
		return Task{}, err
	}

	var t Task
	row := m.DB.QueryRow(ctx, "SELECT id, title, description, status, created_at, updated_at FROM tasks WHERE id = $1", id)
	err = row.Scan(
		&t.Id,
		&t.Title,
		&t.Description,
		&t.Status,
		&t.CreatedAt,
		&t.UpdatedAt)
	if err != nil {
		return Task{}, err
	}

	return t, nil
}

func (m *TaskModel) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := m.DB.Exec(ctx, "DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}
