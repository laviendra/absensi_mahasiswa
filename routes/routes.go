package routes

import (
	"absensi-mahasiswa/controllers"
	"absensi-mahasiswa/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	// Login tidak pakai middleware
	r.POST(
		"/login",
		controllers.Login,
	)


	// Route yang butuh login
	protected := r.Group("/")

	protected.Use(middleware.AuthMiddleware())


	// Dashboard
	protected.GET(
		"/dashboard",
		controllers.Dashboard,
	)


	// Mahasiswa
	protected.GET(
		"/mahasiswa",
		controllers.GetMahasiswa,
	)

	protected.POST(
		"/mahasiswa",
		controllers.CreateMahasiswa,
	)

	protected.PUT(
		"/mahasiswa/:id",
		controllers.UpdateMahasiswa,
	)

	protected.DELETE(
		"/mahasiswa/:id",
		controllers.DeleteMahasiswa,
	)


	// Absensi
	protected.GET(
		"/absensi",
		controllers.GetAbsensi,
	)

	protected.GET(
		"/absensi/filter",
		controllers.FilterAbsensi,
	)

	protected.POST(
		"/absensi/masuk",
		controllers.AbsensiMasuk,
	)

	protected.PUT(
		"/absensi/pulang/:id",
		controllers.AbsensiPulang,
	)

	//rekap
	protected.GET("/rekap", controllers.RekapBulanan())

	// export
	protected.GET("/rekap/pdf", controllers.ExportPDF)

	protected.GET("/rekap/excel", controllers.ExportExcel)
}