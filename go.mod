module kubeception.cloud/kubeception

go 1.12

require (
	github.com/beorn7/perks v1.0.0 // indirect
	github.com/cosmos72/gomacro v0.0.0-20190602202531-a277e261f22e // indirect
	github.com/go-logr/logr v0.1.0
	github.com/go-logr/zapr v0.1.1 // indirect
	github.com/golang/groupcache v0.0.0-20190129154638-5b532d6fd5ef // indirect
	github.com/golang/mock v1.3.0
	github.com/golangci/golangci-lint v1.16.1-0.20190425135923-692dacb773b7
	github.com/google/btree v1.0.0 // indirect
	github.com/google/gofuzz v1.0.0 // indirect
	github.com/google/uuid v1.1.1 // indirect
	github.com/gregjones/httpcache v0.0.0-20190212212710-3befbb6ad0cc // indirect
	github.com/hashicorp/golang-lru v0.5.1 // indirect
	github.com/imdario/mergo v0.3.7 // indirect
	github.com/mattn/go-runewidth v0.0.4 // indirect
	github.com/onsi/ginkgo v1.8.0
	github.com/onsi/gomega v1.5.0
	github.com/peterbourgon/diskv v2.0.1+incompatible // indirect
	github.com/prometheus/client_model v0.0.0-20190129233127-fd36f4220a90 // indirect
	github.com/prometheus/common v0.3.0 // indirect
	github.com/prometheus/procfs v0.0.0-20190425082905-87a4384529e0 // indirect
	github.com/spf13/cobra v0.0.3
	golang.org/x/crypto v0.0.0-20190426145343-a29dc8fdc734 // indirect
	golang.org/x/oauth2 v0.0.0-20190402181905-9f3314589c9a // indirect
	golang.org/x/sys v0.0.0-20190606165138-5da285871e9c // indirect
	golang.org/x/time v0.0.0-20190308202827-9d24e82272b4 // indirect
	k8s.io/api v0.0.0-20190409021203-6e4e0e4f393b
	k8s.io/apimachinery v0.0.0-20190404173353-6a84e37a896d
	k8s.io/apiserver v0.0.0-20190507070644-e9c02aff496d // indirect
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
	k8s.io/code-generator v0.0.0-20190620073620-d55040311883
	k8s.io/klog v0.3.1
	k8s.io/kube-openapi v0.0.0-20190426233423-c5d3b0f4bee0 // indirect
	sigs.k8s.io/cluster-api v0.0.0-20190610203311-5ed76e24e031
	sigs.k8s.io/controller-runtime v0.2.0-beta.2
	sigs.k8s.io/controller-tools v0.2.0-beta.2.0.20190610175510-203d8e8ab133
)

replace (
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20190409022649-727a075fdec8
	sigs.k8s.io/cluster-api => github.com/vincepri/cluster-api v0.0.0-20190621203312-5270eece8091
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.2.0-beta.2
	sigs.k8s.io/controller-tools => sigs.k8s.io/controller-tools v0.2.0-beta.2
)
