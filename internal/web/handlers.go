package web

import (
	_ "embed"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/aureliomalheiros/aragomodoro/internal/pomodoro"
	"github.com/aureliomalheiros/aragomodoro/internal/sound"
	"github.com/gorilla/websocket"
)

//go:embed templates/index.html
var indexHTML string

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type TimerSession struct {
	Active       bool   `json:"active"`
	Type         string `json:"type"`
	Duration     int    `json:"duration"`
	Remaining    int    `json:"remaining"`
	RepeatCount  int    `json:"repeatCount"`
	CurrentCycle int    `json:"currentCycle"`
}

type WebTimerManager struct {
	mu        sync.RWMutex
	session   *TimerSession
	clients   map[*websocket.Conn]bool
	clientsMu sync.RWMutex
	stopChan  chan bool
}

var timerManager = &WebTimerManager{
	clients:  make(map[*websocket.Conn]bool),
	stopChan: make(chan bool),
}

type TimerRequest struct {
	FocusDuration   int  `json:"focusDuration"`
	BreakDuration   int  `json:"breakDuration"`
	RepeatCount     int  `json:"repeatCount"`
	ContinueOnBreak bool `json:"continueOnBreak"`
}

func HandleHome(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("index").Parse(indexHTML)
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	timerManager.mu.RLock()
	session := timerManager.session
	timerManager.mu.RUnlock()

	data := struct {
		Session *TimerSession
	}{
		Session: session,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Template execution error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func HandleStartTimer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req TimerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := pomodoro.ValidateDurations(req.FocusDuration, req.BreakDuration, req.RepeatCount); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Stop any existing timer (without writing response)
	timerManager.mu.Lock()
	if timerManager.session != nil && timerManager.session.Active {
		timerManager.session.Active = false
		select {
		case timerManager.stopChan <- true:
		default:
		}
	}
	timerManager.mu.Unlock()

	go timerManager.startTimerSession(req)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "started"})
}

func HandleStopTimer(w http.ResponseWriter, r *http.Request) {
	timerManager.mu.Lock()
	if timerManager.session != nil && timerManager.session.Active {
		timerManager.session.Active = false
		select {
		case timerManager.stopChan <- true:
		default:
		}
	}
	timerManager.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "stopped"})
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	timerManager.clientsMu.Lock()
	timerManager.clients[conn] = true
	timerManager.clientsMu.Unlock()

	timerManager.mu.RLock()
	if timerManager.session != nil {
		conn.WriteJSON(timerManager.session)
	}
	timerManager.mu.RUnlock()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			timerManager.clientsMu.Lock()
			delete(timerManager.clients, conn)
			timerManager.clientsMu.Unlock()
			break
		}
	}
}

func (tm *WebTimerManager) startTimerSession(req TimerRequest) {
	tm.stopChan = make(chan bool, 1)

	tm.mu.Lock()
	tm.session = &TimerSession{
		Active:       true,
		Type:         "focus",
		Duration:     req.FocusDuration,
		Remaining:    req.FocusDuration * 60,
		RepeatCount:  req.RepeatCount,
		CurrentCycle: 1,
	}
	tm.mu.Unlock()

	for cycle := 1; cycle <= req.RepeatCount; cycle++ {
		tm.mu.Lock()
		tm.session.Type = "focus"
		tm.session.Duration = req.FocusDuration
		tm.session.Remaining = req.FocusDuration * 60
		tm.session.CurrentCycle = cycle
		tm.mu.Unlock()

		if !tm.runTimer(req.FocusDuration * 60) {
			return
		}

		// Play soft sound when focus period completes
		go sound.SoftFocusComplete()

		if cycle < req.RepeatCount || req.ContinueOnBreak {
			tm.mu.Lock()
			tm.session.Type = "break"
			tm.session.Duration = req.BreakDuration
			tm.session.Remaining = req.BreakDuration * 60
			tm.mu.Unlock()

			if !tm.runTimer(req.BreakDuration * 60) {
				return
			}

			// Play soft sound when break period completes
			go sound.SoftBreakComplete()
		}
	}

	tm.mu.Lock()
	tm.session.Active = false
	tm.session.Type = "completed"
	tm.session.Remaining = 0
	tm.mu.Unlock()
	tm.broadcastUpdate()
}

func (tm *WebTimerManager) runTimer(durationSeconds int) bool {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for i := durationSeconds; i > 0; i-- {
		tm.mu.Lock()
		if tm.session != nil {
			tm.session.Remaining = i
		}
		tm.mu.Unlock()
		tm.broadcastUpdate()

		select {
		case <-tm.stopChan:
			return false
		case <-ticker.C:
		}
	}

	return true
}

func (tm *WebTimerManager) broadcastUpdate() {
	tm.mu.RLock()
	session := tm.session
	tm.mu.RUnlock()

	if session == nil {
		return
	}

	tm.clientsMu.Lock()
	defer tm.clientsMu.Unlock()

	for client := range tm.clients {
		err := client.WriteJSON(session)
		if err != nil {
			client.Close()
			delete(tm.clients, client)
		}
	}
}
