package hellostateful

import (
	"github.com/joatmon08/hello-stateful-operator/pkg/apis/hello-stateful/v1alpha1"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Constants for hello-stateful StatefulSet & Volumes
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

// Create generates a StatefulSet and Service for
// hello-stateful
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

	service, err := newService(hs)
	if err != nil {
		logrus.Errorf("Failed to generate service: %v", err)
		return err
	}
	err = sdk.Create(service)
	if err != nil && !errors.IsAlreadyExists(err) {
		logrus.Errorf("Failed to create service: %v", err)
		return err
	}
	err = sdk.Get(service)
	if err != nil {
		logrus.Errorf("Failed to get service: %v", err)
		return err
	}
	return nil
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
	addOwnerRefToObject(statefulset, asOwner(cr))
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
	addOwnerRefToObject(service, asOwner(cr))
	return service, nil
}

// addOwnerRefToObject appends the desired OwnerReference to the object
func addOwnerRefToObject(obj metav1.Object, ownerRef metav1.OwnerReference) {
	obj.SetOwnerReferences(append(obj.GetOwnerReferences(), ownerRef))
}

// asOwner returns an OwnerReference set as the memcached CR
func asOwner(hs *v1alpha1.HelloStateful) metav1.OwnerReference {
	trueVar := true
	return metav1.OwnerReference{
		APIVersion: hs.APIVersion,
		Kind:       hs.Kind,
		Name:       hs.Name,
		UID:        hs.UID,
		Controller: &trueVar,
	}
}
