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
	"net/http"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
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
	logger := log.FromContext(ctx)
	logger.Info("=================== START OF EXECUTE Reconciler: ")
	println(r.LastUpdated.String())

	// Fetch the Execute instance
	var execute lupusv1.Execute
	if err := r.Get(ctx, req.NamespacedName, &execute); err != nil {
		if errors.IsNotFound(err) {
			// Execute resource not found, no need to reconcile
			logger.Info("Execute resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Failed to get Execute")
		return ctrl.Result{}, err
	}
	// If the status of resourcew was not set in cluster, there is no need to reconcile
	if execute.Status.LastUpdated.Time.IsZero() {
		logger.Info("No need to reconcile EXECUTE")
		return ctrl.Result{}, nil
	}
	// If the status last update time is not highger (not after) the time of last update, there no need to reconcile
	if !r.LastUpdated.IsZero() && !execute.Status.LastUpdated.Time.After(r.LastUpdated) {
		logger.Info("No need to reconcile EXECUTE")
		return ctrl.Result{}, nil
	}
	// If this is the first run of controller without any real requests clear the command clear the status
	if r.LastUpdated.IsZero() {
		execute.Status.Input = []lupusv1.MoveCommand{}
	}
	// If we reconcile then set the controller's last updated time same as in the resource spec
	r.LastUpdated = execute.Status.LastUpdated.Time

	// Extract the input list of MoveCommand from Execute's status
	inputCommands := execute.Status.Input

	// Iterate over each MoveCommand and send an HTTP request
	for _, moveCommand := range inputCommands {
		// Convert MoveCommand to JSON
		jsonData, err := json.Marshal(moveCommand)
		if err != nil {
			logger.Error(err, "Failed to marshal MoveCommand to JSON")
			return ctrl.Result{}, err
		}

		// Send HTTP POST request to 'http://localhost:4141/api/move'
		resp, err := http.Post("http://localhost:4141/api/move", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			logger.Error(err, "Failed to send HTTP request")
			return ctrl.Result{}, err
		}
		defer resp.Body.Close()

		// Log the HTTP response
		if resp.StatusCode != http.StatusOK {
			logger.Error(nil, "Received non-OK response from API", "statusCode", resp.StatusCode)
		} else {
			logger.Info("Successfully sent MoveCommand", "From", moveCommand.From, "To", moveCommand.To, "Count", moveCommand.Count)
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ExecuteReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&lupusv1.Execute{}).
		Complete(r)
}
