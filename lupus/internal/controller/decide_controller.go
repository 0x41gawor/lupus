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
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	lupusv1 "github.com/0x41gawor/lupus/api/v1"
)

// DecideReconciler reconciles a Decide object
type DecideReconciler struct {
	client.Client
	Scheme *runtime.Scheme

	LastUpdated time.Time
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
	// Create a logger from the context
	logger := log.FromContext(ctx)
	logger.Info("=================== START OF DECIDE Reconciler: \n")

	// Step 1: Fetch the Decide instance
	var decide lupusv1.Decide
	if err := r.Get(ctx, req.NamespacedName, &decide); err != nil {
		logger.Info("Failed to fetch Decide instance", "error", err)
		// If the resource is not found, return without requeuing
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	logger.Info("Decide instance fetched successfully")

	// Step 2: Check if the status of Decide is empty, if yes, it's the first run
	if decide.Status.LastUpdated.IsZero() {
		logger.Info("No need to reconcile")
		return ctrl.Result{}, nil
	}

	// If the status last update time is not greater (not after) the time of last update, no need to reconcile
	if !r.LastUpdated.IsZero() && !decide.Status.LastUpdated.Time.After(r.LastUpdated) {
		logger.Info("No need to reconcile DECISION")
		return ctrl.Result{}, nil
	}

	// If this is the first run without any real requests, clear the input in status
	if r.LastUpdated.IsZero() {
		decide.Status.Input = runtime.RawExtension{}
	}

	// Set the controller's last updated time to the current time in the resource status
	r.LastUpdated = decide.Status.LastUpdated.Time

	// Step 3: Pass data into OPA by performing an HTTP request
	input := decide.Status.Input

	// Construct the HTTP request for OPA (similar to the curl example)
	url := decide.Spec.Url.Path // Replace with your actual OPA URL

	// Create the body with the input formatted as JSON
	body := map[string]interface{}{
		"input": input,
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		logger.Error(err, "Failed to marshal input into JSON")
		return ctrl.Result{}, err
	}

	// Create the HTTP request
	reqBody := bytes.NewBuffer(bodyBytes)
	httpReq, err := http.NewRequest("POST", url, reqBody)
	if err != nil {
		logger.Error(err, "Failed to create HTTP request")
		return ctrl.Result{}, err
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")

	// Send the HTTP request
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		logger.Error(err, "Failed to send HTTP request to OPA")
		return ctrl.Result{}, err
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err, "Failed to read response body from OPA")
		return ctrl.Result{}, err
	}

	// Check if the response is successful
	if resp.StatusCode != http.StatusOK {
		logger.Error(fmt.Errorf("non-OK HTTP status"), "OPA returned an error", "status", resp.StatusCode)
		return ctrl.Result{}, fmt.Errorf("OPA returned non-OK status: %d", resp.StatusCode)
	}

	// Step 4: Parse the response (assuming the response is also a JSON object)
	var output runtime.RawExtension
	if err := json.Unmarshal(respBody, &output); err != nil {
		logger.Error(err, "Failed to unmarshal OPA response")
		return ctrl.Result{}, err
	}

	logger.Info("Successfully received response from OPA", "data", string(output.Raw))

	// Step 4: Parse the response (assuming the response is a JSON object with a "result" key)
	var outputMap map[string]interface{}
	if err := json.Unmarshal(output.Raw, &outputMap); err != nil {
		logger.Error(err, "Failed to unmarshal OPA response into map")
		return ctrl.Result{}, err
	}

	// Rename the "result" key to "commands"
	if result, ok := outputMap["result"]; ok {
		outputMap["commands"] = result
		delete(outputMap, "result")
	}

	// Marshal the modified map back into runtime.RawExtension
	modifiedOutput, err := json.Marshal(outputMap)
	if err != nil {
		logger.Error(err, "Failed to marshal modified output back to JSON")
		return ctrl.Result{}, err
	}

	// Assign the modified output to the output RawExtension
	output.Raw = modifiedOutput

	// Log the modified output JSON data
	logger.Info("Modified OPA output data", "output", string(output.Raw))

	// Step 5: Fetch nextElement resource
	resourceName := decide.Spec.NextElement
	resourceNamespace := "default"

	var execute lupusv1.Execute
	// With this:
	err = r.Get(ctx, types.NamespacedName{Name: resourceName, Namespace: resourceNamespace}, &execute)
	if err != nil {
		logger.Error(err, "Failed to get Execute resource")
		return ctrl.Result{}, nil
	}

	// Step 6: Set the fields in the nextElement resource
	execute.Status.Input = output
	execute.Status.LastUpdated = decide.Status.LastUpdated

	// Step 7: Update the nextElement resource status
	if err := r.Status().Update(ctx, &execute); err != nil {
		logger.Error(err, "Failed to update Execute status")
		return ctrl.Result{}, err
	}

	logger.Info("Execute resource updated successfully")

	// Step 8: Reconcile again when the status changes
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DecideReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&lupusv1.Decide{}).
		Complete(r)
}
