package main

import (
	"bufio"
	"fmt"
	"github.com/nikolaisalnikov/faraway-test/internal/hashcash"
	"log"
	"math/rand"
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
		log.Println("Error starting server:", err)
		return
	}
	defer listen.Close()
	log.Println("Server started. Listening on :8080")

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
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
		log.Println("Proof of Work failed:", err)
		return
	}

	// Send a random quote
	sendRandomQuote(conn)
}

func performProofOfWork(conn net.Conn) error {
	// Generate a challenge, timestamp, and random nonce
	challenge, timestamp, nonce := hashcash.PerformProofOfWork()

	// Send the challenge, timestamp, and nonce to the client
	conn.Write([]byte(fmt.Sprintf("Challenge: %s %d %d\n", challenge, timestamp, nonce)))

	// Receive the response from the client
	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	// Verify the response, timestamp, and nonce
	if !hashcash.VerifyHashcash(challenge, timestamp, nonce, strings.TrimSpace(response)) {
		return fmt.Errorf("Hashcash verification failed")
	}

	// Hashcash successful
	return nil
}

func sendRandomQuote(conn net.Conn) {
	// Send a random quote to the client
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(quotes))

	conn.Write([]byte(fmt.Sprintf("Quote: %s\n", quotes[randomIndex])))
}
