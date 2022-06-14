package main

import (
	"context"
	"fmt"
	"log"

	"github.com/kulaginds/web-rdp-solution/internal/pkg/rdp"
)

func main() {
	rdpClient, err := rdp.NewClient("192.168.1.2:3389", "Doc", "1qazXSW@", 1280, 800)
	if err != nil {
		log.Println(fmt.Errorf("rdp init: %w", err))

		return
	}
	defer rdpClient.Close()

	err = rdpClient.Connect()
	if err != nil {
		log.Println(fmt.Errorf("rdp connect: %w", err))

		return
	}

	if err = rdpClient.CapabilitiesExchange(); err != nil {
		log.Println(fmt.Errorf("rdp caps: %w", err))

		return
	}

	//if err = rdpClient.ConnectionFinalization(); err != nil {
	//	log.Println(fmt.Errorf("rdp finalization: %w", err))
	//
	//	return
	//}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err = rdpClient.ReaderLoop(ctx); err != nil {
		log.Println(fmt.Errorf("reader loop: %w", err))

		return
	}
}
