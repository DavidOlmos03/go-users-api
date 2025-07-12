# Go Users API

API RESTful para gestión de usuarios desarrollada en Go con MongoDB.

## 🚀 Características

- **Framework**: Gin (Go)
- **Base de datos**: MongoDB
- **Documentación**: Swagger/OpenAPI
- **Arquitectura**: Clean Architecture
- **Testing**: Tests unitarios e integración
- **Docker**: Containerización completa

## 📋 Endpoints

- `POST /api/v1/users/` - Crear usuario
- `GET /api/v1/users/` - Listar usuarios (con paginación)
- `GET /api/v1/users/:id` - Obtener usuario por ID
- `PUT /api/v1/users/:id` - Actualizar usuario
- `DELETE /api/v1/users/:id` - Eliminar usuario
- `GET /api/v1/health` - Health check

## 📥 Instalación

### 1. Clonar el repositorio
```bash
git clone https://github.com/DavidOlmos03/go-users-api.git
```

### 2. Cambiar al directorio del proyecto
```bash
cd go-users-api
```

## 🐳 Ejecutar el proyecto

### Opción 1: Script automatizado (Recomendado)

**Linux/macOS:**
```bash
./scripts/run_docker_tests.sh
```

**Windows:**
```cmd
scripts\run_docker_tests.bat
```

Este script ejecuta los tests y construye la aplicación solo si los tests pasan.

### Opción 2: Comando directo

**Ejecutar tests y construir:**
```bash
docker compose -f docker-compose.test.yml up --build
```

**Solo desarrollo (sin tests):**
```bash
docker compose -f docker-compose.dev.yml up --build
```

## 🌐 Servicios disponibles

| Servicio | URL | Descripción |
|----------|-----|-------------|
| **API Go** | http://localhost:8080 | API REST principal |
| **Swagger Docs** | http://localhost:8080/swagger/index.html | Documentación API |
| **Mongo Express** | http://localhost:8081 | Interfaz web MongoDB |
| **MongoDB** | localhost:27017 | Base de datos |

**Credenciales Mongo Express:**
- Usuario: `admin`
- Contraseña: `admin123`

## 🎨 Frontend Angular

Para una mejor experiencia de trabajo con el CRUD, visita el frontend creado para este proyecto:

**🔗 [Go Users Frontend](https://github.com/DavidOlmos03/go-users-front)**

Desarrollado con Angular, proporciona una interfaz gráfica completa para gestionar usuarios.


# 📘 Qué aprendí trabajando con Go y MongoDB

- Aprendí la **sintaxis básica de Go**, incluyendo:
  - Declaración de funciones con `func`
  - Uso de `nil` como valor nulo
  - Estructura de `go.mod` e importaciones (`import name-module/folder`)

- Conexión con **MongoDB (NoSQL)** y manejo de:
  - Modelos, controladores (`controllers`) y servicios (`services`)
  - Tests básicos

- Uso de **Swagger (Swaggo)** para documentar la API:
  - `swag init -g main.go -o doc` genera la documentación desde anotaciones

- Comandos útiles de Go:
  - `go mod tidy` → actualiza dependencias y limpia `go.sum`
  - `go build -o main` → genera binario para despliegue
  - `go get` → instala dependencias y las registra en `go.mod`
  - `go install` → instala binarios sin afectar `go.mod`
  - `go mod vendor` → crea carpeta `vendor` con dependencias locales

| Comando                      | Afecta `go.mod` | Uso principal                          |
|-----------------------------|------------------|----------------------------------------|
| `go install paquete@versión`| ❌               | Instalar binarios externos             |
| `go get paquete@versión`    | ✅               | Agregar dependencia al proyecto        |

# 🔧 Mejoras futuras

- **Arquitectura**: Refactorizar la estructura del proyecto aplicando mejor los principios SOLID
- **Separación de responsabilidades**: Optimizar la distribución de funciones entre capas
- **Validación**: Implementar validaciones más robustas en los endpoints
- **Seguridad**: Implementar autenticación y autorización

## 📚 Tecnologías

- **Go 1.24.5** - Lenguaje de programación
- **Gin** - Framework web
- **MongoDB** - Base de datos NoSQL
- **Mongo Express** - Interfaz web para MongoDB
- **Swagger** - Documentación de API
- **Docker** - Containerización
- **Docker Compose** - Orquestación de contenedores

