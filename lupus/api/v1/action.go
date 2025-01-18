package v1

import (
	"fmt"

	runtime "k8s.io/apimachinery/pkg/runtime"
)

// Action represents operation that is performed on Data
// Action is used in Element spec. Element has a list of Actions and executes them
// In general each action has an input and output keys that define which Data fields it has to work on
// Each action indicates the name of the next Action in Action Chain
// There is special type - Switch. Actually, it does not perform any operation on Data, but rather controls the flow of Actions chain
type Action struct {
	// Name of the Action, it is for designer to ease the management of the Loop
	Name string `json:"name"`
	// Type of Action
	Type string `json:"type" kubebuilder:"validation:Enum=send,nest,remove,rename,duplicate,print,insert,switch"`
	// One of these fields is not null depending on a Type.
	Send      *SendAction      `json:"send,omitempty" kubebuilder:"validation:Optional"`
	Nest      *NestAction      `json:"nest,omitempty" kubebuilder:"validation:Optional"`
	Remove    *RemoveAction    `json:"remove,omitempty" kubebuilder:"validation:Optional"`
	Rename    *RenameAction    `json:"rename,omitempty" kubebuilder:"validation:Optional"`
	Duplicate *DuplicateAction `json:"duplicate,omitempty" kubebuilder:"validation:Optional"`
	Print     *PrintAction     `json:"print,omitempty" kubebuilder:"validation:Optional"`
	Insert    *InsertAction    `json:"insert,omitempty" kubebuilder:"validation:Optional"`
	Switch    *Switch          `json:"switch,omitempty" kubebuilder:"validation:Optional"`
	// Next is the name of the next action to execute, in the case of Switch-type action it stands as default branch
	Next string `json:"next"`
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

// SendAction is used to make call to external-element
// Element's controller obtains a data field using InputKey,
// and attaches it as a json body when perfoming a call to destination.
// Respnse is saved in data under an OutputKey
type SendAction struct {
	InputKey    string      `json:"inputKey"`
	Destination Destination `json:"destination"`
	OutputKey   string      `json:"outputKey"`
}

// NestAction is used to group a number of data-fields together.
// Element's controllers gathers fields indicates by InputKeys list
// and nests them in a new field under an OutputKey.
type NestAction struct {
	InputKeys []string `json:"inputKeys"`
	OutputKey string   `json:"outputKey"`
}

// RemoveAction is used to delete a data-field.
// Elements's controllers removes fields indicated by the list InputKeys
type RemoveAction struct {
	InputKeys []string `json:"inputKeys"`
}

// RenameAction is used to change name of a data-field.
// InputKey indicates a field to be renamed
// OutputKey is the new field name.
type RenameAction struct {
	InputKey  string `json:"inputKey"`
	OutputKey string `json:"outputKey"`
}

// DuplicateAction is used to make a copy of data-field.
// InputKey indicates the field of which value has to be copied.
// OutputKey indicates the field to which values has to be pasted in.
type DuplicateAction struct {
	InputKey  string `json:"inputKey"`
	OutputKey string `json:"outputKey"`
}

// PrintAction is used to print value of each field indicated by InputKeys in a controller's console.
// It is useful in debugging or logging
type PrintAction struct {
	InputKeys []string `json:"inputKeys"`
}

// InsertAction is used to make a new field and insert value to it
// Normally new fields are created as an outcome of other types of actions
// It is useful in debugging or loggin, e.g. can idicate the path taken by the actions workflow
type InsertAction struct {
	OutputKey string               `json:"outputKey"`
	Value     runtime.RawExtension `json:"value"`
}

// Switch is a special type of action used for flow-control
// When Element's controller encounters switch action on the chain
// it emulates the work of switch known in other programming languages
type Switch struct {
	Conditions []Condition `json:"conditions"`
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
