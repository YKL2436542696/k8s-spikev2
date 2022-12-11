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
	"sync"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	myappv2 "k8s-spikev2/api/v2"
)

var wg sync.WaitGroup
var mutex sync.Mutex
var sucNum = 0
var failNum = 0

// MassesReconciler reconciles a Masses object
type MassesReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=myapp.spike.com,resources=masses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=myapp.spike.com,resources=masses/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=myapp.spike.com,resources=masses/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Masses object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *MassesReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// TODO(user): your logic here
	masses := &myappv2.Masses{}
	if err := r.Get(ctx, req.NamespacedName, masses); err != nil {
		log.Log.Error(err, "获取用户信息失败")
		return ctrl.Result{}, nil
	} else {
		// 如果是第一次创建,更新买家的status字段信息
		if masses.Status.PersonsId == "" {
			masses.Status.PersonsId = uuid.New().String()
			if err = r.Status().Update(ctx, masses); err != nil {
				log.Log.Error(err, "更新status字段信息失败")
				return ctrl.Result{}, nil
			}
			log.Log.Info(fmt.Sprintln("初始化用户群体对象：", masses.String()))
		}

		// 获取商品的库存，
		goodsList := &myappv2.GoodsList{}
		if err := r.List(ctx, goodsList); err != nil {
			log.Log.Error(err, "获取商品信息失败")
			return ctrl.Result{}, nil
		} else {
			item := goodsList.Items
			for _, goods := range item {
				// 与买家的成功数做对比
				if masses.Status.SuccessNum < goods.Spec.Stock {
					wg.Add(masses.Spec.PersonNum)
					for i := 0; i < masses.Spec.PersonNum; i++ {
						go r.buy(ctx)
					}
					wg.Wait()
					masses.Status.SuccessNum = sucNum
					masses.Status.FailNUm = failNum
					if err := r.Status().Update(ctx, masses); err != nil {
						log.Log.Error(err, "更新用户字段信息失败")
						return ctrl.Result{}, nil
					}
				} else if masses.Status.SuccessNum == goods.Spec.Stock {
					log.Log.Info(fmt.Sprintln("商品已售完，用户停止购买！该用户群体信息：", masses.String()))
				}
			}
		}
	}
	return ctrl.Result{}, nil
}

// 购买
func (r *MassesReconciler) buy(ctx context.Context) error {
	mutex.Lock()
	defer mutex.Unlock()
	defer wg.Done()

	// 获取商品数量，若商品数量>0,则商品数量-1。
	goodsList := &myappv2.GoodsList{}
	if err := r.List(ctx, goodsList); err != nil {
		log.Log.Error(err, "获取商品信息失败")
		return nil
	} else {
		item := goodsList.Items
		for _, goods := range item {
			// 先获取商品，判断是否满足条件，准备数据,进行更新
			if goods.Status.CurrentNum > 0 {
				goods.Status.CurrentNum -= 1
				// 如果满足条件
				if err := r.Status().Update(ctx, &goods); err != nil {
					log.Log.Error(err, "更新商品字段信息失败")
					//failNum += 1
					return err
				}
				sucNum += 1
				log.Log.Info(fmt.Sprintln("购买成功！当前购买成功人数：", sucNum))
				return nil
			} else {
				// 不满足条件
				failNum += 1
				log.Log.Info(fmt.Sprintln("购买失败！当前购买失败人数：", failNum))
				return nil
			}
		}
		return nil
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *MassesReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&myappv2.Masses{}).
		Complete(r)
}
