/*
Copyright 2022 ykl.

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
	"fmt"
	"github.com/google/uuid"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	myappv1 "k8s-spikev2/api/v1"
)

// BuyerReconciler reconciles a Buyer object
type BuyerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=myapp.spike.com,resources=buyers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=myapp.spike.com,resources=buyers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=myapp.spike.com,resources=buyers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Buyer object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *BuyerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// TODO(user): your logic here
	log.Log.Info("1. Buyer info start")
	buyer := &myappv1.Buyer{}
	if err := r.Get(ctx, req.NamespacedName, buyer); err != nil {
		log.Log.Error(err, "获取买家信息失败")
		return ctrl.Result{}, nil
	} else {
		// 如果是第一次创建,更新买家的status字段信息
		if buyer.Status.BuyerId == "" {
			buyer.Status.BuyerId = uuid.New().String()
			log.Log.Info("1.1 初始化买家对象")
			if err = r.Status().Update(ctx, buyer); err != nil {
				log.Log.Error(err, "status字段信息失败")
				return ctrl.Result{}, nil
			}
		}

		log.Log.Info(fmt.Sprintln("2. 获得买家对象：", buyer.String()))
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *BuyerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&myappv1.Buyer{}).
		Complete(r)
}
