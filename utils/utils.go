package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateID(size int) (string, error) {
	id := make([]byte, 4)
	_, err := rand.Read(id)

	if err != nil {
		return "", err
	}
	return hex.EncodeToString(id), nil

}
