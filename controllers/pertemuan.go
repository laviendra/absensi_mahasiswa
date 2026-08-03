package controllers

import (
	"net/http"
	"strconv"
	"time"

	"absensi-mahasiswa/database"
	"absensi-mahasiswa/models"

	"github.com/gin-gonic/gin"
)

// BukaPertemuan dipanggil dosen pas mau mulai ngajar hari ini.
// Kalau pertemuan hari ini buat jadwal itu udah pernah dibuka, yang lama dipakai lagi (nggak dobel).
func BukaPertemuan(c *gin.Context) {

	dosenID := c.GetInt("dosen_id")

	jadwalID, err := strconv.Atoi(c.Param("jadwal_id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID jadwal tidak valid"})
		return
	}

	var pemilik int

	err = database.DB.QueryRow("SELECT dosen_id FROM jadwal WHERE id = ?", jadwalID).Scan(&pemilik)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Jadwal tidak ditemukan"})
		return
	}

	if pemilik != dosenID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Ini bukan jadwal Anda"})
		return
	}

	tanggal := time.Now().Format("2006-01-02")

	var pertemuanID int

	err = database.DB.QueryRow(
		"SELECT id FROM pertemuan WHERE jadwal_id = ? AND tanggal = ?",
		jadwalID, tanggal,
	).Scan(&pertemuanID)

	if err == nil {
		c.JSON(http.StatusOK, gin.H{"pertemuan_id": pertemuanID, "tanggal": tanggal})
		return
	}

	result, err := database.DB.Exec(
		"INSERT INTO pertemuan (jadwal_id, tanggal, status) VALUES (?, ?, 'dibuka')",
		jadwalID, tanggal,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()

	c.JSON(http.StatusCreated, gin.H{"pertemuan_id": id, "tanggal": tanggal})
}

// GetAbsensiPertemuan: semua mahasiswa aktif di kelas jadwal ini,
// lengkap status kehadirannya di pertemuan tsb (default "Tidak Hadir" kalau belum diisi dosen)
func GetAbsensiPertemuan(c *gin.Context) {

	dosenID := c.GetInt("dosen_id")

	pertemuanID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID pertemuan tidak valid"})
		return
	}

	var kelasID, pemilik int

	err = database.DB.QueryRow(`
		SELECT j.kelas_id, j.dosen_id
		FROM pertemuan p
		JOIN jadwal j ON j.id = p.jadwal_id
		WHERE p.id = ?
	`, pertemuanID).Scan(&kelasID, &pemilik)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pertemuan tidak ditemukan"})
		return
	}

	if pemilik != dosenID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Ini bukan pertemuan Anda"})
		return
	}

	rows, err := database.DB.Query(`
		SELECT m.id, m.nama,
			COALESCE(ak.jam_hadir, ''),
			COALESCE(ak.status_kehadiran, 'Tidak Hadir')
		FROM mahasiswa m
		LEFT JOIN absensi_kelas ak
			ON ak.mahasiswa_id = m.id AND ak.pertemuan_id = ?
		WHERE m.kelas_id = ? AND m.status = 'aktif'
		ORDER BY m.nama ASC
	`, pertemuanID, kelasID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer rows.Close()

	var data []models.AbsensiKelasItem

	for rows.Next() {

		var item models.AbsensiKelasItem

		rows.Scan(&item.MahasiswaID, &item.Nama, &item.JamHadir, &item.StatusKehadiran)

		data = append(data, item)
	}

	c.JSON(http.StatusOK, gin.H{"pertemuan_id": pertemuanID, "data": data})
}

// SimpanAbsensiKelas: dosen nandain status kehadiran satu mahasiswa di pertemuan ini
func SimpanAbsensiKelas(c *gin.Context) {

	dosenID := c.GetInt("dosen_id")

	pertemuanID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID pertemuan tidak valid"})
		return
	}

	var pemilik int
	var statusPertemuan string

	err = database.DB.QueryRow(`
		SELECT j.dosen_id, p.status
		FROM pertemuan p
		JOIN jadwal j ON j.id = p.jadwal_id
		WHERE p.id = ?
	`, pertemuanID).Scan(&pemilik, &statusPertemuan)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pertemuan tidak ditemukan"})
		return
	}

	if pemilik != dosenID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Ini bukan pertemuan Anda"})
		return
	}

	if statusPertemuan == "selesai" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pertemuan ini sudah ditutup, tidak bisa diubah lagi"})
		return
	}

	var req models.SimpanAbsensiKelasRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jamHadir := ""

	if req.StatusKehadiran == "Hadir" || req.StatusKehadiran == "Terlambat" {
		jamHadir = time.Now().Format("15:04:05")
	}

	_, err = database.DB.Exec(`
		INSERT INTO absensi_kelas (pertemuan_id, mahasiswa_id, jam_hadir, status_kehadiran)
		VALUES (?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE jam_hadir = VALUES(jam_hadir), status_kehadiran = VALUES(status_kehadiran)
	`, pertemuanID, req.MahasiswaID, jamHadir, req.StatusKehadiran)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Kehadiran tersimpan"})
}

// TutupPertemuan mengunci sesi biar nggak keedit-edit lagi
func TutupPertemuan(c *gin.Context) {

	dosenID := c.GetInt("dosen_id")

	pertemuanID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID pertemuan tidak valid"})
		return
	}

	var pemilik int

	err = database.DB.QueryRow(`
		SELECT j.dosen_id
		FROM pertemuan p
		JOIN jadwal j ON j.id = p.jadwal_id
		WHERE p.id = ?
	`, pertemuanID).Scan(&pemilik)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pertemuan tidak ditemukan"})
		return
	}

	if pemilik != dosenID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Ini bukan pertemuan Anda"})
		return
	}

	_, err = database.DB.Exec("UPDATE pertemuan SET status = 'selesai' WHERE id = ?", pertemuanID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pertemuan ditutup"})
}