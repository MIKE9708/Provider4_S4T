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
	s4t "github.com/MIKE9708/s4t-sdk-go/pkg"
	"github.com/MIKE9708/s4t-sdk-go/pkg/api/plugins"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// PluginReconciler reconciles a Plugin object
type PluginReconciler struct {
	client.Client
	Scheme    *runtime.Scheme
	S4tClient *s4t.Client
}

// +kubebuilder:rbac:groups=infrastructure.s4t.example.com,resources=plugins,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.s4t.example.com,resources=plugins/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=infrastructure.s4t.example.com,resources=plugins/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Plugin object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *PluginReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	pluginCR := &infrastructurev1alpha1.Plugin{}

	if err := r.Get(ctx, req.NamespacedName, pluginCR); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)

	}

	if !controllerutil.ContainsFinalizer(pluginCR, boardFinalizer) {
		controllerutil.AddFinalizer(pluginCR, boardFinalizer)
		if err := r.Update(ctx, pluginCR); err != nil {
			return ctrl.Result{}, err
		}
	}

	if !pluginCR.DeletionTimestamp.IsZero() {
		return r.ReconcileDelete(ctx, pluginCR)
	}
	plugin := &plugins.PluginReq{
		Name:       pluginCR.Spec.Plugin.Name,
		Parameters: pluginCR.Spec.Plugin.Parameters,
		Code:       pluginCR.Spec.Plugin.Code,
		Version:    pluginCR.Spec.Plugin.Version,
	}
	plugin_data := &plugins.Plugin{}
	plugin_created, err := plugin_data.CreatePlugin(r.S4tClient, *plugin)

	if err != nil {
		return ctrl.Result{}, err
	}

	// _ = log.FromContext(ctx)
	pluginCR.Status.Name = plugin_created.Name
	pluginCR.Status.Code = plugin_created.UUID

	if err := r.Status().Update(ctx, pluginCR); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *PluginReconciler) ReconcileDelete(ctx context.Context, pluginCR *infrastructurev1alpha1.Plugin) (ctrl.Result, error) {
	plugin := &plugins.Plugin{UUID: pluginCR.Status.Code}
	if err := plugin.DeletePlugin(r.S4tClient); err != nil {
		return ctrl.Result{}, err
	}

	controllerutil.RemoveFinalizer(pluginCR, boardFinalizer)

	if err := r.Update(ctx, pluginCR); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil

}

// SetupWithManager sets up the controller with the Manager.
func (r *PluginReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrastructurev1alpha1.Plugin{}).
		Complete(r)
}
