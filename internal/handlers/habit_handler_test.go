package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockHabitRepository is a mock for habit repository
type MockHabitRepository struct {
	mock.Mock
}

func (m *MockHabitRepository) GetByUserID(userID string) ([]Habit, error) {
	args := m.Called(userID)
	return args.Get(0).([]Habit), args.Error(1)
}

func (m *MockHabitRepository) GetByID(id string) (*Habit, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Habit), args.Error(1)
}

func (m *MockHabitRepository) Create(habit *Habit) error {
	args := m.Called(habit)
	return args.Error(0)
}

func (m *MockHabitRepository) Update(habit *Habit) error {
	args := m.Called(habit)
	return args.Error(0)
}

func (m *MockHabitRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// Habit model for testing
type Habit struct {
	ID            string   `json:"id"`
	UserID        string   `json:"user_id"`
	GoalID        *string  `json:"goal_id"`
	Name          string   `json:"name"`
	Frequency     string   `json:"frequency"`
	ReminderTimes []string `json:"reminder_times"`
	Icon          *string  `json:"icon"`
	CreatedAt     string   `json:"created_at"`
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

func TestGetHabits_Success(t *testing.T) {
	router := setupTestRouter()
	mockRepo := new(MockHabitRepository)

	expectedHabits := []Habit{
		{
			ID:            "1",
			UserID:        "user-1",
			Name:          "Morning Exercise",
			Frequency:     "daily",
			ReminderTimes: []string{"09:00"},
			CreatedAt:     "2025-01-01T00:00:00Z",
		},
	}

	mockRepo.On("GetByUserID", "user-1").Return(expectedHabits, nil)

	router.GET("/habits", func(c *gin.Context) {
		// Mock authentication middleware
		c.Set("user_id", "user-1")

		habits, err := mockRepo.GetByUserID("user-1")
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to fetch habits"})
			return
		}
		c.JSON(200, habits)
	})

	req, _ := http.NewRequest("GET", "/habits", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response []Habit
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(response))
	assert.Equal(t, "Morning Exercise", response[0].Name)

	mockRepo.AssertExpectations(t)
}

func TestGetHabits_EmptyList(t *testing.T) {
	router := setupTestRouter()
	mockRepo := new(MockHabitRepository)

	mockRepo.On("GetByUserID", "user-1").Return([]Habit{}, nil)

	router.GET("/habits", func(c *gin.Context) {
		c.Set("user_id", "user-1")

		habits, err := mockRepo.GetByUserID("user-1")
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to fetch habits"})
			return
		}
		c.JSON(200, habits)
	})

	req, _ := http.NewRequest("GET", "/habits", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response []Habit
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(response))

	mockRepo.AssertExpectations(t)
}

func TestCreateHabit_Success(t *testing.T) {
	router := setupTestRouter()
	mockRepo := new(MockHabitRepository)

	newHabit := Habit{
		Name:          "Read Daily",
		Frequency:     "daily",
		ReminderTimes: []string{"20:00"},
	}

	mockRepo.On("Create", mock.AnythingOfType("*handlers.Habit")).Return(nil)

	router.POST("/habits", func(c *gin.Context) {
		c.Set("user_id", "user-1")

		var habit Habit
		if err := c.ShouldBindJSON(&habit); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		habit.UserID = "user-1"
		habit.ID = "generated-id"
		habit.CreatedAt = "2025-01-13T00:00:00Z"

		err := mockRepo.Create(&habit)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to create habit"})
			return
		}

		c.JSON(201, habit)
	})

	body, _ := json.Marshal(newHabit)
	req, _ := http.NewRequest("POST", "/habits", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)

	var response Habit
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Read Daily", response.Name)
	assert.Equal(t, "user-1", response.UserID)
	assert.NotEmpty(t, response.ID)

	mockRepo.AssertExpectations(t)
}

func TestCreateHabit_ValidationError(t *testing.T) {
	router := setupTestRouter()

	router.POST("/habits", func(c *gin.Context) {
		var habit Habit
		if err := c.ShouldBindJSON(&habit); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		if habit.Name == "" {
			c.JSON(400, gin.H{"error": "Name is required"})
			return
		}

		c.JSON(201, habit)
	})

	invalidHabit := Habit{
		Frequency: "daily",
	}

	body, _ := json.Marshal(invalidHabit)
	req, _ := http.NewRequest("POST", "/habits", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "required")
}

func TestUpdateHabit_Success(t *testing.T) {
	router := setupTestRouter()
	mockRepo := new(MockHabitRepository)

	existingHabit := &Habit{
		ID:        "1",
		UserID:    "user-1",
		Name:      "Old Name",
		Frequency: "daily",
	}

	mockRepo.On("GetByID", "1").Return(existingHabit, nil)
	mockRepo.On("Update", mock.AnythingOfType("*handlers.Habit")).Return(nil)

	router.PATCH("/habits/:id", func(c *gin.Context) {
		c.Set("user_id", "user-1")
		id := c.Param("id")

		habit, err := mockRepo.GetByID(id)
		if err != nil {
			c.JSON(404, gin.H{"error": "Habit not found"})
			return
		}

		var updates map[string]interface{}
		if err := c.ShouldBindJSON(&updates); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		if name, ok := updates["name"].(string); ok {
			habit.Name = name
		}

		err = mockRepo.Update(habit)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to update habit"})
			return
		}

		c.JSON(200, habit)
	})

	updates := map[string]string{"name": "Updated Name"}
	body, _ := json.Marshal(updates)
	req, _ := http.NewRequest("PATCH", "/habits/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response Habit
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Name", response.Name)

	mockRepo.AssertExpectations(t)
}

func TestDeleteHabit_Success(t *testing.T) {
	router := setupTestRouter()
	mockRepo := new(MockHabitRepository)

	mockRepo.On("Delete", "1").Return(nil)

	router.DELETE("/habits/:id", func(c *gin.Context) {
		c.Set("user_id", "user-1")
		id := c.Param("id")

		err := mockRepo.Delete(id)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to delete habit"})
			return
		}

		c.JSON(204, nil)
	})

	req, _ := http.NewRequest("DELETE", "/habits/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 204, w.Code)

	mockRepo.AssertExpectations(t)
}

func TestDeleteHabit_NotFound(t *testing.T) {
	router := setupTestRouter()
	mockRepo := new(MockHabitRepository)

	mockRepo.On("Delete", "999").Return(assert.AnError)

	router.DELETE("/habits/:id", func(c *gin.Context) {
		id := c.Param("id")

		err := mockRepo.Delete(id)
		if err != nil {
			c.JSON(404, gin.H{"error": "Habit not found"})
			return
		}

		c.JSON(204, nil)
	})

	req, _ := http.NewRequest("DELETE", "/habits/999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)

	mockRepo.AssertExpectations(t)
}
