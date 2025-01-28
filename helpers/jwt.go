package helpers

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var JWT_SECRET string
var JWT_EXPIRES_IN uint

// func GenerateJWT(user models.Users) (string, error) {
// 	err :=
// }


// User struct represents the user model.
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"` // Ideally hashed
}

// Claims struct for JWT.
type Claims struct {
	Username string `json:"username"`
	Roles string `json:"roles"`
	jwt.StandardClaims
}

// JWT secret key (use environment variables in production).
var jwtKey = []byte("your_secret_key")

// GenerateJWT generates a JWT token for a given username.
func GenerateJWT(username string, role string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: username,
		Roles: role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// ValidateJWT validates a JWT token and returns the username if valid.
func ValidateJWT(tokenStr string) (Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return *claims, errors.New("invalid token")
	}
	return *claims, nil
}
