# Directorio de Certificados

Este directorio se utiliza para almacenar los certificados digitales necesarios para las operaciones de firma de documentos en la aplicación DTE Signer.

## Propósito

El servicio DTE Signer utiliza los certificados digitales almacenados en este directorio para firmar documentos electrónicos. Cada certificado está asociado a un NIT específico y requiere las credenciales apropiadas para su uso.

## Requisitos de los Certificados

Cada certificado debe:

1. Estar en extension `.crt` para el certificado
2. El nombre del certificado debe ser el mismo que el `NIT` de la organización o individuo propietario del certificado

## Uso en Solicitudes de Firma

Al realizar una solicitud para firmar un documento, deberá proporcionar:

1. El NIT que coincida con el directorio del certificado
2. La contraseña de la clave privada
3. Los datos del documento a firmar

La aplicación localizará y utilizará automáticamente el certificado correcto basándose en el NIT proporcionado en la solicitud.

## Solución de Problemas

Si encuentra errores relacionados con certificados:

- Verifique que el certificado exista en la ubicación correcta
- Compruebe los permisos de archivo
- Asegúrese de que la contraseña proporcionada en la solicitud coincida con el certificado
- Verifique que el certificado esté en el formato correcto `.crt`
- Asegúrese de que el certificado no haya caducado
