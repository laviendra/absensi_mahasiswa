package controllers

import (
	"net/http"
	"time"

	"absensi-mahasiswa/database"
	"absensi-mahasiswa/models"

	"github.com/gin-gonic/gin"
)

func AbsensiMasuk(c *gin.Context) {

	var req models.AbsensiRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	tanggal := time.Now().Format("2006-01-02")
	jamMasuk := time.Now().Format("15:04:05")

	result, err := database.DB.Exec(
		`INSERT INTO absensi
		(mahasiswa_id, tanggal, jam_masuk, status)
		VALUES (?, ?, ?, ?)`,
		req.MahasiswaID,
		tanggal,
		jamMasuk,
		"Hadir",
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	id, _ := result.LastInsertId()

	c.JSON(http.StatusCreated, gin.H{
		"message": "Absensi masuk berhasil",
		"id":      id,
	})
}

func GetAbsensi(c *gin.Context) {

	rows, err := database.DB.Query(`
		SELECT
			a.id,
			m.nama,
			a.tanggal,
			a.jam_masuk,
			IFNULL(a.jam_pulang, '') as jam_pulang,
			a.status
		FROM absensi a
		JOIN mahasiswa m
		ON a.mahasiswa_id = m.id
	`)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	defer rows.Close()

	var data []models.AbsensiResponse

	for rows.Next() {

		var absen models.AbsensiResponse

		err := rows.Scan(
			&absen.ID,
			&absen.Nama,
			&absen.Tanggal,
			&absen.JamMasuk,
			&absen.JamPulang,
			&absen.Status,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		data = append(data, absen)
	}

	c.JSON(http.StatusOK, data)
}

func AbsensiPulang(c *gin.Context) {
	id := c.Param("id")

	jamPulang := time.Now().Format("15:04:05")	
	result, err := database.DB.Exec(
		`UPDATE absensi
		SET jam_pulang = ?, status = 'Pulang'
		WHERE id = ?`,
		jamPulang,
		id,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	rowAffected, _ := result.RowsAffected()

	if rowAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Absensi tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Absensi pulang berhasil",
		"jam_pulang": jamPulang,
	})
}