package utils

import (
	"crud_app/models"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


var JWT_SECRET = os.Getenv("JWT_SECRET")

type Claims struct {
	UserID primitive.ObjectID `json:"user_id"`
	Name   string             `json:"name"`
	Email  string             `json:"email"`
	jwt.RegisteredClaims
}

func GenerateJWT(user models.User) (string, error) {

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: user.ID,
		Name:   user.Name,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(JWT_SECRET))
}

func ValidateJWT(token string) (*Claims, error) {
	claims := &Claims{}

	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWT_SECRET), nil
	})

	if err != nil {
		return nil, err
	}

	if !parsedToken.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}

