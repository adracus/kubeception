#!/bin/bash

KUBERNETES_VERSIONS=(
  "v1.14.1"
  "v1.13.5"
  "v1.12.8"
  "v1.11.10"
)

for kubernetes_version in ${KUBERNETES_VERSIONS[@]}; do
  tag=adracus/dind-kubelet:$kubernetes_version
  docker push $tag
done

