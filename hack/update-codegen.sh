#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

go run "$DIR/../vendor/sigs.k8s.io/controller-tools/cmd/controller-gen/main.go" paths=./pkg/apis/... $@

