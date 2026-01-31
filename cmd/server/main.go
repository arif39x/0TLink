package main

import (
	"0TLink/internal/auth"
	"0TLink/internal/tunnel"
	"context"
	"crypto/tls"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	tlsConfig, err := auth.GetTLSConfig(
		"certs/server.crt",
		"certs/server.key",
		"certs/ca.crt",
		true,
	)
	if err != nil {
		log.Fatalf("TLS config error: %v", err)
	}

	// Enforce modern TLS + mTLS strictly
	tlsConfig.MinVersion = tls.VersionTLS13
	tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert

	controlLn, err := tls.Listen("tcp", ":7000", tlsConfig)
	if err != nil {
		log.Fatalf("Control plane error: %v", err)
	}
	defer controlLn.Close()

	publicLn, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Gateway error: %v", err)
	}
	defer publicLn.Close()

	log.Println("Control Plane (mTLS): :7000")
	log.Println("Public Gateway: :8080")

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down control plane")
			return
		default:
		}

		conn, err := controlLn.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				log.Printf("Temporary accept error: %v", err)
				time.Sleep(100 * time.Millisecond)
				continue
			}
			log.Printf("Accept error: %v", err)
			continue
		}

		go handleClient(conn, publicLn)
	}
}

func handleClient(conn net.Conn, publicLn net.Listener) {
	defer conn.Close()

	tlsConn, ok := conn.(*tls.Conn)
	if !ok {
		log.Println("Non-TLS connection rejected")
		return
	}

	if err := tlsConn.Handshake(); err != nil {
		log.Printf("TLS handshake failed: %v", err)
		return
	}

	state := tlsConn.ConnectionState()
	if len(state.PeerCertificates) == 0 {
		log.Println("Client presented no certificate")
		return
	}

	clientID := state.PeerCertificates[0].Subject.CommonName
	log.Printf("Client authenticated: %s", clientID)

	session, err := tunnel.SetupSession(tlsConn, true)
	if err != nil {
		log.Printf("Yamux setup failed for %s: %v", clientID, err)
		return
	}

	for {
		userConn, err := publicLn.Accept()
		if err != nil {
			if session.IsClosed() {
				log.Printf("Session closed for %s", clientID)
				return
			}
			log.Printf("Public accept error: %v", err)
			continue
		}

		stream, err := session.Open()
		if err != nil {
			userConn.Close()
			if session.IsClosed() {
				return
			}
			log.Printf("Stream open failed for %s: %v", clientID, err)
			continue
		}

		go tunnel.Join(userConn, stream)
	}
}
