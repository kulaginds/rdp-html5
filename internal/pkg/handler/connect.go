package handler

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/kulaginds/web-rdp-solution/internal/pkg/rdp"
	"github.com/kulaginds/web-rdp-solution/internal/pkg/rdp/fastpath"
)

const (
	webSocketReadBufferSize  = 8192
	webSocketWriteBufferSize = 8192 * 2
	rdpMaxPackageSize        = 16 * 1024
	tpktAndX224HeadersLen    = 4 + 3
)

const (
	host     = "192.168.1.2:3389"
	user     = "Doc"
	password = "1qazXSW@"
)

type rdpConn interface {
	GetUpdate() (*fastpath.UpdatePDU, error)
	SendInputEvent(data []byte) error
}

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

	var (
		width  uint16
		height uint16
	)

	width, height, err = wsReadCanvasDimensions(wsConn)
	if err != nil {
		log.Println(fmt.Errorf("ws read canvas dimensions: %w", err))
		return
	}

	rdpClient, err := rdp.NewClient(host, user, password, width, height)
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

func wsReadCanvasDimensions(wsConn *websocket.Conn) (uint16, uint16, error) {
	_, data, err := wsConn.ReadMessage()
	if err != nil {
		return 0, 0, err
	}

	if len(data) != 4 {
		return 0, 0, errors.New("bad data size")
	}

	width := binary.LittleEndian.Uint16(data[0:2])
	height := binary.LittleEndian.Uint16(data[2:4])

	return width, height, nil
}

func wsToRdp(ctx context.Context, wsConn *websocket.Conn, rdpConn rdpConn, cancel context.CancelFunc) {
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

		if err = rdpConn.SendInputEvent(data); err != nil {
			log.Println(fmt.Errorf("failed writing to rdp: %w", err))

			return
		}
	}
}

func rdpToWs(ctx context.Context, rdpConn rdpConn, wsConn *websocket.Conn) {
	defer func() {
		log.Println("rdpToWs done")
	}()

	var (
		update *fastpath.UpdatePDU
		err    error
	)

	for {
		select {
		case <-ctx.Done():
			return
		default: // pass
		}

		update, err = rdpConn.GetUpdate()
		if err != nil {
			log.Println(fmt.Errorf("get update: %w", err))

			return
		}

		if err = wsConn.WriteMessage(2, update.Data); err != nil {
			if err == websocket.ErrCloseSent {
				log.Println("sent to closed websocket")

				return
			}

			log.Println(fmt.Errorf("failed sending message to ws: %w", err))

			return
		}
	}
}
