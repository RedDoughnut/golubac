package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID int64, JWTSecret string) (string, error) {
	claims := UserClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)), // set token expiration to 15 minutes
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "golubac",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString([]byte(JWTSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %s", err.Error())
	}
	return tokenString, nil
}

/*
 * Verifies JWT access/session token
 * return true if its valid and false if it isnt
 * the number it returns is 1 if the JWT is expired and 0 otherwise
 */
func VerifyJWT(JWTToken string, JWTSecret string) (bool, int) {
	parsedToken, err := jwt.ParseWithClaims(
		JWTToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS512 {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(JWTSecret), nil
		},
	)
	if err != nil || parsedToken == nil || !parsedToken.Valid {
		return false, 0
	}
	claims, ok := parsedToken.Claims.(*UserClaims)
	if !ok || claims.ExpiresAt == nil {
		return false, 0
	}
	return claims.ExpiresAt.After(time.Now()), 1
}
