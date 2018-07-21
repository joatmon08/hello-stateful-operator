package stub

import (
	"context"
	"fmt"

	"github.com/joatmon08/hello-stateful-operator/pkg/apis/hello-stateful/v1alpha1"

	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const defaultName = "hello-stateful"
const defaultCapacity = 1 * 1024 * 1024 * 1024

func NewHandler() sdk.Handler {
	return &Handler{}
}

type Handler struct {
	// Fill me
}

func (h *Handler) Handle(ctx context.Context, event sdk.Event) error {
	switch o := event.Object.(type) {
	case *v1alpha1.HelloStateful:
		hellostateful := o
		persistentVolume := newHelloStatefulPersistentVolume(hellostateful)
		err := sdk.Create(persistentVolume)
		if err != nil && !errors.IsAlreadyExists(err) {
			logrus.Errorf("Failed to create %s persistent volume: %v", defaultName, err)
			return err
		}

		err = sdk.Get(persistentVolume)
		if err != nil {
			return fmt.Errorf("Failed to get %s persistent volume: %v", defaultName, err)
		}

		persistentVolumeClaim := newHelloStatefulPersistentVolumeClaim(hellostateful)
		err = sdk.Create(persistentVolumeClaim)
		if err != nil && !errors.IsAlreadyExists(err) {
			logrus.Errorf("Failed to create %s persistent volume claim: %v", defaultName, err)
			return err
		}

		err = sdk.Get(persistentVolumeClaim)
		if err != nil {
			return fmt.Errorf("Failed to get %s persistent volume claim: %v", defaultName, err)
		}
	}
	return nil
}

func newHelloStatefulPersistentVolume(cr *v1alpha1.HelloStateful) *corev1.PersistentVolume {
	return &corev1.PersistentVolume{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PersistentVolume",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-pv", defaultName),
			Namespace: cr.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(cr, schema.GroupVersionKind{
					Group:   v1alpha1.SchemeGroupVersion.Group,
					Version: v1alpha1.SchemeGroupVersion.Version,
					Kind:    "HelloStateful",
				}),
			},
		},
		Spec: cr.Spec.PersistentVolume,
	}
}

func newHelloStatefulPersistentVolumeClaim(cr *v1alpha1.HelloStateful) *corev1.PersistentVolumeClaim {
	persistentVolumeClaimSpec := cr.Spec.PersistentVolumeClaim
	persistentVolumeClaimSpec.VolumeName = fmt.Sprintf("%s-pv", defaultName)
	return &corev1.PersistentVolumeClaim{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PersistentVolumeClaim",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-pv-claim", defaultName),
			Namespace: cr.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(cr, schema.GroupVersionKind{
					Group:   v1alpha1.SchemeGroupVersion.Group,
					Version: v1alpha1.SchemeGroupVersion.Version,
					Kind:    "HelloStateful",
				}),
			},
		},
		Spec: persistentVolumeClaimSpec,
	}
}
