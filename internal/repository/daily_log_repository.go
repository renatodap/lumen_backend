package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/lumen/backend-go/internal/models"
)

type DailyLogRepository interface {
	Create(ctx context.Context, log *models.DailyLog) error
	GetByDate(ctx context.Context, userID uuid.UUID, date time.Time) (*models.DailyLog, error)
	GetByDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]models.DailyLog, error)
	Update(ctx context.Context, log *models.DailyLog) error
	Delete(ctx context.Context, id, userID uuid.UUID) error
}

type dailyLogRepository struct {
	db *Database
}

func NewDailyLogRepository(db *Database) DailyLogRepository {
	return &dailyLogRepository{db: db}
}

func (r *dailyLogRepository) Create(ctx context.Context, log *models.DailyLog) error {
	query := `
		INSERT INTO daily_logs (
			id, user_id, date, morning_routine, evening_routine, water_intake,
			sleep_hours, energy_level, mood_rating, productivity_rating, notes,
			created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		ON CONFLICT (user_id, date) DO UPDATE SET
			morning_routine = EXCLUDED.morning_routine,
			evening_routine = EXCLUDED.evening_routine,
			water_intake = EXCLUDED.water_intake,
			sleep_hours = EXCLUDED.sleep_hours,
			energy_level = EXCLUDED.energy_level,
			mood_rating = EXCLUDED.mood_rating,
			productivity_rating = EXCLUDED.productivity_rating,
			notes = EXCLUDED.notes,
			updated_at = EXCLUDED.updated_at
		RETURNING id, created_at, updated_at
	`

	log.ID = uuid.New()
	log.CreatedAt = time.Now()
	log.UpdatedAt = time.Now()

	err := r.db.Pool.QueryRow(
		ctx,
		query,
		log.ID,
		log.UserID,
		log.Date,
		log.MorningRoutine,
		log.EveningRoutine,
		log.WaterIntake,
		log.SleepHours,
		log.EnergyLevel,
		log.MoodRating,
		log.ProductivityRating,
		log.Notes,
		log.CreatedAt,
		log.UpdatedAt,
	).Scan(&log.ID, &log.CreatedAt, &log.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create daily log: %w", err)
	}

	return nil
}

func (r *dailyLogRepository) GetByDate(ctx context.Context, userID uuid.UUID, date time.Time) (*models.DailyLog, error) {
	query := `
		SELECT id, user_id, date, morning_routine, evening_routine, water_intake,
		       sleep_hours, energy_level, mood_rating, productivity_rating, notes,
		       created_at, updated_at
		FROM daily_logs
		WHERE user_id = $1 AND date = $2
	`

	var log models.DailyLog
	err := r.db.Pool.QueryRow(ctx, query, userID, date).Scan(
		&log.ID,
		&log.UserID,
		&log.Date,
		&log.MorningRoutine,
		&log.EveningRoutine,
		&log.WaterIntake,
		&log.SleepHours,
		&log.EnergyLevel,
		&log.MoodRating,
		&log.ProductivityRating,
		&log.Notes,
		&log.CreatedAt,
		&log.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, models.ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get daily log: %w", err)
	}

	return &log, nil
}

func (r *dailyLogRepository) GetByDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]models.DailyLog, error) {
	query := `
		SELECT id, user_id, date, morning_routine, evening_routine, water_intake,
		       sleep_hours, energy_level, mood_rating, productivity_rating, notes,
		       created_at, updated_at
		FROM daily_logs
		WHERE user_id = $1 AND date BETWEEN $2 AND $3
		ORDER BY date DESC
	`

	rows, err := r.db.Pool.Query(ctx, query, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get daily logs: %w", err)
	}
	defer rows.Close()

	var logs []models.DailyLog
	for rows.Next() {
		var log models.DailyLog
		err := rows.Scan(
			&log.ID,
			&log.UserID,
			&log.Date,
			&log.MorningRoutine,
			&log.EveningRoutine,
			&log.WaterIntake,
			&log.SleepHours,
			&log.EnergyLevel,
			&log.MoodRating,
			&log.ProductivityRating,
			&log.Notes,
			&log.CreatedAt,
			&log.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan daily log: %w", err)
		}
		logs = append(logs, log)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating daily logs: %w", err)
	}

	return logs, nil
}

func (r *dailyLogRepository) Update(ctx context.Context, log *models.DailyLog) error {
	query := `
		UPDATE daily_logs
		SET morning_routine = $3, evening_routine = $4, water_intake = $5,
		    sleep_hours = $6, energy_level = $7, mood_rating = $8,
		    productivity_rating = $9, notes = $10, updated_at = $11
		WHERE id = $1 AND user_id = $2
		RETURNING updated_at
	`

	log.UpdatedAt = time.Now()

	err := r.db.Pool.QueryRow(
		ctx,
		query,
		log.ID,
		log.UserID,
		log.MorningRoutine,
		log.EveningRoutine,
		log.WaterIntake,
		log.SleepHours,
		log.EnergyLevel,
		log.MoodRating,
		log.ProductivityRating,
		log.Notes,
		log.UpdatedAt,
	).Scan(&log.UpdatedAt)

	if err == pgx.ErrNoRows {
		return models.ErrNotFound
	}

	if err != nil {
		return fmt.Errorf("failed to update daily log: %w", err)
	}

	return nil
}

func (r *dailyLogRepository) Delete(ctx context.Context, id, userID uuid.UUID) error {
	query := `DELETE FROM daily_logs WHERE id = $1 AND user_id = $2`

	result, err := r.db.Pool.Exec(ctx, query, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete daily log: %w", err)
	}

	if result.RowsAffected() == 0 {
		return models.ErrNotFound
	}

	return nil
}
