package rdp

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
)

func (c *client) StartTLS() error {
	tlsConn := tls.Client(c.conn, &tls.Config{
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS10,
		MaxVersion:         tls.VersionTLS13,
	})

	log.Println("TPKT: StartTLS")

	if err := tlsConn.Handshake(); err != nil {
		return fmt.Errorf("TLS handshake: %w", err)
	}

	c.conn = tlsConn
	c.buffReader = bufio.NewReaderSize(c.conn, readBufferSize)

	return nil
}
