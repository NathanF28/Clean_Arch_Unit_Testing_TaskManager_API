package services

import (
	"task7/domain"
	"task7/repository/interfaces"
)

type TaskService interface {
	GetAllTasks() ([]domain.Task, error)
	GetTaskById(id int) (domain.Task, error)
	CreateTask(newTask *domain.Task) error
	UpdateTask(id int, updatedTask *domain.Task ) error
	DeleteTaskById(id int) error
}

type taskService struct {
	taskRepo interfaces.TaskRepository
}

func NewTaskService(tr interfaces.TaskRepository) TaskService {
	return &taskService{
		taskRepo: tr,
	}
}

func (s *taskService) GetAllTasks() ([]domain.Task, error) {
	return s.taskRepo.GetAllTasks()
}

func (s *taskService) GetTaskById(id int) (domain.Task, error) {
	return s.taskRepo.GetTaskById(id)
}

func (s *taskService) CreateTask(newTask *domain.Task) error {
	return s.taskRepo.CreateTask(newTask)
}

func (s *taskService) UpdateTask(id int, updatedTask *domain.Task) error {
	return s.taskRepo.UpdateTask(id,updatedTask)
}

func (s *taskService) DeleteTaskById(id int) error {
	return s.taskRepo.DeleteTaskById(id)
}
