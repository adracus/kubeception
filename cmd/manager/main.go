package main

import (
	"flag"
	"github.com/adracus/kubeception/pkg/apis/kubeception/install"
	"github.com/adracus/kubeception/pkg/controller"
	"github.com/go-logr/logr"
	"k8s.io/klog"
	"os"
	clusterapis "sigs.k8s.io/cluster-api/pkg/apis"
	clustercontroller "sigs.k8s.io/cluster-api/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"
)

func logErrorAndExit(log logr.Logger, err error, msg string, keysAndValues ...interface{}) {
	log.Error(err, msg, keysAndValues)
	os.Exit(1)
}

func main() {
	klog.InitFlags(nil)
	flag.Parse()
	log.SetLogger(log.ZapLogger(false))
	logger := log.Log.WithName("main")

	cfg := config.GetConfigOrDie()
	mgr, err := manager.New(cfg, manager.Options{
		LeaderElection: false, // TODO: Dev only, make configurable
	})
	if err != nil {
		logErrorAndExit(logger, err, "Could not initialize manager")
	}

	if err := clusterapis.AddToScheme(mgr.GetScheme()); err != nil {
		logErrorAndExit(logger, err, "Could not modify scheme")
	}

	install.Install(mgr.GetScheme())

	if err := clustercontroller.AddToManager(mgr); err != nil {
		logErrorAndExit(logger, err, "Could add cluster-api controllers")
	}

	if err := controller.AddToManager(mgr); err != nil {
		logErrorAndExit(logger, err, "Could not add controllers")
	}

	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		logErrorAndExit(logger, err, "Error running manager")
	}
}
