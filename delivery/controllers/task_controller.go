package controllers

import (
	"fmt"
	"strconv"
	"task7/domain"
	services "task7/usecases"

	"github.com/gin-gonic/gin"
)

type TaskController struct {
	taskService services.TaskService
}

func NewTaskController(ts services.TaskService) *TaskController {
	return &TaskController{
		taskService: ts,
	}
}

func (t TaskController) GetAllTasks(c *gin.Context) {
	tasks, err := t.taskService.GetAllTasks()
	if err != nil {
		c.JSON(400, gin.H{"message": "Error getting documents"})
		return
	}
	c.JSON(200, tasks)
}

func (t TaskController) GetTasksById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid Task ID"})
		return
	}
	task, err := t.taskService.GetTaskById(id)
	if err != nil {
		c.JSON(404, gin.H{"message": "Task not found"})
		return
	}
	c.JSON(200, task)
}

func (t TaskController) PostTasks(c *gin.Context) {
	var newTask domain.Task
	err := c.BindJSON(&newTask)
	if err != nil {
		c.JSON(400, gin.H{"message": "Error binding JSON"})
		return
	}
	err = t.taskService.CreateTask(&newTask)
	if err != nil {
		fmt.Println(err)
		c.JSON(400, gin.H{"message": fmt.Sprintf("Error %v", err)})
		return
	}
	c.JSON(201, newTask)
}

func (t TaskController) PutTasksById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid Task ID"})
		return
	}
	var updatedTask domain.Task
	err = c.BindJSON(&updatedTask)
	if err != nil {
		c.JSON(400, gin.H{"message": "Error binding JSON"})
		return
	}
	err = t.taskService.UpdateTask(id, &updatedTask)
	if err != nil {
		c.JSON(404, gin.H{"message": "Error updating task"})
		return
	}
	c.JSON(200, gin.H{"message": "Task updated successfully"})
}

func (t TaskController) DeleteTaskById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid Task ID"})
		return
	}
	err = t.taskService.DeleteTaskById(id)
	if err != nil {
		c.JSON(404, gin.H{"message": "Error deleting task"})
		return
	}
	c.Status(204)
}
