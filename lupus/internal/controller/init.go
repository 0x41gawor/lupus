package controller

import (
	"reflect"
)

// UserFunctions struct for user-defined functions
type UserFunctions struct{}

// A global map to store function references
var FunctionRegistry = map[string]func(input interface{}) (interface{}, error){}

// RegisterFunctions dynamically registers user-defined functions
func RegisterFunctions(target interface{}) {
	t := reflect.TypeOf(target)
	v := reflect.ValueOf(target)

	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)

		// Ensure the method matches the required signature
		if method.Type.NumIn() == 2 && // Receiver + input
			method.Type.NumOut() == 2 && // Output + error
			method.Type.In(1).Kind() == reflect.Interface && // Input: interface{}
			method.Type.Out(0).Kind() == reflect.Interface && // Output: interface{}
			method.Type.Out(1).Implements(reflect.TypeOf((*error)(nil)).Elem()) { // Second output: error

			funcName := method.Name
			FunctionRegistry[funcName] = func(input interface{}) (interface{}, error) {
				// Call the user-defined function
				result := method.Func.Call([]reflect.Value{v, reflect.ValueOf(input)})

				// Handle result[1] (error) being nil
				var err error
				if !result[1].IsNil() {
					err = result[1].Interface().(error)
				}

				return result[0].Interface(), err
			}
		}
	}
}

func init() {
	RegisterFunctions(UserFunctions{})
}
