package utils

import (
	"InceptionAnimals/app/models"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

// GenerateNewAccessToken func for generate a new Access token.
func GenerateNewAccessToken(user *models.UserPublic) (string, error) {
	// Set secret key from .env file.
	secret := os.Getenv("JWT_SECRET_KEY")

	// Set expires minutes count for secret key from .env file.
	daysCount, _ := strconv.Atoi(os.Getenv("JWT_SECRET_KEY_EXPIRE_DAYS_COUNT"))

	// Create a new claims.
	claims := jwt.MapClaims{}

	// Set public claims:
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(daysCount*24)).Unix()
	claims["id"] = user.ID
	claims["email"] = user.Email
	claims["flow_address"] = user.FlowAddress

	// TODO: add user information to JWT token

	// Create a new JWT access token with claims.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate token.
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		// Return error, it JWT token generation failed.
		return "", err
	}

	return t, nil
}
