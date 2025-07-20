package main

import "task6/router"
import "task6/data"

func main() {
	data.InitMongo()
	router.StartServer()
}
