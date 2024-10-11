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
	"math"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	lupusv1 "github.com/0x41gawor/lupus/api/v1"
)

// DecisionReconciler reconciles a Decision object
type DecisionReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// distributeLoad takes the Input struct and generates a list of MoveCommand
func distributeLoad(input lupusv1.Input) []lupusv1.MoveCommand {
	// List of MoveCommand that will store the result
	var moves []lupusv1.MoveCommand

	// Map to easily reference the cities
	cities := map[string]*int{
		"Gdansk": &input.Gdansk,
		"Krakow": &input.Krakow,
		"Poznan": &input.Poznan,
		"Warsaw": &input.Warsaw,
	}

	// Calculate the total load and target load per city
	totalLoad := input.Gdansk + input.Krakow + input.Poznan + input.Warsaw
	targetLoad := totalLoad / len(cities)

	// Collect cities with excess and deficit load
	excessCities := make(map[string]int)
	deficitCities := make(map[string]int)

	for city, load := range cities {
		if *load > targetLoad {
			excessCities[city] = *load - targetLoad
		} else if *load < targetLoad {
			deficitCities[city] = targetLoad - *load
		}
	}

	// Process moves from excess cities to deficit cities
	for fromCity, excess := range excessCities {
		for toCity, deficit := range deficitCities {
			if excess == 0 {
				break
			}
			// Calculate the amount to move
			moveCount := int(math.Min(float64(excess), float64(deficit)))
			if moveCount > 0 {
				// Create move command
				moves = append(moves, lupusv1.MoveCommand{
					From:  fromCity,
					To:    toCity,
					Count: moveCount,
				})

				// Update the loads
				excessCities[fromCity] -= moveCount
				deficitCities[toCity] -= moveCount

				// Update cities' actual load
				*cities[fromCity] -= moveCount
				*cities[toCity] += moveCount
			}
		}
	}

	return moves
}

// +kubebuilder:rbac:groups=lupus.gawor.io,resources=decisions,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=lupus.gawor.io,resources=decisions/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=lupus.gawor.io,resources=decisions/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Decision object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *DecisionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Fetch the Decision instance
	var decision lupusv1.Decision
	if err := r.Get(ctx, req.NamespacedName, &decision); err != nil {
		if errors.IsNotFound(err) {
			// Decision resource not found, nothing to reconcile
			logger.Info("Decision resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Failed to get Decision")
		return ctrl.Result{}, err
	}

	// Extract the Input field from Decision's status
	input := decision.Status.Input

	// Pass the Input to the distributeLoad function, which returns a list of MoveCommands
	moveCommands := distributeLoad(input)

	// Update the status of the Decision resource with the generated MoveCommands
	decision.Status.Decision = moveCommands

	// Update the status of the Decision in the cluster
	if err := r.Status().Update(ctx, &decision); err != nil {
		logger.Error(err, "Failed to update Decision status")
		return ctrl.Result{}, err
	}

	logger.Info("Successfully updated Decision with MoveCommands", "MoveCommands", moveCommands)

	// Fetch Execute resource with name "krzysiek" in the "default" namespace
	executeName := "krzysiek"
	executeNamespace := "default"

	var execute lupusv1.Execute
	err := r.Get(ctx, client.ObjectKey{Name: executeName, Namespace: executeNamespace}, &execute)

	if err != nil {
		logger.Error(err, "Failed to get Execute resource.")
		return ctrl.Result{}, err
	}

	execute.Status.Input = moveCommands

	err = r.Status().Update(ctx, &execute)
	if err != nil {
		logger.Error(err, "Failed to update Execute status.")
		return ctrl.Result{}, err
	}

	logger.Info("Successfully updated Execute with MoveCommands", "MoveCommands", moveCommands)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DecisionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&lupusv1.Decision{}).
		Complete(r)
}
