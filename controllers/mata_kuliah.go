package controllers

import (
	"net/http"
	"strconv"

	"absensi-mahasiswa/database"
	"absensi-mahasiswa/models"

	"github.com/gin-gonic/gin"
)

func GetMataKuliah(c *gin.Context) {

	rows, err := database.DB.Query("SELECT id, nama, kode FROM mata_kuliah ORDER BY nama ASC")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer rows.Close()

	var data []models.MataKuliah

	for rows.Next() {

		var mk models.MataKuliah

		rows.Scan(&mk.ID, &mk.Nama, &mk.Kode)

		data = append(data, mk)
	}

	c.JSON(http.StatusOK, data)
}

func CreateMataKuliah(c *gin.Context) {

	var mk models.MataKuliah

	if err := c.ShouldBindJSON(&mk); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if mk.Nama == "" || mk.Kode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama dan kode mata kuliah wajib diisi"})
		return
	}

	result, err := database.DB.Exec("INSERT INTO mata_kuliah (nama, kode) VALUES (?, ?)", mk.Nama, mk.Kode)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()

	c.JSON(http.StatusCreated, gin.H{"message": "Mata kuliah berhasil ditambahkan", "id": id})
}

func UpdateMataKuliah(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var mk models.MataKuliah

	if err := c.ShouldBindJSON(&mk); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if mk.Nama == "" || mk.Kode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama dan kode mata kuliah wajib diisi"})
		return
	}

	result, err := database.DB.Exec("UPDATE mata_kuliah SET nama = ?, kode = ? WHERE id = ?", mk.Nama, mk.Kode, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowAffected, _ := result.RowsAffected()

	if rowAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mata kuliah tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Mata kuliah berhasil diperbarui"})
}

func DeleteMataKuliah(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	result, err := database.DB.Exec("DELETE FROM mata_kuliah WHERE id = ?", id)

	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": fkErrorOr(err, "Mata kuliah masih dipakai jadwal, tidak bisa dihapus")})
		return
	}

	rowAffected, _ := result.RowsAffected()

	if rowAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mata kuliah tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Mata kuliah berhasil dihapus"})
}