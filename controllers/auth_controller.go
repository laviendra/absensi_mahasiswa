package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"absensi-mahasiswa/database"
	"absensi-mahasiswa/models"
	"absensi-mahasiswa/utils"
)

func Login(c *gin.Context) {

	var admin models.Admin

	err := c.ShouldBindJSON(&admin)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Data tidak valid",
		})
		return
	}


	var data models.Admin


	query := `
		SELECT id_admin, username, password
		FROM admin
		WHERE username = ? AND password = ?
	`

	err = database.DB.QueryRow(
		query,
		admin.Username,
		admin.Password,
	).Scan(
		&data.ID,
		&data.Username,
		&data.Password,
	)


	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Username atau password salah",
		})
		return
	}


	token, err := utils.GenerateToken(data.Username)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal membuat token",
		})
		return
	}


	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil",
		"token": token,
	})
}