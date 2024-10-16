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

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	lupusv1 "github.com/0x41gawor/lupus/api/v1"
)

// LoopReconciler reconciles a Loop object
type LoopReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=lupus.gawor.io,resources=loops,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=lupus.gawor.io,resources=loops/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=lupus.gawor.io,resources=loops/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Loop object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *LoopReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("Starting Loop reconciliation")

	// Step 1: Fetch the Loop instance
	var loop lupusv1.Loop
	if err := r.Get(ctx, req.NamespacedName, &loop); err != nil {
		logger.Error(err, "Failed to fetch Loop instance")
		// If the resource is not found, return without requeuing
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	logger.Info("Loop instance fetched successfully")

	if loop.Status.IsActive {
		logger.Info("No need to reconcile")
		return ctrl.Result{}, nil
	}

	// Step 2: Iterate over the elements array and spawn them
	for i, element := range loop.Spec.Elements {
		// Step 3: Construct names for current element and next element
		elementName := fmt.Sprintf("%s-%s", element.Name, loop.Spec.Name)

		// If there is a next element, construct its name, otherwise it's empty
		nextElementName := ""
		if i+1 < len(loop.Spec.Elements) {
			nextElementName = fmt.Sprintf("%s-%s", loop.Spec.Elements[i+1].Name, loop.Spec.Name)
		}

		// Step 4: Create resources based on the element kind (Observe, Decide, Execute)
		switch element.Name {
		case "observe":
			// Create the Observe resource
			observe := &lupusv1.Observe{
				ObjectMeta: metav1.ObjectMeta{
					Name:      elementName,
					Namespace: req.Namespace,
				},
				Spec: lupusv1.ObserveSpec{
					Name:                  elementName,
					NextElement:           nextElementName,
					Url:                   element.Url,
					ReconcileTimeInterval: 60, // Add this only for Observe resources
				},
			}
			if err := r.Create(ctx, observe); err != nil {
				logger.Error(err, "Failed to create Observe resource", "name", elementName)
				return ctrl.Result{}, err
			}
			logger.Info("Observe resource created successfully", "name", elementName)

		case "decide":
			// Create the Decide resource
			decide := &lupusv1.Decide{
				ObjectMeta: metav1.ObjectMeta{
					Name:      elementName,
					Namespace: req.Namespace,
				},
				Spec: lupusv1.DecideSpec{
					Name:        elementName,
					NextElement: nextElementName,
					Url:         element.Url,
				},
			}
			if err := r.Create(ctx, decide); err != nil {
				logger.Error(err, "Failed to create Decide resource", "name", elementName)
				return ctrl.Result{}, err
			}
			logger.Info("Decide resource created successfully", "name", elementName)

		case "execute":
			// Create the Execute resource
			execute := &lupusv1.Execute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      elementName,
					Namespace: req.Namespace,
				},
				Spec: lupusv1.ExecuteSpec{
					Name:        elementName,
					NextElement: nextElementName,
					Url:         element.Url,
				},
			}
			if err := r.Create(ctx, execute); err != nil {
				logger.Error(err, "Failed to create Execute resource", "name", elementName)
				return ctrl.Result{}, err
			}
			logger.Info("Execute resource created successfully", "name", elementName)

		default:
			logger.Info("Unknown element type", "name", element.Name)
		}
	}

	logger.Info("Loop reconciliation finished successfully")
	loop.Status.IsActive = true

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *LoopReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&lupusv1.Loop{}).
		Complete(r)
}
