package idgenerator

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

type IdGeneratorV1 struct{}

func NewIdGeneratorV1() *IdGeneratorV1 {
	return &IdGeneratorV1{}
}

func (t *IdGeneratorV1) Generate(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("length must be > 0")
	}
	if length%2 != 0 {
		return "", fmt.Errorf("length must be an even number (each byte -> two hex chars)")
	}

	bytesLen := length / 2
	traceID := make([]byte, bytesLen)

	for {
		if _, err := rand.Read(traceID); err != nil {
			return "", err
		}

		if !IsAllZero(traceID) {
			return hex.EncodeToString(traceID), nil
		}
	}
}

func IsAllZero(data []byte) bool {
	for _, value := range data {
		if value != 0 {
			return false
		}
	}

	return true
}
