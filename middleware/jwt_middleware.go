package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"absensi-mahasiswa/utils"
)


// AuthMiddleware dipakai buat route admin. Token dosen/mahasiswa ditolak di sini
// biar mereka nggak bisa nyerempet ke endpoint admin.
func AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {


		authHeader := c.GetHeader("Authorization")


		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message":"Token tidak ditemukan",
			})
			c.Abort()
			return
		}


		tokenString := strings.Replace(
			authHeader,
			"Bearer ",
			"",
			1,
		)


		token, err := utils.ValidateToken(tokenString)


		if err != nil || !token.Valid {

			c.JSON(http.StatusUnauthorized, gin.H{
				"message":"Token tidak valid",
			})

			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if role, ada := claims["role"]; ada && role != "" {
				c.JSON(http.StatusForbidden, gin.H{
					"message": "Token ini tidak berlaku di sini",
				})
				c.Abort()
				return
			}
		}


		c.Next()
	}
}


// DosenAuthMiddleware dipakai buat route khusus dosen.
// Menyimpan dosen_id ke context biar controller tau ini dosen siapa.
func DosenAuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Token tidak ditemukan",
			})
			c.Abort()
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		token, err := utils.ValidateToken(tokenString)

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Token tidak valid",
			})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok || claims["role"] != "dosen" {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "Token ini bukan token dosen",
			})
			c.Abort()
			return
		}

		dosenIDFloat, _ := claims["dosen_id"].(float64)

		c.Set("dosen_id", int(dosenIDFloat))

		c.Next()
	}
}


// MahasiswaAuthMiddleware dipakai buat route khusus mahasiswa (lihat rekap sendiri).
// Menyimpan mahasiswa_id ke context biar controller tau ini mahasiswa siapa.
func MahasiswaAuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Token tidak ditemukan",
			})
			c.Abort()
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		token, err := utils.ValidateToken(tokenString)

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Token tidak valid",
			})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok || claims["role"] != "mahasiswa" {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "Token ini bukan token mahasiswa",
			})
			c.Abort()
			return
		}

		mahasiswaIDFloat, _ := claims["mahasiswa_id"].(float64)

		c.Set("mahasiswa_id", int(mahasiswaIDFloat))

		c.Next()
	}
}