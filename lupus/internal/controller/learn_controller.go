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
	"context"
	"fmt"
	"os"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	v1 "github.com/0x41gawor/lupus/api/v1"
	"github.com/go-logr/logr"
)

// LearnReconciler reconciles a Learn object
type LearnReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	// Static fields of Reconciler
	Logger      logr.Logger
	ElementType string
	// Map that holds of reconciler state for each element
	instanceState map[types.NamespacedName]*ElementInstanceState
}

// +kubebuilder:rbac:groups=lupus.gawor.io,resources=learns,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=lupus.gawor.io,resources=learns/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=lupus.gawor.io,resources=learns/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Learn object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *LearnReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// Set up logging context
	r.ElementType = "Learn"
	r.Logger = log.FromContext(ctx)
	r.Logger.Info(fmt.Sprintf("=================== START OF %s Reconciler: \n", strings.ToUpper(r.ElementType)))
	// Step 1 - Fetch the reconciled resource instance (Controller-Runtime nomenclature)
	// Step 1 - Fetch reconciled element 	(Lupus nomenclature)
	var element v1.Learn
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
		r.instanceState[req.NamespacedName] = &ElementInstanceState{}
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

	// Step 4 Send input to destination
	switch element.Spec.Destination.Type {
	case "FILE":
		// Extract the file path from the spec
		filePath := element.Spec.Destination.File.Path
		// Convert input to a JSON string
		inputJSON, err := rawExtensionToString(input)
		if err != nil {
			r.Logger.Error(err, "Failed conversion from rawExtension to JSON string")
			return ctrl.Result{}, nil
		}
		// Append JSON data to the specified file
		err = r.appendToFile(filePath, inputJSON, req)
		if err != nil {
			r.Logger.Error(err, "Failed to append to a file")
			return ctrl.Result{}, nil
		}

	default:
		r.Logger.Info(fmt.Sprintf("Destination %s not yet implemented in Learn", element.Spec.Destination.Type))
	}

	return ctrl.Result{}, nil
}

// Helper function to append JSON data to a file
func (r *LearnReconciler) appendToFile(filePath, data string, req ctrl.Request) error {
	// Open the file in append mode, create it if it doesn't exist
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()
	timeString := r.instanceState[req.NamespacedName].LastUpdated.Format("2006/01/02 15:04:05")
	// Write the JSON data followed by a newline
	if _, err := file.WriteString(timeString + " " + data + "\n"); err != nil {
		return fmt.Errorf("failed to write data to file: %v", err)
	}

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *LearnReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// Initialize instanceState map if it's nil
	if r.instanceState == nil {
		r.instanceState = make(map[types.NamespacedName]*ElementInstanceState)
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1.Learn{}).
		Complete(r)
}
