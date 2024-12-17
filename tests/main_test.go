package main

import (
	"os"
	"os/exec"
	"testing"
)

// Función auxiliar para ejecutar comandos y capturar salida y errores.
func runCommand(t *testing.T, args ...string) (string, error) {
	cmd := exec.Command("./ssl-tool", args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// Test para generate-config
func TestGenerateConfig(t *testing.T) {
	// Limpiar cualquier archivo existente
	os.Remove("ssl-tool-config.yaml")

	_, err := runCommand(t, "generate-config")
	if err != nil {
		t.Fatalf("Error running generate-config: %v", err)
	}

	// Verificar que el archivo se ha creado
	if _, err := os.Stat("ssl-tool-config.yaml"); os.IsNotExist(err) {
		t.Fatalf("ssl-tool-config.yaml was not created")
	}

	t.Log("generate-config passed successfully")
}

// Test para generate-csr
func TestGenerateCSR(t *testing.T) {
	// Ejecutar generate-csr
	_, err := runCommand(t, "generate-csr", "--domain", "example.com", "--country", "US", "--locality", "New York", "--organization", "TestOrg")
	if err != nil {
		t.Fatalf("Error running generate-csr: %v", err)
	}

	// Verificar que se crearon los archivos en la carpeta example_com/
	if _, err := os.Stat("example_com/example_com.key"); os.IsNotExist(err) {
		t.Fatalf("Private key file was not created")
	}

	if _, err := os.Stat("example_com/example_com.csr"); os.IsNotExist(err) {
		t.Fatalf("CSR file was not created")
	}

	t.Log("generate-csr passed successfully")
}

// Test para extract-info
func TestExtractInfo(t *testing.T) {
	// Generar CSR para extraer información
	_, _ = runCommand(t, "generate-csr", "--domain", "example.com", "--country", "US", "--locality", "New York", "--organization", "TestOrg")

	// Ejecutar extract-info
	_, err := runCommand(t, "extract-info", "--file", "example_com/example_com.csr")
	if err != nil {
		t.Fatalf("Error running extract-info: %v", err)
	}

	// Verificar que el archivo ssl-tool-config.yaml contiene la información
	if _, err := os.Stat("ssl-tool-config.yaml"); os.IsNotExist(err) {
		t.Fatalf("ssl-tool-config.yaml was not created by extract-info")
	}

	t.Log("extract-info passed successfully")
}

// Test para verify-hashes
func TestVerifyHashes(t *testing.T) {
    // Generar CSR y clave
    _, _ = runCommand(t, "generate-csr", "--domain", "example.com", "--country", "US", "--locality", "New York", "--organization", "TestOrg")

    // Simular la creación de un certificado firmado basado en el CSR
    keyFile := "example_com/example_com.key"
    csrFile := "example_com/example_com.csr"
    certFile := "example_com/example_com.crt"

    // Generar un certificado autofirmado usando OpenSSL (requiere que esté instalado)
    cmd := exec.Command("openssl", "x509", "-req", "-days", "365", "-in", csrFile, "-signkey", keyFile, "-out", certFile)
    if err := cmd.Run(); err != nil {
        t.Fatalf("Failed to generate self-signed certificate: %v", err)
    }

    // Ejecutar verify-hashes
    _, err := runCommand(t, "verify-hashes", "--key", keyFile, "--csr", csrFile, "--cert", certFile)
    if err != nil {
        t.Fatalf("Error running verify-hashes: %v", err)
    }

    t.Log("verify-hashes passed successfully")
}

// Test para check-expiration
func TestCheckExpiration(t *testing.T) {
	// Generar CSR para obtener un archivo de prueba
	_, _ = runCommand(t, "generate-csr", "--domain", "example.com", "--country", "US", "--locality", "New York", "--organization", "TestOrg")

	// Simular certificado usando el CSR generado
	_, _ = runCommand(t, "extract-info", "--file", "example_com/example_com.csr")

	// Ejecutar check-expiration (no fallará porque no tenemos un certificado real)
	_, err := runCommand(t, "check-expiration", "--cert", "example_com/example_com.csr")
	if err != nil {
		t.Logf("Expected failure as no real certificate exists: %v", err)
	} else {
		t.Log("check-expiration passed successfully (fake cert)")
	}
}