package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lumen/backend-go/internal/models"
	"github.com/lumen/backend-go/internal/repository"
	apperrors "github.com/lumen/backend-go/pkg/errors"
	"github.com/lumen/backend-go/pkg/logger"
	"go.uber.org/zap"
)

type TaskHandler struct {
	repo repository.TaskRepository
}

func NewTaskHandler(repo repository.TaskRepository) *TaskHandler {
	return &TaskHandler{repo: repo}
}

func (h *TaskHandler) Create(c *gin.Context) {
	var req models.CreateTaskRequest
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

	task := &models.Task{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Horizon:     req.Horizon,
		Priority:    req.Priority,
		DueDate:     req.DueDate,
	}

	if err := task.Validate(); err != nil {
		appErr := apperrors.NewValidationError(err.Error())
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	if err := h.repo.Create(c.Request.Context(), task); err != nil {
		logger.Error("Failed to create task", zap.Error(err), zap.String("user_id", userID.String()))
		appErr := apperrors.NewDatabaseError(err)
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	logger.Info("Task created", zap.String("task_id", task.ID.String()), zap.String("user_id", userID.String()))
	c.JSON(http.StatusCreated, task)
}

func (h *TaskHandler) GetAll(c *gin.Context) {
	userID := getUserID(c)
	if userID == uuid.Nil {
		appErr := apperrors.NewUnauthorized("user not authenticated")
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	var filter models.TaskFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		appErr := apperrors.NewBadRequest("invalid query parameters")
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	tasks, err := h.repo.GetByUserID(c.Request.Context(), userID, filter)
	if err != nil {
		logger.Error("Failed to get tasks", zap.Error(err), zap.String("user_id", userID.String()))
		appErr := apperrors.NewDatabaseError(err)
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	if tasks == nil {
		tasks = []models.Task{}
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   tasks,
		"count":  len(tasks),
		"filter": filter,
	})
}

func (h *TaskHandler) GetByID(c *gin.Context) {
	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		appErr := apperrors.NewBadRequest("invalid task ID")
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	userID := getUserID(c)
	if userID == uuid.Nil {
		appErr := apperrors.NewUnauthorized("user not authenticated")
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	task, err := h.repo.GetByID(c.Request.Context(), taskID, userID)
	if err == models.ErrNotFound {
		appErr := apperrors.NewNotFound("task")
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	if err != nil {
		logger.Error("Failed to get task", zap.Error(err), zap.String("task_id", taskID.String()))
		appErr := apperrors.NewDatabaseError(err)
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) Update(c *gin.Context) {
	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		appErr := apperrors.NewBadRequest("invalid task ID")
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	userID := getUserID(c)
	if userID == uuid.Nil {
		appErr := apperrors.NewUnauthorized("user not authenticated")
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	task, err := h.repo.GetByID(c.Request.Context(), taskID, userID)
	if err == models.ErrNotFound {
		appErr := apperrors.NewNotFound("task")
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	if err != nil {
		appErr := apperrors.NewDatabaseError(err)
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	var req models.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := apperrors.NewValidationError(err.Error())
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.Description != nil {
		task.Description = *req.Description
	}
	if req.Horizon != nil {
		task.Horizon = *req.Horizon
	}
	if req.Priority != nil {
		task.Priority = *req.Priority
	}
	if req.Status != nil {
		task.Status = *req.Status
	}
	if req.DueDate != nil {
		task.DueDate = req.DueDate
	}

	if err := task.Validate(); err != nil {
		appErr := apperrors.NewValidationError(err.Error())
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	if err := h.repo.Update(c.Request.Context(), task); err != nil {
		logger.Error("Failed to update task", zap.Error(err), zap.String("task_id", taskID.String()))
		appErr := apperrors.NewDatabaseError(err)
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	logger.Info("Task updated", zap.String("task_id", taskID.String()), zap.String("user_id", userID.String()))
	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) Delete(c *gin.Context) {
	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		appErr := apperrors.NewBadRequest("invalid task ID")
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	userID := getUserID(c)
	if userID == uuid.Nil {
		appErr := apperrors.NewUnauthorized("user not authenticated")
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	if err := h.repo.Delete(c.Request.Context(), taskID, userID); err == models.ErrNotFound {
		appErr := apperrors.NewNotFound("task")
		c.JSON(appErr.StatusCode, appErr)
		return
	} else if err != nil {
		logger.Error("Failed to delete task", zap.Error(err), zap.String("task_id", taskID.String()))
		appErr := apperrors.NewDatabaseError(err)
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	logger.Info("Task deleted", zap.String("task_id", taskID.String()), zap.String("user_id", userID.String()))
	c.JSON(http.StatusNoContent, nil)
}
