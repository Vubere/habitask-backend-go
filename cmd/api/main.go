package main

import "habitask-backend-go/internal/database"

func main() {
	db := database.Connect()
	database.Migrate(db)
}
