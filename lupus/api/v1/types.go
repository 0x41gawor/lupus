package v1

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"k8s.io/apimachinery/pkg/runtime"
)

// Next is used in Observe Spec
// It specifies to which element forward the input
// It allows not to forward the whole input, but also parts of it
type Next struct {
	// Type specifies the type of the element ("Observe", "Decide", "Learn", "Execute", etc.)
	Type string `json:"type" kubebuilder:"validation:Enum=observe;decide;learn;execute"`
	// Kubernetes name of the API Object
	// This is the name that you give in Master CR spec
	Name string `json:"name"`
	// List of input keys (data fields) that have to be forwarded
	// Pass array with single element '*' to forward the whole input
	Keys []string `json:"keys"`
}

// Destination in Action type and in Learn Spec
// Destination is a polymorphic field that can represent different types
type Destination struct {
	// Discriminator: "HTTP", "FILE", "gRPC", etc.
	Type string `json:"type" kubebuilder:"validation:Enum=HTTP;FILE;gRPC;Opa"`

	// HTTP-specific fields
	HTTP *HTTPDestination `json:"http,omitempty" kubebuilder:"validation:Optional"`

	// File-specific fields
	File *FileDestination `json:"file,omitempty" kubebuilder:"validation:Optional"`

	// gRPC-specific fields
	GRPC *GRPCDestination `json:"grpc,omitempty" kubebuilder:"validation:Optional"`
	// Opa-specific fields
	Opa *OpaDestination `json:"opa,omitempty" kubebuilder:"validation:Optional"`
}

// HTTPDestination defines fields specific to HTTP type
type HTTPDestination struct {
	Path   string `json:"path"`
	Method string `json:"method"`
}

// FileDestination defines fields specific to FILE type
type FileDestination struct {
	Path string `json:"path"`
}

// GRPCDestination defines fields specific to gRPC type
type GRPCDestination struct {
	Host    string `json:"host"`
	Service string `json:"service"`
	Method  string `json:"method"`
}

// OpaDestination defines fields specific to Open Policy Agent type
type OpaDestination struct {
	Path string `json:"path"`
}

// Action is used in Decide spec
// It represents the Action that Decide has to perform on Data
type Action struct {
	// Name of the Action, it is for designer to ease the management of the Loop
	Name string `json:"name"`
	// Type of Action one of send;nest,remove,rename,duplicate
	Type string `json:"type" kubebuilder:"validation:Enum=send;nest,remove,rename,duplicate,print,switch"`
	// One of these fields is not null depending on a Type.
	Send      *SendAction      `json:"send,omitempty" kubebuilder:"validation:Optional"`
	Nest      *NestAction      `json:"nest,omitempty" kubebuilder:"validation:Optional"`
	Remove    *RemoveAction    `json:"remove,omitempty" kubebuilder:"validation:Optional"`
	Rename    *RenameAction    `json:"rename,omitempty" kubebuilder:"validation:Optional"`
	Duplicate *DuplicateAction `json:"duplicate,omitempty" kubebuilder:"validation:Optional"`
	Print     *PrintAction     `json:"print,omitempty" kubebuilder:"validation:Optional"`
	Insert    *InsertAction    `json:"insert,omitempty" kubebuilder:"validation:Optional"`
	Switch    *Switch          `json:"switch,omitempty" kubebuilder:"validation:Optional"`
}

func (a *Action) String() string {
	// Start with the Name and Type fields
	result := fmt.Sprintf("Action(Name: %s, Type: %s", a.Name, a.Type)

	// Add details based on the Type of Action using the specific String() methods
	switch a.Type {
	case "send":
		if a.Send != nil {
			result += fmt.Sprintf(", %s", a.Send.String())
		} else {
			result += ", Send: <nil>"
		}
	case "nest":
		if a.Nest != nil {
			result += fmt.Sprintf(", %s", a.Nest.String())
		} else {
			result += ", Nest: <nil>"
		}
	case "remove":
		if a.Remove != nil {
			result += fmt.Sprintf(", %s", a.Remove.String())
		} else {
			result += ", Remove: <nil>"
		}
	case "rename":
		if a.Rename != nil {
			result += fmt.Sprintf(", %s", a.Rename.String())
		} else {
			result += ", Rename: <nil>"
		}
	case "duplicate":
		if a.Duplicate != nil {
			result += fmt.Sprintf(", %s", a.Duplicate.String())
		} else {
			result += ", Duplicate: <nil>"
		}
	case "print":
		if a.Print != nil {
			result += fmt.Sprintf(", %s", a.Print.String())
		} else {
			result += ", Duplicate: <nil>"
		}
	case "insert":
		if a.Insert != nil {
			result += fmt.Sprintf(", %s", a.Insert.String())
		}
	default:
		result += ", <Unknown Action Type>"
	}

	result += ")"
	return result
}

type SendAction struct {
	InputKey    string      `json:"inputKey"`
	Destination Destination `json:"destination"`
	OutputKey   string      `json:"outputKey"`
}

type NestAction struct {
	InputKeys []string `json:"inputKeys"`
	OutputKey string   `json:"outputKey"`
}

type RemoveAction struct {
	InputKeys []string `json:"inputKeys"`
}

type RenameAction struct {
	InputKey  string `json:"inputKey"`
	OutputKey string `json:"outputKey"`
}

type DuplicateAction struct {
	InputKey  string `json:"inputKey"`
	OutputKey string `json:"outputKey"`
}

type PrintAction struct {
	InputKeys []string `json:"inputKeys"`
}

type InsertAction struct {
	OutputKey string               `json:"outputKey"`
	Value     runtime.RawExtension `json:"value"` //value can be of type: int, float, bool, string
}

type Switch struct {
}

func (s *SendAction) String() string {
	return fmt.Sprintf("SendAction(InputKey: %s, Destination: %v, OutputKey: %s)", s.InputKey, s.Destination, s.OutputKey)
}

func (c *NestAction) String() string {
	return fmt.Sprintf("ConcatAction(InputKeys: %v, OutputKey: %s)", c.InputKeys, c.OutputKey)
}

func (r *RemoveAction) String() string {
	return fmt.Sprintf("RemoveAction(InputKeys: %v)", r.InputKeys)
}

func (r *RenameAction) String() string {
	return fmt.Sprintf("RenameAction(InputKey: %s, OutputKey: %s)", r.InputKey, r.OutputKey)
}

func (d *DuplicateAction) String() string {
	return fmt.Sprintf("DuplicateAction(InputKey: %s, OutputKey: %s)", d.InputKey, d.OutputKey)
}

func (d *PrintAction) String() string {
	return fmt.Sprintf("PrintAction(InputKeys: %v)", d.InputKeys)
}

func (d *InsertAction) String() string {
	return fmt.Sprintf("InsertAction(OutputKey: %s, Value: %v)", d.OutputKey, d.Value)
}

// Element is a polymorphic structure that can represent different types of specs
type Element struct {
	// Name is the name of the element
	Name string `json:"name"`
	// Type specifies the type of the element ("Observe", "Decide", "Learn", "Execute", etc.)
	Type string `json:"type" kubebuilder:"validation:Enum=observe;decide;learn;execute"`

	Observe *ObserveSpec `json:"observe,omitempty"`
	Decide  *DecideSpec  `json:"decide,omitempty"`
	Execute *ExecuteSpec `json:"execute,omitempty"`
}

type Condition struct {
	// Key indicates the field of Data that has to be retrieved
	Key string `json:"key"`
	// Operator can be one of: eq (equals to), gt/ls (greater/less than), ne (not equal to)
	Operator string `json:"operator"`
	// Value that will be comapred againts Data field
	Value interface{} `json:"value"`
	// Next specifies name of the next action to execute
	Next string `json:"next"`
}

func (c *Condition) Compare(field interface{}) (bool, error) {
	// Validate the field type
	switch v := field.(type) {
	case int, float64, bool, string:
		// Proceed with validation of c.Value
	default:
		return false, fmt.Errorf("unsupported field type: %T", v)
	}

	// Validate that c.Value has a compatible type with field
	fieldType := reflect.TypeOf(field)
	valueType := reflect.TypeOf(c.Value)

	// Attempt type conversion if types are different
	var value interface{}
	var err error
	if fieldType != valueType {
		value, err = convertValueToType(c.Value, fieldType)
		if err != nil {
			return false, fmt.Errorf("type mismatch and conversion failed: %w", err)
		}
	} else {
		value = c.Value
	}

	// Perform the comparison based on the operator
	switch c.Operator {
	case "eq":
		return reflect.DeepEqual(field, value), nil
	case "ne":
		return !reflect.DeepEqual(field, value), nil
	case "gt":
		return compareNumeric(field, value, ">")
	case "ls":
		return compareNumeric(field, value, "<")
	default:
		return false, fmt.Errorf("unsupported operator: %s", c.Operator)
	}
}

// Helper to convert value to a specific type
func convertValueToType(value interface{}, targetType reflect.Type) (interface{}, error) {
	switch targetType.Kind() {
	case reflect.Int:
		switch v := value.(type) {
		case int:
			return v, nil
		case float64:
			return int(v), nil
		case string:
			i, err := strconv.Atoi(v)
			if err != nil {
				return nil, fmt.Errorf("cannot convert string to int: %w", err)
			}
			return i, nil
		default:
			return nil, fmt.Errorf("cannot convert %T to int", v)
		}
	case reflect.Float64:
		switch v := value.(type) {
		case int:
			return float64(v), nil
		case float64:
			return v, nil
		case string:
			f, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return nil, fmt.Errorf("cannot convert string to float64: %w", err)
			}
			return f, nil
		default:
			return nil, fmt.Errorf("cannot convert %T to float64", v)
		}
	case reflect.Bool:
		switch v := value.(type) {
		case bool:
			return v, nil
		case string:
			if v == "true" {
				return true, nil
			} else if v == "false" {
				return false, nil
			}
			return nil, fmt.Errorf("cannot convert string to bool")
		default:
			return nil, fmt.Errorf("cannot convert %T to bool", v)
		}
	case reflect.String:
		switch v := value.(type) {
		case string:
			return v, nil
		default:
			return fmt.Sprintf("%v", v), nil
		}
	default:
		return nil, fmt.Errorf("unsupported target type: %v", targetType)
	}
}

// Helper for numeric comparisons
func compareNumeric(field interface{}, value interface{}, operator string) (bool, error) {
	var fieldVal, valueVal float64

	switch v := field.(type) {
	case int:
		fieldVal = float64(v)
	case float64:
		fieldVal = v
	default:
		return false, errors.New("field is not numeric")
	}

	switch v := value.(type) {
	case int:
		valueVal = float64(v)
	case float64:
		valueVal = v
	default:
		return false, errors.New("value is not numeric")
	}

	switch operator {
	case ">":
		return fieldVal > valueVal, nil
	case "<":
		return fieldVal < valueVal, nil
	default:
		return false, fmt.Errorf("unsupported numeric operator: %s", operator)
	}
}
