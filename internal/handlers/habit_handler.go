package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lumen/backend/internal/models"
	"github.com/lumen/backend/internal/repository"
	apperrors "github.com/lumen/backend/pkg/errors"
	"github.com/lumen/backend/pkg/logger"
	"go.uber.org/zap"
)

type HabitHandler struct {
	repo repository.HabitRepository
}

func NewHabitHandler(repo repository.HabitRepository) *HabitHandler {
	return &HabitHandler{repo: repo}
}

func (h *HabitHandler) Create(c *gin.Context) {
	var req models.CreateHabitRequest
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

	habit := &models.Habit{
		UserID:      userID,
		Name:        req.Name,
		Color:       req.Color,
		Icon:        req.Icon,
		Frequency:   req.Frequency,
		TargetCount: req.TargetCount,
	}

	if err := habit.Validate(); err != nil {
		appErr := apperrors.NewValidationError(err.Error())
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	if err := h.repo.Create(c.Request.Context(), habit); err != nil {
		logger.Error("Failed to create habit", zap.Error(err), zap.String("user_id", userID.String()))
		appErr := apperrors.NewDatabaseError(err)
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	logger.Info("Habit created", zap.String("habit_id", habit.ID.String()), zap.String("user_id", userID.String()))
	c.JSON(http.StatusCreated, habit)
}

func (h *HabitHandler) GetAll(c *gin.Context) {
	userID := getUserID(c)
	if userID == uuid.Nil {
		appErr := apperrors.NewUnauthorized("user not authenticated")
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	habits, err := h.repo.GetByUserID(c.Request.Context(), userID)
	if err != nil {
		logger.Error("Failed to get habits", zap.Error(err), zap.String("user_id", userID.String()))
		appErr := apperrors.NewDatabaseError(err)
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	if habits == nil {
		habits = []models.Habit{}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": habits,
		"count": len(habits),
	})
}

func (h *HabitHandler) GetByID(c *gin.Context) {
	habitID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		appErr := apperrors.NewBadRequest("invalid habit ID")
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	userID := getUserID(c)
	if userID == uuid.Nil {
		appErr := apperrors.NewUnauthorized("user not authenticated")
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	habit, err := h.repo.GetByID(c.Request.Context(), habitID, userID)
	if err == models.ErrNotFound {
		appErr := apperrors.NewNotFound("habit")
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	if err != nil {
		logger.Error("Failed to get habit", zap.Error(err), zap.String("habit_id", habitID.String()))
		appErr := apperrors.NewDatabaseError(err)
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	c.JSON(http.StatusOK, habit)
}

func (h *HabitHandler) Update(c *gin.Context) {
	habitID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		appErr := apperrors.NewBadRequest("invalid habit ID")
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	userID := getUserID(c)
	if userID == uuid.Nil {
		appErr := apperrors.NewUnauthorized("user not authenticated")
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	habit, err := h.repo.GetByID(c.Request.Context(), habitID, userID)
	if err == models.ErrNotFound {
		appErr := apperrors.NewNotFound("habit")
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	if err != nil {
		appErr := apperrors.NewDatabaseError(err)
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	var req models.UpdateHabitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := apperrors.NewValidationError(err.Error())
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	if req.Name != nil {
		habit.Name = *req.Name
	}
	if req.Color != nil {
		habit.Color = *req.Color
	}
	if req.Icon != nil {
		habit.Icon = *req.Icon
	}
	if req.Frequency != nil {
		habit.Frequency = *req.Frequency
	}
	if req.TargetCount != nil {
		habit.TargetCount = *req.TargetCount
	}
	if req.IsActive != nil {
		habit.IsActive = *req.IsActive
	}

	if err := habit.Validate(); err != nil {
		appErr := apperrors.NewValidationError(err.Error())
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	if err := h.repo.Update(c.Request.Context(), habit); err != nil {
		logger.Error("Failed to update habit", zap.Error(err), zap.String("habit_id", habitID.String()))
		appErr := apperrors.NewDatabaseError(err)
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	logger.Info("Habit updated", zap.String("habit_id", habitID.String()), zap.String("user_id", userID.String()))
	c.JSON(http.StatusOK, habit)
}

func (h *HabitHandler) Delete(c *gin.Context) {
	habitID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		appErr := apperrors.NewBadRequest("invalid habit ID")
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	userID := getUserID(c)
	if userID == uuid.Nil {
		appErr := apperrors.NewUnauthorized("user not authenticated")
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	if err := h.repo.Delete(c.Request.Context(), habitID, userID); err == models.ErrNotFound {
		appErr := apperrors.NewNotFound("habit")
		c.JSON(appErr.StatusCode, appErr)
		return
	} else if err != nil {
		logger.Error("Failed to delete habit", zap.Error(err), zap.String("habit_id", habitID.String()))
		appErr := apperrors.NewDatabaseError(err)
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	logger.Info("Habit deleted", zap.String("habit_id", habitID.String()), zap.String("user_id", userID.String()))
	c.JSON(http.StatusNoContent, nil)
}
