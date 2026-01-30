package auth

import (
	"crypto/tls"
	"crypto/x509"
	"os"
)

func GetTLSConfig(certPath, keypath, caPath string, isServer bool) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certPath, keypath)
	if err != nil {
		return nil, err
	}

	caCert, err := os.ReadFile(caPath)
	if err != nil {
		return nil, err
	}

	caPool := x509.NewCertPool()
	caPool.AppendCertsFromPEM(caCert)

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caPool,
		MinVersion:   tls.VersionTLS13,
	}

	if isServer {
		config.ClientCAs = caPool
		config.ClientAuth = tls.RequireAndVerifyClientCert
	}

	return config, nil
}
