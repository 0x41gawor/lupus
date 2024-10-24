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
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	lupusv1 "github.com/0x41gawor/lupus/api/v1"
	"github.com/go-logr/logr"
)

// LearnReconciler reconciles a Learn object
type LearnReconciler struct {
	client.Client
	Scheme *runtime.Scheme

	Logger      logr.Logger
	ElementType string
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
	r.ElementType = "Learn"
	r.Logger = log.FromContext(ctx)
	r.Logger.Info(fmt.Sprintf("=================== START OF %s Reconciler: \n", strings.ToUpper(r.ElementType)))
	// TODO(user): your logic here

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *LearnReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&lupusv1.Learn{}).
		Complete(r)
}
