/*
Copyright 2020 Mario-Leander Reimer.
*/

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	v1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
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
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Failed to get Microservice.")
		return ctrl.Result{}, err
	}

	logger.Info("Reconcile Microservice.")
	deployment := &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      microservice.Name,
			Namespace: microservice.Namespace,
			Labels:    microservice.Labels,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: microservice.APIVersion,
					Kind:       microservice.Kind,
					Name:       microservice.Name,
					UID:        microservice.UID,
				},
			},
		},
		Spec: v1.DeploymentSpec{
			Replicas: &microservice.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: microservice.Labels,
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: microservice.Labels,
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "microservice",
							Image: microservice.Spec.Image,
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: microservice.Spec.Ports[0],
								},
							},
						},
					},
				},
			},
		},
	}

	error := r.Client.Create(context.TODO(), deployment, &client.CreateOptions{})
	if error != nil {
		logger.Error(nil, error.Error())
	}

	service := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      microservice.Name,
			Namespace: microservice.Namespace,
			Labels:    microservice.Labels,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: microservice.APIVersion,
					Kind:       microservice.Kind,
					Name:       microservice.Name,
					UID:        microservice.UID,
				},
			},
		},
		Spec: apiv1.ServiceSpec{
			Selector: microservice.Labels,
			Ports: []apiv1.ServicePort{
				{
					Name:       "http",
					Port:       microservice.Spec.Ports[0],
					TargetPort: intstr.FromString("http"),
				},
			},
			Type: apiv1.ServiceTypeLoadBalancer,
		},
	}

	error = r.Client.Create(context.TODO(), service, &client.CreateOptions{})
	if error != nil {
		logger.Error(nil, error.Error())
	}

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
