package controllers

import (
	"database/sql"
	"net/http"
	"strconv"

	"absensi-mahasiswa/database"
	"absensi-mahasiswa/models"

	"github.com/gin-gonic/gin"
)

// GetMahasiswa default cuma nampilin mahasiswa yang aktif.
// Tambahkan ?status=semua buat lihat yang nonaktif juga.
func GetMahasiswa(c *gin.Context) {

	statusFilter := c.Query("status")

	query := `
		SELECT m.id, m.nama, m.kelas_id, COALESCE(k.nama, ''), COALESCE(j.nama, ''), m.status
		FROM mahasiswa m
		LEFT JOIN kelas k ON k.id = m.kelas_id
		LEFT JOIN jurusan j ON j.id = k.jurusan_id
	`

	if statusFilter != "semua" {
		query += " WHERE m.status = 'aktif'"
	}

	query += " ORDER BY m.nama ASC"

	rows, err := database.DB.Query(query)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	defer rows.Close()

	var mahasiswa []models.Mahasiswa

	for rows.Next() {

		var mhs models.Mahasiswa
		var kelasID sql.NullInt64

		rows.Scan(
			&mhs.ID,
			&mhs.Nama,
			&kelasID,
			&mhs.NamaKelas,
			&mhs.NamaJurusan,
			&mhs.Status,
		)

		if kelasID.Valid {
			id := int(kelasID.Int64)
			mhs.KelasID = &id
		}

		mahasiswa = append(mahasiswa, mhs)
	}

	c.JSON(http.StatusOK, mahasiswa)
}

func CreateMahasiswa(c *gin.Context) {

	var mhs models.Mahasiswa

	err := c.ShouldBindJSON(&mhs)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if mhs.Nama == "" || mhs.KelasID == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Nama dan kelas wajib diisi",
		})
		return
	}

	result, err := database.DB.Exec(
		"INSERT INTO mahasiswa (nama, kelas_id, status) VALUES (?, ?, 'aktif')",
		mhs.Nama,
		*mhs.KelasID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	id, _ := result.LastInsertId()

	c.JSON(http.StatusCreated, gin.H{
		"message": "Mahasiswa berhasil ditambahkan",
		"id":      id,
	})
}

func UpdateMahasiswa(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID tidak valid",
		})
		return
	}

	var mhs models.Mahasiswa

	err = c.ShouldBindJSON(&mhs)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if mhs.Nama == "" || mhs.KelasID == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Nama dan kelas wajib diisi",
		})
		return
	}

	result, err := database.DB.Exec(
		"UPDATE mahasiswa SET nama = ?, kelas_id = ? WHERE id = ?",
		mhs.Nama,
		*mhs.KelasID,
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
			"error": "Mahasiswa tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Mahasiswa berhasil diperbarui",
	})
}

// DeleteMahasiswa TIDAK menghapus baris dari database (biar riwayat
// absensi yang sudah ada tetap aman / tidak melanggar foreign key).
// Mahasiswa cuma ditandai nonaktif.
func DeleteMahasiswa(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID tidak valid",
		})
		return
	}

	result, err := database.DB.Exec(
		"UPDATE mahasiswa SET status = 'nonaktif' WHERE id = ?",
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
			"error": "Mahasiswa tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Mahasiswa berhasil dinonaktifkan",
	})
}

func AktifkanMahasiswa(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID tidak valid",
		})
		return
	}

	result, err := database.DB.Exec(
		"UPDATE mahasiswa SET status = 'aktif' WHERE id = ?",
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
			"error": "Mahasiswa tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Mahasiswa berhasil diaktifkan",
	})
}