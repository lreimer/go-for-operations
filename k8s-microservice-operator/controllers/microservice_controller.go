/*
Copyright 2020 Mario-Leander Reimer.
*/

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	appsv1 "github.com/lreimer/go-for-operations/k8s-microservice-operator/api/v1"
)

// MicroserviceReconciler reconciles a Microservice object
type MicroserviceReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=apps.qaware.de,resources=microservices,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps.qaware.de,resources=microservices/status,verbs=get;update;patch

// Reconcile loop to apply relevant changes to K8s
func (r *MicroserviceReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	logger := r.Log.WithValues("microservice", req.NamespacedName)

	// lookup the Microservice instance for this reconcile request
	microservice := &appsv1.Microservice{}
	err := r.Get(ctx, req.NamespacedName, microservice)
	if err != nil {
		if errors.IsNotFound(err) {
			logger.Info("Microservice resource not found. Deleting ...")
			// delete all associated resources if required
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Failed to get Microservice.")
		return ctrl.Result{}, err
	}

	logger.Info("Reconcile Microservice.")
	// add the update the associated service, deployment, ...

	return ctrl.Result{}, nil
}

// SetupWithManager watch Microservice and Deployment resources
func (r *MicroserviceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.Microservice{}).
		Owns(&v1.Deployment{}).
		WithOptions(controller.Options{
			MaxConcurrentReconciles: 2,
		}).
		Complete(r)
}
