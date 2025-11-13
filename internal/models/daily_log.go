package models

import (
	"time"

	"github.com/google/uuid"
)

type DailyLog struct {
	ID                 uuid.UUID `json:"id" db:"id"`
	UserID             uuid.UUID `json:"user_id" db:"user_id"`
	Date               time.Time `json:"date" db:"date"`
	MorningRoutine     bool      `json:"morning_routine" db:"morning_routine"`
	EveningRoutine     bool      `json:"evening_routine" db:"evening_routine"`
	WaterIntake        int       `json:"water_intake" db:"water_intake"`
	SleepHours         float64   `json:"sleep_hours" db:"sleep_hours"`
	EnergyLevel        int       `json:"energy_level" db:"energy_level"`
	MoodRating         int       `json:"mood_rating" db:"mood_rating"`
	ProductivityRating int       `json:"productivity_rating" db:"productivity_rating"`
	Notes              string    `json:"notes" db:"notes"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

type CreateDailyLogRequest struct {
	Date               time.Time `json:"date" binding:"required"`
	MorningRoutine     bool      `json:"morning_routine"`
	EveningRoutine     bool      `json:"evening_routine"`
	WaterIntake        int       `json:"water_intake" binding:"min=0,max=20"`
	SleepHours         float64   `json:"sleep_hours" binding:"min=0,max=24"`
	EnergyLevel        int       `json:"energy_level" binding:"min=1,max=5"`
	MoodRating         int       `json:"mood_rating" binding:"min=1,max=5"`
	ProductivityRating int       `json:"productivity_rating" binding:"min=1,max=5"`
	Notes              string    `json:"notes" binding:"max=1000"`
}

type UpdateDailyLogRequest struct {
	MorningRoutine     *bool    `json:"morning_routine"`
	EveningRoutine     *bool    `json:"evening_routine"`
	WaterIntake        *int     `json:"water_intake" binding:"omitempty,min=0,max=20"`
	SleepHours         *float64 `json:"sleep_hours" binding:"omitempty,min=0,max=24"`
	EnergyLevel        *int     `json:"energy_level" binding:"omitempty,min=1,max=5"`
	MoodRating         *int     `json:"mood_rating" binding:"omitempty,min=1,max=5"`
	ProductivityRating *int     `json:"productivity_rating" binding:"omitempty,min=1,max=5"`
	Notes              *string  `json:"notes" binding:"omitempty,max=1000"`
}

type DailyLogStats struct {
	Date               time.Time `json:"date"`
	HabitsCompleted    int       `json:"habits_completed"`
	HabitsTotal        int       `json:"habits_total"`
	TasksCompleted     int       `json:"tasks_completed"`
	TasksTotal         int       `json:"tasks_total"`
	MoodRating         int       `json:"mood_rating"`
	EnergyLevel        int       `json:"energy_level"`
	ProductivityRating int       `json:"productivity_rating"`
}

func (d *DailyLog) Validate() error {
	if d.WaterIntake < 0 || d.WaterIntake > 20 {
		return ErrInvalidWaterIntake
	}

	if d.SleepHours < 0 || d.SleepHours > 24 {
		return ErrInvalidSleepHours
	}

	if d.EnergyLevel < 1 || d.EnergyLevel > 5 {
		return ErrInvalidRating
	}

	if d.MoodRating < 1 || d.MoodRating > 5 {
		return ErrInvalidRating
	}

	if d.ProductivityRating < 1 || d.ProductivityRating > 5 {
		return ErrInvalidRating
	}

	return nil
}
