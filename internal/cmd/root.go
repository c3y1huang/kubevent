package cmd

import (
	"github.com/innobead/kubevent/internal/config"
	"github.com/innobead/kubevent/pkg/engine"
	"github.com/innobead/kubevent/pkg/handler"
	"github.com/innobead/kubevent/pkg/predicater"
	"github.com/innobead/kubevent/pkg/reconciler"
	"github.com/innobead/kubevent/pkg/util"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/thoas/go-funk"
	controllerruntime "sigs.k8s.io/controller-runtime/pkg/handler"
	"time"
)

var (
	cfgFile string
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", config.DefaultConfigFile(), "config file")
}

func initConfig() {
	if _, err := config.Init(cfgFile); err != nil {
		log.Fatalf("Failed to read '%s' or ", cfgFile)
	}
}

func Execute() error {
	return rootCmd.Execute()
}

var rootCmd = &cobra.Command{
	Use:   "kubevent",
	Short: "Kubevent, watch and publish events to external event brokers",
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

		if level, err := log.ParseLevel(cfg.Log.Level); err != nil {
			log.WithField("level", cfg.Log.Level).Warnf("Invalid log level, use %s instead", log.InfoLevel)
			log.SetLevel(log.InfoLevel)
		} else {
			log.SetLevel(level)
		}

		watchedApiTypes := funk.Map(cfg.Resources, func(r config.EventResource) string {
			return r.Kind
		}).([]string)

		var eventHandlers []controllerruntime.EventHandler
		for _, s := range cfg.Brokers {
			switch s.Type {
			case "amqp":
				result := config.AmqpBroker{}
				if err := mapstructure.Decode(s.Value, &result); err != nil {
					log.Errorf("")
				}
				eventHandlers = append(eventHandlers, handler.NewAmqp(result))
			}
		}

		var predictTime time.Time
		if t, err := cfg.Offset.ParsedTime(); err != nil {
			predictTime = time.Now()
			log.WithField("time", cfg.Offset.Time).Warnf("Invalid offset time, set current time (%s) instead", predictTime.Format(time.RFC3339))
		} else {
			predictTime = t
		}

		err = eng.CreateController(
			"kubevent",
			watchedApiTypes,
			eventHandlers,
			predicater.NewTime(predictTime),
			reconciler.NewDummy(),
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
