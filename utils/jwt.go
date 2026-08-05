package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func getSecretKey() []byte {

	secret := os.Getenv("JWT_SECRET")

	if secret == "" {
		secret = "rahasia_absensi"
	}

	return []byte(secret)
}


func GenerateToken(username string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp": time.Now().Add(time.Hour * 24).Unix(),
		})

	return token.SignedString(getSecretKey())
}

// GenerateTokenDosen dipakai buat login dosen, bedanya dari token admin
// ada claim "role":"dosen" dan "dosen_id" biar middleware bisa bedain.
func GenerateTokenDosen(dosenID int, username string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"role":     "dosen",
			"dosen_id": dosenID,
			"username": username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

	return token.SignedString(getSecretKey())
}


func ValidateToken(tokenString string) (*jwt.Token, error) {

	token, err := jwt.Parse(tokenString,
		func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("method tidak valid")
			}

			return getSecretKey(), nil
		})

	return token, err
}