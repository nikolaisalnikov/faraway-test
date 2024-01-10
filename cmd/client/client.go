package main

import (
	"bufio"
	"fmt"
	faraway_test "github.com/nikolaisalnikov/faraway-test"
	"log"
	"net"
	"strings"
	"time"

	"github.com/nikolaisalnikov/faraway-test/internal/hashcash"
)

var config = faraway_test.LoadConfig()

func main() {
	for {
		conn, err := net.Dial("tcp", "word-of-wisdom-server:"+config.Port)
		if err != nil {
			log.Println("Error connecting to server:", err)
			return
		}

		// Receive the challenge, timestamp, and nonce from the server
		challenge, timestamp, nonce, err := receiveChallenge(conn)
		if err != nil {
			log.Println("Error receiving challenge:", err)
			return
		}

		// Solve the Proof of Work
		response := hashcash.SolveProofOfWork(challenge, timestamp, nonce, config.Difficulty)

		// Send the response to the server
		conn.Write([]byte(response + "\n"))

		// Receive and display the quote from the server
		serverResponse, err := receiveResponse(conn)
		if err != nil {
			log.Println("Error receiving quote:", err)
			return
		}

		log.Print("Server response: ", serverResponse)
		
		if strings.HasPrefix(serverResponse, "Error:") {
			log.Println("Server returned an error: ", serverResponse)
			return
		}

		conn.Close()
		time.Sleep(2 * time.Second)
	}
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

func receiveResponse(conn net.Conn) (string, error) {
	// Receive the quote from the server
	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	// Extract and return the quote
	return response, nil
}
