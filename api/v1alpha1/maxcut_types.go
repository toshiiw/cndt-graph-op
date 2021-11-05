/*
Copyright 2021.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MaxCutSpec defines the desired state of MaxCut
type MaxCutSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Iteration count of the randomized algorithm runs
	//+kubebuilder:default:=100
	//+kubebuilder:validation:Optional
	//+kubebuilder:validation:Minimum=1
	Iteration int32 `json:"iteration"`
}

// MaxCutStatus defines the observed state of MaxCut
type MaxCutStatus struct {
	// Total weight of the cut
	Total               int32    `json:"total"`
	VertexSet           []string `json:"vertexset"`
	ComplementVertexSet []string `json:"complementvertexset"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:JSONPath=".status.total",name=Total,type=integer

// MaxCut is the Schema for the maxcuts API
type MaxCut struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MaxCutSpec   `json:"spec,omitempty"`
	Status MaxCutStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MaxCutList contains a list of MaxCut
type MaxCutList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MaxCut `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MaxCut{}, &MaxCutList{})
}
