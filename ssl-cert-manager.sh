#!/bin/bash

# Función para solicitar la ruta de un archivo
function solicitar_archivo() {
  local mensaje="$1"
  local archivo
  while true; do
    read -p "$mensaje " archivo
    if [[ -f "$archivo" ]]; then
      echo "$archivo"
      break
    else
      echo "Error: El archivo especificado no existe. Inténtalo de nuevo."
    fi
  done
}

# Función para crear clave privada y CSR
function crear_key_y_csr() {
  echo "=== Generador de Clave Privada y CSR ==="
  read -p "Introduce el nombre del dominio (ejemplo: example.com): " DOMAIN
  if [[ -z "$DOMAIN" ]]; then
    echo "El dominio no puede estar vacío. Inténtalo de nuevo."
    exit 1
  fi

  # Nombres de archivo
  KEY_FILE="${DOMAIN//./_}.key"
  CSR_FILE="${DOMAIN//./_}.csr"

  # Solicitar información para el CSR
  read -p "País (C, 2 letras, ejemplo: ES): " COUNTRY
  COUNTRY=${COUNTRY:-ES}
  read -p "Localidad (L, ejemplo: Madrid): " LOCALITY
  LOCALITY=${LOCALITY:-Madrid}
  read -p "Organización (O, ejemplo: MiEmpresa): " ORGANIZATION
  ORGANIZATION=${ORGANIZATION:-MiEmpresa}
  COMMON_NAME=$DOMAIN

  # Generar clave privada y CSR
  echo "Generando clave privada (${KEY_FILE}) y CSR (${CSR_FILE})..."
  openssl req -new -newkey rsa:2048 -nodes -out "$CSR_FILE" -keyout "$KEY_FILE" \
    -subj "/C=${COUNTRY}/L=${LOCALITY}/O=${ORGANIZATION}/CN=${COMMON_NAME}" -sha256

  if [[ -f "$KEY_FILE" && -f "$CSR_FILE" ]]; then
    echo "Archivos generados con éxito:"
    echo " - Clave privada: $KEY_FILE"
    echo " - CSR: $CSR_FILE"
  else
    echo "Error: No se pudieron generar los archivos."
  fi
}

# Función para crear certificado final
function crear_certificado_final() {
  echo "=== Generador de Certificado Final ==="
  CERT_FILE=$(solicitar_archivo "Introduce la ruta completa del certificado firmado por la CA (ejemplo: /ruta/a/tu_certificado.crt):")
  INTERMEDIATE_FILE=$(solicitar_archivo "Introduce la ruta completa del certificado intermedio (ejemplo: /ruta/a/intermediate.crt):")

  read -p "Introduce la ruta completa donde quieres guardar el certificado final (ejemplo: /ruta/a/certificado_final.crt): " FINAL_CERT_FILE
  FINAL_CERT_FILE=${FINAL_CERT_FILE:-"./certificado_final_con_cadena.crt"}

  echo "Generando el certificado final con la cadena completa en: ${FINAL_CERT_FILE}..."
  cat "$CERT_FILE" "$INTERMEDIATE_FILE" > "$FINAL_CERT_FILE"

  if [[ -f "$FINAL_CERT_FILE" ]]; then
    echo "Certificado final generado con éxito: ${FINAL_CERT_FILE}"

    # Verificación automática del certificado
    echo "Verificando la validez del certificado generado..."
    VERIFY_OUTPUT=$(openssl verify -CAfile "$INTERMEDIATE_FILE" "$FINAL_CERT_FILE" 2>&1)

    if echo "$VERIFY_OUTPUT" | grep -q ": OK"; then
      echo "=== Verificación exitosa: El certificado es válido y la cadena está completa. ==="
    else
      echo "=== Error en la verificación: ==="
      echo "$VERIFY_OUTPUT"
      echo "Revisa los archivos proporcionados o contacta con la CA si el problema persiste."
    fi
  else
    echo "Error: No se pudo generar el archivo ${FINAL_CERT_FILE}."
  fi
}

# Función para crear archivos Base64 del certificado y la clave privada
function crear_archivos_base64() {
  echo "=== Generador de Certificado y Clave en Base64 ==="
  FINAL_CERT_FILE=$(solicitar_archivo "Introduce la ruta completa del certificado final (ejemplo: /ruta/a/certificado_final.crt):")
  KEY_FILE=$(solicitar_archivo "Introduce la ruta completa de la clave privada (ejemplo: /ruta/a/appweb.key):")

  # Nombres de salida
  BASE64_CERT_FILE="${FINAL_CERT_FILE%.*}_base64.crt"
  BASE64_KEY_FILE="${KEY_FILE%.*}_base64.key"

  # Convertir certificado a Base64
  echo "Generando archivo Base64 del certificado en: ${BASE64_CERT_FILE}..."
  openssl base64 -in "$FINAL_CERT_FILE" -out "$BASE64_CERT_FILE"

  # Convertir clave privada a Base64
  echo "Generando archivo Base64 de la clave privada en: ${BASE64_KEY_FILE}..."
  openssl base64 -in "$KEY_FILE" -out "$BASE64_KEY_FILE"

  # Verificar resultados
  if [[ -f "$BASE64_CERT_FILE" && -f "$BASE64_KEY_FILE" ]]; then
    echo "Archivos Base64 generados con éxito:"
    echo " - Certificado en Base64: $BASE64_CERT_FILE"
    echo " - Clave privada en Base64: $BASE64_KEY_FILE"
  else
    echo "Error: No se pudieron generar los archivos Base64."
  fi
}

# Función para verificar hashes
function verificar_hashes() {
  echo "=== Verificador de Hashes (Clave Privada, CSR y Certificado) ==="
  KEY_FILE=$(solicitar_archivo "Introduce la ruta completa de la clave privada (ejemplo: /ruta/a/appweb.key):")
  CSR_FILE=$(solicitar_archivo "Introduce la ruta completa del CSR (ejemplo: /ruta/a/appweb.csr):")
  CERT_FILE=$(solicitar_archivo "Introduce la ruta completa del certificado final (ejemplo: /ruta/a/appweb.crt):")

  echo "Calculando hashes..."
  KEY_HASH=$(openssl rsa -in "$KEY_FILE" -noout -modulus | openssl md5)
  CSR_HASH=$(openssl req -in "$CSR_FILE" -noout -modulus | openssl md5)
  CERT_HASH=$(openssl x509 -in "$CERT_FILE" -noout -modulus | openssl md5)

  echo "Resultados:"
  echo " - Hash de la clave privada: $KEY_HASH"
  echo " - Hash del CSR: $CSR_HASH"
  echo " - Hash del certificado: $CERT_HASH"

  if [[ "$KEY_HASH" == "$CSR_HASH" && "$CSR_HASH" == "$CERT_HASH" ]]; then
    echo "=== Verificación exitosa: Los hashes coinciden. ==="
  else
    echo "=== Error: Los hashes no coinciden. ==="
    echo "Revisa los archivos proporcionados para posibles errores."
  fi
}

# Menú principal
function menu_principal() {
  echo "=== Gestor de Certificados SSL ==="
  echo "1. Crear Clave Privada y CSR"
  echo "2. Crear Certificado Final (cadena completa)"
  echo "3. Crear Certificado y Clave en Base64"
  echo "4. Verificar Hashes (Clave Privada, CSR y Certificado)"
  echo "5. Salir"
  read -p "Selecciona una opción (1-5): " OPCION

  case $OPCION in
    1)
      crear_key_y_csr
      ;;
    2)
      crear_certificado_final
      ;;
    3)
      crear_archivos_base64
      ;;
    4)
      verificar_hashes
      ;;
    5)
      echo "Saliendo..."
      exit 0
      ;;
    *)
      echo "Opción inválida. Inténtalo de nuevo."
      menu_principal
      ;;
  esac
}

# Ejecutar el menú principal
menu_principal
