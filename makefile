IMAGE=ghc-infra-hello
VERSION=latest
NAME=test-ghc-infra-hello

local:
	kubectl apply -f deploy/crd.yaml
	kubectl apply -f deploy/rbac.yaml
	OPERATOR_NAME=hello-stateful-operator LOCAL=1 operator-sdk up local

run:
	kubectl apply -f deploy/cr.yaml

clean:
	kubectl delete -f deploy