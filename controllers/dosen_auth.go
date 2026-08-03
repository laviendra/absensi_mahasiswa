package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"absensi-mahasiswa/database"
	"absensi-mahasiswa/models"
	"absensi-mahasiswa/utils"
)

func LoginDosen(c *gin.Context) {

	var input models.Dosen

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Data tidak valid"})
		return
	}

	var id int
	var nama string
	var hash string

	err := database.DB.QueryRow(
		"SELECT id, nama, password FROM dosen WHERE username = ?",
		input.Username,
	).Scan(&id, &nama, &hash)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Username atau password salah"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Username atau password salah"})
		return
	}

	token, err := utils.GenerateTokenDosen(id, input.Username)

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