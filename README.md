kubeception
===========
[![Go Report Card](https://goreportcard.com/badge/github.com/adracus/kubeception)](https://goreportcard.com/report/github.com/adracus/kubeception)
[![Documentation](https://godoc.org/github.com/adracus/kubeception?status.svg)](http://godoc.org/github.com/adracus/kubeception)
![GitHub](https://img.shields.io/github/license/adracus/kubeception.svg)
[![Build Status](https://travis-ci.org/adracus/kubeception.svg?branch=master)](https://travis-ci.org/adracus/kubeception)

kubeception is an implementation of the [Kubernetes Cluster API](https://github.com/kubernetes-sigs/cluster-api)
using Kubernetes itself as environment.

> ⚠ Caution ⚠: Kubeception is far from production ready and far from secure:
> It does not care at all about proper authentication / authorization.
> This is merely a PoC of what can be done with the cluster API + the
> kubeception approach.

Intro
-----

The control plane of a cluster is self-hosted by providing the required
components (etcd, Kube API server etc.) via Kubernetes primitives
(Deployments, Services, StatefulSets).

Machines are created by running a privileged pod with kubelet and docker in a
'dind' setup (docker-in-docker). The image for this is built and provided via
[dind-kubelet](/dind-kubelet).

As of now, there is neither a proper overlay network between the nodes nor
a cluster-dns. This may come in the future but as of now this is just a minimal
PoC that allows running a hello-world docker container.

Setup
-----

You'll need a Kubernetes cluster for this. The recommended way of getting one
is by using [kind](https://github.com/kubernetes-sigs/kind), as kubeception is
developed mainly with kind and thus this is the most tested way.

Once you have a cluster, apply all cluster-api CRDs to the cluster. You can do
this by running

```bash
kubectl apply -f config
```

In another window, run

```bash
make start
```

This will run the kubeception cluster and machine controller which reconciles
the respective resources.

### Control Planes

Once you've done that you're ready to create your first cluster by running

```bash
kubectl apply -f example/cluster.yaml
```

This spins up the required components in your current namespace.
After a while, when all components are there, you can connect to the API
server and experiment with it. For quick experiments, there is a hack script
which sets up a container inside your cluster with the kubeconfig already at
the right place and `kubectl` in your path. You can run it via

```bash
./hack/hyper.sh
```

### Machines

To setup a machine, run

```bash
kubectl apply -f example/machine.yaml
```

This runs your 'machine' and connects it to the API server of the referenced
cluster. Once the machine is up and running, you can now connect to your
cluster and run the 'hello-world':

```bash
./hack/hyper.sh

# Verify that the node is there and ready
kubectl get nodes

# Run the hello-world
kubectl run --replicas=1 --restart=Never --image=hello-world -it hello
```
