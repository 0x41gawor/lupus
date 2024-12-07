package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"k8s.io/apimachinery/pkg/runtime"
)

// Generic function to validate an interface{} against a predicate JSON structure
func ValidateJSON(expected interface{}, actual interface{}) (bool, error) {
	// Ensure both inputs are valid JSON-compatible types
	expectedBytes, err := json.Marshal(expected)
	if err != nil {
		return false, fmt.Errorf("error marshaling expected JSON: %w", err)
	}

	var expectedMap map[string]interface{}
	if err := json.Unmarshal(expectedBytes, &expectedMap); err != nil {
		return false, fmt.Errorf("error unmarshaling expected JSON to map: %w", err)
	}

	// Ensure `actual` is a map[string]interface{}
	actualMap, ok := actual.(map[string]interface{})
	if !ok {
		return false, errors.New("actual value is not a JSON-compatible map")
	}

	// Compare the keys and values in expected and actual maps
	for key, expectedValue := range expectedMap {
		actualValue, exists := actualMap[key]
		if !exists {
			return false, fmt.Errorf("missing key: %s", key)
		}

		// Use reflect.DeepEqual to compare values
		if !reflect.DeepEqual(expectedValue, actualValue) {
			return false, fmt.Errorf("mismatched value for key '%s': expected %v, got %v", key, expectedValue, actualValue)
		}
	}

	return true, nil
}

func StringToInterface(str string) (interface{}, error) {
	// Variable to hold the result as an interface{}
	var data interface{}
	// Unmarshal the JSON string into the interface{}
	err := json.Unmarshal([]byte(str), &data)
	if err != nil {
		fmt.Printf("Error unmarshaling JSON: %v\n", err)
		return nil, err
	}
	return data, err
}

func TestData_Get(t *testing.T) {
	// Define the JSON structure
	serverDataJSON := `{
		"cpu": {
			"in_use": 12,
			"license": 8
		},
		"ram": {
			"ram2": {
				"in_use": 10,
				"license": 5
			}
		}
	}`

	// Create a runtime.RawExtension
	rawExtension := &runtime.RawExtension{
		Raw: []byte(serverDataJSON),
	}
	data, err := NewData(*rawExtension)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(data.String())
	// test get non-nested field
	cpu, err := data.Get([]string{"cpu"})
	if err != nil {
		t.Fatal(err.Error())
	}
	cpuStr, err := InterfaceToString(cpu)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(cpuStr)
	// test multiple keys
	mul, err := data.Get([]string{"cpu", "ram.ram2"})
	if err != nil {
		t.Fatal(err.Error())
	}
	mulStr, err := InterfaceToString(mul)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(mulStr)
	// test get nested field
	ram2, err := data.Get([]string{"ram.ram2.in_use"})
	if err != nil {
		t.Fatal(err.Error())
	}
	ram2Str, err := InterfaceToString(ram2)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(ram2Str)
	// test non-nested set
	capitalsStr := `{
		"italy": "Rome",
		"germany": "Berlin"
	}`
	capitals, err := StringToInterface(capitalsStr)
	if err != nil {
		t.Fatal(err.Error())
	}
	err = data.Set("cpu", capitals)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(data.String())
	// test.nested Field
	germanyInfoStr := `{
		"population": 86,
		"language": "German",
		"nested": {
			"elo": 1,
			"freitag": "germany"
		}
	}`
	germanyInfo, err := StringToInterface(germanyInfoStr)
	if err != nil {
		t.Fatal(err.Error())
	}
	err = data.Set("cpu.germany", germanyInfo)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(data.String())
	// test Nest
	err = data.Nest([]string{"ram", "cpu.germany"}, "dump")
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(data.String())
	// test Remove
	err = data.Remove([]string{"cpu", "dump.ram.ram2"})
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(data.String())
	// test Rename
	err = data.Rename("dump", "trump")
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(data.String())
	err = data.Rename("trump.ram", "sram")
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(data.String())
	// test duplicate
	err = data.Duplicate("trump.germany", "trump.italy")
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(data.String())
	err = data.Duplicate("trump.sram", "sram")
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(data.String())
	err = data.Duplicate("sram", "dram")
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(data.String())
	err = data.Duplicate("sram", "elo.kozak")
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(data.String())
	err = data.Print([]string{"*"})
	if err != nil {
		t.Fatal(err.Error())
	}
	err = data.Print([]string{"elo"})
	if err != nil {
		t.Fatal(err.Error())
	}
	err = data.Print([]string{"elo", "trump"})
	if err != nil {
		t.Fatal(err.Error())
	}
	err = data.Print([]string{"trump.germany.language"})
	if err != nil {
		t.Fatal(err.Error())
	}
	data2Str := `
		{
			"elo": "hej",
			"poland": {
				"capital": "Warsaw",
				"population": 38
			}
		}
	`
	// Create a runtime.RawExtension
	rawExtension2 := &runtime.RawExtension{
		Raw: []byte(data2Str),
	}
	data2, err := NewData(*rawExtension2)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(data2.String())
	err = data2.Insert("*", runtime.RawExtension{Raw: []byte(`{"elo2":"siema"}`)})
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(data2.String())
	err = data2.Insert("poland", runtime.RawExtension{Raw: []byte(`{"cities": ["Poznan", "Gdansk", "Krakow"]}`)})
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(data2.String())
	err = data2.Insert("poland.capital", runtime.RawExtension{Raw: []byte(`{"first": "Gniezno"}`)})
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(data2.String())
	err = data2.Insert("poland.capital", runtime.RawExtension{Raw: []byte(`{"second": "Krakow"}`)})
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(data2.String())
	err = data2.Set("*", rawExtension2)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(data2.String())
}
