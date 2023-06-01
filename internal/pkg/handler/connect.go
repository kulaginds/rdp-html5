package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/kulaginds/web-rdp-solution/internal/pkg/rdp/pdu"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"

	"github.com/kulaginds/web-rdp-solution/internal/pkg/rdp"
	"github.com/kulaginds/web-rdp-solution/internal/pkg/rdp/fastpath"
)

const (
	webSocketReadBufferSize  = 8192
	webSocketWriteBufferSize = 8192 * 2
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

		return
	}

	defer func() {
		if err = wsConn.Close(); err != nil {
			log.Println(fmt.Errorf("error closing websocket: %w", err))
		}
	}()

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	width, err := strconv.Atoi(r.URL.Query().Get("width"))
	if err != nil {
		log.Println(fmt.Errorf("get width: %w", err))

		return
	}

	height, err := strconv.Atoi(r.URL.Query().Get("height"))
	if err != nil {
		log.Println(fmt.Errorf("get height: %w", err))

		return
	}

	host := r.URL.Query().Get("host")
	user := r.URL.Query().Get("user")
	password := r.URL.Query().Get("password")

	rdpClient, err := rdp.NewClient(host, user, password, width, height)
	if err != nil {
		log.Println(fmt.Errorf("rdp init: %w", err))

		return
	}
	defer rdpClient.Close()

	// TODO: implement
	//rdpClient.SetRemoteApp("C:\\agent\\agent.exe", ".\\Downloads\\cbct1.zip", "C:\\Users\\Doc")
	//rdpClient.SetRemoteApp("explore", "", "")

	if err = rdpClient.Connect(); err != nil {
		log.Println(fmt.Errorf("rdp connect: %w", err))

		return
	}

	log.Println("begin proxying")

	go wsToRdp(ctx, wsConn, rdpClient, cancel)
	rdpToWs(ctx, rdpClient, wsConn)
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
			if strings.HasSuffix(err.Error(), "use of closed network connection") {
				return
			}

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
		switch {
		case err == nil: // pass
		case errors.Is(err, pdu.ErrDeactiateAll):
			log.Println("deactivate all")

			return

		default:
			log.Println(fmt.Errorf("get update: %w", err))

			return
		}

		if err = wsConn.WriteMessage(websocket.BinaryMessage, update.Data); err != nil {
			if err == websocket.ErrCloseSent {
				log.Println("sent to closed websocket")

				return
			}

			log.Println(fmt.Errorf("failed sending message to ws: %w", err))

			return
		}
	}
}
