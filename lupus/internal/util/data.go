package util

import (
	"fmt"

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

func (d *Data) Get(key string) (interface{}, error) {
	if key == "*" {
		value := d.Body
		return value, nil
	} else {
		value := d.Body[key]
		return value, nil
	}
}

func (d *Data) GetKeys(keys []string) (*runtime.RawExtension, error) {
	// prepare placeholder for output
	outputMap := make(map[string]interface{})
	// cut data parts based on next.Keys
	if len(keys) == 1 && keys[0] == "*" {
		outputMap = d.Body
	} else {
		for _, tag := range keys {
			if value, exists := d.Body[tag]; exists {
				outputMap[tag] = value
			}
		}
	}
	outputRaw, err := MapToRawExtension(outputMap)
	if err != nil {
		return nil, err
	}
	return &outputRaw, nil

}

func (d *Data) Set(key string, value interface{}) error {
	if key == "*" {
		newBody, err := InterfaceToMap(value)
		if err != nil {
			return err
		}
		d.Body = newBody
		return nil
	} else {
		d.Body[key] = value
		return nil
	}
}

func (d *Data) String() string {
	str, _ := MapToString(d.Body)
	return str
}

// Nest creates a combined field with all inputKeys in it and names it outputKey
func (d *Data) Nest(inputKeys []string, outputKey string) error {
	if outputKey == "*" {
		return fmt.Errorf("`*` is not allowed as a key name for Concat")
	}
	// Initialize a map to hold the combined fields
	combinedMap := make(map[string]interface{})
	// Iterate over each field in the inputFields slice
	for _, key := range inputKeys {
		if key == "*" {
			return fmt.Errorf("`*` is not allowed as a key name for Remove")
		}
		// Retrieve the value for the field
		value, exists := d.Body[key]
		if !exists {
			return fmt.Errorf("field not found in body: %s", key)
		}
		// Add the field to the combined map
		combinedMap[key] = value
		// Remove the field from the original map to "move" it
		delete(d.Body, key)
	}
	// Set the combined map as the value of outputField
	d.Body[outputKey] = combinedMap
	return nil
}

// Remove removes fields with given keys in inputKeys
func (d *Data) Remove(inputKeys []string) error {
	for _, key := range inputKeys {
		if key == "*" {
			return fmt.Errorf("`*` is not allowed as a key name for Remove")
		}
		delete(d.Body, key)
	}
	return nil
}

// Rename changes name from inputKey to outputKey
func (d *Data) Rename(inputKey string, outputKey string) error {
	if inputKey == "*" || outputKey == "*" {
		return fmt.Errorf("'*' is not allowed as a key name for Rename")
	}
	// Check if inputKey exists in the map
	value, exists := d.Body[inputKey]
	if !exists {
		return fmt.Errorf("key %s not found in body", inputKey)
	}
	// Set the value of outputKey to the value of inputKey
	d.Body[outputKey] = value
	// Delete the inputKey to complete the "rename"
	delete(d.Body, inputKey)
	return nil
}

// Duplicate function, copies  the value associated with inputKey in d.Body to a new key named outputKey, without removing the original key.
func (d *Data) Duplicate(inputKey string, outputKey string) error {
	if inputKey == "*" || outputKey == "*" {
		return fmt.Errorf("'*' is not allowed as a key name for Duplicate")
	}
	// Check if inputKey exists in the map
	value, exists := d.Body[inputKey]
	if !exists {
		return fmt.Errorf("key %s not found in body", inputKey)
	}

	// Set the value of outputKey to the value of inputKey
	d.Body[outputKey] = value

	return nil
}
