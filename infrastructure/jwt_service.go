package infrastructure

import (
	"os"
	"task7/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func GenerateJWT(user domain.User) (string, error) {
	claims := jwt.MapClaims{
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(2 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
