/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	v1 "github.com/0x41gawor/lupus/api/v1"
	util "github.com/0x41gawor/lupus/internal/util"
	"github.com/go-logr/logr"
)

// DecideReconciler reconciles a Decide object
type DecideReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	// Static fields of Reconciler
	Logger      logr.Logger
	ElementType string
	// Map that holds of reconciler state for each element
	instanceState map[types.NamespacedName]*util.ElementInstanceState
}

// +kubebuilder:rbac:groups=lupus.gawor.io,resources=decides,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=lupus.gawor.io,resources=decides/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=lupus.gawor.io,resources=decides/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Decide object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *DecideReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	r.ElementType = "Decide"
	r.Logger = log.FromContext(ctx)
	r.Logger.Info(fmt.Sprintf("=================== START OF %s Reconciler: \n", strings.ToUpper(r.ElementType)))

	// Step 1 - (k8s) Fetch reconciled resource instance
	// Step 1 - (lupus) Fetch element
	var element v1.Decide
	if err := r.Get(ctx, req.NamespacedName, &element); err != nil {
		r.Logger.Info(fmt.Sprintf("Failed to fetch %s instance", r.ElementType), "error", err)
		// If the resource is not found, we return and don't requeue
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Step 2 - Checks
	// check if we have such element in our state map
	_, exists := r.instanceState[req.NamespacedName]
	if !exists {
		// Initialize instance data if it doesn't exist
		r.instanceState[req.NamespacedName] = &util.ElementInstanceState{}
		// clear status as it can contain some garbage
		element.Status.Input = runtime.RawExtension{}
		element.Status.LastUpdated = metav1.Time{}
		// set the flag
		r.instanceState[req.NamespacedName].IsAfterDryRun = true
		r.Logger.Info("This is the dry run, no need to reconcile")
		return ctrl.Result{}, nil
	}
	// Check for double update in single loop iteration. If r.LastUpdated time is zero it means it is the 2nd run (so double update can't happen)
	// If the Status.LastUpdated time is non-zero we have to check if its not the same as the previous one
	if !r.instanceState[req.NamespacedName].LastUpdated.IsZero() && !element.Status.LastUpdated.Time.After(r.instanceState[req.NamespacedName].LastUpdated) {
		// If this condition is true it means we are reconciling again in the same iteration
		r.Logger.Info("Already reconciled in this loop iteration, no need to reconcile")
		return ctrl.Result{}, nil
	}

	// Step 3 - We reconcile, so let's begin the process with variable settings
	var input runtime.RawExtension = element.Status.Input
	r.instanceState[req.NamespacedName].LastUpdated = element.Status.LastUpdated.Time

	// Step 4 - (Go) Unmarshall input into map[string]interface{}
	// Step 4 - (Lupus) Create Data object
	data, err := util.NewData(input)
	if err != nil {
		r.Logger.Error(err, "Cannot unmarshall the input")
		return ctrl.Result{}, nil
	}

	// Step 5 - Perform actions
	for _, action := range element.Spec.Actions {
		switch action.Type {
		case "send":
			input, err := data.Get([]string{action.Send.InputKey})
			if err != nil {
				r.Logger.Error(err, "cannot get Data inputKey object")
				return ctrl.Result{}, nil
			}
			output, err := sendToDestination(input, action.Send.Destination)
			if err != nil {
				r.Logger.Error(err, "cannot get Data inputKey object")
				return ctrl.Result{}, nil
			}
			if err = data.Set(action.Send.OutputKey, output); err != nil {
				r.Logger.Error(err, "cannot set data field")
			}
		case "nest":
			err := data.Nest(action.Nest.InputKeys, action.Nest.OutputKey)
			if err != nil {
				r.Logger.Error(err, "cannot nest data field")
				return ctrl.Result{}, nil
			}
		case "remove":
			err := data.Remove(action.Remove.InputKeys)
			if err != nil {
				r.Logger.Error(err, "cannot remove data field")
				return ctrl.Result{}, nil
			}
		case "rename":
			err := data.Rename(action.Rename.InputKey, action.Rename.OutputKey)
			if err != nil {
				r.Logger.Error(err, "cannot rename data field")
				return ctrl.Result{}, nil
			}
		case "duplicate":
			err := data.Duplicate(action.Duplicate.InputKey, action.Duplicate.OutputKey)
			if err != nil {
				r.Logger.Error(err, "cannot duplicate data field")
				return ctrl.Result{}, nil
			}
		case "print":
			fmt.Printf("----------------%s-------------------Data:-----------------------------------------\n", action.Name)
			err := data.Print(action.Print.InputKeys)
			if err != nil {
				r.Logger.Error(err, "cannot print data")
			}
		}
	}
	// Step 6 - Send data output to the next elements
	for _, next := range element.Spec.Next {
		// prepare output based on keys
		output, err := data.Get(next.Keys)
		if err != nil {
			r.Logger.Error(err, "Cannot retrieve keys from data")
			return ctrl.Result{}, nil
		}
		// update status of next element with output
		updateTime := metav1.Time{Time: r.instanceState[req.NamespacedName].LastUpdated}
		objectKey := client.ObjectKey{Name: element.Spec.Master + "-" + next.Name, Namespace: "default"}
		err = r.updateStatus(ctx, next.Type, objectKey, updateTime, *output)
		if err != nil {
			r.Logger.Error(err, "Failed to update status of next element")
		}
	}
	r.Logger.Info("Decide sucessfully reconciled", "name", req.Name)
	return ctrl.Result{}, nil
}

func sendToDestination(input interface{}, dest v1.Destination) (interface{}, error) {
	switch dest.Type {
	case "HTTP":
		res, err := sendToHTTP(dest.HTTP.Path, dest.HTTP.Method, input)
		if err != nil {
			return nil, err
		}
		return res, nil
	case "Opa":
		res, err := sendToOpa(dest.Opa.Path, input)
		if err != nil {
			return nil, err
		}
		return res, nil
	default:
		return nil, fmt.Errorf("no such destination type implemented yet")
	}
}

func sendToOpa(path string, reqBody interface{}) (interface{}, error) {
	wrappedBody := map[string]interface{}{
		"input": reqBody,
	}

	// Call sendToHTTP to get the response
	res, err := sendToHTTP(path, "POST", wrappedBody)
	if err != nil {
		return nil, err
	}
	resMap, err := util.InterfaceToMap(res)
	if err != nil {
		return nil, fmt.Errorf("unexpected response format, not a map")
	}
	// Return only the content of "result"
	if result, ok := resMap["result"]; ok {
		return result, nil
	}
	return nil, fmt.Errorf("no 'result' field in response")
}

func sendToHTTP(path string, method string, body interface{}) (interface{}, error) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	reqBody := bytes.NewBuffer(bodyBytes)
	httpReq, err := http.NewRequest(method, path, reqBody)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-ok HTTP Status")
	}

	var res interface{}
	if err := json.Unmarshal(respBody, &res); err != nil {
		return nil, err
	}
	return res, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DecideReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// Initialize instanceState map if it's nil
	if r.instanceState == nil {
		r.instanceState = make(map[types.NamespacedName]*util.ElementInstanceState)
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1.Decide{}).
		Complete(r)
}

func (r *DecideReconciler) updateStatus(ctx context.Context, t string, objKey client.ObjectKey, updateTime metav1.Time, output runtime.RawExtension) error {
	switch t {
	case "decide":
		nextElement := v1.Decide{}
		// (k8s) fetch lupus element from kube-api-server
		err := r.Get(ctx, objKey, &nextElement)
		if err != nil {
			r.Logger.Error(err, "Failed to get next element: Decide")
		}
		// (k8s/go) set fields
		nextElement.Status.Input = output
		nextElement.Status.LastUpdated = updateTime
		// (k8s) update k8s object in kube-api-server
		if err := r.Status().Update(ctx, &nextElement); err != nil {
			r.Logger.Error(err, "Failed to update next element (Decide) status")
		}
	case "execute":
		nextElement := v1.Execute{}
		// (k8s) fetch lupus element from kube-api-server
		err := r.Get(ctx, objKey, &nextElement)
		if err != nil {
			r.Logger.Error(err, "Failed to get next element: Execute")
		}
		// (k8s/go) set fields
		nextElement.Status.Input = output
		nextElement.Status.LastUpdated = updateTime
		// (k8s) update k8s object in kube-api-server
		if err := r.Status().Update(ctx, &nextElement); err != nil {
			r.Logger.Error(err, "Failed to update next element (Execute) status")
		}
	default:
		err := errors.New("cannot pass input to any other element type than Execute or Decide")
		r.Logger.Error(err, "Unrecognized element type")
		return err
	}
	return nil
}
