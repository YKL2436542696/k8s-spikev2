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

package v1

import (
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SpikeOrderSpec defines the desired state of SpikeOrder
type SpikeOrderSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// Foo is an example field of SpikeOrder. Edit spikeorder_types.go to remove/update
	// 买家ID
	BuyerId string `json:"buyerId,omitempty"`
	// 秒杀商品列表
	GoodsMap map[string]int `json:"goodsMap,omitempty"`
	// 是否付款
	IsPay bool `json:"isPay,omitempty"`
}

// SpikeOrderStatus defines the observed state of SpikeOrder
type SpikeOrderStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// 订单ID
	OrderId string `json:"orderId,omitempty"`
	// 订单创建时间
	CreateTime string `json:"createTime,omitempty"`
	// 订单过期时间
	ExpiredTime string `json:"expiredTime,omitempty"`
	// 优惠金额
	Discount int `json:"discount,omitempty"`
	// 需付金额
	Money int `json:"money,omitempty"`
	// 订单状态 0 未支付 1 已支付 2 已过期
	PayStatus string `json:"payStatus,omitempty"`
	// 订单支付时间
	PayTime string `json:"payTime,omitempty"`
	// 收货地址
	Address string `json:"address,omitempty"`
	// 收货电话
	Phone string `json:"phone,omitempty"`
	// 收货人
	Receiver string `json:"receiver,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// SpikeOrder is the Schema for the spikeorders API
type SpikeOrder struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SpikeOrderSpec   `json:"spec,omitempty"`
	Status SpikeOrderStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SpikeOrderList contains a list of SpikeOrder
type SpikeOrderList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SpikeOrder `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SpikeOrder{}, &SpikeOrderList{})
}

func (b *SpikeOrder) String() string {
	return fmt.Sprintf("\n spec: %+v\n status: %+v\n", b.Spec, b.Status)
}
