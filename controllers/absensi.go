package controllers

import (
	"database/sql"
	"net/http"
	"time"

	"absensi-mahasiswa/database"
	"absensi-mahasiswa/models"

	"github.com/gin-gonic/gin"
)

func Dashboard(c *gin.Context) {

	var dashboard models.Dashboard

	//total mahasiswa
	err := database.DB.QueryRow(
		`SELECT COUNT(*) FROM mahasiswa`,
	).Scan(&dashboard.TotalMahasiswa)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	//total hadir hari ini
	err = database.DB.QueryRow(
		`SELECT COUNT(*) FROM absensi WHERE status_kehadiran = 'Hadir' AND tanggal = CURDATE()`,
	).Scan(&dashboard.HadirHariIni)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	//total terlambat hari ini
	err = database.DB.QueryRow(
		`SELECT COUNT(*) FROM absensi WHERE status_kehadiran = 'Terlambat' AND tanggal = CURDATE()`,
	).Scan(&dashboard.Terlambat)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	//total belum absen hari ini
	err = database.DB.QueryRow(
		`SELECT COUNT(*) FROM mahasiswa WHERE id NOT IN (SELECT DISTINCT mahasiswa_id FROM absensi WHERE tanggal = CURDATE())`,
	).Scan(&dashboard.BelumAbsen)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dashboard)
}

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

	status := "Hadir"

	sekarang := time.Now()

	batasMasuk := time.Date(
		sekarang.Year(),
		sekarang.Month(),
		sekarang.Day(),
		8, 0, 0, 0,
		sekarang.Location(),
	)

	batasTutup := time.Date(
		sekarang.Year(),
		sekarang.Month(),
		sekarang.Day(),
		10, 0, 0, 0,
		sekarang.Location(),
	)

	if sekarang.After(batasMasuk) {
		status = "Terlambat"
	}

	if sekarang.After(batasTutup) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Absensi sudah ditutup",
		})
		return
	}

	var jumlah int

	err := database.DB.QueryRow(
		`SELECT COUNT(*) FROM absensi WHERE mahasiswa_id = ? AND tanggal = ?`,
		req.MahasiswaID,
		tanggal,
	).Scan(&jumlah)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if jumlah > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Mahasiswa sudah absen hari ini",
		})
		return
	}

	result, err := database.DB.Exec(
		`INSERT INTO absensi
		(mahasiswa_id, tanggal, jam_masuk, status, status_kehadiran)
		VALUES (?, ?, ?, ?, ?)`,
		req.MahasiswaID,
		tanggal,
		jamMasuk,
		"Masuk",
		status,
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
			a.jam_pulang,
			a.status_kehadiran
		FROM absensi a
		JOIN mahasiswa m
		ON a.mahasiswa_id = m.id
		ORDER BY a.tanggal DESC, m.nama ASC
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

		var jamMasuk sql.NullString
		var jamPulang sql.NullString

		err := rows.Scan(
			&absen.ID,
			&absen.Nama,
			&absen.Tanggal,
			&jamMasuk,
			&jamPulang,
			&absen.Status,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		if jamMasuk.Valid {
			absen.JamMasuk = jamMasuk.String
		}

		if jamPulang.Valid {
			absen.JamPulang = jamPulang.String
		}

		data = append(data, absen)
	}

	c.JSON(http.StatusOK, data)
}

func AbsensiPulang(c *gin.Context) {
	var statusKehadiran string

	id := c.Param("id")

	err := database.DB.QueryRow(
		`SELECT status_kehadiran
		FROM absensi
		WHERE id = ?`,
		id,
	).Scan(&statusKehadiran)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Data absensi tidak ditemukan",
		})
		return
	}

	if statusKehadiran == "Tidak Hadir" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Mahasiswa tidak hadir, tidak dapat melakukan absensi pulang",
		})
		return
	}

	var jamPulangDB sql.NullString

	err = database.DB.QueryRow(
		`SELECT jam_pulang
		FROM absensi
		WHERE id = ?`,
		id,
	).Scan(&jamPulangDB)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if jamPulangDB.Valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Mahasiswa sudah melakukan absensi pulang",
		})
		return
	}

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
		"message":    "Absensi pulang berhasil",
		"jam_pulang": jamPulang,
	})


}

// FilterAbsensi menampilkan absensi PER TANGGAL untuk SEMUA mahasiswa
// (bukan cuma yang sudah absen), sehingga tampilannya jadi satu tabel
// per hari seperti aplikasi absensi pada umumnya.
// Kalau ?tanggal tidak dikirim, otomatis pakai tanggal hari ini.
func FilterAbsensi(c *gin.Context) {

	tanggal := c.Query("tanggal")

	if tanggal == "" {
		tanggal = time.Now().Format("2006-01-02")
	}

	rows, err := database.DB.Query(
		`
		SELECT
			a.id,
			m.id,
			m.nama,
			a.jam_masuk,
			a.jam_pulang,
			a.status_kehadiran
		FROM mahasiswa m
		LEFT JOIN absensi a
		ON m.id = a.mahasiswa_id
		AND a.tanggal = ?
		ORDER BY m.nama ASC
		`,
		tanggal,
	)

	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	defer rows.Close()

	var data []models.AbsensiHarian

	for rows.Next() {

		var d models.AbsensiHarian

		var absensiID sql.NullInt64
		var jamMasuk sql.NullString
		var jamPulang sql.NullString
		var status sql.NullString

		err := rows.Scan(
			&absensiID,
			&d.MahasiswaID,
			&d.Nama,
			&jamMasuk,
			&jamPulang,
			&status,
		)

		if err != nil {
			continue
		}

		if absensiID.Valid {
			id := int(absensiID.Int64)
			d.AbsensiID = &id
		}

		if jamMasuk.Valid {
			d.JamMasuk = jamMasuk.String
		}

		if jamPulang.Valid {
			d.JamPulang = jamPulang.String
		}

		if status.Valid {
			d.Status = status.String
		} else {
			d.Status = "Belum Absen"
		}

		data = append(data, d)
	}

	c.JSON(200, gin.H{
		"tanggal": tanggal,
		"data":    data,
	})
}

func TidakHadir() {

	sekarang := time.Now()

	batasPulang := time.Date(
		sekarang.Year(),
		sekarang.Month(),
		sekarang.Day(),
		10, 0, 0,
		0,
		sekarang.Location(),
	)

	if sekarang.Before(batasPulang) {
		return
	}

	tanggal := sekarang.Format("2006-01-02")

	rows, err := database.DB.Query(`
		SELECT id FROM mahasiswa
		WHERE id NOT IN (
			SELECT mahasiswa_id 
			FROM absensi
			WHERE tanggal = ?
		)
	`, tanggal)

	if err != nil {
		return
	}

	defer rows.Close()

	tx, err := database.DB.Begin()

	if err != nil {
		return
	}

	for rows.Next() {

		var mahasiswaID int

		if err := rows.Scan(&mahasiswaID); err != nil {
			tx.Rollback()
			return
		}

		_, err = tx.Exec(`
			INSERT INTO absensi
			(mahasiswa_id, tanggal, status, status_kehadiran)
			VALUES (?, ?, ?, ?)
		`,
			mahasiswaID,
			tanggal,
			"Tidak Hadir",
			"Tidak Hadir",
		)

		if err != nil {
			tx.Rollback()
			return
		}
	}

	// commit sekali di akhir, setelah semua mahasiswa yang belum
	// absen berhasil di-insert sebagai "Tidak Hadir"
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return
	}
}