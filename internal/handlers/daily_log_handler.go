package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lumen/backend-go/internal/models"
	"github.com/lumen/backend-go/internal/repository"
	apperrors "github.com/lumen/backend-go/pkg/errors"
	"github.com/lumen/backend-go/pkg/logger"
	"go.uber.org/zap"
)

type DailyLogHandler struct {
	repo repository.DailyLogRepository
}

func NewDailyLogHandler(repo repository.DailyLogRepository) *DailyLogHandler {
	return &DailyLogHandler{repo: repo}
}

func (h *DailyLogHandler) Create(c *gin.Context) {
	var req models.CreateDailyLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := apperrors.NewValidationError(err.Error())
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	userID := getUserID(c)
	if userID == uuid.Nil {
		appErr := apperrors.NewUnauthorized("user not authenticated")
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	log := &models.DailyLog{
		UserID:             userID,
		Date:               req.Date,
		MorningRoutine:     req.MorningRoutine,
		EveningRoutine:     req.EveningRoutine,
		WaterIntake:        req.WaterIntake,
		SleepHours:         req.SleepHours,
		EnergyLevel:        req.EnergyLevel,
		MoodRating:         req.MoodRating,
		ProductivityRating: req.ProductivityRating,
		Notes:              req.Notes,
	}

	if err := log.Validate(); err != nil {
		appErr := apperrors.NewValidationError(err.Error())
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	if err := h.repo.Create(c.Request.Context(), log); err != nil {
		logger.Error("Failed to create daily log", zap.Error(err), zap.String("user_id", userID.String()))
		appErr := apperrors.NewDatabaseError(err)
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	logger.Info("Daily log created", zap.String("log_id", log.ID.String()), zap.String("user_id", userID.String()))
	c.JSON(http.StatusCreated, log)
}

func (h *DailyLogHandler) GetByDate(c *gin.Context) {
	dateStr := c.Param("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		appErr := apperrors.NewBadRequest("invalid date format, use YYYY-MM-DD")
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	userID := getUserID(c)
	if userID == uuid.Nil {
		appErr := apperrors.NewUnauthorized("user not authenticated")
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	log, err := h.repo.GetByDate(c.Request.Context(), userID, date)
	if err == models.ErrNotFound {
		appErr := apperrors.NewNotFound("daily log")
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	if err != nil {
		logger.Error("Failed to get daily log", zap.Error(err), zap.String("date", dateStr))
		appErr := apperrors.NewDatabaseError(err)
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	c.JSON(http.StatusOK, log)
}

func (h *DailyLogHandler) GetRange(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		appErr := apperrors.NewBadRequest("start_date and end_date query parameters are required")
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		appErr := apperrors.NewBadRequest("invalid start_date format, use YYYY-MM-DD")
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		appErr := apperrors.NewBadRequest("invalid end_date format, use YYYY-MM-DD")
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	userID := getUserID(c)
	if userID == uuid.Nil {
		appErr := apperrors.NewUnauthorized("user not authenticated")
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	logs, err := h.repo.GetByDateRange(c.Request.Context(), userID, startDate, endDate)
	if err != nil {
		logger.Error("Failed to get daily logs", zap.Error(err), zap.String("user_id", userID.String()))
		appErr := apperrors.NewDatabaseError(err)
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	if logs == nil {
		logs = []models.DailyLog{}
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       logs,
		"count":      len(logs),
		"start_date": startDateStr,
		"end_date":   endDateStr,
	})
}

func (h *DailyLogHandler) Update(c *gin.Context) {
	dateStr := c.Param("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		appErr := apperrors.NewBadRequest("invalid date format, use YYYY-MM-DD")
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	userID := getUserID(c)
	if userID == uuid.Nil {
		appErr := apperrors.NewUnauthorized("user not authenticated")
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	log, err := h.repo.GetByDate(c.Request.Context(), userID, date)
	if err == models.ErrNotFound {
		appErr := apperrors.NewNotFound("daily log")
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	if err != nil {
		appErr := apperrors.NewDatabaseError(err)
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	var req models.UpdateDailyLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := apperrors.NewValidationError(err.Error())
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	if req.MorningRoutine != nil {
		log.MorningRoutine = *req.MorningRoutine
	}
	if req.EveningRoutine != nil {
		log.EveningRoutine = *req.EveningRoutine
	}
	if req.WaterIntake != nil {
		log.WaterIntake = *req.WaterIntake
	}
	if req.SleepHours != nil {
		log.SleepHours = *req.SleepHours
	}
	if req.EnergyLevel != nil {
		log.EnergyLevel = *req.EnergyLevel
	}
	if req.MoodRating != nil {
		log.MoodRating = *req.MoodRating
	}
	if req.ProductivityRating != nil {
		log.ProductivityRating = *req.ProductivityRating
	}
	if req.Notes != nil {
		log.Notes = *req.Notes
	}

	if err := log.Validate(); err != nil {
		appErr := apperrors.NewValidationError(err.Error())
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	if err := h.repo.Update(c.Request.Context(), log); err != nil {
		logger.Error("Failed to update daily log", zap.Error(err), zap.String("date", dateStr))
		appErr := apperrors.NewDatabaseError(err)
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	logger.Info("Daily log updated", zap.String("log_id", log.ID.String()), zap.String("user_id", userID.String()))
	c.JSON(http.StatusOK, log)
}

func getUserID(c *gin.Context) uuid.UUID {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil
	}

	if userID, ok := userIDStr.(uuid.UUID); ok {
		return userID
	}

	if userID, ok := userIDStr.(string); ok {
		parsed, err := uuid.Parse(userID)
		if err != nil {
			return uuid.Nil
		}
		return parsed
	}

	return uuid.Nil
}
