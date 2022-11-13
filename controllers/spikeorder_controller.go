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
	"sort"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	myappv1 "k8s-spikev2/api/v1"
)

// SpikeOrderReconciler reconciles a SpikeOrder object
type SpikeOrderReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=myapp.spike.com,resources=spikeorders,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=myapp.spike.com,resources=spikeorders/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=myapp.spike.com,resources=spikeorders/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the SpikeOrder object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *SpikeOrderReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	// TODO(user): your logic here
	log.Log.Info("1. SpikeOrder info start")

	spikeOrder := &myappv1.SpikeOrder{}
	if err := r.Get(ctx, req.NamespacedName, spikeOrder); err != nil {
		log.Log.Error(err, "获取订单对象出错")
	} else {
		// 如果是第一次创建,更新买家的status字段信息
		if spikeOrder.Status.OrderId == "" {
			log.Log.Info("1.1 初始化订单对象")
			if err = initUpdateOrder(ctx, spikeOrder, r); err != nil {
				return ctrl.Result{}, nil
			}
			// 初始化后，判断当前资源是否为未支付状态，是则开启定时任务
			if spikeOrder.Spec.IsPay == false && spikeOrder.Status.PayStatus == "0" &&
				spikeOrder.Status.ExpiredTime != "" {
				// 定时5分钟后检测订单支付状态
				go TaskOrder(ctx, spikeOrder, r)
			}
		}
		// 判断资源变化是否为支付操作
		if spikeOrder.Spec.IsPay == true && spikeOrder.Status.PayStatus == "0" &&
			spikeOrder.Status.ExpiredTime != "" {
			if err := pay(ctx, spikeOrder, r); err != nil {
				return ctrl.Result{}, nil
			}
		}

		log.Log.Info(fmt.Sprintln("4. 获得订单对象：", spikeOrder.String()))
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SpikeOrderReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&myappv1.SpikeOrder{}).
		Complete(r)
}

// initUpdate 初始化更新status 字段
func initUpdateOrder(ctx context.Context, spikeOrder *myappv1.SpikeOrder, r *SpikeOrderReconciler) error {
	spikeOrder.Status.OrderId = uuid.New().String()
	now := time.Now()
	spikeOrder.Status.CreateTime = now.Format("2006-01-02 15:04:05")
	// 默认过期时间5分钟
	t, _ := time.ParseDuration("5m")
	spikeOrder.Status.ExpiredTime = now.Add(t).Format("2006-01-02 15:04:05")

	// 根据商品ID 计算优惠金额、需付金额
	spikeGoodsMap := spikeOrder.Spec.GoodsMap
	sumMoney := 0
	disMoney := 0
	// 记录下商品列表
	tarGoodsList := []myappv1.SpikeGoods{}
	log.Log.Info("1.2 扣减对应商品剩余库存,更新商品信息")
	spikeGoodsList := &myappv1.SpikeGoodsList{}
	if err := r.List(ctx, spikeGoodsList); err != nil {
		log.Log.Error(err, "获取商品列表对象出错")
		return err
	} else {
		// 遍历所有的商品
		goodItem := spikeGoodsList.Items
		for s, i := range spikeGoodsMap {
			for _, goods := range goodItem {
				if s == goods.Status.GoodsId {
					sumMoney += goods.Spec.Price * i
					disMoney += goods.Spec.SpikePrice * i
					// 先把库存扣减了
					goods.Status.LastStockCount -= i
					if spikeOrder.Spec.IsPay {
						log.Log.Info("1.3 若已支付，增加对应商品已售数")
						goods.Status.Num += i
					}
					if err := r.Status().Update(ctx, &goods); err != nil {
						log.Log.Error(err, "更新goods.status字段信息失败")
						return err
					}
					tarGoodsList = append(tarGoodsList, goods)
					break
				}
			}
		}
	}
	spikeOrder.Status.Discount = sumMoney - disMoney
	spikeOrder.Status.Money = disMoney

	// 订单状态
	if spikeOrder.Spec.IsPay {
		spikeOrder.Status.PayStatus = "1"
		spikeOrder.Status.ExpiredTime = ""
	} else {
		spikeOrder.Status.PayStatus = "0"
	}

	log.Log.Info("1.3  更新买家信息")
	buyerList := &myappv1.BuyerList{}
	if err := r.List(ctx, buyerList); err != nil {
		log.Log.Error(err, "获取买家列表对象出错")
		return err
	} else {
		buyerItem := buyerList.Items
		for i := range buyerItem {
			// 找到目标用户
			if spikeOrder.Spec.BuyerId == buyerItem[i].Status.BuyerId {
				log.Log.Info("1.3.1 更新对应订单的收货信息")
				spikeOrder.Status.Address = buyerItem[i].Spec.Address
				spikeOrder.Status.Phone = buyerItem[i].Spec.Phone
				spikeOrder.Status.Receiver = buyerItem[i].Spec.BuyerName
				log.Log.Info("1.3.2 更新对应买家名下订单列表")
				orderIdList := buyerItem[i].Status.OrderIdList
				orderIdList = append(orderIdList, spikeOrder.Status.OrderId)
				buyerItem[i].Status.OrderIdList = orderIdList

				if spikeOrder.Spec.IsPay {
					log.Log.Info("1.3.3 若已支付，更新对应买家总消费额")
					buyerItem[i].Status.SpendMoney += spikeOrder.Status.Money
				}
				if err := r.Status().Update(ctx, &buyerItem[i]); err != nil {
					log.Log.Error(err, "更新buyer.status字段信息失败")
					return err
				}
				break
			}
		}
	}

	if spikeOrder.Spec.IsPay {
		log.Log.Info("1.4 若已支付，增加对应卖家的总销售额,并更新商家信息")
		sellerList := &myappv1.SellerList{}
		if err := r.List(ctx, sellerList); err != nil {
			log.Log.Error(err, "获取卖家列表对象出错")
			return err
		} else {
			sellerItem := sellerList.Items
			for i, _ := range sellerItem {
				for j, _ := range tarGoodsList {
					if sellerItem[i].Status.SellerId == tarGoodsList[j].Spec.SellerId {
						goodId := tarGoodsList[j].Status.GoodsId
						sellerItem[i].Status.SalesMoney += tarGoodsList[j].Spec.SpikePrice * spikeGoodsMap[goodId]
						if err = r.Status().Update(ctx, &sellerItem[i]); err != nil {
							log.Log.Error(err, "更新seller.status字段信息失败")
							return err
						}
						break
					}
				}
			}
		}
	}

	// 更新订单status
	log.Log.Info("1.5 更新订单信息")
	if err := r.Status().Update(ctx, spikeOrder); err != nil {
		log.Log.Error(err, "更新order.status字段信息失败")
		return err
	}

	return nil
}

// pay 支付操作
func pay(ctx context.Context, spikeOrder *myappv1.SpikeOrder, r *SpikeOrderReconciler) error {
	now := time.Now()
	et := spikeOrder.Status.ExpiredTime
	stamp, _ := time.ParseInLocation("2006-01-02 15:04:05", et, time.Local)
	goodMap := spikeOrder.Spec.GoodsMap
	if now.After(stamp) {
		// 已过期
		log.Log.Info(fmt.Sprintln("3.1 订单已过期，无法支付！"))
		// 更新订单状态
		spikeOrder.Status.PayStatus = "2"
		// 回滚商品的剩余库存数
		log.Log.Info(fmt.Sprintln("3.1.1 回滚商品的剩余库存数"))
		spikeGoodsList := &myappv1.SpikeGoodsList{}
		if err := r.List(ctx, spikeGoodsList); err != nil {
			log.Log.Error(err, "获取商品列表对象出错")
			return err
		} else {
			goodItem := spikeGoodsList.Items
			for s, num := range goodMap {
				for i, _ := range goodItem {
					if s == goodItem[i].Status.GoodsId {
						// 回滚库存
						goodItem[i].Status.LastStockCount += num
						if err := r.Status().Update(ctx, &goodItem[i]); err != nil {
							log.Log.Error(err, "更新goods.status字段信息失败")
							return err
						}
						break
					}
				}
			}
		}
	} else {
		// 没过期 可支付
		log.Log.Info(fmt.Sprintln("3.2 订单未过期，可支付！"))
		spikeOrder.Status.PayStatus = "1"
		spikeOrder.Status.ExpiredTime = ""

		log.Log.Info(fmt.Sprintln("3.2.1 增加对应商品已售数"))
		// 记录下商品列表，以便后续更新商品信息
		tarGoodsList := []myappv1.SpikeGoods{}
		spikeGoodsList := &myappv1.SpikeGoodsList{}
		if err := r.List(ctx, spikeGoodsList); err != nil {
			log.Log.Error(err, "获取商品列表对象出错")
			return err
		} else {
			goodItem := spikeGoodsList.Items
			for s, i := range goodMap {
				for j, _ := range goodItem {
					if goodItem[j].Status.GoodsId == s {
						goodItem[j].Status.Num += i
						if err := r.Status().Update(ctx, &goodItem[j]); err != nil {
							log.Log.Error(err, "更新goods.status字段信息失败")
							return err
						}
						tarGoodsList = append(tarGoodsList, goodItem[j])
						break
					}
				}
			}
		}

		log.Log.Info(fmt.Sprintln("3.2.2 更新对应买家总消费额"))
		buyerId := spikeOrder.Spec.BuyerId
		buyerList := &myappv1.BuyerList{}
		if err := r.List(ctx, buyerList); err != nil {
			log.Log.Error(err, "获取买家列表对象出错")
			return err
		} else {
			buyerItem := buyerList.Items
			for i, _ := range buyerItem {
				if buyerItem[i].Status.BuyerId == buyerId {
					buyerItem[i].Status.SpendMoney += spikeOrder.Status.Money
					if err := r.Status().Update(ctx, &buyerItem[i]); err != nil {
						log.Log.Error(err, "更新buyer.status字段信息失败")
						return err
					}
					break
				}
			}

		}

		log.Log.Info(fmt.Sprintln("3.2.3 增加对应卖家的总销售额"))
		sellerList := &myappv1.SellerList{}
		if err := r.List(ctx, sellerList); err != nil {
			log.Log.Error(err, "获取卖家列表对象出错")
			return err
		} else {
			sellerItem := sellerList.Items
			for i, _ := range tarGoodsList {
				for j, _ := range sellerItem {
					if stringIn(tarGoodsList[i].Status.GoodsId, sellerItem[j].Status.SellerGoodsIdList) {
						goodId := tarGoodsList[i].Status.GoodsId
						sellerItem[j].Status.SalesMoney += tarGoodsList[i].Spec.SpikePrice * goodMap[goodId]
						if err = r.Status().Update(ctx, &sellerItem[i]); err != nil {
							log.Log.Error(err, "更新seller.status字段信息失败")
							return err
						}
						break
					}
				}
			}
		}
	}

	log.Log.Info(fmt.Sprintln("3.3 更新对应订单的状态"))
	if err := r.Status().Update(ctx, spikeOrder); err != nil {
		log.Log.Error(err, "更新order.status字段信息失败")
		return err
	}
	return nil
}

// stringIn 判断string[] 中是否有目标string
func stringIn(target string, strArray []string) bool {
	sort.Strings(strArray)
	index := sort.SearchStrings(strArray, target)
	//index的取值：[0,len(str_array)]
	if index < len(strArray) && strArray[index] == target {
		return true
	}
	return false
}

// TaskOrder 定时任务检测订单是否超时过期
func TaskOrder(ctx context.Context, spikeOrder *myappv1.SpikeOrder, r *SpikeOrderReconciler) {
	log.Log.Info("2.1 启动定时任务...")
	timer := time.NewTimer(5 * time.Minute)
	<-timer.C
	log.Log.Info("2.2 5分钟已到，检查订单是否未支付...")
	if spikeOrder.Spec.IsPay == false && spikeOrder.Status.PayStatus == "0" &&
		spikeOrder.Status.ExpiredTime != "" {
		now := time.Now()
		et := spikeOrder.Status.ExpiredTime
		stamp, _ := time.ParseInLocation("2006-01-02 15:04:05", et, time.Local)
		goodMap := spikeOrder.Spec.GoodsMap
		if now.After(stamp) {
			// 已过期
			// 更新订单状态
			log.Log.Info(fmt.Sprintln("2.2.1 未支付..."))
			log.Log.Info(fmt.Sprintln("2.2.2 回滚商品的剩余库存数..."))
			spikeOrder.Status.PayStatus = "2"
			if err := r.Status().Update(ctx, spikeOrder); err != nil {
				log.Log.Error(err, "更新order.status字段信息失败")
			}

			// 回滚商品的剩余库存数
			log.Log.Info(fmt.Sprintln("2.2.3 回滚商品的剩余库存数..."))
			spikeGoodsList := &myappv1.SpikeGoodsList{}
			if err := r.List(ctx, spikeGoodsList); err != nil {
				log.Log.Error(err, "获取商品列表对象出错")

			} else {
				goodItem := spikeGoodsList.Items
				for s, num := range goodMap {
					for i, _ := range goodItem {
						if s == goodItem[i].Status.GoodsId {
							// 回滚库存
							goodItem[i].Status.LastStockCount += num
							if err := r.Status().Update(ctx, &goodItem[i]); err != nil {
								log.Log.Error(err, "更新goods.status字段信息失败")
							}
							break
						}
					}
				}
			}
		} else {
			log.Log.Info(fmt.Sprintln("2.3 订单未过期..."))
		}
	} else {
		log.Log.Info(fmt.Sprintln("2.4 订单已正常支付..."))
	}
}
