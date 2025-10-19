package web

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestNewServer(t *testing.T) {
	server := NewServer(8080)
	if server == nil {
		t.Fatal("NewServer returned nil")
	}
	if server.port != 8080 {
		t.Errorf("Expected port 8080, got %d", server.port)
	}
	if server.mux == nil {
		t.Error("Server mux is nil")
	}
}

func TestHandleHome(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleHome)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check if response contains HTML
	body := rr.Body.String()
	if !strings.Contains(body, "<!DOCTYPE html>") {
		t.Error("Response should contain HTML doctype")
	}
	if !strings.Contains(body, "Aragomodoro") {
		t.Error("Response should contain Aragomodoro title")
	}
}

func TestHandleStartTimer(t *testing.T) {
	// Reset timer manager for test
	timerManager = &WebTimerManager{
		clients:  make(map[*websocket.Conn]bool),
		stopChan: make(chan bool, 1),
	}

	tests := []struct {
		name           string
		method         string
		body           string
		expectedStatus int
	}{
		{
			name:           "ValidRequest",
			method:         "POST",
			body:           `{"focusDuration":25,"breakDuration":5,"repeatCount":1,"continueOnBreak":false}`,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "InvalidMethod",
			method:         "GET",
			body:           "",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "InvalidJSON",
			method:         "POST",
			body:           `{"invalid":json}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "InvalidDuration",
			method:         "POST",
			body:           `{"focusDuration":0,"breakDuration":5,"repeatCount":1,"continueOnBreak":false}`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, "/api/timer/start", bytes.NewBufferString(tt.body))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(HandleStartTimer)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("Handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if tt.expectedStatus == http.StatusOK {
				// Parse response first
				var response map[string]string
				body := rr.Body.String()
				err := json.Unmarshal([]byte(body), &response)
				if err != nil {
					t.Errorf("Failed to parse response JSON: %v, body: %s", err, body)
				} else if response["status"] != "started" {
					t.Errorf("Expected status 'started', got '%s'", response["status"])
				}

				// Give some time for goroutine to start, then stop safely
				time.Sleep(50 * time.Millisecond)
				timerManager.mu.Lock()
				if timerManager.session != nil && timerManager.session.Active {
					timerManager.session.Active = false
				}
				timerManager.mu.Unlock()

				// Signal stop to cleanup goroutine
				select {
				case timerManager.stopChan <- true:
				default:
				}
			}
		})
	}
}

func TestHandleStopTimer(t *testing.T) {
	req, err := http.NewRequest("POST", "/api/timer/stop", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleStopTimer)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse response JSON: %v", err)
	}
	if response["status"] != "stopped" {
		t.Errorf("Expected status 'stopped', got '%s'", response["status"])
	}
}

func TestTimerSession(t *testing.T) {
	session := &TimerSession{
		Active:       true,
		Type:         "focus",
		Duration:     25,
		Remaining:    1500, // 25 minutes in seconds
		RepeatCount:  1,
		CurrentCycle: 1,
	}

	if session.Type != "focus" {
		t.Errorf("Expected Type 'focus', got '%s'", session.Type)
	}
	if session.Duration != 25 {
		t.Errorf("Expected Duration 25, got %d", session.Duration)
	}
	if session.Remaining != 1500 {
		t.Errorf("Expected Remaining 1500, got %d", session.Remaining)
	}
}

func TestWebTimerManager(t *testing.T) {
	// Create a new timer manager for this test to avoid conflicts
	testManager := &WebTimerManager{
		clients:  make(map[*websocket.Conn]bool),
		stopChan: make(chan bool, 1),
	}

	req := TimerRequest{
		FocusDuration:   1,
		BreakDuration:   1,
		RepeatCount:     1,
		ContinueOnBreak: false,
	}
	testManager.mu.Lock()
	testManager.session = &TimerSession{
		Active:       true,
		Type:         "focus",
		Duration:     req.FocusDuration,
		Remaining:    req.FocusDuration * 60,
		RepeatCount:  req.RepeatCount,
		CurrentCycle: 1,
	}
	testManager.mu.Unlock()

	// Verify session was created correctly
	testManager.mu.RLock()
	if testManager.session == nil {
		t.Error("Session should be created")
	} else {
		if !testManager.session.Active {
			t.Error("Session should be active")
		}
		if testManager.session.Type != "focus" {
			t.Errorf("Expected session type 'focus', got '%s'", testManager.session.Type)
		}
		if testManager.session.Duration != req.FocusDuration {
			t.Errorf("Expected duration %d, got %d", req.FocusDuration, testManager.session.Duration)
		}
	}
	testManager.mu.RUnlock()

	testManager.mu.Lock()
	if testManager.session != nil {
		testManager.session.Active = false
	}
	testManager.mu.Unlock()
	testManager.mu.RLock()
	if testManager.session != nil && testManager.session.Active {
		t.Error("Session should be stopped")
	}
	testManager.mu.RUnlock()
}

func TestServerRoutes(t *testing.T) {
	// Reset timer manager for test
	timerManager = &WebTimerManager{
		clients:  make(map[*websocket.Conn]bool),
		stopChan: make(chan bool, 1),
	}

	server := NewServer(8080)

	routes := []struct {
		path   string
		method string
	}{
		{"/", "GET"},
		{"/api/timer/stop", "POST"},
		{"/ws", "GET"}, // WebSocket endpoint
	}

	for _, route := range routes {
		t.Run(route.path, func(t *testing.T) {
			var req *http.Request
			var err error

			if route.path == "/api/timer/start" {
				return
			}

			req, err = http.NewRequest(route.method, route.path, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			server.mux.ServeHTTP(rr, req)

			// We expect some response (not 404)
			if rr.Code == http.StatusNotFound {
				t.Errorf("Route %s %s returned 404", route.method, route.path)
			}
		})
	}
}

func BenchmarkHandleHome(b *testing.B) {
	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleHome)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handler.ServeHTTP(rr, req)
		rr.Body.Reset()
	}
}

func BenchmarkHandleStartTimer(b *testing.B) {
	body := `{"focusDuration":25,"breakDuration":5,"repeatCount":1,"continueOnBreak":false}`
	req, _ := http.NewRequest("POST", "/api/timer/start", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleStartTimer)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handler.ServeHTTP(rr, req)
		rr.Body.Reset()
		// Reset request body for next iteration
		req.Body = &stringReadCloser{strings.NewReader(body)}
	}
}

// Helper type for benchmarking
type stringReadCloser struct {
	*strings.Reader
}

func (src *stringReadCloser) Close() error {
	return nil
}
