package internal

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"os"
	"time"
)

// ExtractInfo extrae información de un CRT o CSR y la guarda en la estructura de configuración YAML.
func ExtractInfo(filePath, outputPath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return fmt.Errorf("failed to parse PEM block")
	}

	config := Config{}

	switch block.Type {
	case "CERTIFICATE":
		// Procesar archivo de certificado (CRT)
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return fmt.Errorf("error parsing certificate: %v", err)
		}
		config.DefaultDomain = cert.Subject.CommonName
		config.DefaultCountry = firstOrEmpty(cert.Subject.Country)
		config.DefaultLocality = firstOrEmpty(cert.Subject.Locality)
		config.DefaultOrganization = firstOrEmpty(cert.Subject.Organization)
		config.DefaultEmail = firstOrEmpty(cert.EmailAddresses)
		config.DefaultKeySize = 2048 // Valor fijo.

	case "CERTIFICATE REQUEST":
		// Procesar archivo CSR
		csr, err := x509.ParseCertificateRequest(block.Bytes)
		if err != nil {
			return fmt.Errorf("error parsing CSR: %v", err)
		}
		config.DefaultDomain = csr.Subject.CommonName
		config.DefaultCountry = firstOrEmpty(csr.Subject.Country)
		config.DefaultLocality = firstOrEmpty(csr.Subject.Locality)
		config.DefaultOrganization = firstOrEmpty(csr.Subject.Organization)
		config.DefaultKeySize = 2048 // Valor fijo.

	default:
		return fmt.Errorf("unsupported PEM type: %s", block.Type)
	}

	// Guardar la información extraída en YAML
	if err := saveAsYAML(config, outputPath); err != nil {
		return err
	}

	fmt.Printf("Information saved to %s\n", outputPath)
	return nil
}

// firstOrEmpty devuelve el primer elemento de un slice o una cadena vacía si está vacío
func firstOrEmpty(values []string) string {
	if len(values) > 0 {
		return values[0]
	}
	return ""
}

func printCertificateInfo(cert *x509.Certificate) {
	fmt.Println("Certificate Info:")
	fmt.Printf("- Common Name: %s\n", cert.Subject.CommonName)
	fmt.Printf("- Organization: %v\n", cert.Subject.Organization)
	fmt.Printf("- Locality: %v\n", cert.Subject.Locality)
	fmt.Printf("- Country: %v\n", cert.Subject.Country)
	fmt.Printf("- Issuer: %s\n", cert.Issuer.CommonName)
	fmt.Printf("- Valid From: %s\n", cert.NotBefore)
	fmt.Printf("- Valid To: %s\n", cert.NotAfter)
}

func printCSRInfo(csr *x509.CertificateRequest) {
	fmt.Println("CSR Info:")
	fmt.Printf("- Common Name: %s\n", csr.Subject.CommonName)
	fmt.Printf("- Organization: %v\n", csr.Subject.Organization)
	fmt.Printf("- Locality: %v\n", csr.Subject.Locality)
	fmt.Printf("- Country: %v\n", csr.Subject.Country)
}

func DaysUntilExpiration(certFile string) (int, error) {
	data, err := os.ReadFile(certFile)
	if err != nil {
		return 0, fmt.Errorf("error reading certificate: %v", err)
	}
	block, _ := pem.Decode(data)
	if block == nil || block.Type != "CERTIFICATE" {
		return 0, fmt.Errorf("file is not a valid certificate")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return 0, fmt.Errorf("error parsing certificate: %v", err)
	}
	days := int(time.Until(cert.NotAfter).Hours() / 24)
	return days, nil
}

func CertificateFingerprint(certFile string) (string, error) {
	data, err := os.ReadFile(certFile)
	if err != nil {
		return "", fmt.Errorf("error reading certificate: %v", err)
	}
	block, _ := pem.Decode(data)
	if block == nil || block.Type != "CERTIFICATE" {
		return "", fmt.Errorf("file is not a valid certificate")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("error parsing certificate: %v", err)
	}
	hash := sha256.Sum256(cert.Raw)
	return hex.EncodeToString(hash[:]), nil
}

// VerifyHashes asume que está implementada según el código original
func VerifyHashes(keyFile, csrFile, certFile string) error {
	keyHash, err := calculateModulusHash(keyFile, "RSA PRIVATE KEY")
	if err != nil {
		return err
	}
	csrHash, err := calculateModulusHash(csrFile, "CERTIFICATE REQUEST")
	if err != nil {
		return err
	}
	certHash, err := calculateModulusHash(certFile, "CERTIFICATE")
	if err != nil {
		return err
	}

	if keyHash == csrHash && csrHash == certHash {
		fmt.Println("Hashes match! The private key, CSR, and certificate are consistent.")
	} else {
		fmt.Printf("Hashes do not match:\n- Key: %s\n- CSR: %s\n- Certificate: %s\n", keyHash, csrHash, certHash)
	}

	return nil
}

func calculateModulusHash(filePath, pemType string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("error reading file: %v", err)
	}

	block, _ := pem.Decode(data)
	if block == nil || block.Type != pemType {
		return "", fmt.Errorf("file is not a valid %s", pemType)
	}

	hash := md5.Sum(block.Bytes)
	return hex.EncodeToString(hash[:]), nil
}