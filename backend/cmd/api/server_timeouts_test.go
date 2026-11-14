package main

import (
	"net"
	"net/http"
	"testing"
	"time"
)

func TestNewServer_HasEveryTimeout(t *testing.T) {
	server := newServer(http.NotFoundHandler(), defaultTimeouts)

	for name, timeout := range map[string]time.Duration{
		"read header": server.ReadHeaderTimeout,
		"read":        server.ReadTimeout,
		"write":       server.WriteTimeout,
		"idle":        server.IdleTimeout,
	} {
		if timeout <= 0 {
			t.Errorf("%s timeout = %v, want a limit", name, timeout)
		}
	}
}

func TestNewServer_CutsAClientThatNeverFinishesTheHeaders(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	timeouts := defaultTimeouts
	timeouts.readHeader = 100 * time.Millisecond
	server := newServer(http.NotFoundHandler(), timeouts)
	go func() { _ = server.Serve(listener) }()
	t.Cleanup(func() { server.Close() })

	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()
	if _, err := conn.Write([]byte("GET /api/health HTTP/1.1\r\nHost: x\r\n")); err != nil {
		t.Fatalf("write partial header: %v", err)
	}

	if err := conn.SetReadDeadline(time.Now().Add(3 * time.Second)); err != nil {
		t.Fatalf("set deadline: %v", err)
	}
	buffer := make([]byte, 512)
	_, readErr := conn.Read(buffer)
	for readErr == nil {
		_, readErr = conn.Read(buffer)
	}
	if netErr, isNet := readErr.(net.Error); isNet && netErr.Timeout() {
		t.Errorf("the server kept the connection open for 3s, want it cut after %v", timeouts.readHeader)
	}
}
