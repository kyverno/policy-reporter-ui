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

type Display string

const (
	DisplayResources Display = "resources"
	DisplayResults   Display = "results"
)

// +kubebuilder:oneOf:={required:{emails}}
// +kubebuilder:oneOf:={required:{groups}}

// AccessControl defines allowed emails and groups for a resource
type AccessControl struct {
	// Source is an identifier for the policy engine that manages this report
	// +optional
	Emails []string `json:"emails,omitempty"`

	// Rule is the name or identifier of the rule within the policy
	// +optional
	Groups []string `json:"groups,omitempty"`
}

type ValueFilter struct {
	// +optional
	Include []string `json:"include,omitempty"`
	// +optional
	Exclude []string `json:"exclude,omitempty"`
}

type Filter struct {
	// +optional
	NamespaceKinds ValueFilter `json:"namespaceKinds"`
	// +optional
	ClusterKinds ValueFilter `json:"clusterKinds"`
	// +optional
	Results ValueFilter `json:"results"`
	// +optional
	Severities ValueFilter `json:"severities"`
}

// +kubebuilder:oneOf:={required:{labelSelector}}
// +kubebuilder:oneOf:={required:{list}}

// NamespaceSelector allows to select the visualized namespaces via label selector or list
type NamespaceSelector struct {
	// +optional
	LabelSelector map[string]string `json:"labelSelector"`
	// +optional
	List []string `json:"list"`
}

// SourceSelector allows to select the visualized sources via list
type SourceSelector struct {
	// +optional
	List []string `json:"list"`
}

// PolicyReportSelector allows to select the visualized reports via label selector
type PolicyReportSelector struct {
	// +optional
	LabelSelector map[string]string `json:"labelSelector"`
}

type ClusterScope struct {
	// +required
	Enabled bool `json:"enabled"`
}
