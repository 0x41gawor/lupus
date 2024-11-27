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

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	v1 "github.com/0x41gawor/lupus/api/v1"
	"github.com/go-logr/logr"
)

// MasterReconciler reconciles a Master object
type MasterReconciler struct {
	client.Client
	Scheme *runtime.Scheme

	// Prefix that will be added to the k8s object names of loop elements, it's simply the name of Loop
	NamePrefix string
	// Util
	Logger logr.Logger
}

// +kubebuilder:rbac:groups=lupus.gawor.io,resources=masters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=lupus.gawor.io,resources=masters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=lupus.gawor.io,resources=masters/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Master object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *MasterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// Set up logging context
	r.Logger = log.FromContext(ctx)
	r.Logger.Info("Reconciling Master resource")

	// Step 1 - Fetch the reconciled resource instance
	var master v1.Master
	if err := r.Get(ctx, req.NamespacedName, &master); err != nil {
		if errors.IsNotFound(err) {
			r.Logger.Info("Master resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		r.Logger.Error(err, "Failed to get Master")
		return ctrl.Result{}, err
	}

	// Step 2 - Checks
	if master.Status.IsActive {
		r.Logger.Info("Status is active no need to reconcile")
		return ctrl.Result{}, nil
	}
	// Step 3 - We reconcile so let's create loop elements
	r.NamePrefix = master.Spec.Name
	// Loop through each element in the Master spec
	for _, element := range master.Spec.Elements {
		// Check the element's type and create the corresponding resource
		err := r.createElementResource(ctx, *element, master.Namespace, &master)
		if err != nil {
			r.Logger.Error(err, "Failed to create/update Element resource", "ElementName", element.Name)
			return ctrl.Result{}, err
		}
	}
	//Set Master status as active
	master.Status.IsActive = true
	// Step 4 - Update the reconciled resource instance
	if err := r.Status().Update(ctx, &master); err != nil {
		r.Logger.Error(err, "Failed to update Master resource")
		return ctrl.Result{}, nil
	}
	return ctrl.Result{}, nil
}

func (r *MasterReconciler) createElementResource(ctx context.Context, es v1.ElementSpec, namespace string, master *v1.Master) error {
	// Define the desired Observe custom resource
	e := &v1.Element{
		ObjectMeta: metav1.ObjectMeta{
			Name:      r.NamePrefix + "-" + es.Name,
			Namespace: namespace,
		},
		Spec: es,
	}
	e.Spec.Master = r.NamePrefix
	// Set Master as owner of Observe
	if err := setOwnerReference(master, e, r.Scheme); err != nil {
		r.Logger.Error(err, "Failed to set owner reference for Observe")
		return err
	}

	// Simply create the resource without checking if it exists
	return r.Create(ctx, e)
}

// Helper function to set owner reference
func setOwnerReference(owner, object metav1.Object, scheme *runtime.Scheme) error {
	return ctrl.SetControllerReference(owner, object, scheme)
}

// SetupWithManager sets up the controller with the Manager.
func (r *MasterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1.Master{}).
		Complete(r)
}
