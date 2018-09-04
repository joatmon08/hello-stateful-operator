package hellostateful

import (
	"fmt"

	"github.com/joatmon08/hello-stateful-operator/pkg/apis/hello-stateful/v1alpha1"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	DISKSIZE        = 1 * 1000 * 1000 * 1000
	VOLUMENAME      = "log"
	VOLUMEMOUNTPATH = "/usr/share/hello"
	IMAGE           = "joatmon08/hello-stateful:1.0"
	CONTAINERNAME   = "hello-stateful"
)

var (
	storageClassName              = "standard"
	terminationGracePeriodSeconds = int64(10)
)

func Create(hs *v1alpha1.HelloStateful) error {
	statefulSet, err := newStatefulSet(hs)
	if err != nil {
		logrus.Errorf("Failed to generate statefulset: %v", err)
		return err
	}
	err = sdk.Create(statefulSet)
	if err != nil && !errors.IsAlreadyExists(err) {
		logrus.Errorf("Failed to create statefulset: %v", err)
		return err
	}

	err = sdk.Get(statefulSet)
	if err != nil {
		return fmt.Errorf("Failed to get statefulset: %v", err)
	}

	service, err := newService(hs)
	if err != nil {
		return fmt.Errorf("Failed to generate service: %v", err)
	}
	err = sdk.Create(service)
	if err != nil && !errors.IsAlreadyExists(err) {
		logrus.Errorf("Failed to create service: %v", err)
		return err
	}

	err = sdk.Get(service)
	if err != nil {
		return fmt.Errorf("Failed to get service: %v", err)
	}
	return nil
}

func labelsForHelloStateful(name string) map[string]string {
	return map[string]string{"app": name}
}

func newStatefulSet(cr *v1alpha1.HelloStateful) (*appsv1.StatefulSet, error) {
	labels := labelsForHelloStateful(cr.ObjectMeta.Name)
	diskSize := *resource.NewQuantity(DISKSIZE, resource.DecimalSI)
	statefulset := &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "StatefulSet",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.ObjectMeta.Name,
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.StatefulSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			ServiceName: cr.ObjectMeta.Name,
			Replicas:    &cr.Spec.Replicas,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					TerminationGracePeriodSeconds: &terminationGracePeriodSeconds,
					Containers: []corev1.Container{
						corev1.Container{
							Name:  CONTAINERNAME,
							Image: IMAGE,
							VolumeMounts: []corev1.VolumeMount{
								corev1.VolumeMount{
									Name:      VOLUMENAME,
									MountPath: VOLUMEMOUNTPATH,
								},
							},
						},
					},
				},
			},
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{
				corev1.PersistentVolumeClaim{
					ObjectMeta: metav1.ObjectMeta{
						Name:      VOLUMENAME,
						Namespace: cr.Namespace,
					},
					Spec: corev1.PersistentVolumeClaimSpec{
						AccessModes: []corev1.PersistentVolumeAccessMode{
							corev1.ReadWriteOnce,
						},
						StorageClassName: &storageClassName,
						Resources: corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceStorage: diskSize,
							},
						},
					},
				},
			},
		},
	}

	return statefulset, nil
}

func newService(cr *v1alpha1.HelloStateful) (*corev1.Service, error) {
	labels := labelsForHelloStateful(cr.ObjectMeta.Name)
	service := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.ObjectMeta.Name,
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: "None",
			Selector:  labels,
		},
	}
	return service, nil
}
