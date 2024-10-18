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
type Element struct {
	// Name of the element
	Name string `json:"name"`
	// Url (can be monitored-system or some external functional block)
	Url URL `json:"url"`
	// NextElement will be inferred from the array sequence.
}

// LoopSpec defines the desired state of Loop
type LoopSpec struct {
	// Name of the loop e.g. "adam"
	Name string `json:"name"`
	// Loop elements (e.g. "Observe", "Decide", "Execute"). Root Controller will instantiate them one by one.
	// Sequence matters here as the each element will have reference to the next elements on the list as its nextElement.
	// Element is a struct that has two fields: Element kind and url that has to be written in spec.
	Elements []Element `json:"elements"`
	// ObservationTimeInterval specifies the times between data fetches from monitored-system [seconds]
	ObservationTimeInterval int32 `json:"observationTimeInterval"`
}

// LoopStatus defines the observed state of Loop
type LoopStatus struct {
	// Active flag. If set to true, the loop elements are running, if not the controller will instantiate them
	IsActive bool `json:"isActive"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Loop is the Schema for the loops API
type Loop struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LoopSpec   `json:"spec,omitempty"`
	Status LoopStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// LoopList contains a list of Loop
type LoopList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Loop `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Loop{}, &LoopList{})
}
