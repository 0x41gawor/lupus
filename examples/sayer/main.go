package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

// Data represents the JSON structure used in Sayer.
type Data struct {
	Case string `json:"case"`
	Ram  int    `json:"ram"`
}

// Sayer manages periodic data generation and logging.
type Sayer struct {
	mu       sync.Mutex
	data     *Data
	interval time.Duration
	round    int
}

// NewSayer initializes a new Sayer instance.
func NewSayer(interval time.Duration) *Sayer {
	return &Sayer{
		interval: interval,
	}
}

// GenerateRandomData generates a new random Data object.
func (s *Sayer) GenerateRandomData() {
	cases := []string{"A", "B", "C"}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = &Data{
		Case: cases[rand.Intn(len(cases))],
		Ram:  rand.Intn(51), // Random number between 0 and 50
	}
	s.round++
	log.Printf("Round %d: I have data: %s\n", s.round, toJSON(s.data))
}

// DataHandler handles the `/api/data` endpoint.
func (s *Sayer) DataHandler(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	data := s.data
	s.mu.Unlock()
	response := toJSON(data)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(response))
}

// CommandHandler handles the `/api/command` endpoint, logging received data.
func (s *Sayer) CommandHandler(w http.ResponseWriter, r *http.Request) {
	var received interface{}
	if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	receivedStr, err := InterfaceToString(received)
	if err != nil {
		fmt.Print("ERROR")
	}
	log.Printf("Round %d: I've received: %v\n", s.round, receivedStr)
	// Prepare a JSON response
	response := map[string]string{
		"res": "Command received and logged",
	}
	// Encode and send the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, `{"error": "Failed to send response"}`, http.StatusInternalServerError)
	}
}

// Start rounds for periodic data generation.
func (s *Sayer) Start() {
	s.GenerateRandomData()
	ticker := time.NewTicker(s.interval)
	for range ticker.C {
		s.GenerateRandomData()
	}
}

// Convert Data struct to JSON string.
func toJSON(data *Data) string {
	bytes, _ := json.Marshal(data)
	return string(bytes)
}

// InterfaceToString converts any interface{} to its string representation
func InterfaceToString(data interface{}) (string, error) {
	// Use JSON marshaling to handle structured data
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to convert to string: %w", err)
	}

	return string(jsonBytes), nil
}

func main() {
	// Parse the interval flag
	intervalFlag := flag.Int("interval", 30, "Interval between rounds in seconds")
	flag.Parse()
	log.Printf("Round interval time set as: %d seconds", *intervalFlag)

	// Initialize Sayer
	sayer := NewSayer(time.Duration(*intervalFlag) * time.Second)

	// HTTP server setup
	http.HandleFunc("/api/data", sayer.DataHandler)
	http.HandleFunc("/api/command", sayer.CommandHandler)

	// Start the server
	go func() {
		if err := http.ListenAndServe(":7000", nil); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Start periodic rounds
	sayer.Start()
}
