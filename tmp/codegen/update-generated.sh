#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

chmod +x vendor/k8s.io/code-generator/generate-groups.sh
vendor/k8s.io/code-generator/generate-groups.sh \
"deepcopy,client" \
github.com/lyft/flinkk8soperator/pkg/client \
github.com/lyft/flinkk8soperator/pkg/apis \
app:v1beta1 \
--go-header-file "./tmp/codegen/boilerplate.go.txt"

if [[ "$PWD" != "$GOPATH/src/github.com/lyft/flinkk8soperator" ]]; then
  echo "Project not located in $GOPATH subdirectory. Copying generated files from $GOPATH/src/github.com/lyft/flinkk8soperator"
  cp -rf "$GOPATH/src/github.com/lyft/flinkk8soperator/" "$PWD"
fi