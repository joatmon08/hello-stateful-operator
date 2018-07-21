package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
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
	PersistentVolume      corev1.PersistentVolumeSpec      `json:"persistentVolume"`
	PersistentVolumeClaim corev1.PersistentVolumeClaimSpec `json:"persistentVolumeClaim"`
}
type HelloStatefulStatus struct {
	// Fill me
}
