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
	"errors"
	"fmt"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	v1 "github.com/0x41gawor/lupus/api/v1"
	"github.com/go-logr/logr"
)

// ObserveReconciler reconciles a Observe object
type ObserveReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	// Static fields of Reconciler
	Logger      logr.Logger
	ElementType string
	// Dynamic fields of Reconciler
	IsAfterDryRun bool
	LastUpdated   time.Time
}

// +kubebuilder:rbac:groups=lupus.gawor.io,resources=observes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=lupus.gawor.io,resources=observes/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=lupus.gawor.io,resources=observes/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Observe object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *ObserveReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// Set up logging context
	r.ElementType = "Observe"
	r.Logger = log.FromContext(ctx)
	r.Logger.Info(fmt.Sprintf("=================== START OF %s Reconciler: \n", strings.ToUpper(r.ElementType)))

	// Step 1 - Fetch the reconciled resource instance
	var element v1.Observe
	if err := r.Get(ctx, req.NamespacedName, &element); err != nil {
		r.Logger.Info("Failed to fetch Observe instance", "error", err)
		// If the resource is not found, we return and don't requeue
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Step 2 - Checks
	if !r.IsAfterDryRun {
		// clear status as it can contain some garbage
		element.Status.Input = runtime.RawExtension{}
		element.Status.LastUpdated = metav1.Time{}
		// set the flag
		r.IsAfterDryRun = true
		r.Logger.Info("This is the dry run, no need to reconcile.")
		return ctrl.Result{}, nil
	}
	// Check for double update in single loop iteration. If r.LastUpdated time is zero it means it is the 2nd run (so double update can't happen)
	// If the Status.LastUpdated time is non-zero we have to check if its not the same as the previous one
	if !r.LastUpdated.IsZero() && !element.Status.LastUpdated.Time.After(r.LastUpdated) {
		// If this condition is true it means we are reconciling again in the same iteration
		r.Logger.Info("Already reconciled in this loop iteration, no need to reconcile")
	}

	// Step 3 - We reconcile, so let's begin the process with variable settings
	var input runtime.RawExtension = element.Status.Input
	r.LastUpdated = time.Now()
	// Step 4 - Unmarshall input into map[string]interface{} so you can easily access its root elements
	data, err := rawExtensionToMap(input)
	if err != nil {
		r.Logger.Error(err, "Cannot unmarshall the input")
		return ctrl.Result{}, nil
	}

	// Step 5 - Send input to the next elements
	// Iterate over the array of next elements
	for _, next := range element.Spec.Next {
		// prepare placeholder for output
		outputMap := make(map[string]interface{})
		if len(next.Tags) == 1 && next.Tags[0] == "*" {
			outputMap = data
		} else {
			for _, tag := range next.Tags {
				if value, exists := data[tag]; exists {
					outputMap[tag] = value
				}
			}
		}

		outputRaw, err := mapToRawExtension(outputMap)
		if err != nil {
			r.Logger.Error(err, "Cannot convert map to RawExtension")
		}
		// fetch the nextElement
		switch next.Type {
		case "Learn":
			resourceName := element.Spec.Master + "-" + next.Name
			resourceNamespace := "default"

			var nextElement v1.Learn
			err := r.Get(ctx, client.ObjectKey{Name: resourceName, Namespace: resourceNamespace}, &nextElement)
			if err != nil {
				r.Logger.Error(err, "Failed to get next element: Learn")
			}
			nextElement.Status.Input = outputRaw
			nextElement.Status.LastUpdated = metav1.Time{Time: r.LastUpdated}
			if err := r.Status().Update(ctx, &nextElement); err != nil {
				r.Logger.Error(err, "Failed to update next element (Learn) status")
			}
		case "Decide":
			resourceName := element.Spec.Master + "-" + next.Name
			resourceNamespace := "default"

			var nextElement v1.Decide
			err := r.Get(ctx, client.ObjectKey{Name: resourceName, Namespace: resourceNamespace}, &nextElement)
			if err != nil {
				r.Logger.Error(err, "Failed to get next element: Decide")
			}
			nextElement.Status.Input = outputRaw
			nextElement.Status.LastUpdated = metav1.Time{Time: r.LastUpdated}
			if err := r.Status().Update(ctx, &nextElement); err != nil {
				r.Logger.Error(err, "Failed to update next element (Decide) status")
			}
		default:
			r.Logger.Error(errors.New("cannot pass input to any other element type than Lean or Decide"), "Unrecognized element type")
		}
	}
	r.Logger.Info("Observe sucessfully reconciled", "name", req.Name)
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ObserveReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1.Observe{}).
		Complete(r)
}
