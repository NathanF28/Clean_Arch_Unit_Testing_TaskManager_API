package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"task7/delivery/controllers"
	"task7/domain"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockTaskService struct {
	mock.Mock
}

func (m *MockTaskService) GetAllTasks() ([]domain.Task, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Task), args.Error(1)
}

func (m *MockTaskService) GetTaskById(id int) (domain.Task, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return domain.Task{}, args.Error(1)
	}
	return args.Get(0).(domain.Task), args.Error(1)
}

func (m *MockTaskService) CreateTask(newTask *domain.Task) error {
	args := m.Called(newTask)
	return args.Error(0)
}

func (m *MockTaskService) UpdateTask(id int, updatedTask *domain.Task) error {
	args := m.Called(id, updatedTask)
	return args.Error(0)
}

func (m *MockTaskService) DeleteTaskById(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

type TaskControllerSuite struct {
	suite.Suite
	router          *gin.Engine
	mockTaskService *MockTaskService
}

func (s *TaskControllerSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	s.mockTaskService = new(MockTaskService)
	taskController := controllers.NewTaskController(s.mockTaskService)

	s.router = gin.New()
	s.router.GET("/tasks", taskController.GetAllTasks)
	s.router.GET("/tasks/:id", taskController.GetTasksById)
	s.router.POST("/tasks", taskController.PostTasks)
	s.router.PUT("/tasks/:id", taskController.PutTasksById)
	s.router.DELETE("/tasks/:id", taskController.DeleteTaskById)
}

func TestTaskControllerSuite(t *testing.T) {
	suite.Run(t, new(TaskControllerSuite))
}

func (s *TaskControllerSuite) performRequest(method, path string, body interface{}) *httptest.ResponseRecorder {
	var reqBody []byte
	if body != nil {
		reqBody, _ = json.Marshal(body)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, path, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)
	return w
}

func (s *TaskControllerSuite) TestGetAllTasks_SuccessWithTasks() {
	const timePrecision = time.Millisecond

	expectedTasks := []domain.Task{
		{ID: 1, Title: "Task 1", Description: "Desc 1", DueDate: time.Now().Truncate(timePrecision), Status: "pending"},
		{ID: 2, Title: "Task 2", Description: "Desc 2", DueDate: time.Now().Truncate(timePrecision), Status: "completed"},
	}

	s.mockTaskService.On("GetAllTasks").Return(expectedTasks, nil).Once()

	w := s.performRequest("GET", "/tasks", nil)

	s.Equal(http.StatusOK, w.Code)
	var tasks []domain.Task
	err := json.Unmarshal(w.Body.Bytes(), &tasks)
	s.NoError(err)

	for i := range tasks {
		tasks[i].DueDate = tasks[i].DueDate.Truncate(timePrecision)
	}

	s.Equal(expectedTasks, tasks)
	s.mockTaskService.AssertExpectations(s.T())
}

func (s *TaskControllerSuite) TestGetAllTasks_SuccessNoTasks() {
	expectedTasks := []domain.Task{}

	s.mockTaskService.On("GetAllTasks").Return(expectedTasks, nil).Once()

	w := s.performRequest("GET", "/tasks", nil)

	s.Equal(http.StatusOK, w.Code)
	var tasks []domain.Task
	err := json.Unmarshal(w.Body.Bytes(), &tasks)
	s.NoError(err)
	s.Empty(tasks)
	s.mockTaskService.AssertExpectations(s.T())
}

func (s *TaskControllerSuite) TestGetAllTasks_ServiceError() {
	serviceError := errors.New("database connection failed")

	s.mockTaskService.On("GetAllTasks").Return(nil, serviceError).Once()

	w := s.performRequest("GET", "/tasks", nil)

	s.Equal(http.StatusBadRequest, w.Code)
	s.Contains(w.Body.String(), `{"message":"Error getting documents"}`)
	s.mockTaskService.AssertExpectations(s.T())
}

func (s *TaskControllerSuite) TestGetTasksById_Success() {
	const timePrecision = time.Millisecond
	expectedTask := domain.Task{ID: 1, Title: "Test Task", Description: "Desc", DueDate: time.Now().Truncate(timePrecision), Status: "pending"}

	s.mockTaskService.On("GetTaskById", 1).Return(expectedTask, nil).Once()

	w := s.performRequest("GET", "/tasks/1", nil)

	s.Equal(http.StatusOK, w.Code)
	var task domain.Task
	err := json.Unmarshal(w.Body.Bytes(), &task)
	s.NoError(err)
	task.DueDate = task.DueDate.Truncate(timePrecision)
	s.Equal(expectedTask, task)
	s.mockTaskService.AssertExpectations(s.T())
}

func (s *TaskControllerSuite) TestGetTasksById_InvalidID() {
	w := s.performRequest("GET", "/tasks/abc", nil)

	s.Equal(http.StatusBadRequest, w.Code)
	s.Contains(w.Body.String(), `{"message":"Invalid Task ID"}`)
	s.mockTaskService.AssertNotCalled(s.T(), "GetTaskById", mock.Anything)
}

func (s *TaskControllerSuite) TestGetTasksById_NotFound() {
	serviceError := errors.New("task not found")

	s.mockTaskService.On("GetTaskById", 999).Return(domain.Task{}, serviceError).Once()

	w := s.performRequest("GET", "/tasks/999", nil)

	s.Equal(http.StatusNotFound, w.Code)
	s.Contains(w.Body.String(), `{"message":"Task not found"}`)
	s.mockTaskService.AssertExpectations(s.T())
}

func (s *TaskControllerSuite) TestPostTasks_Success() {
	const timePrecision = time.Millisecond
	newTask := domain.Task{Title: "New Task", Description: "Details", DueDate: time.Now().Add(24 * time.Hour).Truncate(timePrecision), Status: "pending"}

	s.mockTaskService.On("CreateTask", mock.AnythingOfType("*domain.Task")).Return(nil).Once()

	w := s.performRequest("POST", "/tasks", newTask)

	s.Equal(http.StatusCreated, w.Code)
	var createdTask domain.Task
	err := json.Unmarshal(w.Body.Bytes(), &createdTask)
	s.NoError(err)

	createdTask.DueDate = createdTask.DueDate.Truncate(timePrecision)
	newTask.ID = createdTask.ID // Assume ID is set by the service/DB, so match it for comparison
	s.Equal(newTask.Title, createdTask.Title)
	s.Equal(newTask.Description, createdTask.Description)
	s.Equal(newTask.Status, createdTask.Status)
	s.Equal(newTask.DueDate, createdTask.DueDate) // Now compare DueDate directly after truncation

	s.mockTaskService.AssertExpectations(s.T())
}

func (s *TaskControllerSuite) TestPostTasks_InvalidJSON() {
	w := s.performRequest("POST", "/tasks", "not json")

	s.Equal(http.StatusBadRequest, w.Code)
	s.Contains(w.Body.String(), `{"message":"Error binding JSON"}`)
	s.mockTaskService.AssertNotCalled(s.T(), "CreateTask", mock.Anything)
}

func (s *TaskControllerSuite) TestPostTasks_ServiceError() {
	newTask := domain.Task{Title: "Failed Task", Description: "Will fail", DueDate: time.Now(), Status: "pending"}
	serviceError := errors.New("database write error")

	s.mockTaskService.On("CreateTask", mock.AnythingOfType("*domain.Task")).Return(serviceError).Once()

	w := s.performRequest("POST", "/tasks", newTask)

	s.Equal(http.StatusBadRequest, w.Code)
	s.Contains(w.Body.String(), fmt.Sprintf(`{"message":"Error %v"}`, serviceError))
	s.mockTaskService.AssertExpectations(s.T())
}

func (s *TaskControllerSuite) TestPutTasksById_Success() {
	const timePrecision = time.Millisecond
	updatedTask := domain.Task{ID: 1, Title: "Updated Task", Description: "New Desc", DueDate: time.Now().Truncate(timePrecision), Status: "completed"}

	s.mockTaskService.On("UpdateTask", 1, mock.AnythingOfType("*domain.Task")).Return(nil).Once()

	w := s.performRequest("PUT", "/tasks/1", updatedTask)

	s.Equal(http.StatusOK, w.Code)
	s.Contains(w.Body.String(), `{"message":"Task updated successfully"}`)
	s.mockTaskService.AssertExpectations(s.T())
}

func (s *TaskControllerSuite) TestPutTasksById_InvalidID() {
	updatedTask := domain.Task{Title: "Invalid ID", Status: "pending"}

	w := s.performRequest("PUT", "/tasks/xyz", updatedTask)

	s.Equal(http.StatusBadRequest, w.Code)
	s.Contains(w.Body.String(), `{"message":"Invalid Task ID"}`)
	s.mockTaskService.AssertNotCalled(s.T(), "UpdateTask", mock.Anything, mock.Anything)
}

func (s *TaskControllerSuite) TestPutTasksById_InvalidJSON() {
	w := s.performRequest("PUT", "/tasks/1", "not json")

	s.Equal(http.StatusBadRequest, w.Code)
	s.Contains(w.Body.String(), `{"message":"Error binding JSON"}`)
	s.mockTaskService.AssertNotCalled(s.T(), "UpdateTask", mock.Anything, mock.Anything)
}

func (s *TaskControllerSuite) TestPutTasksById_ServiceError() {
	updatedTask := domain.Task{ID: 1, Title: "Updated Task", Description: "New Desc", DueDate: time.Now(), Status: "completed"}
	serviceError := errors.New("task not found for update")

	s.mockTaskService.On("UpdateTask", 1, mock.AnythingOfType("*domain.Task")).Return(serviceError).Once()

	w := s.performRequest("PUT", "/tasks/1", updatedTask)

	s.Equal(http.StatusNotFound, w.Code)
	s.Contains(w.Body.String(), `{"message":"Error updating task"}`)
	s.mockTaskService.AssertExpectations(s.T())
}

func (s *TaskControllerSuite) TestDeleteTaskById_Success() {
	s.mockTaskService.On("DeleteTaskById", 1).Return(nil).Once()

	w := s.performRequest("DELETE", "/tasks/1", nil)

	s.Equal(http.StatusNoContent, w.Code)
	s.Empty(w.Body.String())
	s.mockTaskService.AssertExpectations(s.T())
}

func (s *TaskControllerSuite) TestDeleteTaskById_InvalidID() {
	w := s.performRequest("DELETE", "/tasks/abc", nil)

	s.Equal(http.StatusBadRequest, w.Code)
	s.Contains(w.Body.String(), `{"message":"Invalid Task ID"}`)
	s.mockTaskService.AssertNotCalled(s.T(), "DeleteTaskById", mock.Anything)
}

func (s *TaskControllerSuite) TestDeleteTaskById_ServiceError() {
	serviceError := errors.New("task not found for deletion")

	s.mockTaskService.On("DeleteTaskById", 999).Return(serviceError).Once()

	w := s.performRequest("DELETE", "/tasks/999", nil)

	s.Equal(http.StatusNotFound, w.Code)
	s.Contains(w.Body.String(), `{"message":"Error deleting task"}`)
	s.mockTaskService.AssertExpectations(s.T())
}
