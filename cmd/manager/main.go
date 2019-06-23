package main

import (
	"flag"

	"k8s.io/klog"
	"kubeception.cloud/kubeception/cmd/manager/app"
	"kubeception.cloud/kubeception/pkg/util"
	"sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"
)

func main() {
	klog.InitFlags(nil)
	flag.Parse()
	log.SetLogger(log.ZapLogger(false))
	logger := log.Log.WithName("main")

	cmd := app.NewManagerCommand(util.ContextFromStopChannel(signals.SetupSignalHandler()), logger)
	if err := cmd.Execute(); err != nil {
		util.LogErrorAndExit(logger, err, "Error executing command")
	}
}
