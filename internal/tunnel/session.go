package tunnel

import (
	"io"
	"net"

	"github.com/hashicorp/yamux"
)

func Join(src, dst net.Conn) {
	defer src.Close()
	defer dst.Close()

	done := make(chan struct{}, 2)

	go func() {
		io.Copy(src, dst)
		done <- struct{}{}
	}()
	go func() {
		io.Copy(dst, src)
		done <- struct{}{}
	}()

	<-done
}

func SetupSession(conn net.Conn, isServer bool) (*yamux.Session, error) {
	if isServer {
		return yamux.Server(conn, nil)
	}
	return yamux.Client(conn, nil)
}
