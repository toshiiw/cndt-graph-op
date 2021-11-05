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

// MaxFlowSpec defines the desired state of MaxFlow
type MaxFlowSpec struct {
	From string `json:"from,omitempty"`
	To   string `json:"to,omitempty"`
}

// MaxFlowStatus defines the observed state of MaxFlow
type MaxFlowStatus struct {
	Flow  int32 `json:"flow"`
	Stale bool  `json:"stale"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:JSONPath=".status.flow",name=Flow,type=integer
//+kubebuilder:printcolumn:JSONPath=".status.stale",name=Stale,type=boolean

// MaxFlow is the Schema for the maxflows API
type MaxFlow struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MaxFlowSpec   `json:"spec,omitempty"`
	Status MaxFlowStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MaxFlowList contains a list of MaxFlow
type MaxFlowList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MaxFlow `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MaxFlow{}, &MaxFlowList{})
}
