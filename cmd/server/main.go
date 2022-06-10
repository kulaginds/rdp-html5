package main

import (
	"log"
	"net/http"

	"github.com/kulaginds/web-rdp-solution/internal/pkg/handler"
)

func main() {
	if err := startServer(); err != nil {
		log.Fatalln(err)
	}
}

func startServer() error {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("./web")))
	mux.HandleFunc("/connect", handler.Connect)

	log.Println("start web-server on :8080")

	return http.ListenAndServe(":8080", mux)
}
