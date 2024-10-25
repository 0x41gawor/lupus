package controller

import (
	"encoding/json"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
)

// Function to extract and unmarshal the runtime.RawExtension into a map
func rawExtensionToMap(rawExt runtime.RawExtension) (map[string]interface{}, error) {
	// Initialize an empty map to hold the unmarshaled data
	var result map[string]interface{}

	// Unmarshal the RawExtension's Raw field into the map
	err := json.Unmarshal(rawExt.Raw, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal raw extension: %v", err)
	}

	// Return the result map
	return result, nil
}

// Function to marshal a map into runtime.RawExtension
func mapToRawExtension(data map[string]interface{}) (runtime.RawExtension, error) {
	// Marshal the map into a JSON byte slice
	rawBytes, err := json.Marshal(data)
	if err != nil {
		return runtime.RawExtension{}, fmt.Errorf("failed to marshal map to JSON: %v", err)
	}

	// Create a new runtime.RawExtension and assign the JSON byte slice
	return runtime.RawExtension{
		Raw: rawBytes,
	}, nil
}

// Function to convert an interface{} (extracted from map) to RawExtension
func interfaceToRawExtension(value interface{}) (runtime.RawExtension, error) {
	// Marshal the interface{} into a JSON byte slice
	rawBytes, err := json.Marshal(value)
	if err != nil {
		return runtime.RawExtension{}, fmt.Errorf("failed to marshal interface{} to JSON: %v", err)
	}

	// Create a new runtime.RawExtension and assign the marshaled JSON byte slice
	return runtime.RawExtension{
		Raw: rawBytes,
	}, nil
}

// Helper function to convert runtime.RawExtension to a JSON string
func rawExtensionToString(rawExt runtime.RawExtension) (string, error) {
	var jsonData map[string]interface{}
	if err := json.Unmarshal(rawExt.Raw, &jsonData); err != nil {
		return "", fmt.Errorf("failed to unmarshal RawExtension: %v", err)
	}

	// Marshal it back into a JSON string for writing to file
	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON data: %v", err)
	}

	return string(jsonBytes), nil
}
