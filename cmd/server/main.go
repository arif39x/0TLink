package main

import (
	"0TLink/internal/auth"
	"0TLink/internal/tunnel"
	"crypto/tls"
	"log"
	"net"
)

func main() {
	tlsConfig, err := auth.GetTLSConfig("certs/server.crt", "certs/server.key", "certs/ca.crt", true)
	if err != nil {
		log.Fatalf("TLS config error: %v", err)
	}

	clientLn, err := tls.Listen("tcp", ":7000", tlsConfig)
	if err != nil {
		log.Fatalf("Control plane error: %v", err)
	}

	publicLn, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Gateway error: %v", err)
	}

	log.Println("Control Plane: :7000 | Public Gateway: :8080")

	for {
		conn, err := clientLn.Accept()
		if err != nil {
			log.Printf("Accept error: %v", err)
			continue
		}

		session, err := tunnel.SetupSession(conn, true)
		if err != nil {
			log.Printf("Yamux error: %v", err)
			conn.Close()
			continue
		}

		go func() {
			for {
				userConn, err := publicLn.Accept()
				if err != nil {
					return
				}

				stream, err := session.Open()
				if err != nil {
					userConn.Close()
					if session.IsClosed() {
						return
					}
					continue
				}

				go tunnel.Join(userConn, stream)
			}
		}()
	}
}