package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net"
	"strings"
)

const (
	difficulty = 4
)

func main() {
	conn, err := net.Dial("tcp", "word-of-wisdom-server:8080")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	// Receive the challenge, timestamp, and nonce from the server
	challenge, timestamp, nonce, err := receiveChallenge(conn)
	if err != nil {
		fmt.Println("Error receiving challenge:", err)
		return
	}

	// Solve the Proof of Work
	response := solveProofOfWork(challenge, timestamp, nonce)

	// Send the response to the server
	conn.Write([]byte(response + "\n"))

	// Receive and display the quote from the server
	quote, err := receiveQuote(conn)
	if err != nil {
		fmt.Println("Error receiving quote:", err)
		return
	}

	fmt.Println("Received Quote:", quote)
}

func receiveChallenge(conn net.Conn) (string, int64, int, error) {
	// Receive the challenge, timestamp, and nonce from the server
	reader := bufio.NewReader(conn)
	challengeWithNonce, err := reader.ReadString('\n')
	if err != nil {
		return "", 0, 0, err
	}

	// Extract challenge, timestamp, and nonce
	parts := strings.Split(strings.TrimPrefix(challengeWithNonce, "Challenge: "), " ")
	challenge := parts[0]
	timestamp := parseInt64(parts[1])
	nonce := parseInt(parts[2])

	return challenge, timestamp, nonce, nil
}

func parseInt64(s string) int64 {
	var val int64
	_, err := fmt.Sscanf(s, "%d", &val)
	if err != nil {
		return 0
	}
	return val
}

func parseInt(s string) int {
	var val int
	_, err := fmt.Sscanf(s, "%d", &val)
	if err != nil {
		return 0
	}
	return val
}

func solveProofOfWork(challenge string, timestamp int64, nonce int) string {
	for i := 0; ; i++ {
		response := fmt.Sprintf("%s%d", challenge, i)
		if verifyProofOfWork(challenge, timestamp, nonce, response) {
			fmt.Println(response)
			return response
		}
	}
}

func receiveQuote(conn net.Conn) (string, error) {
	// Receive the quote from the server
	reader := bufio.NewReader(conn)
	quote, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	// Extract and return the quote
	return strings.TrimPrefix(quote, "Quote: "), nil
}

func verifyProofOfWork(challenge string, timestamp int64, nonce int, response string) bool {
	// Calculate the hash of the concatenated challenge, timestamp, nonce, and response
	hashInput := fmt.Sprintf("%s:%d:%d:%s", challenge, timestamp, nonce, response)
	hash := sha1.Sum([]byte(hashInput))
	hashString := hex.EncodeToString(hash[:])

	fmt.Printf("Challenge: %s\nTimestamp: %d\nNonce: %d\nResponse: %s\nCalculated Hash: %s\n", challenge, timestamp, nonce, response, hashString)
	// Check if the hash meets the difficulty requirements
	return strings.HasPrefix(hashString, strings.Repeat("0", difficulty))
}
