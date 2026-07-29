package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("rahasia_absensi")


func GenerateToken(username string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp": time.Now().Add(time.Hour * 24).Unix(),
		})

	return token.SignedString(secretKey)
}


func ValidateToken(tokenString string) (*jwt.Token, error) {

	token, err := jwt.Parse(tokenString,
		func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("method tidak valid")
			}

			return secretKey, nil
		})

	return token, err
}