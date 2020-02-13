package main

import (
	kubeventv1alpha1 "github.com/innobead/kubevent/api/v1alpha1"
	"github.com/innobead/kubevent/cmd/kubevent/version"
	"github.com/innobead/kubevent/controllers"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
)

var (
	scheme = runtime.NewScheme()

	logLevel             string
	metricsAddr          string
	enableLeaderElection bool
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)
	_ = kubeventv1alpha1.AddToScheme(scheme)
}

func NewKubeventCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kubevent",
		Short: "Kubevent, watch and publish resource events to external event brokers",
		PreRun: func(cmd *cobra.Command, args []string) {
			// Setup app config
			logLevel, err := logrus.ParseLevel(logLevel)
			if err != nil {
				logrus.WithField("level", logLevel).Fatalf("Invalid log level. Use %s instead\n", logrus.AllLevels)
			}
			logrus.SetLevel(logLevel)
		},
		Run: func(cmd *cobra.Command, args []string) {
			logrus.Infoln("creating an event broker controller manager")
			mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
				Scheme:             scheme,
				MetricsBindAddress: metricsAddr,
				LeaderElection:     enableLeaderElection,
			})
			if err != nil {
				logrus.WithError(err).Fatalln("unable to start the manager")
			}

			logrus.Infoln("creating an event broker controller for resource state reconciliation")
			if err = (&controllers.BrokerReconciler{
				Client: mgr.GetClient(),
				Scheme: mgr.GetScheme(),
				Mgr:    mgr,
			}).SetupWithManager(mgr); err != nil {
				logrus.WithError(err).Fatalln("unable to create the event broker controller")
			}

			logrus.Infoln("starting the manager")
			if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
				logrus.WithError(err).Fatalln("problem running manager")
			}
		},
	}

	cmd.Flags().StringVar(&logLevel, "log_level", string(logrus.InfoLevel), "log level")
	cmd.Flags().StringVar(&metricsAddr, "metricsAddr", "localhost:9000", "metrics address")
	cmd.Flags().BoolVar(&enableLeaderElection, "enableLeaderElection", true, "enable leader election")

	cmd.AddCommand(
		version.NewVersionCmd(),
	)

	return cmd
}
