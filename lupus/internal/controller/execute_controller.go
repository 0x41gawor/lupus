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
	"fmt"
	"net/http"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	lupusv1 "github.com/0x41gawor/lupus/api/v1"
)

// ExecuteReconciler reconciles a Execute object
type ExecuteReconciler struct {
	client.Client
	Scheme *runtime.Scheme

	LastUpdated time.Time
}

// +kubebuilder:rbac:groups=lupus.gawor.io,resources=executes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=lupus.gawor.io,resources=executes/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=lupus.gawor.io,resources=executes/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Execute object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *ExecuteReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// Create a logger from the context
	logger := log.FromContext(ctx)
	logger.Info("=================== START OF EXECUTE Reconciler: \n")
	// Step 1: Fetch the Observe instance
	var execute lupusv1.Execute
	if err := r.Get(ctx, req.NamespacedName, &execute); err != nil {
		logger.Info("Failed to fetch Execute instance", "error", err)
		// If the resource is not found, we return and don't requeue
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	logger.Info("Execute instance fetched successfully")
	// Step 2: Check if status of Execute is empty, if yes it is the first run and reconciler shoudl return
	if execute.Status.LastUpdated.IsZero() {
		logger.Info("No need to reconcile")
		return ctrl.Result{}, nil
	}
	// If the status last update time is not highger (not after) the time of last update, there no need to reconcile
	if !r.LastUpdated.IsZero() && !execute.Status.LastUpdated.Time.After(r.LastUpdated) {
		logger.Info("No need to reconcile EXECUTE")
		return ctrl.Result{}, nil
	}
	// If this is the first run of controller without any real requests clear the command clear the status
	if r.LastUpdated.IsZero() {
		execute.Status.Input = runtime.RawExtension{}
	}
	// If we reconcile then set the controller's last updated time same as in the resource spec
	r.LastUpdated = execute.Status.LastUpdated.Time

	// step 3: Pass input into monitored-system
	input := execute.Status.Input
	// call http client and hit execute.Spec.MonitoredSystemUrl.Path with execute.Spec.MonitoredSystemUrl.Method and input as a json data in body

	// Step 4: Call the HTTP client to hit the monitored system URL
	url := execute.Spec.Url.Path // Assuming this holds the full URL
	method := execute.Spec.Url.Method
	logger.Info("Sending request to monitored system", "url", url, "method", method)

	// Prepare the HTTP request
	reqBody := bytes.NewBuffer(input.Raw)
	httpReq, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		logger.Error(err, "Failed to create HTTP request")
		return ctrl.Result{}, err
	}

	// Set Content-Type to application/json
	httpReq.Header.Set("Content-Type", "application/json")

	// Send the HTTP request
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		logger.Error(err, "Error sending HTTP request to monitored system")
		return ctrl.Result{}, err
	}
	defer resp.Body.Close()

	// Step 5: Handle the response
	if resp.StatusCode != http.StatusOK {
		logger.Error(fmt.Errorf("non-OK HTTP status"), "Monitored system returned an error", "status", resp.StatusCode)
		return ctrl.Result{}, fmt.Errorf("monitored system returned non-OK status: %d", resp.StatusCode)
	}

	logger.Info("Request to monitored system was successful", "status", resp.StatusCode)

	// Step 6: Optionally, update the status or do additional work after successful reconciliation

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ExecuteReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&lupusv1.Execute{}).
		Complete(r)
}
