# Microservicio de Firma Digital para DTE

[![Go Version](https://img.shields.io/badge/Go-1.23-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-GPL%20v3-green.svg)](LICENSE)

> **📚 Nota:** Si ya usabas anteriormente el firmador brindado por Hacienda, este microservicio devuelve el mismo formato de respuesta y error asi como de mensajes y códigos de error.

Microservicio encargado de la firma digital de Documentos Tributarios Electrónicos (DTE) que cumple con los requisitos establecidos por el Ministerio de Hacienda en El Salvador.

## 📋 Características

- Firma digital de documentos utilizando certificados `.crt`
- Soporte para múltiples certificados organizados por NIT
- Internacionalización de mensajes de error (español/inglés)
- API REST simple y de alto rendimiento
- Monitoreo de estado del servicio
- Diseño modular siguiendo principios de arquitectura hexagonal

## 🏗️ Arquitectura

Este proyecto está implementado siguiendo principios de:

- **Arquitectura Hexagonal (Ports & Adapters)**: Separación clara entre lógica de dominio y acceso a recursos externos
- **Diseño guiado por dominio (DDD)**: Estructura organizada de acuerdo al dominio del problema

### Estructura del proyecto:

```
.
├── cmd               # Punto de entrada de la aplicación
├── configs           # Configuraciones y archivos de localización
├── internal          # Código interno
│   ├── application   # Casos de uso y servicios de aplicación
│   ├── domain        # Modelos, puertos y servicios de dominio
│   └── infrastructure # Adaptadores y componentes de infraestructura
├── pkg               # Paquetes reutilizables
│   ├── i18n          # Internacionalización
│   ├── logs          # Logging
│   └── response      # Estructuras de respuesta estandarizadas
└── uploads           # Directorio para almacenar certificados
```

## 🛠️ Tecnologías

- **Go 1.23**: Lenguaje de programación principal
- **Gorilla Mux**: Enrutador HTTP
- **Viper**: Gestión de configuración
- **Docker**: Contenerización para despliegue

## 🔧 Requisitos previos

- Docker y Docker Compose (recomendado para despliegue)
- Go 1.23+ (para desarrollo)
- Certificados digitales de firma en formato PKCS#8

## 📦 Instalación

### Con Docker (Recomendado)

1. Clonar el repositorio:
```bash
git clone https://github.com/chainedpixel/go-dte-signer.git
cd go-dte-signer
```

2. Colocar certificados de firma digital en la carpeta `uploads` siguiendo la estructura especificada en el README de ese directorio.

3. Iniciar el servicio:
```bash
docker-compose up -d
```

### Para desarrollo local

1. Instalar dependencias:
```bash
go mod download
```

2. Ejecutar el servicio:
```bash
go run cmd/signserver/main.go
```

## ⚙️ Configuración

La configuración se maneja mediante un archivo `config.yaml` y variables de entorno:

```yaml
# Server
server:
  port: "8113"
  signerroute: "/sign"
  healthroute: "/health"
  readtimeout: 30
  writetimeout: 30

# Internationalization
locale:
  defaultlocale: "es"
  localesdir: "./configs/locales"

# File system
filesystem:
  certificatesdir: "./uploads/"

# Logging
log:
  level: "info"
  format: "text"
```

## 🚀 Uso

### Endpoints

El servicio expone dos endpoints principales que son configurables a través del archivo `config.yaml`:

#### Firmado de documentos

`POST /sign` (ruta configurable en `server.signerroute`)

Ejemplo de solicitud:
```json
{
  "nit": "06140101780010",
  "passwordPri": "cl4v3-pr1v4d4",
  "dteJson": {
    "identificacion": {
      "version": 3,
      "ambiente": "00",
      "tipoDte": "01"
      // ... más propiedades del DTE
    }
  }
}
```

### Ejemplo de respuesta:
```json
{
  "status": "OK",
  "body": "eyJhbGciOiJSUzUxM..."
}
```

#### Estado de salud del servicio

`GET /health` (ruta configurable en `server.healthroute`)

Proporciona información sobre el estado del servicio, tiempo de ejecución y versión.

### Ejemplo de respuesta:
```json
{
  "status": "OK",
  "body": {
    "status": "UP",
    "uptime": "27.2275698s",
    "timestamp": "2025-04-20T19:39:09.256993-06:00",
    "goVersion": "go1.23.2"
  }
}
```


## 🔌 Integración con API de Facturación Electrónica

Este servicio de firma es un componente esencial para la emisión de DTEs pero no implementa la lógica completa para facturación electrónica. Si estás buscando una solución integral para facturación electrónica, consulta nuestra [API de Facturación Electrónica para El Salvador](https://github.com/chainedpixel/api-facturacion-sv) que integra este servicio de firma con la funcionalidad completa para emisión, validación y transmisión de documentos tributarios electrónicos según normativa vigente.

## 🔒 Seguridad

- Manejo seguro de claves privadas
- Validación de certificados
- Protección contra ataques comunes a APIs

## 📚 Documentación

Para obtener más información sobre el funcionamiento interno del servicio, consulta los siguientes recursos:

- [Manuales técnicos y de usuario oficiales del Ministerio de Hacienda](https://factura.gob.sv/informacion-tecnica-y-funcional/)

## 📝 Licencia

Este proyecto está licenciado bajo la Licencia GNU General Public License v3.0 (GPL-3.0) - ver el archivo [LICENSE](LICENSE) para más detalles.