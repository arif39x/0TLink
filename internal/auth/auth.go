package auth

import (
	"crypto/tls"  // TLS configuration and Handshake
	"crypto/x509" // certificate parsing and validation
	"os"          // read the CA certificate
)

// GetTLSConfig returns a config for tls.Listen or tls.Dial
func GetTLSConfig(certPath, keypath, caPath string, isServer bool) (*tls.Config, error) {
	// certPath: PEM file with certificate, keyPath: PEM file with private key
	cert, err := tls.LoadX509KeyPair(certPath, keypath)
	if err != nil {
		return nil, err
	}

	// Reads CA pem file from disk
	caCert, err := os.ReadFile(caPath)
	if err != nil {
		return nil, err
	}

	// Create an empty certificate trust store - the verification database
	caPool := x509.NewCertPool()
	// Add the CA certificate to the trust pool so we know who to trust
	caPool.AppendCertsFromPEM(caCert)

	config := &tls.Config{
		Certificates: []tls.Certificate{cert}, // Certificates presented during handshake
		RootCAs:      caPool,                  // Used to verify remote certificates
	}

	if isServer {
		// Server-specific: trust CA for inbound connections and mandate client certs
		config.ClientCAs = caPool
		config.ClientAuth = tls.RequireAndVerifyClientCert
	}

	return config, nil
}
