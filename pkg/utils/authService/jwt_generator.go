package authService

import (
	"os"
	"time"

	"bitbucket.org/dyfrag-internal/mass-media-core/pkg/models"
	"github.com/golang-jwt/jwt/v4"
)

func JwtGenerator(user models.User) (string, error) {
	jwtKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	expTime := time.Now().Add(time.Hour * 24)

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     expTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
