package cmd

import (
	"github.com/innobead/kubevent/internal/config"
	"github.com/innobead/kubevent/pkg/engine"
	"github.com/innobead/kubevent/pkg/handler"
	"github.com/innobead/kubevent/pkg/reconciler"
	"github.com/innobead/kubevent/pkg/util"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/thoas/go-funk"
	"os"
	"path/filepath"
	controllerruntime "sigs.k8s.io/controller-runtime/pkg/handler"
)

var (
	cfgFile string
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kubevent/kubevent.yaml)")
}

func initConfig() {
	if _, err := config.Init(cfgFile); err != nil {
		log.Errorf("Failed to read %s", cfgFile)

		os.Exit(1)
	}
}

func Execute() error {
	return rootCmd.Execute()
}

var rootCmd = &cobra.Command{
	Use:   "kubevent",
	Short: "Kubevent, watch and publish events to external event brokers",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		path, _ := filepath.Abs(cfgFile)

		if info, err := os.Stat(path); os.IsNotExist(err) || info.IsDir() {
			return errors.New("%s does not exist or not file")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		eng, err := engine.New()
		if err != nil {
			return err
		}

		util.RegisterShutdownHook(func() {
			if eng != nil {
				_ = eng.Stop()
			}
		})

		cfg := config.Get()

		watchedApiTypes := funk.Map(cfg.Resources, func(r config.EventResource) string {
			return r.Kind
		}).([]string)

		var eventHandlers []controllerruntime.EventHandler
		for _, s := range cfg.Sinks {
			switch s.Type {
			case "amqp":
				result := config.AmqpSink{}
				if err := mapstructure.Decode(s.Value, &result); err != nil {
					log.Errorf("")
				}
				eventHandlers = append(eventHandlers, handler.NewAmqpEventHandler(result))
			}
		}

		err = eng.CreateController(
			watchedApiTypes,
			eventHandlers,
			&reconciler.DummyReconciler{},
		)
		if err != nil {
			return err
		}

		if err := eng.Start(); err != nil {
			return err
		}

		return nil
	},
}
