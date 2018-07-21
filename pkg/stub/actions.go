package stub

import (
	"fmt"

	"github.com/joatmon08/hello-stateful-operator/pkg/apis/hello-stateful/v1alpha1"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
)

func Create(hellostateful *v1alpha1.HelloStateful) error {
	persistentVolume := newPersistentVolume(hellostateful)
	err := sdk.Create(persistentVolume)
	if err != nil && !errors.IsAlreadyExists(err) {
		logrus.Errorf("Failed to create persistent volume: %v", err)
		return err
	}

	err = sdk.Get(persistentVolume)
	if err != nil {
		return fmt.Errorf("Failed to get persistent volume: %v", err)
	}

	persistentVolumeClaim := newPersistentVolumeClaim(hellostateful)
	err = sdk.Create(persistentVolumeClaim)
	if err != nil && !errors.IsAlreadyExists(err) {
		logrus.Errorf("Failed to create persistent volume claim: %v", err)
		return err
	}

	err = sdk.Get(persistentVolumeClaim)
	if err != nil {
		return fmt.Errorf("Failed to get persistent volume claim: %v", err)
	}

	statefulSet := newStatefulSet(hellostateful)
	err = sdk.Create(statefulSet)
	if err != nil && !errors.IsAlreadyExists(err) {
		logrus.Errorf("Failed to create statefulset: %v", err)
		return err
	}

	err = sdk.Get(statefulSet)
	if err != nil {
		return fmt.Errorf("Failed to get statefulset: %v", err)
	}

	service := newService(hellostateful)
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
