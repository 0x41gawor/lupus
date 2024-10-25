package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

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

// interfaceToMap tries to convert any interface{} to map[string]interface{}
func interfaceToMap(data interface{}) (map[string]interface{}, error) {
	// Check if data is already a map[string]interface{}
	if result, ok := data.(map[string]interface{}); ok {
		return result, nil
	}

	// Use reflection to inspect and handle other possible structures
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Map {
		converted := make(map[string]interface{})
		for _, key := range val.MapKeys() {
			// Ensure the key is a string
			if key.Kind() != reflect.String {
				return nil, errors.New("non-string key encountered in map")
			}
			// Set the key and value in the new map
			converted[key.String()] = val.MapIndex(key).Interface()
		}
		return converted, nil
	}

	return nil, errors.New("data is not convertible to map[string]interface{}")
}
