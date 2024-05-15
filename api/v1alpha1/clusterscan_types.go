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
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ClusterScanSpec defines the desired state of ClusterScan
type ClusterScanSpec struct {
	Deployment []DeploymentSpec `json:"deployment"`
}

type DeploymentSpec struct {
	Name        string      `json:"name"`
	Namespace   string      `json:"namespace"`
	Reccuring   bool        `json:"recurring"`
	Schedule    string      `json:"schedule,omitempty"`
	JobTemplate JobTemplate `json:"jobTemplate"`
}

type JobTemplate struct {
	// Containers specifies the list of containers belonging to the pod.
	Name string `json:"name"`

	// Image is the Docker image for the container.
	Image string `json:"image"`

	// Command specifies the command to run in the container.
	Command []string `json:"command"`

	// Args specifies the arguments to pass to the command.
	Args []string `json:"args,omitempty"`

	// RestartPolicy specifies the restart policy for the pod.
	RestartPolicy string `json:"restartPolicy"`
}

// ClusterScanStatus defines the observed state of ClusterScan
type ClusterScanStatus struct {
	// existing fields...
	ToolStatuses []ToolStatus `json:"toolStatuses,omitempty"`
}

type ToolStatus struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ClusterScan is the Schema for the clusterscans API
type ClusterScan struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterScanSpec   `json:"spec,omitempty"`
	Status ClusterScanStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ClusterScanList contains a list of ClusterScan
type ClusterScanList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterScan `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ClusterScan{}, &ClusterScanList{})
}
