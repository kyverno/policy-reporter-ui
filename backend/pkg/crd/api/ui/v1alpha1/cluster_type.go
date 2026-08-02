/*
Copyright 2020 The Kubernetes authors.

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

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +genclient:nonNamespaced
// +kubebuilder:object:root=true
// +kubebuilder:storageversion
// +kubebuilder:resource:path=clusters,scope="Cluster",shortName=cluster
// +kubebuilder:printcolumn:name="Name",type=string,JSONPath=`.metadata.name`,priority=1
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// Cluster is the Schema for the clusters API
type Cluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +required
	Spec ClusterSpec `json:"spec"`
}

type ClusterSpec struct {
	// +optional
	Title string `json:"title,omitempty"`

	// +required
	Host string `json:"host,omitempty"`

	// +optional
	HTTP2 bool `json:"http2,omitempty"`

	// +optional
	Plugins []Plugin `json:"plugins,omitempty"`

	// +optional
	AccessControl *AccessControl `json:"accessControl,omitempty"`

	// +optional
	BasicAuth *BasicAuth `json:"basicAuth,omitempty"`

	// +optional
	SecretRef bool `json:"secretRef,omitempty"`
}

type Plugin struct {
	// +optional
	Title string `json:"title,omitempty"`

	// +required
	Host string `json:"host,omitempty"`

	// +optional
	HTTP2 bool `json:"http2,omitempty"`
}

type BasicAuth struct {
	// +required
	Username string `json:"username,omitempty"`
	// +required
	Password string `json:"password,omitempty"`
}

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterList contains a list of Cluster
type ClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Cluster `json:"items"`
}
