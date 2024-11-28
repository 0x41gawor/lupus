package util

import (
	"encoding/json"
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
)

type Data struct {
	Body map[string]interface{}
}

func NewData(input runtime.RawExtension) (*Data, error) {
	dataMap, err := RawExtensionToMap(input)
	if err != nil {
		return nil, err
	}

	return &Data{Body: dataMap}, nil
}

func (d *Data) Get(keys []string) (*runtime.RawExtension, error) {
	if len(keys) == 0 {
		return nil, fmt.Errorf("no keys provided")
	}

	// Handle wildcard "*"
	if len(keys) == 1 && keys[0] == "*" {
		return InterfaceToRawExtension(d.Body)
	}

	// Handle a single key
	if len(keys) == 1 {
		parsedKeys := ParseKey(keys[0]) // Parse the dotted key
		value, exists := GetNestedValue(d.Body, parsedKeys)
		if !exists {
			return nil, fmt.Errorf("key %s not found", keys[0])
		}

		// Convert the value directly to a runtime.RawExtension
		return InterfaceToRawExtension(value)
	}

	// Handle multiple keys: create a subset object
	outputMap := make(map[string]interface{})
	for _, key := range keys {
		parsedKeys := ParseKey(key)
		value, exists := GetNestedValue(d.Body, parsedKeys)
		if exists {
			// Use only the last part of the key as the field name in the output map
			lastPart := parsedKeys[len(parsedKeys)-1]
			outputMap[lastPart] = value
		}
	}

	// Convert the resulting map to a runtime.RawExtension
	return InterfaceToRawExtension(outputMap)
}

func (d *Data) Set(key string, value interface{}) error {
	if key == "*" {
		// Validate that the value is convertible to map[string]interface{}
		newBody, err := InterfaceToMap2(value)
		if err != nil {
			return fmt.Errorf("failed to set body: value is not convertible to map[string]interface{}: %w", err)
		}

		// Replace the current body with the new map
		d.Body = newBody
		return nil
	}

	// Parse the dotted key into a slice of strings
	parsedKeys := ParseKey(key)

	// Use SetNestedValue to update the nested structure
	err := SetNestedValue(d.Body, parsedKeys, value)
	if err != nil {
		return fmt.Errorf("failed to set value for key %s: %w", key, err)
	}
	return nil
}

func (d *Data) String() string {
	str, _ := MapToString(d.Body)
	return str
}

// Nest creates a combined field with all inputKeys in it and assigns it to outputKey.
func (d *Data) Nest(inputKeys []string, outputKey string) error {
	if outputKey == "*" {
		return fmt.Errorf("`*` is not allowed as a key name for Nest")
	}

	// Initialize a map to hold the combined fields
	combinedMap := make(map[string]interface{})

	// Iterate over each input key
	for _, key := range inputKeys {
		if key == "*" {
			return fmt.Errorf("`*` is not allowed as a key name for Nest")
		}

		// Parse the key and get the value
		parsedKeys := ParseKey(key)
		value, exists := GetNestedValue(d.Body, parsedKeys)
		if !exists {
			return fmt.Errorf("field not found in body: %s", key)
		}

		// Use the last part of the key as the field name in the combined map
		lastKeyPart := parsedKeys[len(parsedKeys)-1]
		combinedMap[lastKeyPart] = value

		// Remove the key from the original map
		if err := DeleteNestedValue(d.Body, parsedKeys); err != nil {
			return fmt.Errorf("failed to remove key %s: %w", key, err)
		}
	}

	// Parse the output key
	parsedOutputKey := ParseKey(outputKey)

	// Set the combined map as the value of outputKey
	if err := SetNestedValue(d.Body, parsedOutputKey, combinedMap); err != nil {
		return fmt.Errorf("failed to set combined value for key %s: %w", outputKey, err)
	}

	return nil
}

// Remove removes fields with given keys in inputKeys
func (d *Data) Remove(inputKeys []string) error {
	for _, key := range inputKeys {
		if key == "*" {
			return fmt.Errorf("`*` is not allowed as a key name for Remove")
		}

		// Parse the key
		parsedKeys := ParseKey(key)

		// Attempt to delete the key; if it doesn't exist, skip it
		err := DeleteNestedValue(d.Body, parsedKeys)
		if err != nil {
			// Log or ignore the error if the key doesn't exist
			fmt.Printf("Warning: failed to remove key %s: %v\n", key, err)
		}
	}
	return nil
}

// Rename changes the last part of inputKey to outputKey, supporting both nested and non-nested inputKey.
func (d *Data) Rename(inputKey string, outputKey string) error {
	if inputKey == "*" || outputKey == "*" {
		return fmt.Errorf("'*' is not allowed as a key name for Rename")
	}

	// Parse the input key
	parsedInputKeys := ParseKey(inputKey)

	// Handle non-nested fields
	if len(parsedInputKeys) == 1 {
		value, exists := d.Body[inputKey]
		if !exists {
			return fmt.Errorf("key %s not found in body", inputKey)
		}
		// Rename by deleting the old key and adding the new key
		delete(d.Body, inputKey)
		d.Body[outputKey] = value
		return nil
	}

	// Handle nested fields
	parentKeys := parsedInputKeys[:len(parsedInputKeys)-1]
	lastKey := parsedInputKeys[len(parsedInputKeys)-1]

	// Navigate to the parent map
	parentMap, ok := GetNestedValue(d.Body, parentKeys)
	if !ok {
		return fmt.Errorf("parent path %v not found in body", parentKeys)
	}

	// Ensure the parent is a map[string]interface{}
	parentMapTyped, ok := parentMap.(map[string]interface{})
	if !ok {
		return fmt.Errorf("parent path %v is not a map[string]interface{}", parentKeys)
	}

	// Check if the last key exists in the parent map
	value, exists := parentMapTyped[lastKey]
	if !exists {
		return fmt.Errorf("key %s not found in body", inputKey)
	}

	// Rename the key by deleting the old key and adding a new one
	delete(parentMapTyped, lastKey)
	parentMapTyped[outputKey] = value

	return nil
}

// Duplicate copies the value associated with inputKey to outputKey, supporting nested keys.
func (d *Data) Duplicate(inputKey string, outputKey string) error {
	if inputKey == "*" || outputKey == "*" {
		return fmt.Errorf("'*' is not allowed as a key name for Duplicate")
	}

	// Parse the input key
	parsedInputKeys := ParseKey(inputKey)

	// Get the value for the inputKey
	value, exists := GetNestedValue(d.Body, parsedInputKeys)
	if !exists {
		return fmt.Errorf("key %s not found in body", inputKey)
	}

	// Parse the output key
	parsedOutputKeys := ParseKey(outputKey)

	// Set the value to the outputKey
	if err := SetNestedValue(d.Body, parsedOutputKeys, value); err != nil {
		return fmt.Errorf("failed to set value for key %s: %w", outputKey, err)
	}

	return nil
}

func (d *Data) Print(keys []string) error {
	// Prepare placeholder for output
	outputMap := make(map[string]interface{})

	// Handle wildcard "*" to print the entire body
	if len(keys) == 1 && keys[0] == "*" {
		outputMap = d.Body
	} else {
		for _, key := range keys {
			// Parse the key to support nested paths
			parsedKeys := ParseKey(key)

			// Retrieve the value for the key
			value, exists := GetNestedValue(d.Body, parsedKeys)
			if exists {
				// Add the value to the output map, using the original key string
				outputMap[key] = value
			}
		}
	}

	// Convert the output map to a JSON string
	str, err := MapToString(outputMap)
	if err != nil {
		return fmt.Errorf("failed to convert map to string: %w", err)
	}

	// Print the result
	fmt.Println(str)
	return nil
}

// Helper function to parse a dotted key into a slice of strings
func ParseKey(key string) []string {
	return strings.Split(key, ".")
}

// Helper function to navigate and get a value from nested maps
func GetNestedValue(data map[string]interface{}, keys []string) (interface{}, bool) {
	current := data
	for i, key := range keys {
		value, exists := current[key]
		if !exists {
			return nil, false
		}
		// If this is the last key, return the value
		if i == len(keys)-1 {
			return value, true
		}
		// Otherwise, continue navigating
		nextMap, ok := value.(map[string]interface{})
		if !ok {
			return nil, false
		}
		current = nextMap
	}
	return nil, false
}

// Helper function to set a value in nested maps
func SetNestedValue(data map[string]interface{}, keys []string, value interface{}) error {
	current := data
	for i, key := range keys {
		// If this is the last key, set the value
		if i == len(keys)-1 {
			current[key] = value
			return nil
		}

		// Navigate or create nested maps
		nextMap, ok := current[key].(map[string]interface{})
		if !ok {
			// If the next level doesn't exist, create it
			nextMap = make(map[string]interface{})
			current[key] = nextMap
		}
		current = nextMap
	}
	return nil
}

// Helper function to delete a nested key
func DeleteNestedValue(data map[string]interface{}, keys []string) error {
	current := data
	for i, key := range keys {
		// If this is the last key, delete it
		if i == len(keys)-1 {
			delete(current, key)
			return nil
		}
		// Otherwise, navigate
		nextMap, ok := current[key].(map[string]interface{})
		if !ok {
			return fmt.Errorf("key path %v does not exist", keys)
		}
		current = nextMap
	}
	return nil
}

// Helper function to navigate to the parent map of the target key
func navigateToParent(data map[string]interface{}, keys []string) (map[string]interface{}, error) {
	current := data
	for i := 0; i < len(keys)-1; i++ {
		key := keys[i]
		nextMap, ok := current[key].(map[string]interface{})
		if !ok {
			// If the next level doesn't exist, create it
			nextMap = make(map[string]interface{})
			current[key] = nextMap
		}
		current = nextMap
	}
	return current, nil
}

func (d *Data) Insert(key string, raw runtime.RawExtension) error {
	// Convert RawExtension to a map[string]interface{}
	var newMap map[string]interface{}
	if err := json.Unmarshal(raw.Raw, &newMap); err != nil {
		return fmt.Errorf("failed to unmarshal RawExtension: %w", err)
	}

	if key == "*" {
		// Add newMap at the root of the JSON structure
		for rootKey, newValue := range newMap {
			if existingValue, exists := d.Body[rootKey]; exists {
				// Merge or replace
				mergedValue, err := mergeOrReplace(existingValue, newValue)
				if err != nil {
					return fmt.Errorf("failed to merge or replace values for root key %s: %w", rootKey, err)
				}
				d.Body[rootKey] = mergedValue
			} else {
				d.Body[rootKey] = newValue
			}
		}
		return nil
	}

	// Parse the dotted key into a slice of strings
	parsedKeys := ParseKey(key)

	// Navigate to the target location
	targetParent, err := navigateToParent(d.Body, parsedKeys)
	if err != nil {
		return fmt.Errorf("failed to navigate to parent key for %s: %w", key, err)
	}

	// Get the last key where the newMap will be inserted
	targetKey := parsedKeys[len(parsedKeys)-1]

	if existingValue, exists := targetParent[targetKey]; exists {
		// Merge or replace
		mergedValue, err := mergeOrReplace(existingValue, newMap)
		if err != nil {
			return fmt.Errorf("failed to merge or replace values for key %s: %w", key, err)
		}
		targetParent[targetKey] = mergedValue
	} else {
		// Insert new field
		targetParent[targetKey] = newMap
	}

	return nil
}

// Helper function to merge or replace values
func mergeOrReplace(existing, incoming interface{}) (interface{}, error) {
	existingMap, err := InterfaceToMap(existing)
	if err != nil {
		// If the existing value is not a map, replace it
		return incoming, nil
	}

	incomingMap, err := InterfaceToMap(incoming)
	if err != nil {
		return nil, fmt.Errorf("failed to convert incoming value to map: %w", err)
	}

	// Perform recursive merging
	for key, incomingValue := range incomingMap {
		if existingValue, exists := existingMap[key]; exists {
			// Merge recursively or replace
			mergedValue, err := mergeOrReplace(existingValue, incomingValue)
			if err != nil {
				return nil, err
			}
			existingMap[key] = mergedValue
		} else {
			// Insert new field
			existingMap[key] = incomingValue
		}
	}

	return existingMap, nil
}
