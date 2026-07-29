package main

import (
	"absensi-mahasiswa/database"
	"absensi-mahasiswa/routes"
	"absensi-mahasiswa/controllers"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {

	database.ConnectDB()

	r := gin.Default()

	r.Static("/frontend", "./frontend")

	routes.SetupRoutes(r)


	// cek otomatis mahasiswa tidak hadir
	go func() {
		for {
			controllers.TidakHadir()
			time.Sleep(1 * time.Minute)
		}
	}()


	r.Run(":8080")
}