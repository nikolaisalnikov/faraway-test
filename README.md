# Word of Wisdom Server

Word of Wisdom Server is a simple TCP server implemented in Go that provides a quote after Proof of Work verification.

## Features

- Protected from DDOS attacks with Proof of Work
- Sends a random quote after successful Proof of Work
- In client there is 50% chance to send correct proof cmd/client/client.go:33:44

## Getting Started

1. Clone the repository:

   ```bash
   git clone https://github.com/nikolaisalnikov/faraway-test.git
   ```

2. Navigate to project firectory

   ```bash
   cd faraway-test
   ```

## Configuration
You can configure the server's port and difficulty by modifying the config.json file in the root directory.

```json 
    {
      "port": "8080",
      "difficulty": 4
    } 
```

## Usage

1. Build the Docker images:

   ```bash
   docker-compose build
   ```

2. Run the Docker containers:

   ```bash
   docker-compose up
   ```

3. For manual running 

   ```bash
   go run cmd/server/server.go
   ```

   ```bash
   go run cmd/client/client.go
   ```

## Dependencies
* Go 1.21.4
* Docker
* Docker Compose

