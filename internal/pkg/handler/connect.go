package handler

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"

	"github.com/kulaginds/web-rdp-solution/internal/pkg/rdp"
)

const (
	webSocketReadBufferSize  = 8192
	webSocketWriteBufferSize = 8192 * 2
	rdpMaxPackageSize        = 16 * 1024
)

const (
	width  = 1280
	height = 800

	host     = "192.168.1.2:3389"
	user     = "Doc"
	password = "1qazXSW@"
)

func Connect(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  webSocketReadBufferSize,
		WriteBufferSize: webSocketWriteBufferSize,
		CheckOrigin: func(r *http.Request) bool {
			return true // TODO: SECURITY: проверить хост
		},
	}
	protocol := r.Header.Get("Sec-Websocket-Protocol")

	wsConn, err := upgrader.Upgrade(w, r, http.Header{
		"Sec-Websocket-Protocol": {protocol},
	})
	if err != nil {
		log.Println(fmt.Errorf("upgrade websocket: %w", err))
		w.Write([]byte(`{"error": "init websocket"}`))
		return
	}

	defer func() {
		if err = wsConn.Close(); err != nil {
			log.Println(fmt.Errorf("error closing websocket: %w", err))
		}
	}()

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	domain := strings.Split(host, ":")

	rdpClient, err := rdp.NewClient(host, domain[0], user, password, width, height)
	if err != nil {
		log.Println(fmt.Errorf("rdp init: %w", err))
		wsConn.WriteMessage(1, []byte(`{"error": "rdp connect"}`))
		return
	}

	defer rdpClient.Close()

	if err = rdpClient.Connect(); err != nil {
		log.Println(fmt.Errorf("rdp connect: %w", err))
		wsConn.WriteMessage(1, []byte(`{"error": "rdp connect"}`))
		return
	}

	log.Println("begin proxying")

	go wsToRdp(ctx, wsConn, rdpClient, cancel)
	rdpToWs(ctx, rdpClient, wsConn)
}

func wsToRdp(ctx context.Context, wsConn *websocket.Conn, rdpConn io.Writer, cancel context.CancelFunc) {
	defer func() {
		log.Println("wsToRdp done")
		cancel()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default: // pass
		}

		_, data, err := wsConn.ReadMessage()
		if err != nil {
			log.Println(fmt.Errorf("error reading message from ws: %w", err))

			return
		}

		if _, err = rdpConn.Write(data); err != nil {
			log.Println(fmt.Errorf("failed writing to guacd: %w", err))

			return
		}
	}
}

func rdpToWs(ctx context.Context, rdpConn io.Reader, wsConn *websocket.Conn) {
	defer func() {
		log.Println("rdpToWs done")
	}()

	var err error

	buf := bytes.NewBuffer(make([]byte, 0, rdpMaxPackageSize))

	for {
		select {
		case <-ctx.Done():
			return
		default: // pass
		}

		if _, err = buf.ReadFrom(rdpConn); err != nil {
			log.Println(fmt.Errorf("failed to buffer guacd to ws: %w", err))

			return
		}

		if err = wsConn.WriteMessage(1, buf.Bytes()); err != nil {
			if err == websocket.ErrCloseSent {
				log.Println("sent to closed websocket")

				return
			}

			log.Println(fmt.Errorf("failed sending message to ws: %w", err))

			return
		}

		buf.Reset()
	}
}
