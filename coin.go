package main

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"time"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const tails = "tails"
const heads = "heads"

// get random server seed
func generateSeed() string {
	length := 20
	result := make([]byte, length)
	rand.Seed(time.Now().UnixNano())

	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

// perform coin flipping return 'tails' or 'heads' as a result
func flipCoin() string {
	rand.Seed(time.Now().UnixNano())
	result := rand.Intn(2)
	if result == 0 {
		return heads
	} else if result == 1 {
		return tails
	}
	return ""
}

// compute hash of flipping result, user can check it when protocol finished
func getFlipHash(flip string) string {
	sum := sha256.Sum256([]byte(flip))
	return hex.EncodeToString(sum[:])
}
