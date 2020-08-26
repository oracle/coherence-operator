/*
 * Copyright (c) 2019, 2020 Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package certs

import (
	"bytes"
	cryptorand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"github.com/oracle/coherence-operator/pkg/operator"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"
	"math/big"
	ctrl "sigs.k8s.io/controller-runtime"
	"strings"
	"time"
)

var (
	log = ctrl.Log.WithName("certificates")

	// SerialNumberLimit is the maximum number used as a certificate serial number
	SerialNumberLimit = new(big.Int).Lsh(big.NewInt(1), 128)
)

// CA is a simple certificate authority
type CA struct {
	// PrivateKey is the CA private key
	PrivateKey *rsa.PrivateKey
	// Cert is the certificate used to issue new certificates
	Cert *x509.Certificate
}

func (c *CA) PopulateSecret(secret *corev1.Secret) {
	secret.Data = map[string][]byte{
		operator.CertFileName: encodePEMCert(c.Cert.Raw),
		operator.KeyFileName:  encodePEMPrivateKey(c.PrivateKey),
	}
}

func (c *CA) ShouldRenew(rotateBefore time.Duration) bool {
	// Read the current certificate used by the server
	if c == nil {
		return true
	}
	// check whether the certs have expired
	if !CanReuseCA(c, rotateBefore) {
		return true
	}
	// check the DNS names
	if !c.hasCorrectDNS() {
		return true
	}
	return false
}

func (c *CA) hasCorrectDNS() bool {
	dns := operator.GetWebhookServiceDNSNames()
	for _, name := range c.Cert.DNSNames {
		if strings.HasPrefix(name, dns[0]) {
			return true
		}
	}
	return false
}


// WebhookCertificates holds the artifacts used by the webhook server and the webhook configuration.
type WebhookCertificates struct {
	CaCert []byte
	ServerKey  []byte
	ServerCert []byte
}

func CreateSelfSignedCA() (*CA, error) {
	// generate a serial number
	serial, err := cryptorand.Int(cryptorand.Reader, SerialNumberLimit)
	if err != nil {
		return nil, err
	}

	privateKey, err := rsa.GenerateKey(cryptorand.Reader, 2048)
	if err != nil {
		return nil, errors.Wrap(err, "unable to generate the private key")
	}

	expireIn :=  viper.GetDuration(operator.FlagCACertValidity)
	notAfter := time.Now().Add(expireIn)

	subject := pkix.Name{
		CommonName:         "coherence-webhook-ca",
		OrganizationalUnit: []string{"coherence-webhook"},
	}

	certTemplate := x509.Certificate{
		SerialNumber:          serial,
		Subject:               subject,
		NotBefore:             time.Now().Add(-10 * time.Minute),
		NotAfter:              notAfter,
		SignatureAlgorithm:    x509.SHA256WithRSA,
		IsCA:                  true,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		DNSNames:              operator.GetWebhookServiceDNSNames(),
	}

	certData, err := x509.CreateCertificate(cryptorand.Reader, &certTemplate, &certTemplate, privateKey.Public(), privateKey)
	if err != nil {
		return nil, err
	}

	cert, err := x509.ParseCertificate(certData)
	if err != nil {
		return nil, err
	}

	return &CA{PrivateKey: privateKey, Cert: cert}, nil
}

// BuildCAFromSecret parses the given secret into a CA.
// It returns nil if the secrets could not be parsed into a CA.
func BuildCAFromSecret(caInternalSecret corev1.Secret) *CA {
	if caInternalSecret.Data == nil {
		return nil
	}
	caBytes, exists := caInternalSecret.Data[operator.CertFileName]
	if !exists || len(caBytes) == 0 {
		return nil
	}
	certs, err := parsePEMCerts(caBytes)
	if err != nil {
		log.Error(err, "Cannot parse PEM cert from CA secret, will create a new one", "namespace", caInternalSecret.Namespace, "secret_name", caInternalSecret.Name)
		return nil
	}
	if len(certs) == 0 {
		return nil
	}
	if len(certs) > 1 {
		log.Info(
			"More than 1 certificate in the CA secret, continuing with the first one",
			"namespace", caInternalSecret.Namespace,
			"secret_name", caInternalSecret.Name,
		)
	}
	cert := certs[0]

	privateKeyBytes, exists := caInternalSecret.Data[operator.KeyFileName]
	if !exists || len(privateKeyBytes) == 0 {
		return nil
	}
	privateKey, err := parsePEMPrivateKey(privateKeyBytes)
	if err != nil {
		log.Error(err, "Cannot parse PEM private key from CA secret, will create a new one", "namespace", caInternalSecret.Namespace, "secret_name", caInternalSecret.Name)
		return nil
	}

	return &CA{
		PrivateKey: privateKey,
		Cert:       cert,
	}
}

// encodePEMCert encodes the given certificate blocks as a PEM certificate
func encodePEMCert(certBlocks ...[]byte) []byte {
	var buf bytes.Buffer
	for _, block := range certBlocks {
		_, _ = buf.Write(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: block}))
	}
	return buf.Bytes()
}

// encodePEMPrivateKey encodes the given private key in the PEM format
func encodePEMPrivateKey(privateKey *rsa.PrivateKey) []byte {
	return pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
}

func parsePEMCerts(pemData []byte) ([]*x509.Certificate, error) {
	certs := []*x509.Certificate{}
	for len(pemData) > 0 {
		var block *pem.Block
		block, pemData = pem.Decode(pemData)
		if block == nil {
			break
		}
		if block.Type != "CERTIFICATE" || len(block.Headers) != 0 {
			continue
		}

		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		certs = append(certs, cert)
	}
	return certs, nil
}

func parsePEMPrivateKey(pemData []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing private key")
	}

	switch {
	case block.Type == "PRIVATE KEY":
		return parsePKCS8PrivateKey(block.Bytes)
	case block.Type == "RSA PRIVATE KEY" && len(block.Headers) == 0:
		return x509.ParsePKCS1PrivateKey(block.Bytes)
	default:
		return nil, errors.New("expected PEM block to contain an RSA private key")
	}
}

func parsePKCS8PrivateKey(block []byte) (*rsa.PrivateKey, error) {
	key, err := x509.ParsePKCS8PrivateKey(block)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse private key")
	}

	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.Errorf("expected an RSA private key but got %t", key)
	}

	return rsaKey, nil
}

// RotateIn determines when a cert should be rotated
func RotateIn(now time.Time, certExpiration time.Time, certRotateBefore time.Duration) time.Duration {
	// make sure we are past the safety margin when rotating, by making it a little bit shorter
	safetyMargin := certRotateBefore - 1*time.Second
	requeueTime := certExpiration.Add(-safetyMargin)
	requeueIn := requeueTime.Sub(now)
	if requeueIn < 0 {
		// requeue asap
		requeueIn = 0
	}
	return requeueIn
}

func CanReuseCA(ca *CA, expirationSafetyMargin time.Duration) bool {
	return privateMatchesPublicKey(ca.Cert.PublicKey, *ca.PrivateKey) && CertIsValid(*ca.Cert, expirationSafetyMargin)
}

func CertIsValid(cert x509.Certificate, expirationSafetyMargin time.Duration) bool {
	now := time.Now()
	if now.Before(cert.NotBefore) {
		log.Info("CA cert is not valid yet", "subject", cert.Subject)
		return false
	}
	if now.After(cert.NotAfter.Add(-expirationSafetyMargin)) {
		log.Info("CA cert expired or soon to expire", "subject", cert.Subject, "expiration", cert.NotAfter)
		return false
	}
	return true
}

func privateMatchesPublicKey(publicKey interface{}, privateKey rsa.PrivateKey) bool {
	pubKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		log.Error(errors.New("Public key is not an RSA public key"), "")
		return false
	}
	// check that public and private keys share the same modulus and exponent
	if pubKey.N.Cmp(privateKey.N) != 0 || pubKey.E != privateKey.E {
		return false
	}
	return true
}

