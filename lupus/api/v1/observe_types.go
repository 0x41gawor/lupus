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

// URL defines the structure for monitored_system_url with path and method
type URL struct {
	Path   string `json:"path"`   // URL path
	Method string `json:"method"` // HTTP method (GET, POST, etc.)
}

// ObserveSpec defines the desired state of Observe
type ObserveSpec struct {
	// Name is a string field
	Name string `json:"name"`

	// NextElement is a string field
	NextElement string `json:"nextElement"`

	// Url contains URL details with path and method
	Url URL `json:"url"`

	// ObservationTimeInterval specifies the times between data fetches in seconds
	ObservationTimeInterval int32 `json:"observationTimeInterval"`
}

// ObserveStatus defines the observed state of Observe
type ObserveStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Observe is the Schema for the observes API
type Observe struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ObserveSpec   `json:"spec,omitempty"`
	Status ObserveStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ObserveList contains a list of Observe
type ObserveList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Observe `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Observe{}, &ObserveList{})
}
