# hello-stateful-operator

This is a Kubernetes operator that creates HelloStateful resources
using the Operator framework. HelloStateful writes the date and time
to a `stateful.log` file every 10 seconds. It it useful for practicing
with PersistentVolumes.

The operator itself creates:
* A StatefulSet for the application.
* A CronJob for backup of the stateful.log file.
* A Job run if you choose to do a restore because it existed before.

The CustomResource should be defined as follows:
```
apiVersion: "hello-stateful.example.com/v1alpha1"
kind: "HelloStateful"
metadata:
  name: <name here>
spec:
  replicas: <number of replicas>
  restoreFromExisting: <true | false>
  backupSchedule: <cron schedule notation for backup>
```

## Pre-Requisites
* Minikube
* Kubernetes 1.10

## Run
1. Deploy the CustomResource Definition & RBAC.
    ```
    make setup
    ```
1. Deploy the Operator.
    ```
    kubectl create -f deploy/operator.yaml
    ```
1. Deploy the HelloStateful Custom Resource.
    ```
    kubectl create -f deploy/cr.yaml
    ```

## Test
```
make tests
```

## To Consider
In order for us to immutably restore a PV & PVC, we need to create & delete them separately
from the StatefulSet.

This is a bit annoying to statically define, especially since StatefulSets should probably
have their own PVs and PVCs anyway. But statically defining a PVC means tying replicas to
the same PV, which isn't good.