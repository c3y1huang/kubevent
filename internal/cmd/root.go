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

		// Setup app config
		if level, err := log.ParseLevel(cfg.Log.Level); err != nil {
			log.WithField("level", cfg.Log.Level).Warnf("Invalid log level, use %s instead", log.InfoLevel)
			log.SetLevel(log.InfoLevel)
		} else {
			log.SetLevel(level)
		}

		// Prepare watched resources
		watchedApiTypes := funk.Map(cfg.Resources, func(r config.EventResource) string {
			return r.Kind
		}).([]string)

		// Prepare event handler
		var eventHandlers []controllerruntime.EventHandler
		for _, s := range cfg.Brokers {
			switch s.Type {
			case "amqp":
				broker := config.AmqpBroker{}
				if err := mapstructure.Decode(s.Value, &broker); err != nil {
					log.Errorf("")
				}
				eventHandlers = append(eventHandlers, handler.NewAmqp(broker))

			case "kafka":
				broker := config.KafkaBroker{}
				if err := mapstructure.Decode(s.Value, &broker); err != nil {
					log.Errorf("")
				}
				eventHandlers = append(eventHandlers, handler.NewKafka(broker))

			}
		}

		// Prepare offset events start to receive
		var predictTime time.Time
		if t, err := cfg.Offset.ParsedTime(); err != nil {
			predictTime = time.Now()
			log.WithField("time", cfg.Offset.Time).Warnf("Invalid offset time, set current time (%s) instead", predictTime.Format(time.RFC3339))
		} else {
			predictTime = t
		}

		// Create controller
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

		// Start event handler
		go func() {
			for {
				for _, h := range eventHandlers {
					op := h.(handler.Operation)

					if err := op.Start(); err != nil {
						log.Warnf("Failed to connect broker, %v", err)
						_ = op.Stop()
					}
				}

				time.Sleep(time.Second * time.Duration(cfg.ReconnectPeriod))
			}

		}()

		defer func() {
			for _, h := range eventHandlers {
				op := h.(handler.Operation)

				if err := op.Stop(); err != nil {
					log.Warnf("Failed to disconnect broker, %v", err)
				}
			}
		}()

		// Start controller engine
		if err := eng.Start(); err != nil {
			return err
		}

		return nil
	},
}
