/*
Copyright 2024.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// LearnSpec defines the desired state of Learn
type LearnSpec struct {
	// This is the Destination in which Learn will store its Input
	// +kubebuilder:validation:Required
	Destination Destination `json:"destination,omitempty"`
	// Name of master element
	Master string `json:"master,omitempty"`
}

// LearnStatus defines the observed state of Learn
type LearnStatus struct {
	// Input contains operational data
	Input runtime.RawExtension `json:"input"`
	// Timestamp of the last update
	LastUpdated metav1.Time `json:"lastUpdated"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Learn is the Schema for the learns API
type Learn struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LearnSpec   `json:"spec,omitempty"`
	Status LearnStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// LearnList contains a list of Learn
type LearnList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Learn `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Learn{}, &LearnList{})
}
