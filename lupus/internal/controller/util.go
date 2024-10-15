package controller

import (
	"encoding/json"
	"fmt"
	"math"

	"k8s.io/apimachinery/pkg/runtime"
)

// Input defines the 4 integer fields for different locations
type Input struct {
	Gdansk int `json:"gdansk"`
	Krakow int `json:"krakow"`
	Poznan int `json:"poznan"`
	Warsaw int `json:"warsaw"`
}

// MoveCommand represents a move command with source, destination, and count
type MoveCommand struct {
	From  string `json:"from"`  // Source location
	To    string `json:"to"`    // Destination location
	Count int    `json:"count"` // Number of items to move
}

// distributeLoad takes the Input struct and generates a list of MoveCommand
func distributeLoad(input Input) []MoveCommand {
	// List of MoveCommand that will store the result
	var moves []MoveCommand

	// Map to easily reference the cities
	cities := map[string]*int{
		"Gdansk": &input.Gdansk,
		"Krakow": &input.Krakow,
		"Poznan": &input.Poznan,
		"Warsaw": &input.Warsaw,
	}

	// Calculate the total load and target load per city
	totalLoad := input.Gdansk + input.Krakow + input.Poznan + input.Warsaw
	targetLoad := totalLoad / len(cities)

	// Collect cities with excess and deficit load
	excessCities := make(map[string]int)
	deficitCities := make(map[string]int)

	for city, load := range cities {
		if *load > targetLoad {
			excessCities[city] = *load - targetLoad
		} else if *load < targetLoad {
			deficitCities[city] = targetLoad - *load
		}
	}

	// Process moves from excess cities to deficit cities
	for fromCity, excess := range excessCities {
		for toCity, deficit := range deficitCities {
			if excess == 0 {
				break
			}
			// Calculate the amount to move
			moveCount := int(math.Min(float64(excess), float64(deficit)))
			if moveCount > 0 {
				// Create move command
				moves = append(moves, MoveCommand{
					From:  fromCity,
					To:    toCity,
					Count: moveCount,
				})

				// Update the loads
				excessCities[fromCity] -= moveCount
				deficitCities[toCity] -= moveCount

				// Update cities' actual load
				*cities[fromCity] -= moveCount
				*cities[toCity] += moveCount
			}
		}
	}

	return moves
}

// parseRawExtensionToInput parses the RawExtension into the Input struct
func parseRawExtensionToInput(input runtime.RawExtension) (Input, error) {
	// Define a variable of type Input to hold the parsed data
	var parsedInput Input

	// Unmarshal the raw JSON into the Input struct
	if err := json.Unmarshal(input.Raw, &parsedInput); err != nil {
		return Input{}, fmt.Errorf("failed to unmarshal RawExtension into Input struct: %w", err)
	}

	// Return the successfully parsed Input struct
	return parsedInput, nil
}

// convertMoveCommandSliceToRawExtension converts a slice of MoveCommand into a runtime.RawExtension
func convertMoveCommandSliceToRawExtension(commands []MoveCommand) (runtime.RawExtension, error) {
	// Marshal the slice of MoveCommand into JSON
	rawData, err := json.Marshal(commands)
	if err != nil {
		return runtime.RawExtension{}, fmt.Errorf("failed to marshal []MoveCommand to JSON: %w", err)
	}

	// Create a RawExtension and assign the raw JSON data
	return runtime.RawExtension{
		Raw: rawData, // Assign the marshaled JSON bytes to Raw
	}, nil
}

func opaSimulation(input runtime.RawExtension) (runtime.RawExtension, error) {
	data, err := parseRawExtensionToInput(input)
	if err != nil {
		return runtime.RawExtension{}, err
	}
	moveCommands := distributeLoad(data)
	return convertMoveCommandSliceToRawExtension(moveCommands)
}
