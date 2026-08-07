package routes

import (
	"absensi-mahasiswa/controllers"
	"absensi-mahasiswa/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	// Login admin & dosen & mahasiswa (tidak pakai middleware)
	r.POST("/login", controllers.Login)
	r.POST("/login-dosen", controllers.LoginDosen)
	r.POST("/login-mahasiswa", controllers.LoginMahasiswa)
	r.POST("/login-universal", controllers.LoginUniversal)


	// ================= ROUTE ADMIN =================
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())


	// Dashboard
	protected.GET("/dashboard", controllers.Dashboard)


	// Mahasiswa
	protected.GET("/mahasiswa", controllers.GetMahasiswa)
	protected.POST("/mahasiswa", controllers.CreateMahasiswa)
	protected.PUT("/mahasiswa/:id", controllers.UpdateMahasiswa)
	protected.DELETE("/mahasiswa/:id", controllers.DeleteMahasiswa)
	protected.PUT("/mahasiswa/:id/aktifkan", controllers.AktifkanMahasiswa)


	// Jurusan
	protected.GET("/jurusan", controllers.GetJurusan)
	protected.POST("/jurusan", controllers.CreateJurusan)
	protected.PUT("/jurusan/:id", controllers.UpdateJurusan)
	protected.DELETE("/jurusan/:id", controllers.DeleteJurusan)


	// Kelas
	protected.GET("/kelas", controllers.GetKelas)
	protected.POST("/kelas", controllers.CreateKelas)
	protected.PUT("/kelas/:id", controllers.UpdateKelas)
	protected.DELETE("/kelas/:id", controllers.DeleteKelas)


	// Dosen (data master, dikelola admin)
	protected.GET("/dosen", controllers.GetDosen)
	protected.POST("/dosen", controllers.CreateDosen)
	protected.PUT("/dosen/:id", controllers.UpdateDosen)
	protected.DELETE("/dosen/:id", controllers.DeleteDosen)


	// Mata kuliah
	protected.GET("/mata-kuliah", controllers.GetMataKuliah)
	protected.POST("/mata-kuliah", controllers.CreateMataKuliah)
	protected.PUT("/mata-kuliah/:id", controllers.UpdateMataKuliah)
	protected.DELETE("/mata-kuliah/:id", controllers.DeleteMataKuliah)
	protected.GET("/mata-kuliah/:id/kelas", controllers.GetKelasByMataKuliah)


	// Jadwal
	protected.GET("/jadwal", controllers.GetJadwal)
	protected.POST("/jadwal", controllers.CreateJadwal)
	protected.DELETE("/jadwal/:id", controllers.DeleteJadwal)
	protected.GET("/dosen/:id/mata-kuliah", controllers.GetMataKuliahByDosen)


	// Rekap absensi kelas (pantau & koreksi hasil absen dosen)
	protected.GET("/rekap-kelas/pertemuan", controllers.GetPertemuanAdmin)
	protected.GET("/rekap-kelas/pertemuan/:id/absensi", controllers.GetAbsensiPertemuanAdmin)
	protected.PUT("/rekap-kelas/pertemuan/:id/absensi", controllers.UpdateAbsensiKelasAdmin)
	protected.GET("/rekap-kelas/jadwal/:id/export/pdf", controllers.ExportRekapKelasPDF)
	protected.GET("/rekap-kelas/jadwal/:id/export/excel", controllers.ExportRekapKelasExcel)


	// ================= ROUTE DOSEN =================
	dosenRoute := r.Group("/dosen-area")
	dosenRoute.Use(middleware.DosenAuthMiddleware())

	dosenRoute.GET("/jadwal-saya", controllers.GetJadwalSaya)
	dosenRoute.POST("/jadwal/:jadwal_id/buka-pertemuan", controllers.BukaPertemuan)
	dosenRoute.GET("/pertemuan/:id/absensi", controllers.GetAbsensiPertemuan)
	dosenRoute.POST("/pertemuan/:id/absensi", controllers.SimpanAbsensiKelas)
	dosenRoute.PUT("/pertemuan/:id/tutup", controllers.TutupPertemuan)


	// ================= ROUTE MAHASISWA =================
	mahasiswaRoute := r.Group("/mahasiswa-area")
	mahasiswaRoute.Use(middleware.MahasiswaAuthMiddleware())

	mahasiswaRoute.GET("/rekap", controllers.GetRekapSaya)
}