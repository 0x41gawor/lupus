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
	"time"

	"k8s.io/apimachinery/pkg/runtime"
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
	// Step 1: Fetch the Observe instance
	var decide lupusv1.Decide
	if err := r.Get(ctx, req.NamespacedName, &decide); err != nil {
		logger.Info("Failed to fetch Decide instance", "error", err)
		// If the resource is not found, we return and don't requeue
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	logger.Info("Decide instance fetched successfully")
	// Step 2: Check if status of Decide is empty, if yes it is the first run and reconciler shoudl return
	if decide.Status.LastUpdated.IsZero() {
		logger.Info("No need to reconcile")
		return ctrl.Result{}, nil
	}

	// If the status last update time is not highger (not after) the time of last update, there no need to reconcile
	if !r.LastUpdated.IsZero() && !decide.Status.LastUpdated.Time.After(r.LastUpdated) {
		logger.Info("No need to reconcile DECISION")
		return ctrl.Result{}, nil
	}
	// If this is the first run of controller without any real requests clear the command clear the status
	if r.LastUpdated.IsZero() {
		decide.Status.Input = runtime.RawExtension{}
	}
	// If we reconcile then set the controller's last updated time same as in the resource spec
	r.LastUpdated = decide.Status.LastUpdated.Time

	// Step 3: Pass data into OPA
	input := decide.Status.Input
	// here the call to OPA will occur
	output, err := opaSimulation(input)
	if err != nil {
		logger.Error(err, "Failed to distribute the load")
	}

	// Step 4: Fetch nextElement resource
	resourceName := decide.Spec.NextElement
	resourceNamespace := "default"

	var execute lupusv1.Execute
	err = r.Get(ctx, client.ObjectKey{Name: resourceName, Namespace: resourceNamespace}, &execute)
	if err != nil {
		logger.Error(err, "Failed to get Execute resource")
		return ctrl.Result{}, nil
	}
	// Step 5: Set the fields in nextElement resource
	execute.Status.Input = output
	execute.Status.LastUpdated = decide.Status.LastUpdated

	// Step 6: Update the nextElement resource status
	if err := r.Status().Update(ctx, &execute); err != nil {
		logger.Error(err, "Failed to update Execute status")

	}
	// Step 7: Next time reconcile when the status changes
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DecideReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&lupusv1.Decide{}).
		Complete(r)
}
