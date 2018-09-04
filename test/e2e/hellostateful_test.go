package e2e

import (
	goctx "context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/ghodss/yaml"
	hellov1alpha1 "github.com/joatmon08/hello-stateful-operator/pkg/apis/hello-stateful/v1alpha1"

	framework "github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	retryInterval         = time.Second * 10
	timeout               = time.Second * 60
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

func readYAMLtoHelloStateful(filename string) (*hellov1alpha1.HelloStateful, error) {
	yamlData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	jsonBytes, err := yaml.YAMLToJSON(yamlData)
	if err != nil {
		return nil, err
	}
	// unmarshal the json into the kube struct
	var helloStateful = hellov1alpha1.HelloStateful{}
	err = json.Unmarshal(jsonBytes, &helloStateful)
	if err != nil {
		return nil, err
	}
	return &helloStateful, nil
}

func createTest(t *testing.T, f *framework.Framework, ctx framework.TestCtx) error {
	namespace, err := ctx.GetNamespace()
	if err != nil {
		return fmt.Errorf("could not get namespace: %v", err)
	}

	exampleHelloStateful, err := readYAMLtoHelloStateful(customResourceFixture)
	if err != nil {
		return err
	}

	exampleHelloStateful.ObjectMeta.Namespace = namespace

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

	if err = createTest(t, f, ctx); err != nil {
		t.Fatal(err)
	}
}
