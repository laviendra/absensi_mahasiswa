package controllers

import (
	"net/http"

	"absensi-mahasiswa/database"
	"absensi-mahasiswa/models"

	"github.com/gin-gonic/gin"
)

func Dashboard(c *gin.Context) {

	var dashboard models.Dashboard

	err := database.DB.QueryRow(
		`SELECT COUNT(*) FROM mahasiswa WHERE status = 'aktif'`,
	).Scan(&dashboard.TotalMahasiswa)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = database.DB.QueryRow(
		`SELECT COUNT(*) FROM dosen`,
	).Scan(&dashboard.TotalDosen)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = database.DB.QueryRow(
		`SELECT COUNT(*) FROM kelas`,
	).Scan(&dashboard.TotalKelas)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = database.DB.QueryRow(
		`SELECT COUNT(*) FROM pertemuan WHERE tanggal = CURDATE()`,
	).Scan(&dashboard.PertemuanHariIni)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dashboard)
}
