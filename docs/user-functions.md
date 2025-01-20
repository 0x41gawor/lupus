# User-functions

In Lupus we adopted the requirement of [data-driven](defs.md#data-driven) design, by which [lupus-elements](defs.md#lupus-element) do not perform [compuation-part](defs.md#computing-part) of the [loop-logic](defs.md#loop-logic). Instead it is delegated to the [external-element](defs.md#external-element). As for now, [external-elements](defs.md#external-element) can be any HTTP Server implemented by the user. Especially, we recommend use of [Open Policy Agents](open-policy-agents.md). 

We have two aspects that arise here:
- What if someone needs to perform small and simple operation on Data (like addition of two fields) and it is uneconomic to deploy a special HTTP server for that
- Any framework platform should be extensible. We should given some mechanism of extendind our platform

Fot this two reasons we are deploying a feature called "User-defined, internal Go functions" or shortly - "User functions".

User can define its own snippets of Go code as functions and call them as one of the Destination in Send Action.

User-Functions act completely the same as [external-elements](defs.md#external-element), with the difference that they run inside a Kubernetes cluster, not outside of it. 

## Usage
Ok, but where do user "writes" its functions?

In [lupus/internal/controller/user-functions.go](../lupus/internal/controller/user-functions.go).

This file already contains one func that states as an example.

```go
// Exemplary user-function. It just returns the input
func (UserFunctions) Echo(input interface{}) (interface{}, error) {
	return input, nil
}
```

What is `UserFunctions` here? It is a struct, that will gather user-functions as its methods.

```go
// UserFunctions struct for user-defined, internal functions
type UserFunctions struct{}
```

But how we will convert string representing the function name from the Destination spec to a given user-defined function?

Fortunately in Go func can be treated as a type and referred as value, thus it can be a value in key-value map. We will simply make aforementioned strings a keys.

```go
// A global map to store function references
var FunctionRegistry = map[string]func(input interface{}) (interface{}, error){}
```

Now we will use the reflect library to fill the FuntionRegistry map.

```go
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
```

And such function will be called at the initialization of the controller package and with UserFunctions{} strcut as input.

```go
func init() {
    // Fill in the FunctionRegistry map with functions defined as a method of UserFunctions{}
	RegisterFunctions(UserFunctions{})
}
```

All of the code can be found in:
- [lupus/internal/controller/user-functions.go](../lupus/internal/controller/user-functions.go)
- [lupus/internal/controller/init.go](../lupus/internal/controller/init.go)

And from now on, we have a map that has functions names as keys and functions itself as values. Map is initialized at the start of operator runtime.

Now, we just simply need to call apropriate function based on the func name from the spec.

```go
func sendToGoFunc(funcName string, body interface{}) (interface{}, error) {
	if fn, exists := FunctionRegistry[funcName]; exists {
		return fn(body)
	} else {
		return nil, fmt.Errorf("no such UserFunction defined")
	}
}
```

In LupN user just writes:
```yaml
  - name: "bounce"
    type: send
    send:
        inputKey: "field1"
        destination:
        type: gofunc
        gofunc:
            name: echo
        outputKey: "field2"
```