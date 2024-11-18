package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
)

type ElementInstanceState struct {
	IsAfterDryRun bool
	LastUpdated   time.Time
}

// Function to extract and unmarshal the runtime.RawExtension into a map
func RawExtensionToMap(rawExt runtime.RawExtension) (map[string]interface{}, error) {
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
func MapToRawExtension(data map[string]interface{}) (runtime.RawExtension, error) {
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

// Helper function to convert runtime.RawExtension to a JSON string
func RawExtensionToString(rawExt runtime.RawExtension) (string, error) {
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
func InterfaceToMap(data interface{}) (map[string]interface{}, error) {
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

// Function to convert a map[string]interface{} to a JSON string
func MapToString(data map[string]interface{}) (string, error) {
	// Marshal the map into a JSON byte slice
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal map to JSON string: %v", err)
	}

	// Convert JSON bytes to a string and return
	return string(jsonBytes), nil
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

// InterfaceToRawExtension converts any interface{} to runtime.RawExtension
func InterfaceToRawExtension(data interface{}) (*runtime.RawExtension, error) {
	// Marshal the interface into a JSON byte slice
	rawBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal interface to JSON: %v", err)
	}

	// Create and return a runtime.RawExtension
	return &runtime.RawExtension{
		Raw: rawBytes,
	}, nil
}
