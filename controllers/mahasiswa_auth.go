package controllers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"absensi-mahasiswa/database"
	"absensi-mahasiswa/utils"
)

func LoginMahasiswa(c *gin.Context) {

	var input struct {
		Nim      string `json:"nim"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Data tidak valid"})
		return
	}

	var id int
	var nama string
	var hash sql.NullString

	err := database.DB.QueryRow(
		"SELECT id, nama, password FROM mahasiswa WHERE nim = ? AND status = 'aktif'",
		input.Nim,
	).Scan(&id, &nama, &hash)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "NIM atau password salah"})
		return
	}

	if !hash.Valid || hash.String == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Akun ini belum diaktifkan, hubungi admin buat setting password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash.String), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "NIM atau password salah"})
		return
	}

	token, err := utils.GenerateTokenMahasiswa(id, input.Nim)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal membuat token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil",
		"token":   token,
		"nama":    nama,
	})
}