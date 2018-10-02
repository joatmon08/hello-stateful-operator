package hellostateful

import (
	"fmt"

	"github.com/joatmon08/hello-stateful-operator/pkg/apis/hello-stateful/v1alpha1"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Update should change the replicas for a StatefulSet.
func Update(hs *v1alpha1.HelloStateful) error {
	statefulSet := &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "StatefulSet",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      hs.GetName(),
			Namespace: hs.GetNamespace(),
		},
	}
	err := sdk.Get(statefulSet)
	if err != nil {
		return fmt.Errorf("Failed to get statefulset: %v", err)
	}

	if *statefulSet.Spec.Replicas != hs.Spec.Replicas {
		statefulSet.Spec.Replicas = &(hs.Spec.Replicas)

		err = sdk.Update(statefulSet)
		if err != nil {
			return fmt.Errorf("failed to update hello-stateful statefulset: %v", err)
		}
	}
	return nil
}
