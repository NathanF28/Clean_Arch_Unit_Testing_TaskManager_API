package main

import (
	"task7/data"
	"task7/delivery/controllers"
	"task7/delivery/router"
	mongoRepo "task7/repository/mongo"
	services "task7/usecases"
)

func main() {
	userCol, taskCol := data.InitMongo()
	userRepo := mongoRepo.NewMongoUserRepository(userCol)
	taskRepo := mongoRepo.NewMongoTaskRepository(taskCol)
	userService := services.NewUserService(userRepo)
	taskService := services.NewTaskService(taskRepo)
	authController := controllers.NewAuthController(userService)
	taskController := controllers.NewTaskController(taskService)
	r := router.SetupRouter(authController, taskController)
	r.Run(":8080")
}
