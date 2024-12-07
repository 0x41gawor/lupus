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

// ElementSpec defines the desired state of Element
type ElementSpec struct {
	// Name is the name of the element, its distinct from Kubernetes API Object name, but rather serves ease of managemenet aspect for loop-designer
	Name string `json:"name"`
	// Descr is the description of the element, same as Name is serves as ease of management aspect for loop-designer
	Descr string `json:"descr"`
	// Actions is a list of Actions that Element has to perform
	Actions []Action `json:"actions,omitempty"`
	// Next is a list of element to which send the Data final form
	Next []Next `json:"next,omitempty"`
	// Name of master element
	Master string `json:"master,omitempty"`
}

// ElementStatus defines the observed state of Element
type ElementStatus struct {
	// Input contains operational data
	Input runtime.RawExtension `json:"input"`
	// Timestamp of the last update
	LastUpdated metav1.Time `json:"lastUpdated"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Element is the Schema for the elements API
type Element struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ElementSpec   `json:"spec,omitempty"`
	Status ElementStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ElementList contains a list of Element
type ElementList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Element `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Element{}, &ElementList{})
}
