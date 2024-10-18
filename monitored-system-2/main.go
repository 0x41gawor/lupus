package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

// Command represents a single command to update either CPU or RAM license
type Command struct {
	CPU *int `json:"cpu,omitempty"`
	RAM *int `json:"ram,omitempty"`
}

// CommandsRequest represents the JSON structure with a set of commands
type CommandsRequest struct {
	Commands []Command `json:"commands"`
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

	// Log current usage data to the screen
	logUsage()
}

// Log current data to the screen in the specified format
func logUsage() {
	round++
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
	var request CommandsRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	dataMutex.Lock()
	defer dataMutex.Unlock()

	// Apply the commands and log the old and new values
	for _, command := range request.Commands {
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
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "License values updated")
}

func main() {
	// Parse the interval from command line arguments
	intervalArg := "5" // default interval of 5 seconds
	if len(os.Args) > 1 {
		intervalArg = os.Args[1]
	}

	parsedInterval, err := strconv.Atoi(intervalArg)
	if err != nil {
		log.Fatal("Invalid interval argument:", err)
	}
	interval = time.Duration(parsedInterval) * time.Second

	// Initialize the server data
	data.CPU.License = 20
	data.RAM.License = 8
	data.CPU.InUse = 10 // Starting with 10% CPU usage
	data.RAM.InUse = 5  // Starting with 5% RAM usage

	// Create a ticker to update server usage periodically
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Handle graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			select {
			case <-ticker.C:
				updateUsage()
			case <-stop:
				log.Println("Shutting down server...")
				os.Exit(0)
			}
		}
	}()

	// Set up the HTTP server and routes
	http.HandleFunc("/api/data", dataHandler)
	http.HandleFunc("/api/set-license", setLicenseHandler)

	log.Println("Server is running...")
	log.Fatal(http.ListenAndServe(":4242", nil))
}
