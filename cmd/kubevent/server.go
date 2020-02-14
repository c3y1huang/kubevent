package main

import (
	kubeventv1alpha1 "github.com/innobead/kubevent/api/v1alpha1"
	"github.com/innobead/kubevent/controllers"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
)

var (
	scheme      = runtime.NewScheme()
	metricsAddr string
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)
	_ = kubeventv1alpha1.AddToScheme(scheme)
}

func NewServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Server, run event broker controller",
		Run: func(cmd *cobra.Command, args []string) {
			logrus.Infoln("creating event broker controller manager")

			mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
				Scheme:             scheme,
				MetricsBindAddress: metricsAddr,
			})
			if err != nil {
				logrus.WithError(err).Fatalln("unable to start manager")
			}

			logrus.Infoln("creating event broker controller for managing event broker watcher for supported event brokers")

			if _, err := controllers.NewEventBrokerController(mgr); err != nil {
				logrus.WithError(err).Fatalln("unable to create event broker controller")
			}

			logrus.Infoln("starting the manager")

			if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
				logrus.WithError(err).Fatalln("failed to run event broker controller manager")
			}
		},
	}

	cmd.Flags().StringVar(&metricsAddr, "metricsAddr", "localhost:9000", "Metrics address")

	return cmd
}
