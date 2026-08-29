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
	saltLength = 16

	argonTime    = 1
	argonMemory  = 64 * 1024
	argonThreads = 4
	argonKeyLen  = 32
)

func GenerateSalt() ([]byte, error) {
	salt := make([]byte, saltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
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

func HashPassword(password string, salt []byte) []byte {
	return argon2.IDKey(
		[]byte(password),
		salt,
		argonTime,
		argonMemory,
		argonThreads,
		argonKeyLen,
	)
}
func VerifyPassword(password string, storedHash string, storedSalt string) (bool, error) {
	salt, err := base64.RawStdEncoding.DecodeString(storedSalt)
	if err != nil {
		return false, errors.New("invalid stored salt")
	}

	hash, err := base64.RawStdEncoding.DecodeString(storedHash)
	if err != nil {
		return false, errors.New("invalid stored hash")
	}

	newHash := HashPassword(password, salt)

	return subtle.ConstantTimeCompare(hash, newHash) == 1, nil
}
