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
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MasterSpec defines the desired state of Master
type MasterSpec struct {
	// Name of the Master CR
	Name string `json:"name"`

	// Elements is a list of elements, each with its name, type, and corresponding spec
	Elements []Element `json:"elements"`
}

// MasterStatus defines the observed state of Master
type MasterStatus struct {
	// Active flag. If set to true, the loop elements are running, if not the controller will instantiate them
	IsActive bool `json:"isActive"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Master is the Schema for the masters API
type Master struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MasterSpec   `json:"spec,omitempty"`
	Status MasterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MasterList contains a list of Master
type MasterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Master `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Master{}, &MasterList{})
}
