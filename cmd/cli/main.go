package main

import (
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

	err = rdpClient.Connect()
	if err != nil {
		log.Println(fmt.Errorf("rdp connect: %w", err))

		return
	}
}
