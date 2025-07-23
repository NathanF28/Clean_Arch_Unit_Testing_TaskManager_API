package interfaces 

import (
	"task7/domain"
)

type TaskRepository interface {            // choose any db that implements register and login
	GetAllTasks() ([]domain.Task,error)
	GetTaskById(id int) (domain.Task,error)
	CreateTask(newTask *domain.Task)  error
	UpdateTask(id int, updatedTask *domain.Task)  error
	DeleteTaskById(id int) error
}


 




