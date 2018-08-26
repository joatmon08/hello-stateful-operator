package e2e

import (
	goctx "context"
	"fmt"
	"testing"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	v1alpha1 "github.com/joatmon08/hello-stateful-operator/pkg/apis/hello-stateful/v1alpha1"

	framework "github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	retryInterval = time.Second * 5
	timeout       = time.Second * 30
)

func TestHelloStateful(t *testing.T) {
	helloStatefulList := &v1alpha1.HelloStatefulList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "HelloStateful",
			APIVersion: "hello-stateful.example.com/v1alpha1",
		},
	}
	err := framework.AddToFrameworkScheme(v1alpha1.AddToScheme, helloStatefulList)
	if err != nil {
		t.Fatalf("failed to add custom resource scheme to framework: %v", err)
	}
	// run subtests
	t.Run("hello-stateful-group", func(t *testing.T) {
		t.Run("Instance", HelloStatefulInstance)
	})
}

func createTest(t *testing.T, f *framework.Framework, ctx framework.TestCtx) error {
	namespace, err := ctx.GetNamespace()
	if err != nil {
		return fmt.Errorf("could not get namespace: %v", err)
	}
	// create hello-stateful custom resource
	exampleHelloStateful := &v1alpha1.HelloStateful{
		TypeMeta: metav1.TypeMeta{
			Kind:       "HelloStateful",
			APIVersion: "hello-stateful.example.com/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "example-hello-stateful",
			Namespace: namespace,
		},
		Spec: v1alpha1.HelloStatefulSpec{
			PersistentVolume:      corev1.PersistentVolume{},
			PersistentVolumeClaim: corev1.PersistentVolumeClaim{},
			StatefulSet:           appsv1.StatefulSet{},
			Service:               corev1.Service{},
		},
	}
	err = f.DynamicClient.Create(goctx.TODO(), exampleHelloStateful)
	if err != nil {
		return err
	}
	// wait for example-hello-stateful to reach 1 replica
	return e2eutil.WaitForDeployment(t, f.KubeClient, namespace, "example-hello-stateful", 1, retryInterval, timeout)
}

func HelloStatefulInstance(t *testing.T) {
	t.Parallel()
	ctx := framework.NewTestCtx(t)
	defer ctx.Cleanup(t)
	err := ctx.InitializeClusterResources()
	if err != nil {
		t.Fatalf("failed to initialize cluster resources: %v", err)
	}
	t.Log("Initialized cluster resources")
	namespace, err := ctx.GetNamespace()
	if err != nil {
		t.Fatal(err)
	}
	// get global framework variables
	f := framework.Global
	// wait for memcached-operator to be ready
	err = e2eutil.WaitForDeployment(t, f.KubeClient, namespace, "hello-stateful-operator", 1, retryInterval, timeout)
	if err != nil {
		t.Fatal(err)
	}

	if err = createTest(t, f, ctx); err != nil {
		t.Fatal(err)
	}
}
