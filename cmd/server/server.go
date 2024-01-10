package main

import (
	"bufio"
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"math/big"
	"net"
	"strings"
	"time"
)

const (
	difficulty = 4
)

var quotes = []string{
	"The only true wisdom is in knowing you know nothing. - Socrates",
	"Life is what happens when you're busy making other plans. - John Lennon",
	"Success is not final, failure is not fatal: It is the courage to continue that counts. - Winston Churchill",
	// Add more quotes as needed
}

func main() {
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listen.Close()
	fmt.Println("Server started. Listening on :8080")

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Perform Proof of Work with Hashcash
	err := performProofOfWork(conn)
	if err != nil {
		fmt.Println("Proof of Work failed:", err)
		return
	}

	// Send a random quote
	sendRandomQuote(conn)
}

func performProofOfWork(conn net.Conn) error {
	// Generate a challenge, timestamp, and random nonce
	challenge := generateChallenge()
	timestamp := time.Now().Unix()
	nonce := generateRandomNonce()

	// Send the challenge, timestamp, and nonce to the client
	conn.Write([]byte(fmt.Sprintf("Challenge: %s %d %d\n", challenge, timestamp, nonce)))

	// Receive the response from the client
	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	// Verify the response, timestamp, and nonce
	if !verifyHashcash(challenge, timestamp, nonce, strings.TrimSpace(response)) {
		return fmt.Errorf("Hashcash verification failed")
	}

	// Hashcash successful
	return nil
}

func generateChallenge() string {
	// Generate a random challenge string for illustration purposes
	return "abcd"
}

func generateRandomNonce() int {
	// Generate a random nonce
	nonce, err := rand.Int(rand.Reader, big.NewInt(10000))
	if err != nil {
		panic(err) // handle error appropriately in a real application
	}
	return int(nonce.Int64())
}

func verifyHashcash(challenge string, timestamp int64, nonce int, response string) bool {
	// Calculate the hash of the concatenated challenge, timestamp, nonce, and response
	hashInput := fmt.Sprintf("%s:%d:%d:%s", challenge, timestamp, nonce, response)
	hash := sha1.Sum([]byte(hashInput))
	hashString := hex.EncodeToString(hash[:])

	fmt.Printf("Challenge: %s\nTimestamp: %d\nNonce: %d\nResponse: %s\nCalculated Hash: %s\n", challenge, timestamp, nonce, strings.TrimSpace(response), hashString)
	// Check if the hash meets the difficulty requirements
	return strings.HasPrefix(hashString, strings.Repeat("0", difficulty))
}

func sendRandomQuote(conn net.Conn) {
	// Send a random quote to the client
	randomIndex := 0 // You can use a random number generator for a real application
	conn.Write([]byte(fmt.Sprintf("Quote: %s\n", quotes[randomIndex])))
}
