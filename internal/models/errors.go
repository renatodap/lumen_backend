package models

import "errors"

var (
	ErrInvalidFrequency    = errors.New("invalid frequency: must be daily, weekly, or monthly")
	ErrInvalidTargetCount  = errors.New("invalid target count: must be at least 1")
	ErrInvalidHorizon      = errors.New("invalid horizon: must be now, next, later, or someday")
	ErrInvalidPriority     = errors.New("invalid priority: must be low, medium, high, or urgent")
	ErrInvalidStatus       = errors.New("invalid status: must be todo, in_progress, done, or archived")
	ErrInvalidWaterIntake  = errors.New("invalid water intake: must be between 0 and 20")
	ErrInvalidSleepHours   = errors.New("invalid sleep hours: must be between 0 and 24")
	ErrInvalidRating       = errors.New("invalid rating: must be between 1 and 5")
	ErrNotFound            = errors.New("resource not found")
	ErrUnauthorized        = errors.New("unauthorized access")
	ErrForbidden           = errors.New("forbidden: insufficient permissions")
	ErrConflict            = errors.New("resource conflict")
	ErrInternalServer      = errors.New("internal server error")
	ErrBadRequest          = errors.New("bad request")
	ErrValidationFailed    = errors.New("validation failed")
	ErrDatabaseConnection  = errors.New("database connection error")
	ErrDatabaseQuery       = errors.New("database query error")
)
