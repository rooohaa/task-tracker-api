package main

import (
	database "task-tracker-api/src/database"
	router "task-tracker-api/src/router"
)

func main() {
	database.ConnectDB()
	router.InitRouter()
}
