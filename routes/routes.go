package routes

import (
	"absensi-mahasiswa/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	r.GET(
		"/mahasiswa",
		controllers.GetMahasiswa,
	)

	r.POST(
		"/mahasiswa",
		controllers.CreateMahasiswa,
	)

	r.PUT(
		"/mahasiswa",
		controllers.UpdateMahasiswa,
	)


	r.DELETE(
		"/mahasiswa",
		controllers.DeleteMahasiswa,
	)
}