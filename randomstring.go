package tools

import (
	"math/rand"
)

// RandomString creates a random string of alphanumeric characters of length strlen
func RandomString(strlen int) string {
	// golang 1.20 or higher is required for this function, otherwise, the string isn't random,
	//  older golangs required setting the seed.
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}
