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

// DiEdgeSpec defines the desired state of DiEdge
type DiEdgeSpec struct {
	Capacity int32 `json:"capacity"`
	//+kubebuilder:default:=0
	//+kubebuilder:validation:Optional
	Allocated int32  `json:"allocated"`
	From      string `json:"from"`
	To        string `json:"to"`
}

// DiEdgeStatus defines the observed state of DiEdge
type DiEdgeStatus struct {
	// nothing to declare...
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:JSONPath=".spec.from",name=From,type=string
//+kubebuilder:printcolumn:JSONPath=".spec.to",name=To,type=string
//+kubebuilder:printcolumn:JSONPath=".spec.allocated",name=Allocated,type=integer
//+kubebuilder:printcolumn:JSONPath=".spec.capacity",name=Capacity,type=integer

// DiEdge is the Schema for the diedges API
type DiEdge struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DiEdgeSpec   `json:"spec,omitempty"`
	Status DiEdgeStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DiEdgeList contains a list of DiEdge
type DiEdgeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DiEdge `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DiEdge{}, &DiEdgeList{})
}
