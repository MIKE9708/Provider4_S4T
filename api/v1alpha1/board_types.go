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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/MIKE9708/s4t-sdk-go/pkg/api/boards"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// BoardSpec defines the desired state of Board
type BoardSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Board. Edit board_types.go to remove/update
	Board boards.Board 
}

// BoardStatus defines the observed state of Board
// this are the fields visible in kubernetes
type BoardStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Status string  `json:"status,omitempty"`
	UUID string  `json:"uuid,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Board is the Schema for the boards API
type Board struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BoardSpec   `json:"spec,omitempty"`
	Status BoardStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// BoardList contains a list of Board
type BoardList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Board `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Board{}, &BoardList{})
}
