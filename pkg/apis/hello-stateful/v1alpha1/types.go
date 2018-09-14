package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type HelloStatefulList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []HelloStateful `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type HelloStateful struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              HelloStatefulSpec   `json:"spec"`
	Status            HelloStatefulStatus `json:"status,omitempty"`
}

type HelloStatefulSpec struct {
	Replicas int32 `json:"replicas"`
}
type HelloStatefulStatus struct {
	BackendVolumes []string `json:"backendVolumes"`
}
