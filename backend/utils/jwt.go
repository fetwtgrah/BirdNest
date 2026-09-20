package utils

import (
	"errors"
	"fmt"
	"github/fetwtgrah/BirdNest/backend/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JWT = []byte("JWT_SECRET")

type CustomClaims struct {
	Username  string `json:"username"`
	UserID    uint   `json:"userid"`
	UserEmail string `json:"useremail"`
	jwt.RegisteredClaims
}

func GenerateToken(user model.User) (string, error) {
	claims := CustomClaims{
		UserID:    user.ID,
		UserEmail: user.Email,
		Username:  user.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWT)
}

func ParseToken(tokenString string) (*CustomClaims, error) {
	claims := &CustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return JWT, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
