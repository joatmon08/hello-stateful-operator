package hellostateful

import (
	"github.com/joatmon08/hello-stateful-operator/pkg/apis/hello-stateful/v1alpha1"
)

func Get(hellostateful *v1alpha1.HelloStateful) error {
	return nil
}

func getPersistentVolumeIds() []string {
	return nil
}

func getPersistentVolumeClaimIds() []string {
	return nil
}

func getBackendVolumes() []string {
	return nil
}
