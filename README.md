Todo:
1. Convert statefulset to using volume claim template
2. Update status spec to get volume claim name & volume claim information
3. figure out how to time backups (snapshots)

Run in local mode for exploratory testing:
1.  Create the CustomResourceDefintion.
    ```
    $ kubectl create -f deploy/crd.yaml
    ```
1.  Apply the RBAC.
    ```
    $ kubectl create -f deploy/rbac.yaml
    ```
    ```
1. Start it up in local mode.
    ```
    $ OPERATOR_NAME=hello-stateful-operator LOCAL=1 operator-sdk up local
    ```
1. Create the CustomResource.
    ```
    $ kubectl create -f deploy/cr.yaml
    ```

## To Consider
In order for us to immutably restore a PV & PVC, we need to create & delete them separately
from the StatefulSet.

This is a bit annoying to statically define, especially since StatefulSets should probably
have their own PVs and PVCs anyway. But statically defining a PVC means tying replicas to
the same PV, which isn't good.