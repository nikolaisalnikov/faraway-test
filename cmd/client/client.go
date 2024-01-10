package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/nikolaisalnikov/faraway-test/internal/hashcash"
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
		log.Println("Error receiving challenge:", err)
		return
	}

	// Solve the Proof of Work
	response := hashcash.SolveProofOfWork(challenge, timestamp, nonce)

	// Send the response to the server
	conn.Write([]byte(response + "\n"))

	// Receive and display the quote from the server
	quote, err := receiveQuote(conn)
	if err != nil {
		log.Println("Error receiving quote:", err)
		return
	}

	log.Println("Received Quote:", quote)
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
