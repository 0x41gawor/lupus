package controller

import (
	"fmt"

	v1 "github.com/0x41gawor/lupus/api/v1"
	util "github.com/0x41gawor/lupus/internal/util"
)

func PerformAction(data *util.Data, action v1.Action) (string, error) {
	if action.Name == "" {
		return "", fmt.Errorf("empty action")
	}
	switch action.Type {
	case "send":
		input, err := data.Get([]string{action.Send.InputKey})
		if err != nil {
			err = fmt.Errorf("cannot get Data inputKey object: %w", err)
			return "exit", err
		}
		output, err := sendToDestination(input, action.Send.Destination)
		if err != nil {
			err = fmt.Errorf("send to destination failed: %w", err)
			return "exit", err
		}
		if err = data.Set(action.Send.OutputKey, output); err != nil {
			err = fmt.Errorf("cannot set data field: %w", err)
			return "exit", err
		}
	case "nest":
		err := data.Nest(action.Nest.InputKeys, action.Nest.OutputKey)
		if err != nil {
			err = fmt.Errorf("cannot nest data field: %w", err)
			return "exit", err
		}
	case "remove":
		err := data.Remove(action.Remove.InputKeys)
		if err != nil {
			err = fmt.Errorf("cannot remove data field: %w", err)
			return "exit", err
		}
	case "rename":
		err := data.Rename(action.Rename.InputKey, action.Rename.OutputKey)
		if err != nil {
			err = fmt.Errorf("cannot rename data field: %w", err)
			return "exit", err
		}
	case "duplicate":
		err := data.Duplicate(action.Duplicate.InputKey, action.Duplicate.OutputKey)
		if err != nil {
			err = fmt.Errorf("cannot duplicate data field: %w", err)
			return "exit", err
		}
	case "insert":
		err := data.Insert(action.Insert.OutputKey, action.Insert.Value)
		if err != nil {
			err = fmt.Errorf("cannot insert data field: %w", err)
			return "exit", err
		}
	case "print":
		fmt.Printf("----------------%s-------------------Data:-----------------------------------------\n", action.Name)
		err := data.Print(action.Print.InputKeys)
		if err != nil {
			err = fmt.Errorf("cannot print data: %w", err)
			return "exit", err
		}
	case "switch":
		for _, condition := range action.Switch.Conditions {
			field, err := data.Get([]string{condition.Key})
			if err != nil {
				err = fmt.Errorf("could not retrieve data field for evaluating condition: %w", err)
				return "exit", err
			}
			eval, err := condition.Evaluate(*field)
			if err != nil {
				err = fmt.Errorf("error during condition evaluation: %w", err)
				return "exit", err
			}
			if eval {
				return condition.Next, nil
			}
		}
	}
	return action.Next, nil
}

// ActionMap is an auxilliary type used by Element Controller's Reconcile func
// As routing between actions is done by name, we need a map that maps names to specific Actions
type ActionMap map[string]v1.Action

func ConvertActionsToMap(actions []v1.Action) (ActionMap, error) {
	actionMap := make(ActionMap)
	for _, action := range actions {
		if action.Name == "" {
			return nil, fmt.Errorf("action with empty name cannot be added to the map")
		}
		if _, exists := actionMap[action.Name]; exists {
			return nil, fmt.Errorf("duplicate action name found: %s", action.Name)
		}
		actionMap[action.Name] = action
	}
	return actionMap, nil
}
