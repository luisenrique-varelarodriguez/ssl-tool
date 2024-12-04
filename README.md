# Gestor de Certificados SSL

Este script en Bash es una herramienta interactiva para gestionar certificados SSL. Proporciona una serie de funcionalidades que permiten crear claves privadas, CSRs, certificados finales y convertirlos al formato Base64, además de verificar la coherencia entre los componentes de los certificados.

---

## Características

1. **Crear Clave Privada y CSR**  
   Genera una clave privada y un CSR (Certificate Signing Request) con los datos del usuario.

2. **Crear Certificado Final (Cadena Completa)**  
   Combina un certificado firmado por la autoridad certificadora (CA) con un certificado intermedio para generar un certificado final válido.

3. **Crear Certificado y Clave en Base64**  
   Convierte el certificado final y la clave privada al formato Base64 para casos en los que se requiera este formato.

4. **Verificar Hashes (Clave Privada, CSR y Certificado)**  
   Comprueba que los hashes de la clave privada, CSR y certificado coincidan, asegurando que todos los componentes están relacionados.

---

## Requisitos

- **Sistema operativo:** Linux/MacOS (o WSL en Windows con un entorno Bash).
- **OpenSSL:** Debe estar instalado en tu sistema.  
  Puedes verificarlo ejecutando:
  ```bash
  openssl version
  ```

---

## Cómo usar el script

### 1. Descarga y prepara el script

1. Guarda el script como `ssl_cert_manager.sh`.
2. Dale permisos de ejecución:
   ```bash
   chmod +x ssl_cert_manager.sh
   ```

### 2. Ejecuta el script

```bash
./ssl_cert_manager.sh
```

### 3. Selecciona una opción del menú principal

Cuando ejecutes el script, verás el siguiente menú:

```plaintext
=== Gestor de Certificados SSL ===
1. Crear Clave Privada y CSR
2. Crear Certificado Final (cadena completa)
3. Crear Certificado y Clave en Base64
4. Verificar Hashes (Clave Privada, CSR y Certificado)
5. Salir
```

Introduce el número correspondiente a la acción que deseas realizar.

---

## Funcionalidades detalladas

### Opción 1: Crear Clave Privada y CSR

Esta opción genera:
- Una clave privada (`.key`).
- Un CSR (`.csr`) basado en los datos proporcionados por el usuario.

**Pasos:**
1. Introduce el nombre del dominio (por ejemplo, `example.com`).
2. Proporciona información para el CSR:
   - País (`C`).
   - Localidad (`L`).
   - Organización (`O`).
   - Nombre común (dominio).

**Salida:**
- Clave privada: `example_com.key`.
- CSR: `example_com.csr`.

---

### Opción 2: Crear Certificado Final (Cadena Completa)

Esta opción combina un certificado firmado por la CA con un certificado intermedio para generar un certificado final.

**Pasos:**
1. Introduce la ruta del certificado firmado (`tu_certificado.crt`).
2. Introduce la ruta del certificado intermedio (`intermediate.crt`).
3. Proporciona una ubicación para guardar el certificado final.

**Salida:**
- Certificado final con cadena completa: `certificado_final_con_cadena.crt`.

**Verificación automática:**
El script verifica que el certificado final sea válido y muestra mensajes de éxito o error.

---

### Opción 3: Crear Certificado y Clave en Base64

Convierte:
- El certificado final (`.crt`) al formato Base64.
- La clave privada (`.key`) al formato Base64.

**Pasos:**
1. Introduce la ruta del certificado final (`certificado_final.crt`).
2. Introduce la ruta de la clave privada (`example_com.key`).

**Salida:**
- Certificado en Base64: `certificado_final_base64.crt`.
- Clave privada en Base64: `example_com_base64.key`.

---

### Opción 4: Verificar Hashes (Clave Privada, CSR y Certificado)

Comprueba que los hashes de la clave privada, CSR y certificado coincidan. Esto asegura que todos los componentes están relacionados entre sí.

**Pasos:**
1. Introduce la ruta de la clave privada (`example_com.key`).
2. Introduce la ruta del CSR (`example_com.csr`).
3. Introduce la ruta del certificado final (`certificado_final.crt`).

**Salida:**
- Muestra los hashes calculados de cada componente.
- Verifica si coinciden y muestra un mensaje de éxito o error.

**Ejemplo de salida:**

```plaintext
Calculando hashes...
 - Hash de la clave privada: MD5(stdin)= 86b79ebd9ef04f39b904043d9a65bfcd
 - Hash del CSR: MD5(stdin)= 86b79ebd9ef04f39b904043d9a65bfcd
 - Hash del certificado: MD5(stdin)= 86b79ebd9ef04f39b904043d9a65bfcd

=== Verificación exitosa: Los hashes coinciden. ===
```

---

## Errores comunes y soluciones

1. **Error: OpenSSL no está instalado.**
   - Solución: Instala OpenSSL en tu sistema.
     ```bash
     sudo apt install openssl  # En distribuciones basadas en Debian
     brew install openssl     # En macOS
     ```

2. **Error: El archivo no existe.**
   - Solución: Verifica que la ruta del archivo proporcionado sea correcta.

3. **Error: Los hashes no coinciden.**
   - Solución: Asegúrate de usar los archivos correctos que corresponden a la misma clave privada y CSR.

---

## Ejemplo de flujo completo

1. **Generar clave privada y CSR:**
   ```bash
   ./ssl_cert_manager.sh
   Selecciona una opción (1-5): 1
   ```
   Proporciona el dominio y los datos requeridos.

2. **Crear certificado final:**
   ```bash
   ./ssl_cert_manager.sh
   Selecciona una opción (1-5): 2
   ```
   Introduce las rutas de los certificados y guarda el certificado final.

3. **Convertir a Base64:**
   ```bash
   ./ssl_cert_manager.sh
   Selecciona una opción (1-5): 3
   ```
   Proporciona las rutas del certificado y la clave privada.

4. **Verificar hashes:**
   ```bash
   ./ssl_cert_manager.sh
   Selecciona una opción (1-5): 4
   ```
   Introduce las rutas de la clave privada, CSR y certificado.

---
