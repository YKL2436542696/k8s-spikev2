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

// SellerSpec defines the desired state of Seller
type SellerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Seller. Edit seller_types.go to remove/update
	// 卖家电话
	SellerPhone string `json:"sellerPhone,omitempty"`
}

// SellerStatus defines the observed state of Seller
type SellerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// 卖家ID
	SellerId string `json:"sellerId,omitempty"`
	// 卖家发布的秒杀商品ID列表
	SellerGoodsIdList []string `json:"sellerGoodsList,omitempty"`
	// 卖家总销售额
	SalesMoney int `json:"salesMoney,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Seller is the Schema for the sellers API
type Seller struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SellerSpec   `json:"spec,omitempty"`
	Status SellerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SellerList contains a list of Seller
type SellerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Seller `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Seller{}, &SellerList{})
}

func (b *Seller) String() string {
	return fmt.Sprintf("\n spec: %+v\n status: %+v\n", b.Spec, b.Status)
}
