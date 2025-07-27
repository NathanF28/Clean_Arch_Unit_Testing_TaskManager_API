package router

import (
	"task7/delivery/controllers"
	"task7/infrastructure"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	authController *controllers.AuthController,
	taskController *controllers.TaskController,
) *gin.Engine {
	router := gin.Default()
	router.POST("/register", authController.RegisterUser)
	router.POST("/login", authController.LoginUser)

	router.PUT("/promote", infrastructure.AuthMiddleware(), infrastructure.AdminAuth(), authController.PromoteUser)

	r := router.Group("/tasks")
	r.Use(infrastructure.AuthMiddleware())
	{
		r.GET("", taskController.GetAllTasks)
		r.GET("/:id", taskController.GetTasksById)
		r.POST("", infrastructure.AdminAuth(), taskController.PostTasks)
		r.PUT("/:id", infrastructure.AdminAuth(), taskController.PutTasksById)
		r.DELETE("/:id", infrastructure.AdminAuth(), taskController.DeleteTaskById)
	}
	return router
}
