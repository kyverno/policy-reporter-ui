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
// +kubebuilder:resource:path=customboards,scope="Cluster",shortName=cb
// +kubebuilder:printcolumn:name="Name",type=string,JSONPath=`.metadata.name`,priority=1
// +kubebuilder:printcolumn:name="Namespaces",type=array,JSONPath=`.summary.namespaces`
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// CustomBoard is the Schema for the customboards API
type CustomBoard struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// AccessControl is an optional reference to the custom boards access control
	// +optional
	AccessControl *AccessControl `json:"accessControl,omitempty"`

	// Allowed values are "resources" or "results".
	// +optional
	Display Display `json:"display,omitempty"`

	// NamespaceSelector allows to select visualized namespaces
	// +required
	NamespaceSelector NamespaceSelector `json:"namespaces"`

	// SourceSelector allows to select visualized sources
	// +optional
	SourceSelector *SourceSelector `json:"sources,omitempty"`

	// PolicyReportSelector allows to select visualized reports
	// +optional
	PolicyReportSelector *PolicyReportSelector `json:"policyReports,omitempty"`

	// +optional
	ClusterScope *ClusterScope `json:"clusterScope,omitempty"`
}

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CustomBoardList contains a list of ClusterPolicyReport
type CustomBoardList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CustomBoard `json:"items"`
}
