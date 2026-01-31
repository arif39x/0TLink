package auth

import (
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"errors"
	"math/big"
	"time"
)

func SignCSR(
	csrBytes []byte,
	expectedCommonName string,
	caCert *x509.Certificate,
	caPriv crypto.Signer,
) ([]byte, error) {

	if !caCert.IsCA {
		return nil, errors.New("provided certificate is not a CA")
	}

	if caCert.KeyUsage&x509.KeyUsageCertSign == 0 {
		return nil, errors.New("CA certificate cannot sign certificates")
	}

	csr, err := x509.ParseCertificateRequest(csrBytes)
	if err != nil {
		return nil, err
	}

	if err := csr.CheckSignature(); err != nil {
		return nil, err
	}

	if csr.Subject.CommonName != expectedCommonName {
		return nil, errors.New("CSR common name not authorized")
	}

	serial, err := rand.Int(
		rand.Reader,
		new(big.Int).Lsh(big.NewInt(1), 128),
	)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	template := &x509.Certificate{
		SerialNumber: serial,
		Subject:      csr.Subject,
		NotBefore:    now.Add(-5 * time.Minute),
		NotAfter:     now.Add(24 * time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}

	certBytes, err := x509.CreateCertificate(
		rand.Reader,
		template,
		caCert,
		csr.PublicKey,
		caPriv,
	)
	if err != nil {
		return nil, err
	}

	return certBytes, nil
}
