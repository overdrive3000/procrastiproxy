package procrastiproxy_test

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"procrastiproxy"
	"testing"
	"time"
)

func TestInitServerListenConnections(t *testing.T) {
	t.Parallel()

	var conErr error
	var conn net.Conn

	c, err := procrastiproxy.New()
	if err != nil {
		t.Fatal(err)
	}
	c.Start()

	for i := 0; i < 5; i++ {
		conn, conErr = net.Dial("tcp", c.FormatAddr())
		if conErr != nil {
			time.Sleep(10 * time.Millisecond)
			continue
		}
		conn.Close()
		break
	}

	if conErr != nil {
		t.Error("connection failed", err)
	}

}

func TestInitServerListenConnectionsFail(t *testing.T) {
	t.Parallel()

	p, _ := procrastiproxy.New()
	p.Port = 80 // Should fail as there is no access to privileged ports

	var conErr error
	var conn net.Conn

	p.Start()
	for i := 0; i < 5; i++ {
		conn, conErr = net.Dial("tcp", p.FormatAddr())
		if conErr != nil {
			time.Sleep(10 * time.Millisecond)
			continue
		}
		conn.Close()
		break
	}

	if conErr == nil {
		t.Error("want error, got nil")
	}
}

func TestFormatAddress(t *testing.T) {
	t.Parallel()

	p, err := procrastiproxy.New()
	if err != nil {
		t.Fatal(err)
	}
	p.Port = 8080

	want := "127.0.0.1:8080"
	got := p.FormatAddr()
	if got != want {
		t.Errorf("want %s, got %s", want, got)
	}

}

func TestProxyServer(t *testing.T) {
	t.Parallel()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))

	defer ts.Close()

	p, _ := procrastiproxy.New()
	p.Start()
	proxyURL, _ := url.Parse(fmt.Sprintf("http://%s", p.FormatAddr()))

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}

	want := http.StatusOK

	var got *http.Response
	var conErr error
	for i := 0; i < 5; i++ {
		got, conErr = client.Get(ts.URL)
		if conErr != nil {
			time.Sleep(10 * time.Millisecond)
			continue
		}
	}
	if conErr != nil {
		t.Fatal(conErr)
	}

	if got.StatusCode != want {
		t.Errorf("want %d, got %d", want, got.StatusCode)
	}

}
