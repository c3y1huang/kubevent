package util

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/innobead/kubevent/internal/config"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
)

func RegisterShutdownHook(callback func()) {
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-stopCh

		log.Infof("Receiving shutdown hook")
		callback()
	}()
}

func CreateTLSConfig(cfg config.TLS) (*tls.Config, error) {
	if cfg.InsecureSkipVerify == false && (cfg.Cert == "" || cfg.Key == "") {
		return nil, errors.New("no tls")
	}

	if cfg.InsecureSkipVerify {
		return &tls.Config{
			InsecureSkipVerify: true,
		}, nil
	}

	cer, err := tls.LoadX509KeyPair(cfg.Cert, cfg.Key)
	if err != nil {
		return nil, err
	}

	tlsCfg := tls.Config{Certificates: []tls.Certificate{cer}}

	if cfg.CaCert != "" {
		caCert, err := ioutil.ReadFile(cfg.CaCert)
		if err != nil {
			return nil, err
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		tlsCfg.RootCAs = caCertPool
	}

	return &tlsCfg, nil
}
