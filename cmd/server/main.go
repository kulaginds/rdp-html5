package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/kulaginds/rdp-html5/internal/pkg/handler"
)

func main() {
	if err := startServer(); err != nil {
		log.Fatalln(err)
	}
}

func startServer() error {
	http.Handle("/", http.FileServer(http.Dir("./web")))
	http.HandleFunc("/connect", handler.Connect)

	log.Println("start web-server on :8080")

	return http.ListenAndServe(":8080", nil)
}
