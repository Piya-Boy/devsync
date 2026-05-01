package transport

import (
	"fmt"
	"net"
	"time"

	"github.com/Piya-Boy/devsync/config"
)

const (
	sshPort    = 22
	dialTimeout = 3 * time.Second
)

// Detect returns the appropriate Transport for the given server.
// It checks if port 22 is reachable; if so, SSH is used; otherwise SMB is used.
func Detect(server *config.Server) (Transport, error) {
	if IsPortOpen(server.Host, sshPort) {
		return SSHTransport{}, nil
	}
	return SMBTransport{}, nil
}

// IsPortOpen returns true if the given host:port is reachable within the timeout.
func IsPortOpen(host string, port int) bool {
	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", addr, dialTimeout)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}
