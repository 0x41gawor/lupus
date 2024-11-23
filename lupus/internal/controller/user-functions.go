package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	util "github.com/0x41gawor/lupus/internal/util"
	"k8s.io/apimachinery/pkg/runtime"
)

func (UserFunctions) Echo(input interface{}) (interface{}, error) {
	fmt.Println("\n\n ---------- WYKON -----------")
	inputStr, err := util.InterfaceToString(input)
	if err != nil {
		return nil, err
	}
	fmt.Println("input: ", inputStr)
	return input, nil
}

func (UserFunctions) AddTen(input interface{}) (interface{}, error) {
	fmt.Println("\n\n ---------- EXECUTING AddTen -----------")
	inputStr, err := util.InterfaceToString(input)
	if err != nil {
		return nil, err
	}
	fmt.Println("input: ", inputStr, " type: ", reflect.TypeOf(input))
	// Check if the input is of type *runtime.RawExtension
	rawExt, ok := input.(*runtime.RawExtension)
	if !ok {
		return nil, errors.New("input is not of type *runtime.RawExtension")
	}
	// Deserialize the RawExtension into an int
	var value int
	if err := json.Unmarshal(rawExt.Raw, &value); err != nil {
		return nil, fmt.Errorf("failed to unmarshal input into int: %w", err)
	}

	// Add 10 to the deserialized value
	newValue := value + 10
	fmt.Println("Original input:", value, "Modified value:", newValue)

	return newValue, nil
}
