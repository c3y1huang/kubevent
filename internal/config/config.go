package config

import (
	"bytes"
	"github.com/spf13/viper"
	"io/ioutil"
)

var (
	cfg Config
)

type Config struct {
	Resources []EventResource
	Sinks     []Sink
}

type EventResource struct {
	Group   string
	Version string
	Kind    string
}

type Sink struct {
	Type  string // amqp
	Value interface{}
}

type AmqpSink struct {
	Uri      string
	Exchange string
}

func init() {
	viper.AutomaticEnv()
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.kubevent")
}

func Init(cfgFile string) (*Config, error) {
	data, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return nil, err
	}

	if err := viper.ReadConfig(bytes.NewBuffer(data)); err != nil {
		return nil, err
	}

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func Get() *Config {
	return &cfg
}
