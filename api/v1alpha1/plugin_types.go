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
	"k8s.io/apimachinery/pkg/runtime"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// PluginSpec defines the desired state of Plugin

type PluginReq struct {
	Name       string               `json:"name"`
	Parameters runtime.RawExtension `json:"parameters,omitempty"`
	Code       string               `json:"code"`
	Version    string               `json:"version,omitempty"`
}

type PluginSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Plugin. Edit plugin_types.go to remove/update
	Plugin PluginReq `json:"plugin"`
}

// PluginStatus defines the observed state of Plugin
type PluginStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Name string `json:"name"`
	Code string `json:"code"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Plugin is the Schema for the plugins API
type Plugin struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PluginSpec   `json:"spec,omitempty"`
	Status PluginStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// PluginList contains a list of Plugin
type PluginList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Plugin `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Plugin{}, &PluginList{})
}
