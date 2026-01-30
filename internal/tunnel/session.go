package tunnel

import (
	"io"
	"net"

	"github.com/hashicorp/yamux"
)

// Join bridges two network connections together
func Join(src, dst net.Conn) {
	defer src.Close()
	defer dst.Close()

	// Channel to coordinate the two goroutines
	done := make(chan struct{}, 2)

	// Start bidirectional copy
	go func() {
		io.Copy(src, dst) // Data from dst to src
		done <- struct{}{}
	}()
	go func() {
		io.Copy(dst, src) // Data from src to dst
		done <- struct{}{}
	}()

	// Wait for at least one direction to close/fail
	<-done
}

// SetupSession initializes a Yamux multiplexing session
func SetupSession(conn net.Conn, isServer bool) (*yamux.Session, error) {
	if isServer {
		// Server-side waits for client to initiate streams
		return yamux.Server(conn, nil)
	}
	// Client-side can initiate or accept streams
	return yamux.Client(conn, nil)
}
