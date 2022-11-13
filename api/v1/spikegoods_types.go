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

// SpikeGoodsSpec defines the desired state of SpikeGoods
type SpikeGoodsSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of SpikeGoods. Edit spikegoods_types.go to remove/update
	// 商品标题
	GoodsName string `json:"goodsName,omitempty"`
	// 商家ID
	SellerId string `json:"sellerId,omitempty"`
	// 原价
	Price int `json:"price,omitempty"`
	// 秒杀价格
	SpikePrice int `json:"spikePrice,omitempty"`
	// 活动持续时间
	Day int `json:"day,omitempty"`
	// 总库存数
	StockCount int `json:"stockCount,omitempty"`
}

// SpikeGoodsStatus defines the observed state of SpikeGoods
type SpikeGoodsStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// 商品ID
	GoodsId string `json:"goodsId,omitempty"`
	// 活动开始时间
	StartTime string `json:"start_time,omitempty"`
	// 活动结束时间
	EndTime string `json:"endTime,omitempty"`
	// 已售商品数
	Num int `json:"num,omitempty"`
	// 剩余库存数
	LastStockCount int `json:"lastCount,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// SpikeGoods is the Schema for the spikegoods API
type SpikeGoods struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SpikeGoodsSpec   `json:"spec,omitempty"`
	Status SpikeGoodsStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SpikeGoodsList contains a list of SpikeGoods
type SpikeGoodsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SpikeGoods `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SpikeGoods{}, &SpikeGoodsList{})
}

func (b *SpikeGoods) String() string {
	return fmt.Sprintf("\n spec: %+v\n status: %+v\n", b.Spec, b.Status)
}
