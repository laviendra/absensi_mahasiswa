package controllers

import (
	"net/http"
	"strconv"

	"absensi-mahasiswa/database"
	"absensi-mahasiswa/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// GetDosen TIDAK mengembalikan password (hash sekalipun) ke client
func GetDosen(c *gin.Context) {

	rows, err := database.DB.Query("SELECT id, nama, username FROM dosen ORDER BY nama ASC")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer rows.Close()

	var data []models.Dosen

	for rows.Next() {

		var d models.Dosen

		rows.Scan(&d.ID, &d.Nama, &d.Username)

		data = append(data, d)
	}

	c.JSON(http.StatusOK, data)
}

func CreateDosen(c *gin.Context) {

	var d models.Dosen

	if err := c.ShouldBindJSON(&d); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if d.Nama == "" || d.Username == "" || d.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama, username, dan password wajib diisi"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(d.Password), bcrypt.DefaultCost)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result, err := database.DB.Exec(
		"INSERT INTO dosen (nama, username, password) VALUES (?, ?, ?)",
		d.Nama, d.Username, string(hash),
	)

	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": fkErrorOr(err, "Username sudah dipakai")})
		return
	}

	id, _ := result.LastInsertId()

	c.JSON(http.StatusCreated, gin.H{"message": "Dosen berhasil ditambahkan", "id": id})
}

// UpdateDosen: password cuma diganti kalau field password dikirim (nggak kosong)
func UpdateDosen(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var d models.Dosen

	if err := c.ShouldBindJSON(&d); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if d.Nama == "" || d.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama dan username wajib diisi"})
		return
	}

	var execErr error
	var rowAffected int64

	if d.Password != "" {

		hash, hashErr := bcrypt.GenerateFromPassword([]byte(d.Password), bcrypt.DefaultCost)

		if hashErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": hashErr.Error()})
			return
		}

		result, err := database.DB.Exec(
			"UPDATE dosen SET nama = ?, username = ?, password = ? WHERE id = ?",
			d.Nama, d.Username, string(hash), id,
		)

		execErr = err

		if err == nil {
			rowAffected, _ = result.RowsAffected()
		}

	} else {

		result, err := database.DB.Exec(
			"UPDATE dosen SET nama = ?, username = ? WHERE id = ?",
			d.Nama, d.Username, id,
		)

		execErr = err

		if err == nil {
			rowAffected, _ = result.RowsAffected()
		}
	}

	if execErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fkErrorOr(execErr, "Username sudah dipakai")})
		return
	}

	if rowAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dosen tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Dosen berhasil diperbarui"})
}

func DeleteDosen(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	result, err := database.DB.Exec("DELETE FROM dosen WHERE id = ?", id)

	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": fkErrorOr(err, "Dosen masih punya jadwal mengajar, tidak bisa dihapus")})
		return
	}

	rowAffected, _ := result.RowsAffected()

	if rowAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dosen tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Dosen berhasil dihapus"})
}