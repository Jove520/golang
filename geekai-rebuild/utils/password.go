package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"math/big"
)

func GenSalt(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes)[:length], nil
}

func GenPassword(password, salt string) string {
	sum := sha256.Sum256([]byte(password + salt))
	return hex.EncodeToString(sum[:])
}

func CheckPassword(rawPassword, salt, hashedPassword string) bool {
	return GenPassword(rawPassword, salt) == hashedPassword
}

func GenNumericCode(length int) (string, error) {
	code := make([]byte, length)

	for i := range code {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}

		code[i] = byte('0' + n.Int64())
	}

	return string(code), nil
}
