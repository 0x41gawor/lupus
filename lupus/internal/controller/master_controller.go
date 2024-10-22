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

	NamePrefix string
	Logger     logr.Logger
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

	// Fetch the Master instance
	var master v1.Master
	if err := r.Get(ctx, req.NamespacedName, &master); err != nil {
		if errors.IsNotFound(err) {
			r.Logger.Info("Master resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		r.Logger.Error(err, "Failed to get Master")
		return ctrl.Result{}, err
	}

	if master.Status.IsActive {
		r.Logger.Info("Status is active no need to reconcile")
		return ctrl.Result{}, nil
	}

	r.NamePrefix = master.Spec.Name
	// Loop through each element in the Master spec
	for _, element := range master.Spec.Elements {
		// Check the element's type and create the corresponding resource
		switch element.Type {
		case "Observe":
			err := r.createObserveResource(ctx, element, master.Namespace, &master)
			if err != nil {
				r.Logger.Error(err, "Failed to create/update Observe resource", "ElementName", element.Name)
				return ctrl.Result{}, err
			}
		case "Decide":
			err := r.createDecideResource(ctx, element, master.Namespace, &master)
			if err != nil {
				r.Logger.Error(err, "Failed to create/update Decide resource", "ElementName", element.Name)
				return ctrl.Result{}, err
			}
		case "Learn":
			err := r.createLearnResource(ctx, element, master.Namespace, &master)
			if err != nil {
				r.Logger.Error(err, "Failed to create/update Learn resource", "ElementName", element.Name)
				return ctrl.Result{}, err
			}
		case "Execute":
			err := r.createExecuteResource(ctx, element, master.Namespace, &master)
			if err != nil {
				r.Logger.Error(err, "Failed to create/update Execute resource", "ElementName", element.Name)
				return ctrl.Result{}, err
			}
		default:
			r.Logger.Info("Unknown element type, skipping", "ElementName", element.Name)
		}
	}
	master.Status.IsActive = true
	if err := r.Status().Update(ctx, &master); err != nil {
		r.Logger.Error(err, "Failed to update Master resource")
		return ctrl.Result{}, nil
	}
	// Reconciliation successful
	return ctrl.Result{}, nil
}

func (r *MasterReconciler) createObserveResource(ctx context.Context, element v1.Element, namespace string, master *v1.Master) error {
	// Define the desired Observe custom resource
	observe := &v1.Observe{
		ObjectMeta: metav1.ObjectMeta{
			Name:      r.NamePrefix + "-" + element.Name,
			Namespace: namespace,
		},
		Spec: *element.ObserveSpec,
	}
	// Set Master as owner of Observe
	if err := setOwnerReference(master, observe, r.Scheme); err != nil {
		r.Logger.Error(err, "Failed to set owner reference for Observe")
		return err
	}

	// Simply create the resource without checking if it exists
	return r.Create(ctx, observe)
}

func (r *MasterReconciler) createDecideResource(ctx context.Context, element v1.Element, namespace string, master *v1.Master) error {
	// Define the desired Observe custom resource
	decide := &v1.Decide{
		ObjectMeta: metav1.ObjectMeta{
			Name:      r.NamePrefix + "-" + element.Name,
			Namespace: namespace,
		},
		Spec: *element.DecideSpec,
	}
	// Set Master as owner of Decide
	if err := setOwnerReference(master, decide, r.Scheme); err != nil {
		r.Logger.Error(err, "Failed to set owner reference for Decide")
		return err
	}

	// Simply create the resource without checking if it exists
	return r.Create(ctx, decide)
}

func (r *MasterReconciler) createLearnResource(ctx context.Context, element v1.Element, namespace string, master *v1.Master) error {
	// Define the desired Observe custom resource
	learn := &v1.Learn{
		ObjectMeta: metav1.ObjectMeta{
			Name:      r.NamePrefix + "-" + element.Name,
			Namespace: namespace,
		},
		Spec: *element.LearnSpec,
	}

	// Set Master as owner of Learn
	if err := setOwnerReference(master, learn, r.Scheme); err != nil {
		r.Logger.Error(err, "Failed to set owner reference for Learn")
		return err
	}

	// Simply create the resource without checking if it exists
	return r.Create(ctx, learn)
}

func (r *MasterReconciler) createExecuteResource(ctx context.Context, element v1.Element, namespace string, master *v1.Master) error {
	// Define the desired Observe custom resource
	execute := &v1.Execute{
		ObjectMeta: metav1.ObjectMeta{
			Name:      r.NamePrefix + "-" + element.Name,
			Namespace: namespace,
		},
		Spec: *element.ExecuteSpec,
	}

	// Set Master as owner of Execute
	if err := setOwnerReference(master, execute, r.Scheme); err != nil {
		r.Logger.Error(err, "Failed to set owner reference for Execute")
		return err
	}

	// Simply create the resource without checking if it exists
	return r.Create(ctx, execute)
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
