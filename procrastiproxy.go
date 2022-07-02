package procrastiproxy

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/phayes/freeport"
)

type Proxy struct {
	HTTPClient *http.Client
	Port       int
	Addr       string
}

func New() (*Proxy, error) {
	port, err := freeport.GetFreePort()
	if err != nil {
		return nil, err
	}
	return &Proxy{
		HTTPClient: &http.Client{},
		Port:       port,
		Addr:       "127.0.0.1",
	}, nil
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.RequestURI = ""
	log.Println("Requesting", r.Host)
	resp, err := p.HTTPClient.Do(r)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		log.Fatal(err)
	}
	defer resp.Body.Close()
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func (p *Proxy) FormatAddr() string {
	return fmt.Sprintf("%s:%d", p.Addr, p.Port)
}

func (p *Proxy) Start() {
	addr := p.FormatAddr()
	log.Printf("Starting proxy on %s", addr)
	go func() error {
		if err := http.ListenAndServe(addr, p); err != nil {
			return err
		}
		return nil
	}()
}
