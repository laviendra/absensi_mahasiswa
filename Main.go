package main

import (
	"os"

	"absensi-mahasiswa/database"
	"absensi-mahasiswa/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()

	database.ConnectDB()

	r := gin.Default()

	r.Static("/frontend", "./frontend")

	routes.SetupRoutes(r)

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}