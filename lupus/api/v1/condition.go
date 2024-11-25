package v1

import (
	"encoding/json"
	"fmt"

	runtime "k8s.io/apimachinery/pkg/runtime"
)

// Condition represent signle condition present in Switch action.
type Condition struct {
	// Key indicates the field of Data that has to be retrieved
	Key string `json:"key"`
	// Operator defines the comparison operation, e.g., eq, ne, gt, lt
	Operator string `json:"operator" kubebuilder:"validation:Enum=eq,ne,gt,lt"`
	// Type specifies the type of the value: string, int, float, bool
	Type string `json:"type" kubebuilder:"validation:Enum=string,int,float,bool"`
	// BoolCondition specifies the condition for boolean values
	BoolCondition *BoolCondition `json:"bool,omitempty" kubebuilder:"validation:Optional"`
	// IntCondition specifies the condition for integer values
	IntCondition *IntCondition `json:"int,omitempty" kubebuilder:"validation:Optional"`
	// StringCondition specifies the condition for string values
	StringCondition *StringCondition `json:"string,omitempty" kubebuilder:"validation:Optional"`
	// Next specifies the name of the next action to execute
	Next string `json:"next"`
}

// BoolCondition defines a boolean-specific condition
type BoolCondition struct {
	Value bool `json:"value"`
}

// IntCondition defines an integer-specific condition
type IntCondition struct {
	Value int `json:"value"`
}

// StringCondition defines a string-specific condition
type StringCondition struct {
	Value string `json:"value"`
}

func (c *Condition) Evaluate(field runtime.RawExtension) (bool, error) {
	var fieldValue interface{}

	// Convert RawExtension to interface{}
	if err := json.Unmarshal(field.Raw, &fieldValue); err != nil {
		return false, fmt.Errorf("failed to unmarshal field: %w", err)
	}

	// Determine the type and evaluate the condition
	switch c.Type {
	case "bool":
		if c.BoolCondition == nil {
			return false, fmt.Errorf("expected BoolCondition for type bool")
		}
		value, ok := fieldValue.(bool)
		if !ok {
			return false, fmt.Errorf("expected bool field, got %T", fieldValue)
		}
		return compareEquality(value, c.BoolCondition.Value, c.Operator)
	case "int":
		if c.IntCondition == nil {
			return false, fmt.Errorf("expected IntCondition for type int")
		}
		// Handle JSON numbers, which may be parsed as float64
		var value int
		if err := parseJSONNumber(fieldValue, &value); err != nil {
			return false, err
		}
		return compareOrdered(value, c.IntCondition.Value, c.Operator)
	case "string":
		if c.StringCondition == nil {
			return false, fmt.Errorf("expected StringCondition for type string")
		}
		value, ok := fieldValue.(string)
		if !ok {
			return false, fmt.Errorf("expected string field, got %T", fieldValue)
		}
		return compareEquality(value, c.StringCondition.Value, c.Operator)
	default:
		return false, fmt.Errorf("unsupported type: %s", c.Type)
	}
}

// Helper function to handle JSON numbers
func parseJSONNumber(fieldValue interface{}, target *int) error {
	switch v := fieldValue.(type) {
	case float64:
		*target = int(v)
		return nil
	case int:
		*target = v
		return nil
	default:
		return fmt.Errorf("expected number, got %T", fieldValue)
	}
}

// Compare for equality-based operators (eq, ne)
func compareEquality[T comparable](field, value T, operator string) (bool, error) {
	switch operator {
	case "eq":
		return field == value, nil
	case "ne":
		return field != value, nil
	default:
		return false, fmt.Errorf("unsupported operator for equality: %s", operator)
	}
}

// Compare for ordered operators (gt, lt, eq, ne)
func compareOrdered[T int | float64](field, value T, operator string) (bool, error) {
	switch operator {
	case "eq":
		return field == value, nil
	case "ne":
		return field != value, nil
	case "gt":
		return field > value, nil
	case "lt":
		return field < value, nil
	default:
		return false, fmt.Errorf("unsupported operator for ordered comparison: %s", operator)
	}
}
