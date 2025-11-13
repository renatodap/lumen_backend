package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/lumen/backend-go/internal/models"
)

type HabitRepository interface {
	Create(ctx context.Context, habit *models.Habit) error
	GetByID(ctx context.Context, id, userID uuid.UUID) (*models.Habit, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]models.Habit, error)
	Update(ctx context.Context, habit *models.Habit) error
	Delete(ctx context.Context, id, userID uuid.UUID) error
}

type habitRepository struct {
	db *Database
}

func NewHabitRepository(db *Database) HabitRepository {
	return &habitRepository{db: db}
}

func (r *habitRepository) Create(ctx context.Context, habit *models.Habit) error {
	query := `
		INSERT INTO habits (id, user_id, name, color, icon, frequency, target_count, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at
	`

	habit.ID = uuid.New()
	habit.IsActive = true
	habit.CreatedAt = time.Now()
	habit.UpdatedAt = time.Now()

	err := r.db.Pool.QueryRow(
		ctx,
		query,
		habit.ID,
		habit.UserID,
		habit.Name,
		habit.Color,
		habit.Icon,
		habit.Frequency,
		habit.TargetCount,
		habit.IsActive,
		habit.CreatedAt,
		habit.UpdatedAt,
	).Scan(&habit.ID, &habit.CreatedAt, &habit.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create habit: %w", err)
	}

	return nil
}

func (r *habitRepository) GetByID(ctx context.Context, id, userID uuid.UUID) (*models.Habit, error) {
	query := `
		SELECT id, user_id, name, color, icon, frequency, target_count, is_active, created_at, updated_at
		FROM habits
		WHERE id = $1 AND user_id = $2
	`

	var habit models.Habit
	err := r.db.Pool.QueryRow(ctx, query, id, userID).Scan(
		&habit.ID,
		&habit.UserID,
		&habit.Name,
		&habit.Color,
		&habit.Icon,
		&habit.Frequency,
		&habit.TargetCount,
		&habit.IsActive,
		&habit.CreatedAt,
		&habit.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, models.ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get habit: %w", err)
	}

	return &habit, nil
}

func (r *habitRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]models.Habit, error) {
	query := `
		SELECT id, user_id, name, color, icon, frequency, target_count, is_active, created_at, updated_at
		FROM habits
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get habits: %w", err)
	}
	defer rows.Close()

	var habits []models.Habit
	for rows.Next() {
		var habit models.Habit
		err := rows.Scan(
			&habit.ID,
			&habit.UserID,
			&habit.Name,
			&habit.Color,
			&habit.Icon,
			&habit.Frequency,
			&habit.TargetCount,
			&habit.IsActive,
			&habit.CreatedAt,
			&habit.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan habit: %w", err)
		}
		habits = append(habits, habit)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating habits: %w", err)
	}

	return habits, nil
}

func (r *habitRepository) Update(ctx context.Context, habit *models.Habit) error {
	query := `
		UPDATE habits
		SET name = $3, color = $4, icon = $5, frequency = $6, target_count = $7, is_active = $8, updated_at = $9
		WHERE id = $1 AND user_id = $2
		RETURNING updated_at
	`

	habit.UpdatedAt = time.Now()

	err := r.db.Pool.QueryRow(
		ctx,
		query,
		habit.ID,
		habit.UserID,
		habit.Name,
		habit.Color,
		habit.Icon,
		habit.Frequency,
		habit.TargetCount,
		habit.IsActive,
		habit.UpdatedAt,
	).Scan(&habit.UpdatedAt)

	if err == pgx.ErrNoRows {
		return models.ErrNotFound
	}

	if err != nil {
		return fmt.Errorf("failed to update habit: %w", err)
	}

	return nil
}

func (r *habitRepository) Delete(ctx context.Context, id, userID uuid.UUID) error {
	query := `DELETE FROM habits WHERE id = $1 AND user_id = $2`

	result, err := r.db.Pool.Exec(ctx, query, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete habit: %w", err)
	}

	if result.RowsAffected() == 0 {
		return models.ErrNotFound
	}

	return nil
}
