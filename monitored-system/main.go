package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

// MoveCommand represents a request to move sessions between nodes.
type MoveCommand struct {
	Count int    `json:"count"`
	From  string `json:"from"`
	To    string `json:"to"`
}

// Override String for MoveCommand struct
func (mc *MoveCommand) String() string {
	return "{From: " + mc.From + ", To: " + mc.To + ", Count: " + fmt.Sprintf("%d", mc.Count) + "}"
}

// NodeSessions holds the session counts for each node.
type NodeSessions struct {
	mu     sync.Mutex // Ensure thread-safe access
	Counts map[string]int
}

// Override String for NodeSessions struct
func (ns *NodeSessions) String() string {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	// Collect cities into a slice
	cities := make([]string, 0, len(ns.Counts))
	for city := range ns.Counts {
		cities = append(cities, city)
	}

	// Sort the cities alphabetically
	sort.Strings(cities)

	// Create parts of the final string
	parts := make([]string, 0, len(cities))
	for _, city := range cities {
		parts = append(parts, fmt.Sprintf("%s: %2d", city, ns.Counts[city]))
	}

	return "{" + strings.Join(parts, ", ") + "}"
}

// Applies random change on traffic distribution
// Adds a value betwee -2 and 2 to each node
// This func is called each round
func (ns *NodeSessions) ApplyRandomChange() {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	for city := range ns.Counts {
		change := rand.Intn(5) - 2 // Random number between -2 and 2
		newCount := ns.Counts[city] + change
		if newCount < 0 {
			newCount = 0 // Ensure the count never goes below zero
		}
		ns.Counts[city] = newCount
	}
}

// Handler of /api/move endpoint
// This endpoits has to receive MoveCommand json
func (ns *NodeSessions) MoveHandler(w http.ResponseWriter, r *http.Request) {
	var cmd MoveCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Apply the move command if valid
	if _, ok := ns.Counts[cmd.From]; ok && ns.Counts[cmd.From] >= cmd.Count {
		ns.mu.Lock()
		ns.Counts[cmd.From] -= cmd.Count
		ns.Counts[cmd.To] += cmd.Count
		ns.mu.Unlock()
		log.Printf("Got move command: %+v", cmd.String())
		log.Printf("%-10s %s", fmt.Sprintf("Round %d:", round), ns.String())
		// Return success response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Move command executed successfully")) // Optionally send a message
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Move command cannot be executed")) // Optionally send a message
	}
}

// ToJSON converts the Counts map to a JSON string
func (ns *NodeSessions) ToJSON() ([]byte, error) {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	// Marshal the Counts map to JSON
	jsonData, err := json.Marshal(ns.Counts)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

// SendJSON sends the Counts map as JSON to the provided URL
func (ns *NodeSessions) SendJSON(url string) error {
	// Convert the Counts map to JSON
	jsonData, err := ns.ToJSON()
	if err != nil {
		return err
	}

	// Send the JSON data using HTTP POST
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check if the server responded with a non-2xx status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send JSON: server responded with status %d", resp.StatusCode)
	}

	return nil
}

// Handler for /api/data endpoint
// This endpoint will return NodeSessions as JSON
func (ns *NodeSessions) DataHandler(w http.ResponseWriter, r *http.Request) {
	// Convert the Counts map to JSON
	jsonData, err := ns.ToJSON()
	if err != nil {
		http.Error(w, "Failed to marshal data", http.StatusInternalServerError)
		return
	}

	// Set response header to JSON content type
	w.Header().Set("Content-Type", "application/json")
	// Write the JSON response
	w.Write(jsonData)
}

var round = 0

func main() {
	// Command-line flag to get the interval value
	interval := flag.Int("interval", 5, "Interval in seconds for applying random changes")
	flag.Parse()
	log.Printf("Round interval time set as: %d seconds", *interval)

	time.Now().UnixNano()

	sessions := NodeSessions{
		Counts: map[string]int{"Gdansk": 10, "Poznan": 12, "Warsaw": 25, "Krakow": 4},
	}

	// HTTP Server for move commands and data retrieval
	http.HandleFunc("/api/move", sessions.MoveHandler)
	http.HandleFunc("/api/data", sessions.DataHandler)

	go func() {
		if err := http.ListenAndServe(":4040", nil); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Main loop
	ticker := time.NewTicker(1 * time.Second * time.Duration(*interval))
	for ; true; <-ticker.C {
		// apply random change on nodes
		sessions.ApplyRandomChange()
		// incerement the round number
		round++
		// // send the monitor data (node distribution)
		// err := sessions.SendJSON("http://localhost:4141/api/monitor")
		// if err != nil {
		// 	// log.Printf("Error sending JSON: %v", err)
		// }
		// log out the current node distribution
		log.Printf("%-10s %s", fmt.Sprintf("Round %d:", round), sessions.String())
	}
}
