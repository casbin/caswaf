package object

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"time"
)

func getCertificateExpiry(certificatePEM string) (time.Time, error) {
	block, _ := pem.Decode([]byte(certificatePEM))
	if block == nil || block.Type != "CERTIFICATE" {
		return time.Time{}, errors.New("failed to decode PEM block containing certificate")
	}

	certificate, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return time.Time{}, err
	}

	return certificate.NotAfter, nil
}

func getCertExpireTime(certificate string) string {
	t, err := getCertificateExpiry(certificate)
	if err != nil {
		panic(err)
	}

	return t.Local().Format(time.RFC3339)
}
