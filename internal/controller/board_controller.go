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
	"github.com/MIKE9708/s4t-sdk-go/pkg"
	"github.com/MIKE9708/s4t-sdk-go/pkg/api/boards"
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
	Scheme *runtime.Scheme
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
	
	if err := r.Get(ctx, req.NamespacedName, boardCR); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	} 
	// If no UUID is provided means the 
	// resource has to be created	
	if boardCR.Status.UUID == "" {
		return r.ReconcileCreate(ctx, boardCR)		
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
	board_exist, err := boardCR.Spec.Board.GetBoardDetail(r.S4tClient, boardCR.Spec.Board.Uuid)
	if (err != nil || board_exist.Uuid == "") {
		return ctrl.Result{}, err
	}
	
	return r.ReconcileUpdate(ctx, boardCR)
}

func (r *BoardReconciler) ReconcileCreate(ctx context.Context, boardCR *infrastructurev1alpha1.Board) (ctrl.Result, error) {
	board := boards.Board{}
	board.Name = boardCR.Spec.Board.Name
	board.Code = boardCR.Spec.Board.Code
	board.Location = boardCR.Spec.Board.Location
	
	resp,err:= board.CreateBoard(r.S4tClient, board)
	if err != nil {
		return ctrl.Result{}, err
	}

	boardCR.Spec.Board = *resp 
	boardCR.Status.UUID = resp.Uuid
	boardCR.Status.Status = resp.Status

	if err := r.Update(ctx, boardCR); err != nil{
		return ctrl.Result{}, err
	} 
	return ctrl.Result{}, nil
}

func (r *BoardReconciler) ReconcileUpdate(ctx context.Context, boardCR *infrastructurev1alpha1.Board) (ctrl.Result, error) {
	board := &boards.Board{Uuid: boardCR.Status.UUID}	
	resp,err:= board.PatchBoard(r.S4tClient, board.Uuid, map[string]interface{}{"code":boardCR.Spec.Board.Code})
	if err != nil {
		return ctrl.Result{}, err
	}

	boardCR.Spec.Board = *resp 

	if err := r.Update(ctx, boardCR); err != nil{
		return ctrl.Result{}, err
	} 

	return ctrl.Result{}, nil
}

func (r *BoardReconciler) ReconcileDelete(ctx context.Context, boardCR *infrastructurev1alpha1.Board) (ctrl.Result, error) {
	board := &boards.Board{Uuid: boardCR.Status.UUID}	
	if err:=board.DeleteBoard(r.S4tClient, board.Uuid); err != nil {
		return ctrl.Result{}, err
	}
	
	controllerutil.RemoveFinalizer(boardCR, boardFinalizer)
	if err := r.Update(ctx, boardCR); err != nil{
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
