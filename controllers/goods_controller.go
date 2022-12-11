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

	myappv2 "k8s-spikev2/api/v2"
)

// GoodsReconciler reconciles a Goods object
type GoodsReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=myapp.spike.com,resources=goods,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=myapp.spike.com,resources=goods/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=myapp.spike.com,resources=goods/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Goods object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *GoodsReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// 如果是第一次创建,更新买家的status字段信息
	goods := &myappv2.Goods{}
	if err := r.Get(ctx, req.NamespacedName, goods); err != nil {
		log.Log.Error(err, "获取商品信息失败")
		return ctrl.Result{}, nil
	} else {
		// 如果是第一次创建,更新买家的status字段信息
		if goods.Status.GoodsId == "" {
			goods.Status.GoodsId = uuid.New().String()
			goods.Status.CurrentNum = goods.Spec.Stock
			if err = r.Status().Update(ctx, goods); err != nil {
				log.Log.Error(err, "更新status字段信息失败")
				return ctrl.Result{}, nil
			}
			log.Log.Info(fmt.Sprintln("初始化商品对象：", goods.String()))
		}
		if goods.Status.CurrentNum == 0 {
			log.Log.Info(fmt.Sprintln("商品售卖完毕！商品信息：", goods.String()))
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *GoodsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&myappv2.Goods{}).
		Complete(r)
}
