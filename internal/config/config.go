package config

import (
	"bytes"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

var (
	cfg Config
)

type Config struct {
	Log             Log
	Offset          Offset
	ReconnectPeriod int
	Resources       []EventResource
	Brokers         []Broker
}

type Log struct {
	Level string
}

type Offset struct {
	Time string
}

type TLS struct {
	InsecureSkipVerify bool
	CaCert             string
	Cert               string
	Key                string
}

type EventResource struct {
	Group   string
	Version string
	Kind    string
}

type Broker struct {
	Type  string // amqp
	Value interface{}
}

type AmqpBroker struct {
	Uri      string
	Exchange string
	Tls      TLS
}

type KafkaBroker struct {
	Servers []string
	Topic   string
	Tls     TLS
}

var defaultConfigFile string

func init() {
	viper.AutomaticEnv()

	configType := "yaml"
	viper.SetConfigType(configType)

	homeDir, _ := os.UserHomeDir()
	configPath := filepath.Join(homeDir, ".kubevent")
	viper.AddConfigPath(configPath)

	configName := "config"
	viper.SetConfigName(configName)

	defaultConfigFile = filepath.Join(configPath, configName) + "." + configType

	viper.SetDefault("ReconnectPeriod", 5)
}

func DefaultConfigFile() string {
	return defaultConfigFile
}

func Init(cfgFile string) (*Config, error) {
	path, _ := filepath.Abs(cfgFile)

	if info, err := os.Stat(path); !os.IsNotExist(err) && !info.IsDir() {
		data, err := ioutil.ReadFile(cfgFile)
		if err != nil {
			return nil, err
		}

		if err := viper.ReadConfig(bytes.NewBuffer(data)); err != nil {
			return nil, err
		}
	} else {
		if err := viper.ReadInConfig(); err != nil {
			return nil, err
		}
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func Get() *Config {
	return &cfg
}

func (receiver *Offset) ParsedTime() (time.Time, error) {
	return time.Parse(time.RFC3339, receiver.Time)
}
