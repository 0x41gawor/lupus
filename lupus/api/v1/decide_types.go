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

// DecideSpec defines the desired state of Decide
type DecideSpec struct {
	// Actions is a list of Actions that Decide has to perform see Action in types.go
	Actions []Action `json:"actions,omitempty"`
	// Next is a list of elements to which send the output
	Next []Next `json:"next,omitempty"`
	// Name of master element
	Master string `json:"master,omitempty"`
}

// DecideStatus defines the observed state of Decide
type DecideStatus struct {
	// Input contains operational data
	Input runtime.RawExtension `json:"input"`
	// Timestamp of the last update
	LastUpdated metav1.Time `json:"lastUpdated"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Decide is the Schema for the decides API
type Decide struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DecideSpec   `json:"spec,omitempty"`
	Status DecideStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DecideList contains a list of Decide
type DecideList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Decide `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Decide{}, &DecideList{})
}
