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
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	myappv1 "k8s-spikev2/api/v1"
)

// SpikeGoodsReconciler reconciles a SpikeGoods object
type SpikeGoodsReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=myapp.spike.com,resources=spikegoods,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=myapp.spike.com,resources=spikegoods/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=myapp.spike.com,resources=spikegoods/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the SpikeGoods object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *SpikeGoodsReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// TODO(user): your logic here
	// TODO(user): your logic here
	log.Log.Info("1. SpikeGoods info start")
	spikeGoods := &myappv1.SpikeGoods{}
	if err := r.Get(ctx, req.NamespacedName, spikeGoods); err != nil {
		log.Log.Error(err, "获取商品对象出错")
	} else {
		// 如果是第一次创建,更新买家的status字段信息
		if spikeGoods.Status.GoodsId == "" {
			log.Log.Info("1.1 初始化商品对象")
			if err = initUpdateGoods(ctx, spikeGoods, r); err != nil {
				return ctrl.Result{}, nil
			}
		}
		log.Log.Info(fmt.Sprintln("2. 获得商品对象：", spikeGoods.String()))
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SpikeGoodsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&myappv1.SpikeGoods{}).
		Complete(r)
}

// initUpdate 初始化更新status 字段
func initUpdateGoods(ctx context.Context, spikeGoods *myappv1.SpikeGoods, r *SpikeGoodsReconciler) error {
	spikeGoods.Status.GoodsId = uuid.New().String()
	now := time.Now()
	spikeGoods.Status.StartTime = now.Format("2006-01-02 15:04:05")
	day := spikeGoods.Spec.Day
	spikeGoods.Status.EndTime = now.AddDate(0, 0, day).Format("2006-01-02 15:04:05")
	spikeGoods.Status.LastStockCount = spikeGoods.Spec.StockCount

	if err := r.Status().Update(ctx, spikeGoods); err != nil {
		log.Log.Error(err, "更新goods.status字段信息失败")
		return err
	} else {
		// 更新卖家名下的商品列表
		log.Log.Info("1.2 更新卖家名下的商品列表")
		sellerList := &myappv1.SellerList{}
		if err := r.List(ctx, sellerList); err != nil {
			log.Log.Error(err, "获取卖家列表对象出错")
			return err
		} else {
			// 找到卖家
			log.Log.Info("1.3 查找卖家名下的商品列表")
			item := sellerList.Items
			for i, _ := range item {
				if spikeGoods.Spec.SellerId == item[i].Status.SellerId {
					// 更新卖家信息
					idList := item[i].Status.SellerGoodsIdList
					idList = append(idList, spikeGoods.Status.GoodsId)
					item[i].Status.SellerGoodsIdList = idList
					if err := r.Status().Update(ctx, &item[i]); err != nil {
						log.Log.Error(err, "更新seller.status字段信息失败")
						return err
					}
					break
				}
			}
		}
	}
	return nil
}
