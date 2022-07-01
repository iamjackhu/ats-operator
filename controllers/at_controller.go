/*
Copyright 2022 Jack Hu.

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

package controllers

import (
	"context"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	cnatv1 "ats-operator.programming-kubernetes.info/api/v1"
)

// AtReconciler reconciles a At object
type AtReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=cnat.programming-kubernetes.info,resources=ats,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cnat.programming-kubernetes.info,resources=ats/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cnat.programming-kubernetes.info,resources=ats/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the At object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.1/pkg/reconcile
func (r *AtReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := log.Log.WithValues("namespace", req.Namespace, "at", req.Name)

	reqLogger.Info("==== Reconciling At")

	instance := &cnatv1.At{}
	err := r.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	reqLogger.Info("Get At", "instance", instance)

	if instance.Status.Phase == "" {
		instance.Status.Phase = cnatv1.PhasePending
	}

	switch instance.Status.Phase {
	case cnatv1.PhasePending:
		reqLogger.Info("Phase PENDING")
		time.Sleep(time.Second * 3)
		instance.Status.Phase = cnatv1.PhaseRunning
	case cnatv1.PhaseRunning:
		reqLogger.Info("Phase RUNNING")
		time.Sleep(time.Second * 3)
		instance.Status.Phase = cnatv1.PhaseDone
	case cnatv1.PhaseDone:
		reqLogger.Info("Phase DONE")
		return ctrl.Result{}, nil
	default:
		reqLogger.Info("NOP")
		return ctrl.Result{}, nil
	}

	rerr := r.Status().Update(context.TODO(), instance)
	if rerr != nil {
		return reconcile.Result{}, rerr
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AtReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cnatv1.At{}).
		Complete(r)
}
