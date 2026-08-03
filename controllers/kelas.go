package controllers

import (
	"net/http"
	"strconv"

	"absensi-mahasiswa/database"
	"absensi-mahasiswa/models"

	"github.com/gin-gonic/gin"
)

func GetKelas(c *gin.Context) {

	rows, err := database.DB.Query(`
		SELECT k.id, k.nama, k.jurusan_id, j.nama
		FROM kelas k
		JOIN jurusan j ON j.id = k.jurusan_id
		ORDER BY k.nama ASC
	`)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer rows.Close()

	var data []models.Kelas

	for rows.Next() {

		var k models.Kelas

		rows.Scan(&k.ID, &k.Nama, &k.JurusanID, &k.NamaJurusan)

		data = append(data, k)
	}

	c.JSON(http.StatusOK, data)
}

func CreateKelas(c *gin.Context) {

	var k models.Kelas

	if err := c.ShouldBindJSON(&k); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if k.Nama == "" || k.JurusanID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama kelas dan jurusan wajib diisi"})
		return
	}

	result, err := database.DB.Exec("INSERT INTO kelas (nama, jurusan_id) VALUES (?, ?)", k.Nama, k.JurusanID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()

	c.JSON(http.StatusCreated, gin.H{"message": "Kelas berhasil ditambahkan", "id": id})
}

func UpdateKelas(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var k models.Kelas

	if err := c.ShouldBindJSON(&k); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if k.Nama == "" || k.JurusanID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama kelas dan jurusan wajib diisi"})
		return
	}

	result, err := database.DB.Exec("UPDATE kelas SET nama = ?, jurusan_id = ? WHERE id = ?", k.Nama, k.JurusanID, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowAffected, _ := result.RowsAffected()

	if rowAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kelas tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Kelas berhasil diperbarui"})
}

func DeleteKelas(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	result, err := database.DB.Exec("DELETE FROM kelas WHERE id = ?", id)

	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": fkErrorOr(err, "Kelas masih dipakai mahasiswa/jadwal, tidak bisa dihapus")})
		return
	}

	rowAffected, _ := result.RowsAffected()

	if rowAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kelas tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Kelas berhasil dihapus"})
}