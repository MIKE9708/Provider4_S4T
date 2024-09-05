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

	"github.com/MIKE9708/Provider4S4T.git/api/v1alpha1"
	infrastructurev1alpha1 "github.com/MIKE9708/Provider4S4T.git/api/v1alpha1"
	"github.com/MIKE9708/s4t-sdk-go/pkg/api"
	"github.com/MIKE9708/s4t-sdk-go/pkg/api/data/board"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const boardFinalizer = "board.openstack-fork.example.com/finalizer"

// BoardReconciler reconciles a Board object
type BoardReconciler struct {
	client.Client
	Scheme    *runtime.Scheme
	S4tClient *s4t.Client
}

// +kubebuilder:rbac:groups=infrastructure.s4t.example.com,resources=boards,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.s4t.example.com,resources=boards/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=infrastructure.s4t.example.com,resources=boards/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Board object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *BoardReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	boardCR := &infrastructurev1alpha1.Board{}
	s4t := s4t.Client{}
	s4t_client, err := s4t.GetClientConnection()

	if err != nil {
		return ctrl.Result{}, err
	}

	r.S4tClient = s4t_client

	if err := r.Get(ctx, req.NamespacedName, boardCR); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	// If no UUID is provided means the
	// resource has to be created
	if boardCR.Status.UUID == "" {
		return r.ReconcileCreate(ctx, s4t_client, boardCR)
	}
	// If UUID is provided Board exist so
	// Check if board should be deleted
	if !controllerutil.ContainsFinalizer(boardCR, boardFinalizer) {
		controllerutil.AddFinalizer(boardCR, boardFinalizer)
		if err := r.Update(ctx, boardCR); err != nil {
			return ctrl.Result{}, err
		}
	}

	if !boardCR.DeletionTimestamp.IsZero() {
		return r.ReconcileDelete(ctx, boardCR)
	}

	// Board don't have to be deleted,
	// so the board exist and has an ID
	// test if board UUID exist
	board_exist, err := r.S4tClient.GetBoardDetail(boardCR.Spec.Board.Uuid)
	if err != nil || board_exist.Uuid == "" {
		return ctrl.Result{}, err
	}

	return r.ReconcileUpdate(ctx, boardCR)
}

func (r *BoardReconciler) ReconcileCreate(ctx context.Context, s4t_client *s4t.Client, boardCR *infrastructurev1alpha1.Board) (ctrl.Result, error) {
	board := v1alpha1.BoardData{}
	board.Name = boardCR.Spec.Board.Name
	board.Code = boardCR.Spec.Board.Code
	// board.Location = boardCR.Spec.Board.Location
	board_data := boards.Board{}

	board_data.Name = boardCR.Spec.Board.Name
	board_data.Code = boardCR.Spec.Board.Code

	for index, location := range boardCR.Spec.Board.Location {
		board_data.Location[index].Altitude = location.Altitude
		board_data.Location[index].Latitude = location.Latitude
		board_data.Location[index].Longitude = location.Longitude
	}

	resp, err := s4t_client.CreateBoard(board_data)
	if err != nil {
		return ctrl.Result{}, err
	}

	boardCR.Spec.Board.Uuid = resp.Uuid
	boardCR.Spec.Board.Code = resp.Code
	boardCR.Spec.Board.Status = resp.Status
	boardCR.Spec.Board.Name = resp.Name
	boardCR.Spec.Board.Session = resp.Session
	boardCR.Spec.Board.Wstunip = resp.Wstunip
	boardCR.Spec.Board.Type = resp.Type
	boardCR.Spec.Board.LRversion = resp.LRversion

	boardCR.Status.UUID = resp.Uuid
	boardCR.Status.Status = resp.Status

	if err := r.Update(ctx, boardCR); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func (r *BoardReconciler) ReconcileUpdate(ctx context.Context, boardCR *infrastructurev1alpha1.Board) (ctrl.Result, error) {
	resp, err := r.S4tClient.PatchBoard(boardCR.Status.UUID, map[string]interface{}{"code": boardCR.Spec.Board.Code})
	if err != nil {
		return ctrl.Result{}, err
	}

	boardCR.Spec.Board.Uuid = resp.Uuid
	boardCR.Spec.Board.Code = resp.Code
	boardCR.Spec.Board.Status = resp.Status
	boardCR.Spec.Board.Name = resp.Name
	boardCR.Spec.Board.Session = resp.Session
	boardCR.Spec.Board.Wstunip = resp.Wstunip
	boardCR.Spec.Board.Type = resp.Type
	boardCR.Spec.Board.LRversion = resp.LRversion

	boardCR.Status.UUID = resp.Uuid
	boardCR.Status.Status = resp.Status

	if err := r.Update(ctx, boardCR); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *BoardReconciler) ReconcileDelete(ctx context.Context, boardCR *infrastructurev1alpha1.Board) (ctrl.Result, error) {
	if err := r.S4tClient.DeleteBoard(boardCR.Status.UUID); err != nil {
		return ctrl.Result{}, err
	}

	controllerutil.RemoveFinalizer(boardCR, boardFinalizer)
	if err := r.Update(ctx, boardCR); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil

}

// SetupWithManager sets up the controller with the Manager.
func (r *BoardReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrastructurev1alpha1.Board{}).
		Complete(r)
}
