package main

import (
	"bufio"
	"fmt"
	faraway_test "github.com/nikolaisalnikov/faraway-test"
	"github.com/nikolaisalnikov/faraway-test/internal/hashcash"
	"log"
	"math/rand"
	"net"
	"strings"
	"time"
)

var (
	quotes = []string{
		"The only true wisdom is in knowing you know nothing. - Socrates\n",
		"Life is what happens when you're busy making other plans. - John Lennon\n",
		"Success is not final, failure is not fatal: It is the courage to continue that counts. - Winston Churchill\n",
		"The greatest glory in living lies not in never falling, but in rising every time we fall. - Nelson Mandela\n",
		"The way to get started is to quit talking and begin doing. - Walt Disney\n",
		"Don't cry because it's over, smile because it happened. - Dr. Seuss\n",
		"Life is either a daring adventure or nothing at all. - Helen Keller\n",
		"The only limit to our realization of tomorrow will be our doubts of today. - Franklin D. Roosevelt\n",
		"The future belongs to those who believe in the beauty of their dreams. - Eleanor Roosevelt\n",
		"The purpose of our lives is to be happy. - Dalai Lama\n",
	}
	config = faraway_test.LoadConfig()
)

func main() {
	listen, err := net.Listen("tcp", ":"+config.Port)
	if err != nil {
		log.Println("Error starting server:", err)
		return
	}
	defer listen.Close()
	log.Println("Server started. Listening on 0.0.0.0:" + config.Port)

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
		conn.Write([]byte(fmt.Sprintf("Error: Proof of Work verification failed\n")))
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
	if !hashcash.VerifyHashcash(challenge, timestamp, nonce, strings.TrimSpace(response), config.Difficulty) {
		return fmt.Errorf("Hashcash verification failed")
	}

	// Hashcash successful
	return nil
}

func sendRandomQuote(conn net.Conn) {
	// Send a random quote to the client
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(quotes))

	conn.Write([]byte(quotes[randomIndex]))
}
