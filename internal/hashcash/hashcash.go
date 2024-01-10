package hashcash

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"math/big"
	mathrand "math/rand"
	"strings"
	"time"
)

// GenerateChallenge generates a random challenge string.
func GenerateChallenge() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, 4)
	for i := range result {
		result[i] = charset[mathrand.Intn(len(charset))]
	}
	return string(result)
}

// GenerateRandomNonce generates a random nonce.
func GenerateRandomNonce() int {
	// Generate a random nonce
	nonce, err := rand.Int(rand.Reader, big.NewInt(10000))
	if err != nil {
		panic(err) // handle error appropriately in a real application
	}
	return int(nonce.Int64())
}

// PerformProofOfWork performs Hashcash proof of work and returns the challenge, timestamp, and nonce.
func PerformProofOfWork() (string, int64, int) {
	// Generate a challenge, timestamp, and random nonce
	challenge := GenerateChallenge()
	timestamp := time.Now().Unix()
	nonce := GenerateRandomNonce()

	return challenge, timestamp, nonce
}

// VerifyHashcash verifies the Hashcash proof of work.
func VerifyHashcash(challenge string, timestamp int64, nonce int, response string, difficulty int) bool {
	// Calculate the hash of the concatenated challenge, timestamp, nonce, and response
	hashInput := fmt.Sprintf("%s:%d:%d:%s", challenge, timestamp, nonce, response)
	hash := sha1.Sum([]byte(hashInput))
	hashString := hex.EncodeToString(hash[:])

	// Check if the hash meets the difficulty requirements
	return strings.HasPrefix(hashString, strings.Repeat("0", difficulty))
}

func SolveProofOfWork(challenge string, timestamp int64, nonce int, difficulty int) string {
	for i := 0; ; i++ {
		response := fmt.Sprintf("%s%d", challenge, i)
		if VerifyHashcash(challenge, timestamp, nonce, response, difficulty) {
			return response
		}
	}
}
