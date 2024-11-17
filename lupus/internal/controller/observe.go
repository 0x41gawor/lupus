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
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	v1 "github.com/0x41gawor/lupus/api/v1"
	util "github.com/0x41gawor/lupus/internal/util"
	"github.com/go-logr/logr"
)

// ObserveReconciler reconciles a Observe object
type ObserveReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	// Static fields of Reconciler
	Logger      logr.Logger
	ElementType string
	// Map that holds of reconciler state for each element
	instanceState map[types.NamespacedName]*util.ElementInstanceState
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

	// Step 1 - (k8s) Fetch reconciled resource instance
	// Step 1 - (lupus) Fetch element
	var element v1.Observe
	if err := r.Get(ctx, req.NamespacedName, &element); err != nil {
		r.Logger.Info("Failed to fetch Observe instance", "error", err)
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
	r.instanceState[req.NamespacedName].LastUpdated = time.Now()

	// Step 4 - (Go) Unmarshall input into map[string]interface{}
	// Step 4 - (Lupus) Create Data object
	data, err := util.NewData(input)
	if err != nil {
		r.Logger.Error(err, "Cannot unmarshall the input")
		return ctrl.Result{}, nil
	}

	// Step 5 - (Lupus) Send input to the next elements
	// Step 5 - (K8s) Update Status.Input field of next element
	for _, next := range element.Spec.Next {
		// prepare output based on keys
		output, err := data.GetKeys(next.Keys)
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
	r.Logger.Info("Observe sucessfully reconciled", "name", req.Name)
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ObserveReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// Initialize instanceState map if it's nil
	if r.instanceState == nil {
		r.instanceState = make(map[types.NamespacedName]*util.ElementInstanceState)
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1.Observe{}).
		Complete(r)
}

func (r *ObserveReconciler) updateStatus(ctx context.Context, t string, objKey client.ObjectKey, updateTime metav1.Time, output runtime.RawExtension) error {
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
	default:
		err := errors.New("cannot pass input to any other element type than Decide")
		r.Logger.Error(err, "Unrecognized element type")
		return err
	}
	return nil
}
