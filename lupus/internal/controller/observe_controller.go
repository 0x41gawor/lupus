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
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	lupusv1 "github.com/0x41gawor/lupus/api/v1"
)

// ObserveReconciler reconciles a Observe object
type ObserveReconciler struct {
	client.Client
	Scheme *runtime.Scheme
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

// Reconcile is the main logic for the controller
func (r *ObserveReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// Create a logger from the context
	logger := log.FromContext(ctx)
	logger.Info("=================== START OF OBSERVE Reconciler: \n")
	// Step 1: Fetch the Observe instance
	var observe lupusv1.Observe
	if err := r.Get(ctx, req.NamespacedName, &observe); err != nil {
		logger.Info("Failed to fetch Observe instance", "error", err)
		// If the resource is not found, we return and don't requeue
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	logger.Info("Observe instance fetched successfully")

	// Step 2: Fetch data from the URL in spec.monitored_system_url
	url := observe.Spec.MonitoredSystemURL.Path
	logger.Info(fmt.Sprintf("Fetching data from: %s", url))
	resp, err := http.Get(url)
	if err != nil {
		logger.Info("Error fetching data", "err", err)
		return ctrl.Result{}, err
	}
	defer resp.Body.Close()

	// Step 3: Read and parse the response body into map[string]interface{}
	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		logger.Info("Error decoding response", "body", err)
		return ctrl.Result{}, err
	}
	logger.Info(fmt.Sprintf("Data fetched successfully: %+v", data), "data", data)

	// Step 4: Log success and requeue after 1 minute
	logger.Info(fmt.Sprintf("Reconciliation success. Requeue after %d seconds", observe.Spec.ReconcileTimeInterval))
	return ctrl.Result{RequeueAfter: time.Second * time.Duration(observe.Spec.ReconcileTimeInterval)}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ObserveReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&lupusv1.Observe{}).
		Complete(r)
}
