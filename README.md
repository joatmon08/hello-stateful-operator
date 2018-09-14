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