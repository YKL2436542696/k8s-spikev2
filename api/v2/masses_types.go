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

// MassesSpec defines the desired state of Masses
type MassesSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// 买家所属地
	PersonAddr string `json:"personAddr,omitempty"`
	// 买家人数
	PersonNum int `json:"personNum,omitempty"`
}

// MassesStatus defines the observed state of Masses
type MassesStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// 用户群体ID
	PersonsId string `json:"persons_id,omitempty"`
	// 购买成功数
	SuccessNum int `json:"SuccessNum,omitempty"`
	// 购买失败数
	FailNUm int `json:"failNUm,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Masses is the Schema for the masses API
type Masses struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MassesSpec   `json:"spec,omitempty"`
	Status MassesStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MassesList contains a list of Masses
type MassesList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Masses `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Masses{}, &MassesList{})
}

func (b *Masses) String() string {
	return fmt.Sprintf("\n spec: %+v\n status: %+v\n", b.Spec, b.Status)
}
