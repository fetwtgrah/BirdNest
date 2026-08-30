package utils

import (
	"errors"
	"github/fetwtgrah/BirdNest/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JWT = []byte("JWT_SECRET")

func GenerateToken(user model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userid":    user.ID,
		"useremail": user.Email,
		"username":  user.Name,
		"exp":       time.Now().Add(time.Hour * 1).Unix(),
	})
	return token.SignedString(JWT)
}

func ParseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {
			return JWT, nil
		})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
