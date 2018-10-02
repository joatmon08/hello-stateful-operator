package hellostateful

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/joatmon08/hello-stateful-operator/pkg/apis/hello-stateful/v1alpha1"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// UpdateStatus updates the status of HelloStateful CustomResource
// with hostpath and restoration state.
func UpdateStatus(cr *v1alpha1.HelloStateful) error {
	pvcList := &corev1.PersistentVolumeClaimList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PersistentVolumeClaim",
			APIVersion: "v1",
		},
	}
	labelSelector := labels.SelectorFromSet(labelsForHelloStateful(cr.Name)).String()
	listOps := &metav1.ListOptions{LabelSelector: labelSelector}
	err := sdk.List(cr.Namespace, pvcList, sdk.WithListOptions(listOps))
	if err != nil {
		return err
	}
	pvNames := getPersistentVolumes(pvcList.Items)
	backendVolumes, err := getBackendVolumes(pvNames)
	if err != nil {
		return err
	}
	if !reflect.DeepEqual(backendVolumes, cr.Status.BackendVolumes) {
		cr.Status.BackendVolumes = backendVolumes
		err := sdk.Update(cr)
		if err != nil {
			return fmt.Errorf("failed to update hellostateful status: %v", err)
		}
	}
	return nil
}

func getPersistentVolumes(pvcs []corev1.PersistentVolumeClaim) []string {
	var pvNames []string
	for _, pvc := range pvcs {
		pvNames = append(pvNames, pvc.Spec.VolumeName)
	}
	return pvNames
}

func getPersistentVolumeClaims(pvcs []corev1.PersistentVolumeClaim) []string {
	var pvcNames []string
	for _, pvc := range pvcs {
		pvcNames = append(pvcNames, pvc.Name)
	}
	return pvcNames
}

func getBackendVolumes(persistentVolumes []string) ([]string, error) {
	var backendVolumes []string
	var config *rest.Config
	var err error
	if os.Getenv("LOCAL") == "1" {
		kubeconfig := filepath.Join(
			os.Getenv("HOME"), ".kube", "config",
		)
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	} else {
		config, err = rest.InClusterConfig()
	}
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	for _, pv := range persistentVolumes {
		name := fmt.Sprintf("metadata.name=%s", pv)
		volumes, err := clientset.CoreV1().PersistentVolumes().List(metav1.ListOptions{FieldSelector: name})
		if err != nil {
			return nil, err
		}
		if len(volumes.Items) != 1 {
			return backendVolumes, fmt.Errorf("Searched for %s PersistentVolume, found %d volumes", name, len(volumes.Items))
		}
		backendVolumes = append(backendVolumes, volumes.Items[0].Spec.HostPath.Path)
	}
	return backendVolumes, nil
}
