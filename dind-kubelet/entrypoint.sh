#!/bin/sh

dockerd $ADDITIONAL_DOCKERD_ARGS > /var/log/dockerd.log 2>&1 &

while [ ! -f /var/run/docker.pid ]; do echo "Docker pid not available"; sleep 1; done

/kubelet "$@"
