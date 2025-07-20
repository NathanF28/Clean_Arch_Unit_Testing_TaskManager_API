package router

import (
    "task6/controllers"
    "task6/middleware"
    "github.com/gin-gonic/gin"
)

func StartServer() {
    router := gin.Default()

    router.POST("/register", controllers.RegisterUser)
    router.POST("/login", controllers.LoginUser)
	router.PUT("/promote", middleware.AuthMiddleware(),middleware.AdminAuth(), controllers.PromoteUser)
	
    r := router.Group("/tasks")
    r.Use(middleware.AuthMiddleware())
    {
        r.GET("", controllers.GetTasks)
        r.GET("/:id", controllers.GetTasksById)
        r.POST("", middleware.AdminAuth(), controllers.PostTasks)
        r.PUT("/:id", middleware.AdminAuth(), controllers.PutTasksById)
        r.DELETE("/:id", middleware.AdminAuth(), controllers.DeleteTaskById)
    }
    router.Run("localhost:8080")
}