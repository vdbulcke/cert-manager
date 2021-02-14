package data

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"strings"
)

// ParseX509FromPEM return a x509.Certificate from pem input string
func ParseX509FromPEM(rawpem string) (*x509.Certificate, error) {

	certPem, _ := pem.Decode([]byte(rawpem))
	if certPem == nil {
		return nil, errors.New("Error Parsing PEM: ")
	}

	cert, err := x509.ParseCertificate(certPem.Bytes)
	if err != nil {
		return nil, err
	}

	return cert, nil
}

// GetSHA256FingerprintFromX509Cert compute the sha256 fingerprint of X509 Cert
func GetSHA256FingerprintFromX509Cert(cert *x509.Certificate) string {

	fingerprint := sha256.Sum256(cert.Raw)
	return fmt.Sprintf("%x", fingerprint)

}

// GetSubjectFromX509Cert return X509 Cert Subject as string
func GetSubjectFromX509Cert(cert *x509.Certificate) string {
	return cert.Subject.String()
}

// GetIssuerFromX509Cert return X509 Cert Issuer as string
func GetIssuerFromX509Cert(cert *x509.Certificate) string {
	return cert.Issuer.String()
}

// GetSerialNumberFromX509Cert return X509 Cert SerialNumber as string
func GetSerialNumberFromX509Cert(cert *x509.Certificate) string {
	return toHexInt(cert.SerialNumber)
}

// GetSignatureAlgorithmFromX509Cert return X509 Cert SignatureAlgorithm as string
func GetSignatureAlgorithmFromX509Cert(cert *x509.Certificate) string {
	return cert.SignatureAlgorithm.String()
}

// GetAKIFromX509Cert return X509 Cert Authority Key ID as string
func GetAKIFromX509Cert(cert *x509.Certificate) string {
	return formatKeyID(cert.AuthorityKeyId)
}

// GetSKIFromX509Cert return X509 Cert Subject Key ID as string
func GetSKIFromX509Cert(cert *x509.Certificate) string {
	return formatKeyID(cert.SubjectKeyId)
}

// GetSNASFromX509Cert return X509 Certlist of Subject Alternative Names as list of string
func GetSNASFromX509Cert(cert *x509.Certificate) string {

	sans := cert.DNSNames
	// append IP address to SANS
	for _, ip := range cert.IPAddresses {
		sans = append(sans, ip.String())
	}

	return strings.Join(sans, "|")
}

// GetOCSPFromX509Cert return X509 OCSP Server as list of  string
func GetOCSPFromX509Cert(cert *x509.Certificate) string {
	return strings.Join(cert.OCSPServer, "|")
}

// GetIsCAFromX509Cert return if the x509 cert is CA
func GetIsCAFromX509Cert(cert *x509.Certificate) bool {
	return cert.IsCA
}

// GetCRLFromX509Cert return X509 CRL Server as list of  string
func GetCRLFromX509Cert(cert *x509.Certificate) string {
	return strings.Join(cert.CRLDistributionPoints, "|")
}

// GetIssuingCAFromX509Cert return X509 Issueing CA Url as list of  string
func GetIssuingCAFromX509Cert(cert *x509.Certificate) string {
	return strings.Join(cert.IssuingCertificateURL, "|")
}

// from https://github.com/cloudflare/cfssl/blob/master/certinfo/certinfo.go
func formatKeyID(id []byte) string {
	var s string

	for i, c := range id {
		if i > 0 {
			s += ":"
		}
		s += fmt.Sprintf("%02X", c)
	}

	return s
}

func toHexInt(n *big.Int) string {
	return fmt.Sprintf("%x", n) // or %X or upper case
}
