package models

import (
	"time"

	"github.com/google/uuid"
)

type Habit struct {
	ID          uuid.UUID `json:"id" db:"id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	Name        string    `json:"name" db:"name" binding:"required"`
	Color       string    `json:"color" db:"color" binding:"required"`
	Icon        string    `json:"icon" db:"icon" binding:"required"`
	Frequency   string    `json:"frequency" db:"frequency" binding:"required"`
	TargetCount int       `json:"target_count" db:"target_count" binding:"required,min=1"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type CreateHabitRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=100"`
	Color       string `json:"color" binding:"required,hexcolor"`
	Icon        string `json:"icon" binding:"required,min=1,max=50"`
	Frequency   string `json:"frequency" binding:"required,oneof=daily weekly monthly"`
	TargetCount int    `json:"target_count" binding:"required,min=1,max=100"`
}

type UpdateHabitRequest struct {
	Name        *string `json:"name" binding:"omitempty,min=1,max=100"`
	Color       *string `json:"color" binding:"omitempty,hexcolor"`
	Icon        *string `json:"icon" binding:"omitempty,min=1,max=50"`
	Frequency   *string `json:"frequency" binding:"omitempty,oneof=daily weekly monthly"`
	TargetCount *int    `json:"target_count" binding:"omitempty,min=1,max=100"`
	IsActive    *bool   `json:"is_active"`
}

type HabitCompletion struct {
	ID          uuid.UUID `json:"id" db:"id"`
	HabitID     uuid.UUID `json:"habit_id" db:"habit_id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	CompletedAt time.Time `json:"completed_at" db:"completed_at"`
	Notes       string    `json:"notes" db:"notes"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

func (h *Habit) Validate() error {
	validFrequencies := map[string]bool{
		"daily":   true,
		"weekly":  true,
		"monthly": true,
	}

	if !validFrequencies[h.Frequency] {
		return ErrInvalidFrequency
	}

	if h.TargetCount < 1 {
		return ErrInvalidTargetCount
	}

	return nil
}
