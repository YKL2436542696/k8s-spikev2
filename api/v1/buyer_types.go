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

// BuyerSpec defines the desired state of Buyer
type BuyerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Buyer. Edit buyer_types.go to remove/update
	// 买家账号名
	BuyerName string `json:"buyerName,omitempty"`
	// 买家收货地址
	Address string `json:"address,omitempty"`
	// 买家电话
	Phone string `json:"phone,omitempty"`
}

// BuyerStatus defines the observed state of Buyer
type BuyerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// 买家ID
	BuyerId string `json:"buyerId,omitempty"`
	// 买家订单ID列表
	OrderIdList []string `json:"orderList,omitempty"`
	// 买家总消费额
	SpendMoney int `json:"spendMoney,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Buyer is the Schema for the buyers API
type Buyer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BuyerSpec   `json:"spec,omitempty"`
	Status BuyerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// BuyerList contains a list of Buyer
type BuyerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Buyer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Buyer{}, &BuyerList{})
}

func (b *Buyer) String() string {
	return fmt.Sprintf("\n spec: %+v\n status: %+v\n", b.Spec, b.Status)
}
