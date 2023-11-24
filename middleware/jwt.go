package middleware

import (
	"github.com/golang-jwt/jwt/v5"
)

var secretkey = "rahasia"
var err error

func GenerateToken(ID uint, Email string, Role string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    ID,
		"email": Email,
		"role":  Role,
	})

	signedToken, err := token.SignedString([]byte(secretkey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
