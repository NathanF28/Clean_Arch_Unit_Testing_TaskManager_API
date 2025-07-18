package controllers

import (
	"fmt"
	"strconv"
	"task4/data"
	"task4/models"

	"github.com/gin-gonic/gin"
)

func GetTasks(c *gin.Context) {
	tasks, err := data.GetAllTasks()
	if err != nil {
		c.JSON(400, gin.H{"message": "Error getting documents"})
		return
	}
	c.JSON(200, tasks)
}

func GetTasksById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid Task ID"})
		return
	}
	task, err := data.GetTaskById(id)
	if err != nil {
		c.JSON(404, gin.H{"message": "Task not found"})
		return
	}
	c.JSON(200, task)
}

func PostTasks(c *gin.Context) {
	var newTask models.Task
	err := c.BindJSON(&newTask)
	if err != nil {
		c.JSON(400, gin.H{"message": "Error binding JSON"})
		return
	}
	err = data.CreateTask(&newTask)
	if err != nil {
		fmt.Println(err)
		c.JSON(400, gin.H{"message": fmt.Sprintf("Error %v", err)})
		return
	}
	c.JSON(201, newTask)
}

func PutTasksById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid Task ID"})
		return
	}
	var updatedTask models.Task
	err = c.BindJSON(&updatedTask)
	if err != nil {
		c.JSON(400, gin.H{"message": "Error binding JSON"})
		return
	}
	err = data.UpdateTask(id, &updatedTask)
	if err != nil {
		c.JSON(404, gin.H{"message": "Error updating task"})
		return
	}
	c.JSON(200, gin.H{"message": "Task updated successfully"})
}

func DeleteTaskById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid Task ID"})
		return
	}
	err = data.RemoveTasks(id)
	if err != nil {
		c.JSON(404, gin.H{"message": "Error deleting task"})
		return
	}
	c.Status(204)
}
