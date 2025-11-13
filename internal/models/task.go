package models

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	UserID      uuid.UUID  `json:"user_id" db:"user_id"`
	Title       string     `json:"title" db:"title" binding:"required"`
	Description string     `json:"description" db:"description"`
	Horizon     string     `json:"horizon" db:"horizon" binding:"required"`
	Priority    string     `json:"priority" db:"priority" binding:"required"`
	Status      string     `json:"status" db:"status" binding:"required"`
	DueDate     *time.Time `json:"due_date" db:"due_date"`
	CompletedAt *time.Time `json:"completed_at" db:"completed_at"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

type CreateTaskRequest struct {
	Title       string     `json:"title" binding:"required,min=1,max=200"`
	Description string     `json:"description" binding:"max=1000"`
	Horizon     string     `json:"horizon" binding:"required,oneof=now next later someday"`
	Priority    string     `json:"priority" binding:"required,oneof=low medium high urgent"`
	DueDate     *time.Time `json:"due_date"`
}

type UpdateTaskRequest struct {
	Title       *string    `json:"title" binding:"omitempty,min=1,max=200"`
	Description *string    `json:"description" binding:"omitempty,max=1000"`
	Horizon     *string    `json:"horizon" binding:"omitempty,oneof=now next later someday"`
	Priority    *string    `json:"priority" binding:"omitempty,oneof=low medium high urgent"`
	Status      *string    `json:"status" binding:"omitempty,oneof=todo in_progress done archived"`
	DueDate     *time.Time `json:"due_date"`
}

type TaskFilter struct {
	Horizon  string     `form:"horizon"`
	Status   string     `form:"status"`
	Priority string     `form:"priority"`
	FromDate *time.Time `form:"from_date"`
	ToDate   *time.Time `form:"to_date"`
}

func (t *Task) Validate() error {
	validHorizons := map[string]bool{
		"now":     true,
		"next":    true,
		"later":   true,
		"someday": true,
	}

	validPriorities := map[string]bool{
		"low":    true,
		"medium": true,
		"high":   true,
		"urgent": true,
	}

	validStatuses := map[string]bool{
		"todo":        true,
		"in_progress": true,
		"done":        true,
		"archived":    true,
	}

	if !validHorizons[t.Horizon] {
		return ErrInvalidHorizon
	}

	if !validPriorities[t.Priority] {
		return ErrInvalidPriority
	}

	if !validStatuses[t.Status] {
		return ErrInvalidStatus
	}

	return nil
}
