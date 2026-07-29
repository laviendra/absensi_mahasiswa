package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"absensi-mahasiswa/database"
	"absensi-mahasiswa/models"
)

func RekapBulanan() gin.HandlerFunc {
	return func(c *gin.Context) {

		bulan, _ := strconv.Atoi(c.Query("bulan"))
		tahun, _ := strconv.Atoi(c.Query("tahun"))

		var rekap []models.Rekap

		query := `
		SELECT 
			mahasiswa.nama,
			SUM(CASE WHEN absensi.status_kehadiran='Hadir' THEN 1 ELSE 0 END) AS hadir,
			SUM(CASE WHEN absensi.status_kehadiran='Terlambat' THEN 1 ELSE 0 END) AS terlambat,
			SUM(CASE WHEN absensi.status_kehadiran='Tidak Hadir' THEN 1 ELSE 0 END) AS tidak_hadir
		FROM absensi
		JOIN mahasiswa 
		ON mahasiswa.id = absensi.mahasiswa_id
		WHERE MONTH(absensi.tanggal)=?
		AND YEAR(absensi.tanggal)=?
		GROUP BY mahasiswa.nama
		`

		rows, err := database.DB.Query(query, bulan, tahun)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		defer rows.Close()

		for rows.Next() {
			var r models.Rekap

			err := rows.Scan(
				&r.Nama,
				&r.Hadir,
				&r.Terlambat,
				&r.TidakHadir,
			)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			rekap = append(rekap, r)
		}

		c.JSON(http.StatusOK, rekap)
	}
}