package services_test

import (
	"errors"
	"testing"
	"time"

	"task7/domain"
	services "task7/usecases"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) GetAllTasks() ([]domain.Task, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Task), args.Error(1)
}

func (m *MockTaskRepository) GetTaskById(id int) (domain.Task, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return domain.Task{}, args.Error(1)
	}
	return args.Get(0).(domain.Task), args.Error(1)
}

func (m *MockTaskRepository) CreateTask(newTask *domain.Task) error {
	args := m.Called(newTask)
	return args.Error(0)
}

func (m *MockTaskRepository) UpdateTask(id int, updatedTask *domain.Task) error {
	args := m.Called(id, updatedTask)
	return args.Error(0)
}

func (m *MockTaskRepository) DeleteTaskById(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

type TaskServiceSuite struct {
	suite.Suite
	mockRepo    *MockTaskRepository
	taskService services.TaskService
}

func (s *TaskServiceSuite) SetupTest() {
	s.mockRepo = new(MockTaskRepository)
	s.taskService = services.NewTaskService(s.mockRepo)
}

func TestTaskServiceSuite(t *testing.T) {
	suite.Run(t, new(TaskServiceSuite))
}

func (s *TaskServiceSuite) TestGetAllTasks_SuccessWithTasks() {
	expectedTasks := []domain.Task{
		{ID: 1, Title: "Task 1", Description: "Desc 1", Status: "pending"},
		{ID: 2, Title: "Task 2", Description: "Desc 2", Status: "completed"},
	}

	s.mockRepo.On("GetAllTasks").Return(expectedTasks, nil).Once()

	tasks, err := s.taskService.GetAllTasks()
	s.NoError(err, "GetAllTasks should not return an error on success")
	assert.Len(s.T(), tasks, 2, "Should return 2 tasks")
	assert.Equal(s.T(), expectedTasks, tasks, "Returned tasks should match expected tasks")
	s.mockRepo.AssertExpectations(s.T())
}

func (s *TaskServiceSuite) TestGetAllTasks_SuccessNoTasks() {
	expectedTasks := []domain.Task{}

	s.mockRepo.On("GetAllTasks").Return(expectedTasks, nil).Once()

	tasks, err := s.taskService.GetAllTasks()
	s.NoError(err, "GetAllTasks should not return an error on success with no tasks")
	assert.Len(s.T(), tasks, 0, "Should return 0 tasks")
	assert.Equal(s.T(), expectedTasks, tasks, "Returned tasks should be empty slice")
	s.mockRepo.AssertExpectations(s.T())
}

func (s *TaskServiceSuite) TestGetAllTasks_RepositoryError() {
	repoError := errors.New("database error fetching tasks")

	s.mockRepo.On("GetAllTasks").Return(nil, repoError).Once()

	tasks, err := s.taskService.GetAllTasks()
	s.Error(err, "GetAllTasks should return an error when repository fails")
	s.Equal(repoError, err, "Error returned should be the repository error")
	s.Nil(tasks, "Tasks should be nil on error")
	s.mockRepo.AssertExpectations(s.T())
}

func (s *TaskServiceSuite) TestGetTaskById_Success() {
	expectedTask := domain.Task{ID: 1, Title: "Test Task", Description: "Description", Status: "pending"}

	s.mockRepo.On("GetTaskById", 1).Return(expectedTask, nil).Once()

	task, err := s.taskService.GetTaskById(1)
	s.NoError(err, "GetTaskById should not return an error on success")
	assert.Equal(s.T(), expectedTask, task, "Returned task should match expected task")
	s.mockRepo.AssertExpectations(s.T())
}

func (s *TaskServiceSuite) TestGetTaskById_NotFound() {
	repoError := errors.New("task not found")

	s.mockRepo.On("GetTaskById", 999).Return(domain.Task{}, repoError).Once()

	task, err := s.taskService.GetTaskById(999)
	s.Error(err, "GetTaskById should return an error when task is not found")
	s.Equal(repoError, err, "Error returned should indicate task not found")
	s.Equal(domain.Task{}, task, "Task should be empty on error")
	s.mockRepo.AssertExpectations(s.T())
}

func (s *TaskServiceSuite) TestCreateTask_Success() {
	newTask := &domain.Task{Title: "New Task", Description: "To be created", DueDate: time.Now(), Status: "pending"}

	s.mockRepo.On("CreateTask", newTask).Return(nil).Once()

	err := s.taskService.CreateTask(newTask)
	s.NoError(err, "CreateTask should not return an error on success")
	s.mockRepo.AssertExpectations(s.T())
}

func (s *TaskServiceSuite) TestCreateTask_RepositoryError() {
	newTask := &domain.Task{Title: "New Task", Description: "To be created", DueDate: time.Now(), Status: "pending"}
	repoError := errors.New("database creation failed")

	s.mockRepo.On("CreateTask", newTask).Return(repoError).Once()

	err := s.taskService.CreateTask(newTask)
	s.Error(err, "CreateTask should return an error when repository fails")
	s.Equal(repoError, err, "Error returned should be the repository error")
	s.mockRepo.AssertExpectations(s.T())
}

func (s *TaskServiceSuite) TestUpdateTask_Success() {
	updatedTask := &domain.Task{ID: 1, Title: "Updated Task", Status: "completed"}

	s.mockRepo.On("UpdateTask", 1, updatedTask).Return(nil).Once()

	err := s.taskService.UpdateTask(1, updatedTask)
	s.NoError(err, "UpdateTask should not return an error on success")
	s.mockRepo.AssertExpectations(s.T())
}

func (s *TaskServiceSuite) TestUpdateTask_NotFound() {
	updatedTask := &domain.Task{ID: 999, Title: "Non-existent", Status: "pending"}
	repoError := errors.New("task not found for update")

	s.mockRepo.On("UpdateTask", 999, updatedTask).Return(repoError).Once()

	err := s.taskService.UpdateTask(999, updatedTask)
	s.Error(err, "UpdateTask should return an error when task is not found")
	s.Equal(repoError, err, "Error returned should indicate task not found")
	s.mockRepo.AssertExpectations(s.T())
}

func (s *TaskServiceSuite) TestDeleteTaskById_Success() {
	s.mockRepo.On("DeleteTaskById", 1).Return(nil).Once()

	err := s.taskService.DeleteTaskById(1)
	s.NoError(err, "DeleteTaskById should not return an error on success")
	s.mockRepo.AssertExpectations(s.T())
}

func (s *TaskServiceSuite) TestDeleteTaskById_NotFound() {
	repoError := errors.New("task not found for deletion")

	s.mockRepo.On("DeleteTaskById", 999).Return(repoError).Once()

	err := s.taskService.DeleteTaskById(999)
	s.Error(err, "DeleteTaskById should return an error when task is not found")
	s.Equal(repoError, err, "Error returned should indicate task not found")
	s.mockRepo.AssertExpectations(s.T())
}
