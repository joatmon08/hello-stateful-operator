package v1alpha1

import (
	appsv1 "k8s.io/api/apps/v1"
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
	PersistentVolume      corev1.PersistentVolume      `json:"persistentVolume"`
	PersistentVolumeClaim corev1.PersistentVolumeClaim `json:"persistentVolumeClaim"`
	StatefulSet           appsv1.StatefulSet           `json:"statefulSet"`
	Service               corev1.Service               `json:"service"`
}
type HelloStatefulStatus struct {
	// Fill me
}
