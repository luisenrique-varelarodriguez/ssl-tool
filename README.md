
# SSL Tool

**SSL Tool** es una herramienta de línea de comandos escrita en Go para gestionar certificados SSL, generar CSRs, verificar la consistencia entre clave privada/CSR/certificado, y realizar otras operaciones relacionadas con SSL. Está diseñada para funcionar tanto de manera no interactiva (ideal para entornos automatizados o CI/CD) como en modo interactivo (tipo asistente), facilitando su uso a personas que no recuerden todos los flags o valores necesarios.

## Características principales

- **Generación de CSR y clave privada:**  
  Crea una clave RSA y un CSR con información personalizada (dominio, país, localidad, organización).  
  El CSR se genera en una carpeta con el nombre derivado del dominio (por ejemplo: `example_com/`) dentro del directorio actual.

- **Extracción de información de CRT/CSR:**  
  Imprime detalles de un certificado o CSR (Common Name, Organización, Fechas de validez, etc.).  
  Además, guarda la información en el mismo archivo YAML generado por `generate-config` (`ssl-tool-config.yaml`).

- **Verificación de hashes (private key, CSR, cert):**  
  Comprueba que el trío clave privada, CSR y certificado concuerden, calculando el hash del módulo de cada uno.

- **Verificación de expiración:**  
  Muestra cuántos días faltan para que un certificado caduque.

- **Fingerprint (huella digital):**  
  Muestra el fingerprint SHA256 de un certificado.

- **Modo interactivo:**  
  Si se activa `--interactive`, la herramienta funciona como un asistente que pregunta todos los datos necesarios, mostrando los valores por defecto si existen en flags o en el archivo de configuración. Esto permite no tener que recordar todos los parámetros.

- **Archivo de configuración YAML:**  
  Permite establecer valores por defecto (país, localidad, organización, etc.) en un archivo `ssl-tool-config.yaml` en el directorio actual.  
  Si existe, la herramienta lo carga antes de ejecutar comandos, facilitando la reutilización de configuraciones.  
  Además, comandos como `extract-info` guardan su salida en este archivo.

## Requisitos

- Go 1.16 o superior.
- (Opcional) Un archivo de configuración `ssl-tool-config.yaml` en el directorio actual, generado con `ssl-tool generate-config`.

## Instalación

### Compilar el binario

Para compilar el programa en un binario ejecutable, ejecuta lo siguiente desde la raíz del proyecto:

```bash
go build -o ssl-tool cmd/main.go
```

Esto generará un archivo **`ssl-tool`** en el directorio actual.

### Mover el binario al PATH del sistema

Para ejecutar `ssl-tool` como un comando desde cualquier lugar, mueve el binario a un directorio incluido en el `PATH` del sistema.

#### En Linux/macOS:

```bash
sudo mv ssl-tool /usr/local/bin/
```

Verifica que esté correctamente instalado:

```bash
ssl-tool --help
```

#### En Windows:

1. Compila el binario para Windows:

   ```bash
   GOOS=windows GOARCH=amd64 go build -o ssl-tool.exe cmd/main.go
   ```

2. Copia `ssl-tool.exe` a un directorio incluido en el `PATH` del sistema, como `C:\Windows\`.

3. Verifica que esté correctamente instalado:

   ```bash
   ssl-tool.exe --help
   ```

## Uso general

```bash
ssl-tool [flags] [command]
```

Para ver la lista de comandos disponibles:

```bash
ssl-tool --help
```

## Modo interactivo

Si se usa `--interactive`, la herramienta preguntará todos los datos necesarios. Por ejemplo:

```bash
ssl-tool generate-csr --interactive
```

La herramienta mostrará prompts para el dominio, país, localidad, organización, etc. Si se han pasado flags o hay valores por defecto en la configuración, se mostrarán como sugerencias. El usuario puede presionar Enter para aceptar el valor por defecto o escribir uno nuevo.

En modo no interactivo, si faltan parámetros, el comando falla y muestra un mensaje de error indicando qué falta.

## Comandos principales

### `generate-config`

Genera un archivo YAML con valores por defecto.

```bash
ssl-tool generate-config
```

Esto crea (o sobrescribe) `ssl-tool-config.yaml` en el directorio actual.

### `generate-csr`

Genera un CSR y una clave privada.

```bash
ssl-tool generate-csr --domain example.com --country US --locality "New York" --organization ExampleOrg
```

- **Modo interactivo**:
```bash
ssl-tool generate-csr --interactive
```

El CSR y la clave se generan en una carpeta `example_com/` dentro del directorio actual.

### `extract-info`

Extrae información de un certificado o CSR y la guarda en `ssl-tool-config.yaml`.

```bash
ssl-tool extract-info --file path/to/cert.crt
```

### `verify-hashes`

Verifica que la clave privada, el CSR y el certificado coincidan.

```bash
ssl-tool verify-hashes --key path/to/key.key --csr path/to/req.csr --cert path/to/cert.crt
```

### `check-expiration`

Muestra cuántos días quedan hasta la expiración del certificado.

```bash
ssl-tool check-expiration --cert path/to/cert.crt
```

### `fingerprint`

Muestra el fingerprint SHA256 de un certificado.

```bash
ssl-tool fingerprint --cert path/to/cert.crt
```

## Ejemplo de flujo completo

1. Generar un archivo de configuración YAML predeterminado:

   ```bash
   ssl-tool generate-config
   ```

2. Generar un CSR y una clave privada:

   ```bash
   ssl-tool generate-csr --domain example.com
   ```

3. Extraer información del CSR y actualizar el YAML:

   ```bash
   ssl-tool extract-info --file example_com/example_com.csr
   ```

4. Verificar consistencia de la clave, CSR y certificado:

   ```bash
   ssl-tool verify-hashes --key example_com/example_com.key --csr example_com/example_com.csr --cert example_com/example_com.crt
   ```

5. Comprobar expiración del certificado:

   ```bash
   ssl-tool check-expiration --cert example_com/example_com.crt
   ```

## Conclusión

SSL Tool proporciona un flujo de trabajo flexible para la gestión de certificados:

- **No interactivo:** Para entornos automatizados, CI/CD, o usuarios que conocen los flags.
- **Interactivo:** Para usuarios que prefieren una experiencia guiada paso a paso.

La configuración YAML, junto con la exportación automática de datos en `extract-info`, hace que la herramienta sea cómoda, segura y fácil de usar.
