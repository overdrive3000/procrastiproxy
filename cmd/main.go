package main

import (
	"flag"
	"log"
	"net/http"
	"procrastiproxy"
)

func main() {
	var addr = flag.String("addr", "127.0.0.1:8080", "Proxy bind address")
	flag.Parse()

	handler := procrastiproxy.New()
	log.Println("Starting proxy on", *addr)

	if err := http.ListenAndServe(*addr, handler); err != nil {
		log.Fatal(err)
	}
}
