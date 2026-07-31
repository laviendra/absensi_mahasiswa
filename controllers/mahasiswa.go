package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"absensi-mahasiswa/database"
	"absensi-mahasiswa/models"

	"github.com/gin-gonic/gin"
)

func GetMahasiswa(c *gin.Context) {

	rows, err := database.DB.Query(
		"SELECT id,nama,jurusan FROM mahasiswa",
	)

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

		rows.Scan(
			&mhs.ID,
			&mhs.Nama,
			&mhs.Jurusan,
		)

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

	if mhs.Nama == "" || mhs.Jurusan == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Nama dan jurusan wajib diisi",
		})
		return
	}

	result, err := database.DB.Exec(
		"INSERT INTO mahasiswa (nama, jurusan) VALUES (?, ?)",
		mhs.Nama,
		mhs.Jurusan,
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

	if mhs.Nama == "" || mhs.Jurusan == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Nama dan jurusan wajib diisi",
		})
		return
	}

	result, err := database.DB.Exec(
		"UPDATE mahasiswa SET nama = ?, jurusan = ? WHERE id = ?",
		mhs.Nama,
		mhs.Jurusan,
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

func DeleteMahasiswa(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID tidak valid",
		})
		return
	}

	result, err := database.DB.Exec(
		"DELETE FROM mahasiswa WHERE id = ?",
		id,
	)

	if err != nil {

		if strings.Contains(err.Error(), "foreign key constraint") {
			c.JSON(http.StatusConflict, gin.H{
				"error": "Mahasiswa ini masih punya data absensi, tidak bisa dihapus",
			})
			return
		}

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
		"message": "Mahasiswa berhasil dihapus",
	})
}