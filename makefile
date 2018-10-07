OPERATOR_NAME=hello-stateful-operator

local: setup
	OPERATOR_NAME=$(OPERATOR_NAME) LOCAL=1 operator-sdk up local

generate-types:
	operator-sdk generate k8s

run-local:
	kubectl create -f deploy/cr.yaml

setup:
	kubectl apply -f deploy/crd.yaml
	kubectl apply -f deploy/rbac.yaml

run: setup
	kubectl create -f deploy/operator.yaml
	kubectl create -f deploy/cr.yaml

build:
	operator-sdk build joatmon08/$(OPERATOR_NAME):latest
	docker push joatmon08/$(OPERATOR_NAME):latest

tests:
	operator-sdk test -t ./test/e2e

clean:
	kubectl delete -f deploy/cr.yaml --ignore-not-found
	kubectl delete -f deploy/operator.yaml --ignore-not-found
	kubectl delete pvc --all
	kubectl delete pv --all

clean-all:
	kubectl delete -f deploy/ --ignore-not-found