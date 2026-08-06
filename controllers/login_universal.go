package controllers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"absensi-mahasiswa/database"
	"absensi-mahasiswa/utils"
)

// LoginUniversal: satu form buat admin, dosen, maupun mahasiswa.
// "identifier" bisa diisi username (admin/dosen) atau NIM (mahasiswa).
// Dicoba berurutan: admin dulu, dosen, baru mahasiswa.
func LoginUniversal(c *gin.Context) {

	var input struct {
		Identifier string `json:"identifier"`
		Password   string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Data tidak valid"})
		return
	}

	// ---------- coba sebagai ADMIN ----------
	var adminUsername, adminHash string

	err := database.DB.QueryRow(
		"SELECT username, password FROM admin WHERE username = ?",
		input.Identifier,
	).Scan(&adminUsername, &adminHash)

	if err == nil && bcrypt.CompareHashAndPassword([]byte(adminHash), []byte(input.Password)) == nil {

		token, tErr := utils.GenerateToken(adminUsername)

		if tErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal membuat token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Login berhasil",
			"token":   token,
			"role":    "admin",
			"nama":    adminUsername,
		})
		return
	}

	// ---------- coba sebagai DOSEN ----------
	var dosenID int
	var dosenNama, dosenUsername, dosenHash string

	err = database.DB.QueryRow(
		"SELECT id, nama, username, password FROM dosen WHERE username = ?",
		input.Identifier,
	).Scan(&dosenID, &dosenNama, &dosenUsername, &dosenHash)

	if err == nil && bcrypt.CompareHashAndPassword([]byte(dosenHash), []byte(input.Password)) == nil {

		token, tErr := utils.GenerateTokenDosen(dosenID, dosenUsername)

		if tErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal membuat token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Login berhasil",
			"token":   token,
			"role":    "dosen",
			"nama":    dosenNama,
		})
		return
	}

	// ---------- coba sebagai MAHASISWA (identifier = NIM) ----------
	var mhsID int
	var mhsNama, mhsNim string
	var mhsHash sql.NullString

	err = database.DB.QueryRow(
		"SELECT id, nama, nim, password FROM mahasiswa WHERE nim = ? AND status = 'aktif'",
		input.Identifier,
	).Scan(&mhsID, &mhsNama, &mhsNim, &mhsHash)

	if err == nil && mhsHash.Valid && mhsHash.String != "" &&
		bcrypt.CompareHashAndPassword([]byte(mhsHash.String), []byte(input.Password)) == nil {

		token, tErr := utils.GenerateTokenMahasiswa(mhsID, mhsNim)

		if tErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal membuat token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Login berhasil",
			"token":   token,
			"role":    "mahasiswa",
			"nama":    mhsNama,
		})
		return
	}

	// nggak ketemu di admin, dosen, maupun mahasiswa
	c.JSON(http.StatusUnauthorized, gin.H{"message": "Username/NIM atau password salah"})
}