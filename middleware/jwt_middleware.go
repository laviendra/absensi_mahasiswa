package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"absensi-mahasiswa/utils"
)


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


		c.Next()
	}
}