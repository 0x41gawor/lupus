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

// ExecuteSpec defines the desired state of Execute
type ExecuteSpec struct {
	// Destination specifies where the input of Execute has to be sent
	// It will be always of type HTTP since this is what Lupout interface specifies it to be, but we will reuse the known
	// +kubebuilder:validation:Required
	Destination Destination `json:"destination,omitempty"`
	// Name of master element
	Master string `json:"master,omitempty"`
}

// ExecuteStatus defines the observed state of Execute
type ExecuteStatus struct {
	// Input contains operational data
	Input runtime.RawExtension `json:"input"`
	// Timestamp of the last update
	LastUpdated metav1.Time `json:"lastUpdated"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Execute is the Schema for the executes API
type Execute struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ExecuteSpec   `json:"spec,omitempty"`
	Status ExecuteStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ExecuteList contains a list of Execute
type ExecuteList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Execute `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Execute{}, &ExecuteList{})
}
