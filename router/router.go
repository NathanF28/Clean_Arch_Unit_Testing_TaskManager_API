package router

import (
	"task4/controllers"

	"github.com/gin-gonic/gin"
)

func StartServer() {
	router := gin.Default()
	r := router.Group("tasks") // best practice to group routes
	{
		r.GET("", controllers.GetTasks)
		r.POST("", controllers.PostTasks)
		r.GET("/:id", controllers.GetTasksById)
		r.PUT("/:id", controllers.PutTasksById)
		r.DELETE("/:id", controllers.DeleteTaskById)
	}
	router.Run("localhost:8080")
}
