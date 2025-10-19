package web

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestServer_NewServer(t *testing.T) {
	tests := []struct {
		name string
		port int
	}{
		{"StandardPort", 8080},
		{"AlternativePort", 3000},
		{"HighPort", 9999},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := NewServer(tt.port)
			if server.port != tt.port {
				t.Errorf("Expected port %d, got %d", tt.port, server.port)
			}
			if server.mux == nil {
				t.Error("Server mux should not be nil")
			}
		})
	}
}

func TestServer_SetupRoutes(t *testing.T) {
	server := NewServer(8080)

	// Test that routes are properly configured
	testRoutes := []string{
		"/",
		"/api/timer/start",
		"/api/timer/stop",
		"/ws",
	}

	for _, route := range testRoutes {
		t.Run(route, func(t *testing.T) {
			// Try to match the route pattern
			req, err := http.NewRequest("GET", route, nil)
			if err != nil {
				t.Fatal(err)
			}

			// Use a custom ResponseWriter to capture the response
			rw := &testResponseWriter{}
			server.mux.ServeHTTP(rw, req)

			// Check if route exists (not 404)
			if rw.statusCode == http.StatusNotFound {
				t.Errorf("Route %s should be configured", route)
			}
		})
	}
}

func TestServer_StartAndStop(t *testing.T) {
	// Find an available port
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()

	server := NewServer(port)

	// Start server in goroutine
	serverErr := make(chan error, 1)
	go func() {
		serverErr <- server.Start()
	}()

	// Wait for server to start
	time.Sleep(100 * time.Millisecond)

	// Test that server is responding
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d", port))
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Note: In a real scenario, you'd implement graceful shutdown
	// For testing, we'll just verify the server started successfully
}

func TestServer_Integration(t *testing.T) {
	// Find an available port
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()

	server := NewServer(port)

	// Create a context for server shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start server
	go func() {
		httpServer := &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: server.mux,
		}

		go func() {
			<-ctx.Done()
			httpServer.Shutdown(context.Background())
		}()

		httpServer.ListenAndServe()
	}()

	// Wait for server to start
	time.Sleep(200 * time.Millisecond)

	// Test different endpoints
	testCases := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{"Home", "GET", "/", http.StatusOK},
		{"StartTimer_NoBody", "POST", "/api/timer/start", http.StatusBadRequest},
		{"StopTimer", "POST", "/api/timer/stop", http.StatusOK},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("http://localhost:%d%s", port, tc.path)
			req, err := http.NewRequest(tc.method, url, nil)
			if err != nil {
				t.Fatal(err)
			}

			client := &http.Client{Timeout: 5 * time.Second}
			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Request failed: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, resp.StatusCode)
			}
		})
	}

	// Stop server
	cancel()
}

func BenchmarkServer_HandleRequests(b *testing.B) {
	server := NewServer(8080)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", "/", nil)
		rw := &testResponseWriter{}
		server.mux.ServeHTTP(rw, req)
	}
}

// testResponseWriter is a simple implementation for testing
type testResponseWriter struct {
	headers    http.Header
	body       []byte
	statusCode int
}

func (w *testResponseWriter) Header() http.Header {
	if w.headers == nil {
		w.headers = make(http.Header)
	}
	return w.headers
}

func (w *testResponseWriter) Write(data []byte) (int, error) {
	w.body = append(w.body, data...)
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}
	return len(data), nil
}

func (w *testResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}
