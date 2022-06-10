package tpkt

import (
	"crypto/tls"
	"fmt"
	"log"
)

func (p *protocol) StartTLS() error {
	tlsConn := tls.Client(p.conn, &tls.Config{
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS10,
		MaxVersion:         tls.VersionTLS13,
	})

	log.Println("TPKT: StartTLS")

	if err := tlsConn.Handshake(); err != nil {
		return fmt.Errorf("TLS handshake: %w", err)
	}

	p.conn = tlsConn

	return nil
}
