package main

import (
	"crypto/tls"
	"0TLink/internal/auth"
	"0TLink/internal/tunnel"
	"log"
	"net"
)

func main() {
	// 1. Load mTLS config
	tlsConfig, err := auth.GetTLSConfig("certs/server.crt", "certs/server.key", "certs/ca.crt", true)
	if err != nil {
		log.Fatalf("Failed to load TLS config: %v", err)
	}

	// 2. Listen for the Agent (Client)
	clientLn, err := tls.Listen("tcp", ":7000", tlsConfig)
	if err != nil {
		log.Fatalf("Failed to start TLS listener: %v", err)
	}
	log.Println("Control Plane listening on :7000")

	// 3. Wait for Agent
	conn, err := clientLn.Accept()
	if err != nil {
		log.Printf("Accept error: %v", err)
		return
	}
	log.Println("Agent connected via mTLS")

	// 4. Setup Yamux Session
	session, err := tunnel.SetupSession(conn, true)
	if err != nil {
		log.Fatalf("Yamux error: %v", err)
	}

	// 5. Public Gateway
	publicLn, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Public listener error: %v", err)
	}
	log.Println("Public Gateway open on :8080")

	for {
		userConn, err := publicLn.Accept()
		if err != nil {
			continue
		}

		stream, err := session.Open()
		if err != nil {
			userConn.Close()
			continue
		}

		log.Println("Tunneling new request...")
		go tunnel.Join(userConn, stream)
	}
}
