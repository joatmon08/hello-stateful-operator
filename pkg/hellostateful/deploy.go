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

// Constants for hello-stateful StatefulSet & Volumes
const (
	DiskSize            = 1 * 1000 * 1000 * 1000
	AppVolumeName       = "app"
	AppVolumeMountPath  = "/usr/share/hello"
	HostProvisionerPath = "/tmp/hostpath-provisioner"
	AppImage            = "joatmon08/hello-stateful:latest"
	AppContainerName    = "hello-stateful"
	ImagePullPolicy     = corev1.PullAlways
)

var (
	storageClassName              = "standard"
	diskSize                      = *resource.NewQuantity(DiskSize, resource.DecimalSI)
	terminationGracePeriodSeconds = int64(10)
	accessMode                    = []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce}
	resourceList                  = corev1.ResourceList{corev1.ResourceStorage: diskSize}
)

// CreateVolume generates a PersistentVolume and Claim.
func CreateVolume(hs *v1alpha1.HelloStateful) error {
	pv, err := newPersistentVolume(hs)
	if err != nil {
		logrus.Errorf("Failed to generate persistentVolume: %v", err)
		return err
	}
	err = sdk.Create(pv)
	if err != nil && !errors.IsAlreadyExists(err) {
		logrus.Errorf("Failed to create persistentVolume: %v", err)
		return err
	}
	err = sdk.Get(pv)
	if err != nil {
		logrus.Errorf("Failed to get persistentVolume: %v", err)
		return err
	}

	pvc, err := newPersistentVolumeClaim(hs)
	if err != nil {
		logrus.Errorf("Failed to generate persistentVolumeClaim: %v", err)
		return err
	}
	err = sdk.Create(pvc)
	if err != nil && !errors.IsAlreadyExists(err) {
		logrus.Errorf("Failed to create persistentVolumeClaim: %v", err)
		return err
	}
	err = sdk.Get(pvc)
	if err != nil {
		logrus.Errorf("Failed to get persistentVolumeClaim: %v", err)
		return err
	}
	return nil
}

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
			Selector:    labelSelector(labels),
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
							Name:            AppContainerName,
							Image:           AppImage,
							ImagePullPolicy: ImagePullPolicy,
							VolumeMounts: []corev1.VolumeMount{
								corev1.VolumeMount{
									Name:      AppVolumeName,
									MountPath: AppVolumeMountPath,
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						corev1.Volume{
							Name: AppVolumeName,
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: cr.ObjectMeta.Name,
								},
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

func newPersistentVolume(cr *v1alpha1.HelloStateful) (*corev1.PersistentVolume, error) {
	labels := labelsForHelloStateful(cr.ObjectMeta.Name)
	pv := &corev1.PersistentVolume{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "PersistentVolume",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.ObjectMeta.Name,
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PersistentVolumeSpec{
			StorageClassName: storageClassName,
			AccessModes:      accessMode,
			Capacity:         resourceList,
			PersistentVolumeSource: corev1.PersistentVolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: fmt.Sprintf("%s/%s", HostProvisionerPath, cr.ObjectMeta.Name),
				},
			},
		},
	}
	addOwnerRefToObject(pv, asOwner(cr))
	return pv, nil
}

func newPersistentVolumeClaim(cr *v1alpha1.HelloStateful) (*corev1.PersistentVolumeClaim, error) {
	labels := labelsForHelloStateful(cr.ObjectMeta.Name)
	pvc := &corev1.PersistentVolumeClaim{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "PersistentVolumeClaim",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.ObjectMeta.Name,
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: accessMode,
			Selector:    labelSelector(labels),
			VolumeName:  cr.ObjectMeta.Name,
			Resources:   corev1.ResourceRequirements{Requests: resourceList},
		},
	}
	addOwnerRefToObject(pvc, asOwner(cr))
	return pvc, nil
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
			ClusterIP: corev1.ClusterIPNone,
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
