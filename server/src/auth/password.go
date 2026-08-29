package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"errors"

	"golang.org/x/crypto/argon2"
)

const (
	saltLength   = 16
	argonTime    = 1
	argonMemory  = 64 * 1024
	argonThreads = 4
	argonKeyLen  = 32
)

func GenerateSalt() (string, error) {
	salt := make([]byte, saltLength)

	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	return base64.RawStdEncoding.EncodeToString(salt), nil
}

func GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}

func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func HashPassword(password string, salt string) (string, error) {
	saltBytes, err := base64.RawStdEncoding.DecodeString(salt)
	if err != nil {
		return "", errors.New("invalid salt")
	}

	hash := argon2.IDKey(
		[]byte(password),
		saltBytes,
		argonTime,
		argonMemory,
		argonThreads,
		argonKeyLen,
	)

	return base64.RawStdEncoding.EncodeToString(hash), nil
}

func VerifyPassword(password, storedHash, storedSalt string) (bool, error) {
	newHash, err := HashPassword(password, storedSalt)
	if err != nil {
		return false, err
	}

	return subtle.ConstantTimeCompare(
		[]byte(storedHash),
		[]byte(newHash),
	) == 1, nil
}
