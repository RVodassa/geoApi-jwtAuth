package main

import "test/server"

func main() {
	s := server.NewServer()
	s.Start()
}

// curl -X POST http://localhost:8080/api/register -d '{"email": "user@example.com", "password": "secret"}' -H "Content-Type: application/json"
