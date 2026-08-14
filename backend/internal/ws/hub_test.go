package ws

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func dialTestServer(t *testing.T, server *httptest.Server) *websocket.Conn {
	t.Helper()
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("failed to dial test server: %v", err)
	}
	t.Cleanup(func() { conn.Close() })
	return conn
}

func TestHubBroadcastDeliversToConnectedClient(t *testing.T) {
	hub := NewHub()
	server := httptest.NewServer(http.HandlerFunc(hub.ServeWS))
	t.Cleanup(server.Close)

	conn := dialTestServer(t, server)

	time.Sleep(50 * time.Millisecond)
	hub.Broadcast([]byte("hello"))

	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, msg, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("expected to receive broadcast message, got error: %v", err)
	}
	if string(msg) != "hello" {
		t.Errorf("got message %q, want %q", msg, "hello")
	}
}

func TestHubBroadcastReachesMultipleClients(t *testing.T) {
	hub := NewHub()
	server := httptest.NewServer(http.HandlerFunc(hub.ServeWS))
	t.Cleanup(server.Close)

	conn1 := dialTestServer(t, server)
	conn2 := dialTestServer(t, server)

	time.Sleep(50 * time.Millisecond)
	hub.Broadcast([]byte("hi all"))

	for i, conn := range []*websocket.Conn{conn1, conn2} {
		conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		_, msg, err := conn.ReadMessage()
		if err != nil {
			t.Fatalf("client %d: expected broadcast message, got error: %v", i, err)
		}
		if string(msg) != "hi all" {
			t.Errorf("client %d: got message %q, want %q", i, msg, "hi all")
		}
	}
}

func TestHubRemovesClientAfterDisconnect(t *testing.T) {
	hub := NewHub()
	server := httptest.NewServer(http.HandlerFunc(hub.ServeWS))
	t.Cleanup(server.Close)

	conn := dialTestServer(t, server)
	time.Sleep(50 * time.Millisecond)

	conn.Close()
	// Give readPump time to notice the closed connection and remove the client.
	time.Sleep(50 * time.Millisecond)

	hub.mu.Lock()
	clientCount := len(hub.clients)
	hub.mu.Unlock()

	if clientCount != 0 {
		t.Errorf("got %d clients after disconnect, want 0", clientCount)
	}
}
