package utils

import (
	"crypto/rand"
	"encoding/hex"
	"os"
)

// generateRandomString generates a random string of specified length
func GenerateRandomString(length int) string {
	bytes := make([]byte, length/2)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)[:length]
}

// writeToFile writes content to a file
func WriteToFile(filename, content string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	return err
}

// fileExists checks if a file exists
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}