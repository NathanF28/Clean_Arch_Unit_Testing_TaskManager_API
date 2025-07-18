package main

import "task4/router"
import "task4/data"

func main() {
	data.InitMongo()
	router.StartServer()
}
