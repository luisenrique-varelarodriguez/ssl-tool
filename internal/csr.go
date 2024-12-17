package internal

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func ValidateCSRParams(domain, country, locality, organization string) error {
	if domain == "" {
		return errors.New("domain cannot be empty")
	}
	if len(country) != 2 {
		return errors.New("country must be 2 letters")
	}
	matched, _ := regexp.MatchString("^[A-Za-z]{2}$", country)
	if !matched {
		return errors.New("country must be alphabetic 2-letter code")
	}
	if strings.TrimSpace(locality) == "" {
		return errors.New("locality cannot be empty")
	}
	if strings.TrimSpace(organization) == "" {
		return errors.New("organization cannot be empty")
	}
	return nil
}

func GenerateCSR(domain, country, locality, organization string) error {
	dirName := strings.ReplaceAll(domain, ".", "_")
	if err := os.MkdirAll(dirName, 0755); err != nil {
		return fmt.Errorf("error creating directory: %v", err)
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("error generating private key: %v", err)
	}

	keyFilePath := filepath.Join(dirName, fmt.Sprintf("%s.key", dirName))
	keyFile, err := os.Create(keyFilePath)
	if err != nil {
		return fmt.Errorf("error creating key file: %v", err)
	}
	defer keyFile.Close()

	if err := pem.Encode(keyFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}); err != nil {
		return fmt.Errorf("error writing private key: %v", err)
	}

	subject := pkix.Name{
		CommonName:   domain,
		Country:      []string{country},
		Locality:     []string{locality},
		Organization: []string{organization},
	}

	csrTemplate := &x509.CertificateRequest{
		Subject:            subject,
		SignatureAlgorithm: x509.SHA256WithRSA,
	}

	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, csrTemplate, privateKey)
	if err != nil {
		return fmt.Errorf("error creating CSR: %v", err)
	}

	csrFilePath := filepath.Join(dirName, fmt.Sprintf("%s.csr", dirName))
	csrFile, err := os.Create(csrFilePath)
	if err != nil {
		return fmt.Errorf("error creating CSR file: %v", err)
	}
	defer csrFile.Close()

	if err := pem.Encode(csrFile, &pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrBytes}); err != nil {
		return fmt.Errorf("error writing CSR: %v", err)
	}

	fmt.Printf("Files generated successfully:\n- Private Key: %s\n- CSR: %s\n", keyFilePath, csrFilePath)
	return nil
}