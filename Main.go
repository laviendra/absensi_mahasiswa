package main

import (
	"absensi-mahasiswa/database"
	"absensi-mahasiswa/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	database.ConnectDB()

	r := gin.Default()

	r.Static("/frontend", "./frontend")

	routes.SetupRoutes(r)

	r.Run(":8080")
}