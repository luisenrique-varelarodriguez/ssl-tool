package internal

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config estructura para el archivo de configuración YAML
type Config struct {
	DefaultDomain            string `yaml:"default_domain,omitempty"`            // Dominio (Common Name)
	DefaultCountry           string `yaml:"default_country,omitempty"`           // País
	DefaultLocality          string `yaml:"default_locality,omitempty"`          // Localidad
	DefaultOrganization      string `yaml:"default_organization,omitempty"`      // Organización
	DefaultOrganizationalUnit string `yaml:"default_organizational_unit,omitempty"`
	DefaultEmail             string `yaml:"default_email,omitempty"`             // Correo electrónico
	DefaultKeySize           int    `yaml:"default_key_size,omitempty"`          // Tamaño de clave
}

// GenerateConfigTemplate genera un archivo de configuración YAML predeterminado.
func GenerateConfigTemplate(outputPath string) error {
	defaultTemplate := Config{
		DefaultDomain:            "example.com",
		DefaultCountry:           "US",
		DefaultLocality:          "New York",
		DefaultOrganization:      "DefaultOrg",
		DefaultOrganizationalUnit: "IT",
		DefaultEmail:             "admin@example.com",
		DefaultKeySize:           2048,
	}

	return saveAsYAML(defaultTemplate, outputPath)
}

// saveAsYAML guarda una estructura en formato YAML en el archivo indicado
func saveAsYAML(data interface{}, path string) error {
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshalling to YAML: %v", err)
	}

	// Crear el archivo YAML
	if err := os.WriteFile(path, yamlData, 0644); err != nil {
		return fmt.Errorf("error writing YAML file: %v", err)
	}
	return nil
}

func LoadConfig(path string) (Config, error) {
	var cfg Config
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("error parsing config file: %v", err)
	}
	return cfg, nil
}