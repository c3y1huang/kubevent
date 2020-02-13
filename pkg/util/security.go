package util

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"github.com/innobead/kubevent/api/v1alpha1"
	"io/ioutil"
)

func CreateTLSConfig(cfg *v1alpha1.TlsConfig) (*tls.Config, error) {
	if cfg.Insecure {
		return &tls.Config{
			InsecureSkipVerify: true,
		}, nil
	}

	if cfg.CACert == "" || cfg.ClientCert == "" || cfg.ClientKey == "" {
		return nil, errors.New("missing certificates or key info")
	}

	if cer, err := tls.LoadX509KeyPair(cfg.ClientCert, cfg.ClientKey); err != nil {
		return nil, err
	} else {
		caCert, err := ioutil.ReadFile(cfg.CACert)
		if err != nil {
			return nil, err
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		tlsConfig := tls.Config{Certificates: []tls.Certificate{cer}}
		tlsConfig.RootCAs = caCertPool

		return &tlsConfig, nil
	}
}
