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
	infrastructurev1alpha1 "github.com/MIKE9708/Provider4S4T.git/api/v1alpha1"
	"github.com/MIKE9708/s4t-sdk-go/pkg/api"
	"github.com/MIKE9708/s4t-sdk-go/pkg/api/data/service"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// ServiceReconciler reconciles a Service object
type ServiceReconciler struct {
	client.Client
	Scheme    *runtime.Scheme
	S4tClient *s4t.Client
}

// +kubebuilder:rbac:groups=infrastructure.s4t.example.com,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.s4t.example.com,resources=services/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=infrastructure.s4t.example.com,resources=services/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Service object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *ServiceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	serviceCR := &infrastructurev1alpha1.Service{}

	if err := r.Get(ctx, req.NamespacedName, serviceCR); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)

	}

	if !controllerutil.ContainsFinalizer(serviceCR, boardFinalizer) {
		controllerutil.AddFinalizer(serviceCR, boardFinalizer)
		if err := r.Update(ctx, serviceCR); err != nil {
			return ctrl.Result{}, err
		}
	}

	if !serviceCR.DeletionTimestamp.IsZero() {
		return r.ReconcileDelete(ctx, serviceCR)
	}

	service_data := services.Service{
		Name:     serviceCR.Spec.Service.Name,
		Port:     serviceCR.Spec.Service.Port,
		Protocol: serviceCR.Spec.Service.Protocol,
	}

	service_created, err := r.S4tClient.CreateService(service_data)

	if err != nil {
		return ctrl.Result{}, err
	}
	serviceCR.Status.Name = service_created.Name
	serviceCR.Status.Port = service_created.Port

	// TODO(user): your logic here

	return ctrl.Result{}, nil
}

func (r *ServiceReconciler) ReconcileDelete(ctx context.Context, serviceCR *infrastructurev1alpha1.Service) (ctrl.Result, error) {
	if err := r.S4tClient.DeleteService(serviceCR.Status.UUID); err != nil {
		return ctrl.Result{}, err
	}

	controllerutil.RemoveFinalizer(serviceCR, boardFinalizer)
	if err := r.Update(ctx, serviceCR); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrastructurev1alpha1.Service{}).
		Complete(r)
}
