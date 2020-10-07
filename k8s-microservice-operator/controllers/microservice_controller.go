/*
Copyright 2020 Mario-Leander Reimer.
*/

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	v1 "k8s.io/api/apps/v1"
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
	_ = context.Background()
	_ = r.Log.WithValues("microservice", req.NamespacedName)

	// your logic here

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
