package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/lumen/backend-go/internal/models"
)

type TaskRepository interface {
	Create(ctx context.Context, task *models.Task) error
	GetByID(ctx context.Context, id, userID uuid.UUID) (*models.Task, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, filter models.TaskFilter) ([]models.Task, error)
	Update(ctx context.Context, task *models.Task) error
	Delete(ctx context.Context, id, userID uuid.UUID) error
}

type taskRepository struct {
	db *Database
}

func NewTaskRepository(db *Database) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) Create(ctx context.Context, task *models.Task) error {
	query := `
		INSERT INTO tasks (id, user_id, title, description, horizon, priority, status, due_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at
	`

	task.ID = uuid.New()
	task.Status = "todo"
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	err := r.db.Pool.QueryRow(
		ctx,
		query,
		task.ID,
		task.UserID,
		task.Title,
		task.Description,
		task.Horizon,
		task.Priority,
		task.Status,
		task.DueDate,
		task.CreatedAt,
		task.UpdatedAt,
	).Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	return nil
}

func (r *taskRepository) GetByID(ctx context.Context, id, userID uuid.UUID) (*models.Task, error) {
	query := `
		SELECT id, user_id, title, description, horizon, priority, status, due_date, completed_at, created_at, updated_at
		FROM tasks
		WHERE id = $1 AND user_id = $2
	`

	var task models.Task
	err := r.db.Pool.QueryRow(ctx, query, id, userID).Scan(
		&task.ID,
		&task.UserID,
		&task.Title,
		&task.Description,
		&task.Horizon,
		&task.Priority,
		&task.Status,
		&task.DueDate,
		&task.CompletedAt,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, models.ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return &task, nil
}

func (r *taskRepository) GetByUserID(ctx context.Context, userID uuid.UUID, filter models.TaskFilter) ([]models.Task, error) {
	query := `
		SELECT id, user_id, title, description, horizon, priority, status, due_date, completed_at, created_at, updated_at
		FROM tasks
		WHERE user_id = $1
	`

	args := []interface{}{userID}
	argCount := 1

	if filter.Horizon != "" {
		argCount++
		query += fmt.Sprintf(" AND horizon = $%d", argCount)
		args = append(args, filter.Horizon)
	}

	if filter.Status != "" {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, filter.Status)
	}

	if filter.Priority != "" {
		argCount++
		query += fmt.Sprintf(" AND priority = $%d", argCount)
		args = append(args, filter.Priority)
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(
			&task.ID,
			&task.UserID,
			&task.Title,
			&task.Description,
			&task.Horizon,
			&task.Priority,
			&task.Status,
			&task.DueDate,
			&task.CompletedAt,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tasks: %w", err)
	}

	return tasks, nil
}

func (r *taskRepository) Update(ctx context.Context, task *models.Task) error {
	setClauses := []string{
		"title = $3",
		"description = $4",
		"horizon = $5",
		"priority = $6",
		"status = $7",
		"due_date = $8",
		"updated_at = $9",
	}

	if task.Status == "done" && task.CompletedAt == nil {
		now := time.Now()
		task.CompletedAt = &now
		setClauses = append(setClauses, "completed_at = $10")
	}

	query := fmt.Sprintf(`
		UPDATE tasks
		SET %s
		WHERE id = $1 AND user_id = $2
		RETURNING updated_at
	`, strings.Join(setClauses, ", "))

	task.UpdatedAt = time.Now()

	err := r.db.Pool.QueryRow(
		ctx,
		query,
		task.ID,
		task.UserID,
		task.Title,
		task.Description,
		task.Horizon,
		task.Priority,
		task.Status,
		task.DueDate,
		task.UpdatedAt,
		task.CompletedAt,
	).Scan(&task.UpdatedAt)

	if err == pgx.ErrNoRows {
		return models.ErrNotFound
	}

	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	return nil
}

func (r *taskRepository) Delete(ctx context.Context, id, userID uuid.UUID) error {
	query := `DELETE FROM tasks WHERE id = $1 AND user_id = $2`

	result, err := r.db.Pool.Exec(ctx, query, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	if result.RowsAffected() == 0 {
		return models.ErrNotFound
	}

	return nil
}
