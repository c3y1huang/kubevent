package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	logLevel string
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kubevent",
		Short: "kubevent, a solution for watching resource events then forward to external event brokers",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Setup app config
			logLevel, err := logrus.ParseLevel(logLevel)
			if err != nil {
				logrus.WithField("level", logLevel).Fatalf("invalid log level. Use %s instead\n", logrus.AllLevels)
			}
			logrus.SetLevel(logLevel)
		},
	}

	cmd.PersistentFlags().StringVar(&logLevel, "log-level", logrus.InfoLevel.String(), "log level")

	cmd.AddCommand(
		NewServerCmd(),
		NewVersionCmd(),
		NewSchemaCmd(),
	)

	return cmd
}
