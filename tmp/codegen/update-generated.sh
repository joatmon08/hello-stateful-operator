#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

vendor/k8s.io/code-generator/generate-groups.sh \
deepcopy \
github.com/joatmon08/hello-stateful-operator/pkg/generated \
github.com/joatmon08/hello-stateful-operator/pkg/apis \
hello-stateful:v1alpha1 \
--go-header-file "./tmp/codegen/boilerplate.go.txt"
