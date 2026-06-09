package controllers

import (
	"net/http"

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

	var mhs models.Mahasiswa

	err := c.ShouldBindJSON(&mhs)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	_, err = database.DB.Exec(
		"UPDATE mahasiswa SET nama = ?, jurusan = ? WHERE id = ?",
		mhs.Nama,
		mhs.Jurusan,
		mhs.ID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Mahasiswa berhasil diperbarui",
	})
}

func DeleteMahasiswa(c *gin.Context) {

	var mhs models.Mahasiswa
	err := c.ShouldBindJSON(&mhs)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}	

	_, err = database.DB.Exec(
		"DELETE FROM mahasiswa WHERE id = ?",
		mhs.ID,
	)	

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Mahasiswa berhasil dihapus",
	})
}