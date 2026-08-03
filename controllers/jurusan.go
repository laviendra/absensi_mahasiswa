package controllers

import (
	"net/http"
	"strconv"

	"absensi-mahasiswa/database"
	"absensi-mahasiswa/models"

	"github.com/gin-gonic/gin"
)

func GetJurusan(c *gin.Context) {

	rows, err := database.DB.Query("SELECT id, nama FROM jurusan ORDER BY nama ASC")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer rows.Close()

	var data []models.Jurusan

	for rows.Next() {

		var j models.Jurusan

		rows.Scan(&j.ID, &j.Nama)

		data = append(data, j)
	}

	c.JSON(http.StatusOK, data)
}

func CreateJurusan(c *gin.Context) {

	var j models.Jurusan

	if err := c.ShouldBindJSON(&j); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if j.Nama == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama jurusan wajib diisi"})
		return
	}

	result, err := database.DB.Exec("INSERT INTO jurusan (nama) VALUES (?)", j.Nama)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()

	c.JSON(http.StatusCreated, gin.H{"message": "Jurusan berhasil ditambahkan", "id": id})
}

func UpdateJurusan(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var j models.Jurusan

	if err := c.ShouldBindJSON(&j); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if j.Nama == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama jurusan wajib diisi"})
		return
	}

	result, err := database.DB.Exec("UPDATE jurusan SET nama = ? WHERE id = ?", j.Nama, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowAffected, _ := result.RowsAffected()

	if rowAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Jurusan tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Jurusan berhasil diperbarui"})
}

func DeleteJurusan(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	result, err := database.DB.Exec("DELETE FROM jurusan WHERE id = ?", id)

	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": fkErrorOr(err, "Jurusan masih dipakai kelas, tidak bisa dihapus")})
		return
	}

	rowAffected, _ := result.RowsAffected()

	if rowAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Jurusan tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Jurusan berhasil dihapus"})
}