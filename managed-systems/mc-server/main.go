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

// Command represents a single command to update either CPU or RAM license
type Command struct {
	CPU *int `json:"cpu,omitempty"`
	RAM *int `json:"ram,omitempty"`
}

// ServerData holds the CPU and RAM usage information
type ServerData struct {
	CPU struct {
		InUse   int `json:"in_use"`
		License int `json:"license"`
	} `json:"cpu"`
	RAM struct {
		InUse   int `json:"in_use"`
		License int `json:"license"`
	} `json:"ram"`
}

var (
	data      ServerData
	dataMutex sync.Mutex
	interval  time.Duration
	round     int
)

// Update server usage, with random changes for in_use values
func updateUsage() {
	dataMutex.Lock()
	defer dataMutex.Unlock()

	// Randomly update in_use values for CPU and RAM
	data.CPU.InUse += rand.Intn(4) - 1 // -1, 0, 1 or 2
	data.RAM.InUse += rand.Intn(4) - 1 // -1, 0, 1 or 2

	// Keep values within a logical range (e.g., 0 to license limit)
	if data.CPU.InUse < 0 {
		data.CPU.InUse = 0
	}
	if data.CPU.InUse > data.CPU.License {
		data.CPU.InUse = data.CPU.License
	}
	if data.RAM.InUse < 0 {
		data.RAM.InUse = 0
	}
	if data.RAM.InUse > data.RAM.License {
		data.RAM.InUse = data.RAM.License
	}
}

// Log current data to the screen in the specified format
func logUsage() {
	dataJSON, err := json.Marshal(data)
	if err != nil {
		log.Println("Error marshaling JSON:", err)
		return
	}

	// Print the JSON data to the console in a formatted manner
	log.Printf("%-10s %s", fmt.Sprintf("Round %d:", round), string(dataJSON))
}

// Handler for the `/api/data` endpoint
func dataHandler(w http.ResponseWriter, r *http.Request) {
	dataMutex.Lock()
	defer dataMutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// Handler for the `/api/set-license` endpoint
func setLicenseHandler(w http.ResponseWriter, r *http.Request) {
	var command Command

	err := json.NewDecoder(r.Body).Decode(&command)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	dataMutex.Lock()
	defer dataMutex.Unlock()

	// Apply the commands and log the old and new values
	if command.CPU != nil {
		oldCPU := data.CPU.License
		data.CPU.License = *command.CPU
		log.Printf("CPU License changed from: %d to: %d", oldCPU, data.CPU.License)
	}
	if command.RAM != nil {
		oldRAM := data.RAM.License
		data.RAM.License = *command.RAM
		log.Printf("RAM License changed from: %d to:  %d", oldRAM, data.RAM.License)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "License values updated")
}

func main() {
	// Parse the interval from command line arguments
	interval := flag.Int("interval", 5, "Interval in seconds for applying random changes")
	flag.Parse()
	log.Printf("Round interval time set as: %d seconds", *interval)

	// Initialize the server data
	data.CPU.License = 20
	data.RAM.License = 8
	data.CPU.InUse = 10 // Starting with 10% CPU usage
	data.RAM.InUse = 5  // Starting with 5% RAM usage

	// Create a ticker to update server usage periodically
	http.HandleFunc("/api/data", dataHandler)
	http.HandleFunc("/api/set-license", setLicenseHandler)
	go func() {
		if err := http.ListenAndServe(":6000", nil); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Main loop
	ticker := time.NewTicker(1 * time.Second * time.Duration(*interval))
	for ; true; <-ticker.C {
		// apply random change on nodes
		updateUsage()
		// increment the round number
		round++
		// log out the current node distribution
		logUsage()
	}
}
