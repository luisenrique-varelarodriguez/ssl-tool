package main

import (
    "errors"
    "fmt"
    "os"
    "strings"

    "github.com/spf13/cobra"
    "github.com/luisenrique-varelarodriguez/ssl-tool/internal"
)

var (
    filePath     string
    certFile     string
    keyFile      string
    csrFile      string
    outputDir    string
    domain       string
    country      string
    locality     string
    organization string
    configPath   string
    interactive  bool
    config       internal.Config
)

func main() {
    rootCmd := &cobra.Command{
        Use:   "ssl-tool",
        Short: "SSL Tool is a CLI for managing SSL certificates",
        PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
            // Cargar config si existe
            if _, err := os.Stat(configPath); err == nil {
                cfg, err := internal.LoadConfig(configPath)
                if err != nil {
                    return fmt.Errorf("error loading config: %v", err)
                }
                config = cfg
            }
            return nil
        },
        RunE: func(cmd *cobra.Command, args []string) error {
            return cmd.Help()
        },
    }

    rootCmd.CompletionOptions.DisableDefaultCmd = true

    rootCmd.PersistentFlags().BoolVarP(&interactive, "interactive", "i", false, "Enable interactive mode")
    configPath = "ssl-tool-config.yaml"
    rootCmd.PersistentFlags().StringVar(&configPath, "config", "ssl-tool-config.yaml", "Path to the configuration file")

    // Comando: generate-config
    generateConfigCmd := &cobra.Command{
        Use:   "generate-config",
        Short: "Generate a default YAML configuration file",
        RunE: func(cmd *cobra.Command, args []string) error {
            return internal.GenerateConfigTemplate(configPath)
        },
    }
    generateConfigCmd.Flags().StringVar(&configPath, "output", "ssl-tool-config.yaml", "Path to save the configuration file")

    // Comando: generate-csr
    generateCSRCmd := &cobra.Command{
        Use:   "generate-csr",
        Short: "Generate a new CSR and private key",
        RunE: func(cmd *cobra.Command, args []string) error {
            // Usar defaults de config si faltan
            if country == "" && config.DefaultCountry != "" {
                country = config.DefaultCountry
            }
            if locality == "" && config.DefaultLocality != "" {
                locality = config.DefaultLocality
            }
            if organization == "" && config.DefaultOrganization != "" {
                organization = config.DefaultOrganization
            }

            requiredParams := []string{"domain", "country", "locality", "organization"}

            if interactive {
                // En modo interactivo, preguntar por todos los datos
                domain = promptFor("Domain", domain)
                country = promptFor("Country (2 letters)", country)
                locality = promptFor("Locality (City)", locality)
                organization = promptFor("Organization", organization)
            } else {
                // Modo no interactivo: deben estar todos los datos
                for _, p := range requiredParams {
                    val := getParamValue(p)
                    if val == "" {
                        return fmt.Errorf("missing required parameter: %s. Provide flags or use --interactive", p)
                    }
                }
            }

            if err := internal.ValidateCSRParams(domain, country, locality, organization); err != nil {
                return err
            }

            return internal.GenerateCSR(domain, country, locality, organization)
        },
    }
    generateCSRCmd.Flags().StringVar(&domain, "domain", "", "Domain name for the CSR")
    generateCSRCmd.Flags().StringVar(&country, "country", "", "Country (2 letters)")
    generateCSRCmd.Flags().StringVar(&locality, "locality", "", "Locality (City)")
    generateCSRCmd.Flags().StringVar(&organization, "organization", "", "Organization")

    // Comando: extract-info
    extractInfoCmd := &cobra.Command{
        Use:   "extract-info",
        Short: "Extract information from a CRT or CSR and save it to ssl-tool-config.yaml",
        RunE: func(cmd *cobra.Command, args []string) error {
            if interactive {
                filePath = promptFor("Path to CRT or CSR file", filePath)
            } else if filePath == "" {
                return errors.New("please provide --file or use --interactive")
            }

            if filePath == "" || !fileExists(filePath) {
                return fmt.Errorf("file does not exist: %s", filePath)
            }

            // Utiliza ssl-tool-config.yaml como destino
            outputPath := "ssl-tool-config.yaml"
            return internal.ExtractInfo(filePath, outputPath)
        },
    }
    extractInfoCmd.Flags().StringVar(&filePath, "file", "", "Path to the CRT or CSR file")

    // Comando: check-expiration
    checkExpirationCmd := &cobra.Command{
        Use:   "check-expiration",
        Short: "Check how many days until a certificate expires",
        RunE: func(cmd *cobra.Command, args []string) error {
            if interactive {
                certFile = promptFor("Path to certificate (.crt)", certFile)
            } else {
                if certFile == "" {
                    return errors.New("please provide --cert or use --interactive")
                }
            }

            if certFile == "" {
                return errors.New("certificate path cannot be empty")
            }

            if !fileExists(certFile) {
                return fmt.Errorf("certificate file does not exist: %s", certFile)
            }

            days, err := internal.DaysUntilExpiration(certFile)
            if err != nil {
                return err
            }
            fmt.Printf("Certificate expires in %d days.\n", days)
            return nil
        },
    }
    checkExpirationCmd.Flags().StringVar(&certFile, "cert", "", "Path to the certificate file")

    // Comando: fingerprint
    fingerprintCmd := &cobra.Command{
        Use:   "fingerprint",
        Short: "Show the SHA256 fingerprint of a certificate",
        RunE: func(cmd *cobra.Command, args []string) error {
            if interactive {
                certFile = promptFor("Path to certificate (.crt)", certFile)
            } else {
                if certFile == "" {
                    return errors.New("please provide --cert or use --interactive")
                }
            }

            if certFile == "" {
                return errors.New("certificate path cannot be empty")
            }

            if !fileExists(certFile) {
                return fmt.Errorf("certificate file does not exist: %s", certFile)
            }

            fp, err := internal.CertificateFingerprint(certFile)
            if err != nil {
                return err
            }
            fmt.Printf("SHA256 Fingerprint: %s\n", fp)
            return nil
        },
    }
    fingerprintCmd.Flags().StringVar(&certFile, "cert", "", "Path to the certificate")

    // Comando: verify-hashes
    verifyHashesCmd := &cobra.Command{
        Use:   "verify-hashes",
        Short: "Verify that the private key, CSR, and certificate hashes match",
        RunE: func(cmd *cobra.Command, args []string) error {
            if interactive {
                keyFile = promptFor("Path to private key (.key)", keyFile)
                csrFile = promptFor("Path to CSR (.csr)", csrFile)
                certFile = promptFor("Path to certificate (.crt)", certFile)
            } else {
                // No interactivo: deben estar todos
                if keyFile == "" || csrFile == "" || certFile == "" {
                    return errors.New("missing required parameters: --key, --csr, --cert. Provide flags or use --interactive")
                }
            }

            if keyFile == "" || csrFile == "" || certFile == "" {
                return errors.New("key, csr, and cert cannot be empty in interactive mode")
            }

            if !fileExists(keyFile) {
                return fmt.Errorf("key file does not exist: %s", keyFile)
            }
            if !fileExists(csrFile) {
                return fmt.Errorf("csr file does not exist: %s", csrFile)
            }
            if !fileExists(certFile) {
                return fmt.Errorf("certificate file does not exist: %s", certFile)
            }

            return internal.VerifyHashes(keyFile, csrFile, certFile)
        },
    }
    verifyHashesCmd.Flags().StringVar(&keyFile, "key", "", "Path to the private key file")
    verifyHashesCmd.Flags().StringVar(&csrFile, "csr", "", "Path to the CSR file")
    verifyHashesCmd.Flags().StringVar(&certFile, "cert", "", "Path to the certificate file")

    rootCmd.AddCommand(generateConfigCmd)
    rootCmd.AddCommand(generateCSRCmd)
    rootCmd.AddCommand(extractInfoCmd)
    rootCmd.AddCommand(checkExpirationCmd)
    rootCmd.AddCommand(fingerprintCmd)
    rootCmd.AddCommand(verifyHashesCmd)

    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}

func readLine() string {
    var line string
    fmt.Scanln(&line)
    return strings.TrimSpace(line)
}

// promptFor muestra un prompt con un valor por defecto. Si el usuario presiona Enter, se mantiene el valor por defecto.
func promptFor(label, defaultVal string) string {
    if defaultVal != "" {
        fmt.Printf("%s [%s]: ", label, defaultVal)
    } else {
        fmt.Printf("%s: ", label)
    }
    input := readLine()
    if input == "" {
        return defaultVal
    }
    return input
}

func fileExists(path string) bool {
    info, err := os.Stat(path)
    if os.IsNotExist(err) {
        return false
    }
    return !info.IsDir()
}

func getParamValue(param string) string {
    switch param {
    case "domain":
        return domain
    case "country":
        return country
    case "locality":
        return locality
    case "organization":
        return organization
    }
    return ""
}