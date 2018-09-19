IMAGE=ghc-infra-hello
VERSION=latest
NAME=test-ghc-infra-hello

local:
	kubectl apply -f deploy/crd.yaml
	kubectl apply -f deploy/rbac.yaml
	OPERATOR_NAME=hello-stateful-operator LOCAL=1 operator-sdk up local

run:
	kubectl create -f deploy/operator.yaml
	kubectl create -f deploy/cr.yaml

build:
	operator-sdk build joatmon08/hello-stateful-operator:latest
	docker push joatmon08/hello-stateful-operator:latest

clean:
	kubectl delete -f deploy/cr.yaml --ignore-not-found
	kubectl delete -f deploy/operator.yaml --ignore-not-found
	kubectl delete pvc --all
	kubectl delete pv --all