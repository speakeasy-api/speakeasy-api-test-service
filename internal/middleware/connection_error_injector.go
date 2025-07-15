package middleware

import (
	"net"
	"net/http"

	"github.com/lingrino/go-fault"
)

var _ fault.Injector = (*ConnectionErrorInjector)(nil)

// Injects a connection error by closing the connection immediately.
// This simulates a connection reset or close error.
type ConnectionErrorInjector struct {
	// Enable to set SO_LINGER to 0 before closing, which will cause a TCP RST
	// packet on most platforms when closing.
	Reset bool
}

func (i *ConnectionErrorInjector) Handler(_ http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hijacker, ok := w.(http.Hijacker)

		if !ok {
			http.Error(w, "connection hijacking not supported", http.StatusInternalServerError)
			return
		}

		conn, _, err := hijacker.Hijack()

		if err != nil {
			http.Error(w, "failed to hijack connection", http.StatusInternalServerError)
			return
		}

		if tcpConn, ok := conn.(*net.TCPConn); ok && i.Reset {
			_ = tcpConn.SetLinger(0) // Best effort RST on close
		}

		conn.Close()
	})
}
