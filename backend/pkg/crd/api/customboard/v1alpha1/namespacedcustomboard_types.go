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
// +kubebuilder:storageversion
// +kubebuilder:resource:path=namespacecustomboards,scope="Namespaced",shortName=ncb
// +kubebuilder:printcolumn:name="Name",type=string,JSONPath=`.metadata.name`,priority=1
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// NamespaceCustomBoard is the Schema for the namespace customboards API
type NamespaceCustomBoard struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +required
	Spec NamespaceCustomBoardSpec `json:"spec"`
}

type NamespaceCustomBoardSpec struct {
	// +optional
	Title string `json:"title,omitempty"`

	// AccessControl is an optional reference to the custom boards access control
	// +optional
	AccessControl *AccessControl `json:"accessControl,omitempty"`

	// Allowed values are "resources" or "results".
	// +optional
	Display Display `json:"display,omitempty"`

	// SourceSelector allows to select visualized sources
	// +optional
	SourceSelector *SourceSelector `json:"sources,omitempty"`

	// PolicyReportSelector allows to select visualized reports
	// +optional
	PolicyReportSelector *PolicyReportSelector `json:"policyReports,omitempty"`

	// +optional
	Filter *Filter `json:"filter,omitempty"`
}

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NamespaceCustomBoardList contains a list of NamespaceCustomBoard
type NamespaceCustomBoardList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NamespaceCustomBoard `json:"items"`
}
