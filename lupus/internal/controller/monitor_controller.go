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

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	lupusv1 "github.com/0x41gawor/lupus/api/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MonitorReconciler reconciles a Monitor object
type MonitorReconciler struct {
	client.Client
	Scheme *runtime.Scheme

	LastUpdated time.Time
}

// +kubebuilder:rbac:groups=lupus.gawor.io,resources=monitors,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=lupus.gawor.io,resources=monitors/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=lupus.gawor.io,resources=monitors/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Monitor object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *MonitorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// Create a logger from the context
	logger := log.FromContext(ctx)
	logger.Info("=================== START OF MONITOR Reconciler: ")
	println(r.LastUpdated.String())

	// Fetch the Monitor instance
	var monitor lupusv1.Monitor
	if err := r.Get(ctx, req.NamespacedName, &monitor); err != nil {
		if errors.IsNotFound(err) {
			logger.Info("Monitor resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Failed to get Monitor resource.")
		return ctrl.Result{}, err
	}
	// If the status of resource was not set in cluster, there is no need to reconcile
	if monitor.Status.LastUpdated.Time.IsZero() {
		logger.Info("No need to reconcile MONITOR")
		return ctrl.Result{}, nil
	}
	// If the status last update time is not highger (not after) the time of last update, there no need to reconcile
	if !r.LastUpdated.IsZero() && !monitor.Status.LastUpdated.Time.After(r.LastUpdated) {
		logger.Info("No need to reconcile MONITOR")
		return ctrl.Result{}, nil
	}
	// If we reconcile then set the controller's last updated time same as in the resource spec
	r.LastUpdated = monitor.Status.LastUpdated.Time

	// Extract Gdansk, Krakow, Poznan, and Warsaw values from the Monitor's status
	gdansk := monitor.Status.Gdansk
	krakow := monitor.Status.Krakow
	poznan := monitor.Status.Poznan
	warsaw := monitor.Status.Warsaw

	// Fetch the Decision resource with name "piotrek" in the "default" namespace
	decisionName := "piotrek"
	decisionNamespace := "default"

	var decision lupusv1.Decision
	err := r.Get(ctx, client.ObjectKey{Name: decisionName, Namespace: decisionNamespace}, &decision)
	if err != nil && errors.IsNotFound(err) {
		// Decision resource not found, so let's create a new one
		decision = lupusv1.Decision{
			ObjectMeta: metav1.ObjectMeta{
				Name:      decisionName,
				Namespace: decisionNamespace,
			},
			Spec: lupusv1.DecisionSpec{
				// Add any relevant Spec fields here if required
			},
		}
	} else if err != nil {
		// Handle errors other than NotFound
		logger.Error(err, "Failed to get Decision resource.")
		return ctrl.Result{}, err
	}

	// Set the fields in the Decision resource's status
	decision.Status.Input = lupusv1.Input{
		Gdansk: gdansk,
		Krakow: krakow,
		Poznan: poznan,
		Warsaw: warsaw,
	}
	decision.Status.LastUpdated = monitor.Status.LastUpdated

	// Update or create the Decision resource
	if err == nil && errors.IsNotFound(err) {
		logger.Error(err, "Failed to create new Decision resource.")
		return ctrl.Result{}, err
	} else {
		// Update the existing Decision resource
		if err := r.Status().Update(ctx, &decision); err != nil {
			logger.Error(err, "Failed to update Decision status.")
			return ctrl.Result{}, err
		}
	}

	logger.Info("MONITOR: Successfully reconciled Monitor and updated Decision resource.")
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MonitorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&lupusv1.Monitor{}).
		Complete(r)
}
