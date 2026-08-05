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

	kelasFilter := c.Query("kelas_id")

	query := `
		SELECT m.id, m.nama, COALESCE(m.nim, ''), m.kelas_id, COALESCE(k.nama, ''), k.jurusan_id, COALESCE(j.nama, ''), m.status
		FROM mahasiswa m
		LEFT JOIN kelas k ON k.id = m.kelas_id
		LEFT JOIN jurusan j ON j.id = k.jurusan_id
		WHERE 1=1
	`

	var args []interface{}

	if statusFilter != "semua" {
		query += " AND m.status = 'aktif'"
	}

	if kelasFilter != "" {
		query += " AND m.kelas_id = ?"
		args = append(args, kelasFilter)
	}

	query += " ORDER BY m.nama ASC"

	rows, err := database.DB.Query(query, args...)

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
		var jurusanID sql.NullInt64

		rows.Scan(
			&mhs.ID,
			&mhs.Nama,
			&mhs.Nim,
			&kelasID,
			&mhs.NamaKelas,
			&jurusanID,
			&mhs.NamaJurusan,
			&mhs.Status,
		)

		if kelasID.Valid {
			id := int(kelasID.Int64)
			mhs.KelasID = &id
		}

		if jurusanID.Valid {
			id := int(jurusanID.Int64)
			mhs.JurusanID = &id
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

	if mhs.Nama == "" || mhs.Nim == "" || mhs.KelasID == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Nama, NIM, dan kelas wajib diisi",
		})
		return
	}

	result, err := database.DB.Exec(
		"INSERT INTO mahasiswa (nama, nim, kelas_id, status) VALUES (?, ?, ?, 'aktif')",
		mhs.Nama,
		mhs.Nim,
		*mhs.KelasID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fkErrorOr(err, "NIM sudah dipakai mahasiswa lain"),
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

	if mhs.Nama == "" || mhs.Nim == "" || mhs.KelasID == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Nama, NIM, dan kelas wajib diisi",
		})
		return
	}

	result, err := database.DB.Exec(
		"UPDATE mahasiswa SET nama = ?, nim = ?, kelas_id = ? WHERE id = ?",
		mhs.Nama,
		mhs.Nim,
		*mhs.KelasID,
		id,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fkErrorOr(err, "NIM sudah dipakai mahasiswa lain"),
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