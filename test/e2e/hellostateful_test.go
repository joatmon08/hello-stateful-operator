package e2e

import (
	goctx "context"
	"fmt"
	"testing"
	"time"

	hellov1alpha1 "github.com/joatmon08/hello-stateful-operator/pkg/apis/hello-stateful/v1alpha1"

	framework "github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	retryInterval         = time.Second * 10
	timeout               = time.Second * 75
	helloStatefulName     = "test-hello-stateful"
	customResourceFixture = "./test/e2e/fixtures/cr.yaml"
)

func TestHelloStateful(t *testing.T) {
	helloStatefulList := &hellov1alpha1.HelloStatefulList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "HelloStateful",
			APIVersion: "hello-stateful.example.com/v1alpha1",
		},
	}
	err := framework.AddToFrameworkScheme(hellov1alpha1.AddToScheme, helloStatefulList)
	if err != nil {
		t.Fatalf("failed to add custom resource scheme to framework: %v", err)
	}
	// run subtests
	t.Run("hello-stateful-group", func(t *testing.T) {
		t.Run("Instance", HelloStatefulInstance)
	})
}

func createHelloStatefulCustomResource(namespace string) *hellov1alpha1.HelloStateful {
	return &hellov1alpha1.HelloStateful{
		TypeMeta: metav1.TypeMeta{
			Kind:       "HelloStateful",
			APIVersion: "hello-stateful.example.com/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      helloStatefulName,
			Namespace: namespace,
		},
		Spec: hellov1alpha1.HelloStatefulSpec{
			Replicas: 1,
		},
	}
}

func createHelloStatefulTest(t *testing.T, f *framework.Framework, ctx *framework.TestCtx) error {
	namespace, err := ctx.GetNamespace()
	if err != nil {
		return fmt.Errorf("could not get namespace: %v", err)
	}

	exampleHelloStateful := createHelloStatefulCustomResource(namespace)
	if err != nil {
		return err
	}

	err = f.DynamicClient.Create(goctx.TODO(), exampleHelloStateful)
	if err != nil {
		return err
	}
	// wait for test-hello-stateful to reach 1 replica
	if err = WaitForStatefulSet(t, f.KubeClient, namespace, exampleHelloStateful.Name, 1, retryInterval, timeout); err != nil {
		return err
	}

	jobName := fmt.Sprintf("backup-%s", exampleHelloStateful.Name)
	return WaitForCronJob(t, f.KubeClient, namespace, jobName, retryInterval, timeout)
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
	t.Log(namespace)
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

	if err = createHelloStatefulTest(t, f, ctx); err != nil {
		t.Fatal(err)
	}
}
