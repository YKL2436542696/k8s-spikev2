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

package v2

import (
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// GoodsSpec defines the desired state of Goods
type GoodsSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	GoodsName string `json:"goodsName,omitempty"`
	// 初始库存
	Stock int `json:"stock,omitempty"`
}

// GoodsStatus defines the observed state of Goods
type GoodsStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// 商品ID
	GoodsId string `json:"goodsId,omitempty"`
	// 当前商品数
	CurrentNum int `json:"currentNum,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Goods is the Schema for the goods API
type Goods struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GoodsSpec   `json:"spec,omitempty"`
	Status GoodsStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// GoodsList contains a list of Goods
type GoodsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Goods `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Goods{}, &GoodsList{})
}

func (b *Goods) String() string {
	return fmt.Sprintf("\n spec: %+v\n status: %+v\n", b.Spec, b.Status)
}
