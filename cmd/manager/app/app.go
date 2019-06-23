package app

import (
	"context"
	"flag"

	certificateinstall "kubeception.cloud/kubeception/pkg/apis/certificate/install"

	"github.com/go-logr/logr"
	"github.com/spf13/cobra"
	kubeceptioninstall "kubeception.cloud/kubeception/pkg/apis/kubeception/install"
	"kubeception.cloud/kubeception/pkg/controller"
	"kubeception.cloud/kubeception/pkg/util"
	clusterapis "sigs.k8s.io/cluster-api/pkg/apis"
	clustercontroller "sigs.k8s.io/cluster-api/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func NewManagerCommand(ctx context.Context, logger logr.Logger) *cobra.Command {
	var (
		leaderElection bool
	)

	cmd := &cobra.Command{
		Use:   "kubeception",
		Short: "Kubeception allows creating and managing Kubernetes clusters with Kubernetes itself.",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := config.GetConfigOrDie()
			mgr, err := manager.New(cfg, manager.Options{
				LeaderElection: leaderElection,
			})
			if err != nil {
				util.LogErrorAndExit(logger, err, "Could not initialize manager")
			}

			if err := clusterapis.AddToScheme(mgr.GetScheme()); err != nil {
				util.LogErrorAndExit(logger, err, "Could not modify scheme")
			}

			kubeceptioninstall.Install(mgr.GetScheme())
			certificateinstall.Install(mgr.GetScheme())

			if err := clustercontroller.AddToManager(mgr); err != nil {
				util.LogErrorAndExit(logger, err, "Could add cluster-api controllers")
			}

			if err := controller.AddToManager(mgr); err != nil {
				util.LogErrorAndExit(logger, err, "Could not add controllers")
			}

			if err := mgr.Start(ctx.Done()); err != nil {
				util.LogErrorAndExit(logger, err, "Error running manager")
			}
		},
	}
	cmd.Flags().AddGoFlagSet(flag.CommandLine)
	cmd.Flags().BoolVar(&leaderElection, "leader-election", false, "Whether to do leader")

	return cmd
}
