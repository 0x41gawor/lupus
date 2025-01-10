package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func bounceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var jsonData map[string]interface{}
	if err := json.Unmarshal(body, &jsonData); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body) // Echoes the received JSON back as the response
}

func main() {
	http.HandleFunc("/api/bounce", bounceHandler)
	log.Println("Starting server on :7000...")
	log.Fatal(http.ListenAndServe(":7000", nil))
}
