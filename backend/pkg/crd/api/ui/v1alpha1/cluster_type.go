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
// +kubebuilder:resource:path=clusters,scope="Cluster",shortName=cl
// +kubebuilder:printcolumn:name="Name",type=string,JSONPath=`.metadata.name`,priority=1
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// Cluster is the Schema for the external clusters API
// It is used to configure additional external clusters for the PolicyReporter UI
type Cluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +required
	Spec ClusterSpec `json:"spec"`
}

// ClusterSpec defines the desired state of Cluster
type ClusterSpec struct {
	// Optional title for the cluster in the Dropdown menu. If not set, the cluster name will be used.
	// +optional
	Title string `json:"title,omitempty"`

	// Host of the external Policy Reporter API. This is required to connect to the api and needs to be reachable from Policy Reporter UI.
	// +required
	Host string `json:"host,omitempty"`

	// Forces HTTP2 for the connection to the external Policy Reporter API.
	// +optional
	HTTP2 bool `json:"http2,omitempty"`

	// SkipTLS disables TLS verification for the connection to the external Policy Reporter API.
	// +optional
	SkipTLS bool `json:"skipTLS,omitempty"`

	// Path to a certificate file for the connection to the external Policy Reporter API.
	// +optional
	Certificate string `json:"certificate,omitempty"`

	// List of plugins running in the external cluster.
	// +optional
	Plugins []Plugin `json:"plugins,omitempty"`

	// Access control configuration for the external cluster. If not set, the cluster will be accessible for all users.
	// Requires Authentication to be enabled in the Policy Reporter UI.
	// +optional
	AccessControl *AccessControl `json:"accessControl,omitempty"`

	// BasicAuth configuration for the external API. If not set, no authentication will be used.
	// +optional
	BasicAuth *BasicAuth `json:"basicAuth,omitempty"`

	// SecretRef is the name of a secret containing additional configuration for the external cluster API.
	// The secret needs to be in the same namespace as the Policy Reporter UI deployment.
	// +optional
	SecretRef string `json:"secretRef,omitempty"`
}

// Plugin is the Schema for plugins running in the external cluster
type Plugin struct {
	// The Report source/engine the plugin is referenced to
	// +required
	Source string `json:"source,omitempty"`

	// Host of the plugin API. This is required to connect to the plugin and needs to be reachable from Policy Reporter UI.
	// +required
	Host string `json:"host,omitempty"`

	// Forces HTTP2 for the connection to the plugin API.
	// +optional
	HTTP2 bool `json:"http2,omitempty"`

	// SkipTLS disables TLS verification for the connection to the plugin API.
	// +optional
	SkipTLS bool `json:"skipTLS,omitempty"`

	// Path to a certificate file for the connection to the plugin API.
	// +optional
	Certificate string `json:"certificate,omitempty"`

	// BasicAuth configuration for the plugin API. If not set, no authentication will be used.
	// +optional
	BasicAuth *BasicAuth `json:"basicAuth,omitempty"`

	// SecretRef is the name of a secret containing additional configuration for the plugin API.
	// The secret needs to be in the same namespace as the Policy Reporter UI deployment.
	// +optional
	SecretRef string `json:"secretRef,omitempty"`
}

// BasicAuth is the Schema for basic authentication configuration for external APIs
type BasicAuth struct {
	// Username for basic authentication to the external API.
	// +required
	Username string `json:"username,omitempty"`
	// Password for basic authentication to the external API.
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
