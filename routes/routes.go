package routes

import (
	"absensi-mahasiswa/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	//Mahasiswa
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
	//Absensi
	r.GET(
		"/absensi",
		controllers.GetAbsensi,
	)
	r.POST(
		"/absensi/masuk",
		controllers.AbsensiMasuk,
	)
	r.PUT(
		"/absensi/pulang/:id",
		controllers.AbsensiPulang,
	)
}