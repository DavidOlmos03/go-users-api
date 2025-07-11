# Go Users API

API RESTful para gestión de usuarios desarrollada en Go con MongoDB.

## 🚀 Características

- **Framework**: Gin (Go)
- **Base de datos**: MongoDB
- **Documentación**: Swagger/OpenAPI
- **Arquitectura**: Clean Architecture (Controllers, Services, Repository)
- **Validación**: Binding validation
- **Logging**: Middleware personalizado
- **CORS**: Configurado para desarrollo

## 📋 Endpoints disponibles

- `POST /api/v1/users/` - Crear usuario
- `GET /api/v1/users/` - Listar usuarios (con paginación)
- `GET /api/v1/users/:id` - Obtener usuario por ID
- `PUT /api/v1/users/:id` - Actualizar usuario
- `DELETE /api/v1/users/:id` - Eliminar usuario
- `GET /api/v1/health` - Health check

## 🐳 Ejecutar con Docker

### Prerrequisitos

- Docker
- Docker Compose

### Pasos para ejecutar

1. **Clonar el repositorio**
   ```bash
   git clone https://github.com/DavidOlmos03/go_users_api
   cd go-users-api
   ```

2. **Regenerar dependencias (si es necesario)**
   ```bash
   go mod tidy
   ```

3. **Ejecutar con Docker Compose**
   ```bash
   docker compose -f docker-compose.dev.yml up --build
   ```

### Servicios disponibles

El proyecto incluye tres servicios:

- **API Go** (Puerto 8080) - Aplicación principal
- **MongoDB** (Puerto 27017) - Base de datos
- **Mongo Express** (Puerto 8081) - Interfaz web para MongoDB

## 🌐 Acceso a los servicios

### API REST
- **URL**: http://localhost:8080
- **Base Path**: `/api/v1`
- **Health Check**: http://localhost:8080/api/v1/health

### Documentación Swagger
- **URL**: http://localhost:8080/swagger/index.html
- **API Docs JSON**: http://localhost:8080/swagger/doc.json

### MongoDB Express (Interfaz web)
- **URL**: http://localhost:8081
- **Usuario**: admin
- **Contraseña**: admin123

### MongoDB (Conexión directa)
- **Host**: localhost
- **Puerto**: 27017
- **Base de datos**: users_brm_dev

## 📝 Ejemplos de uso

### Crear un usuario
```bash
curl -X POST http://localhost:8080/api/v1/users/ \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john.doe@example.com",
    "age": 30,
    "phone": "+1234567890",
    "address": "123 Main St, City, Country"
  }'
```

### Obtener usuarios
```bash
curl http://localhost:8080/api/v1/users/
```

### Obtener usuario por ID
```bash
curl http://localhost:8080/api/v1/users/{user_id}
```

### Actualizar usuario
```bash
curl -X PUT http://localhost:8080/api/v1/users/{user_id} \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Updated",
    "email": "john.updated@example.com"
  }'
```

### Eliminar usuario
```bash
curl -X DELETE http://localhost:8080/api/v1/users/{user_id}
```

## 🛠️ Desarrollo

### Estructura del proyecto
```
go-users-api/
├── config/          # Configuración de la aplicación
├── controllers/     # Controladores HTTP
├── middleware/      # Middlewares personalizados
├── models/          # Modelos de datos
├── repository/      # Capa de acceso a datos
├── services/        # Lógica de negocio
├── docs/           # Documentación Swagger generada
├── tests/          # Tests de la API
├── main.go         # Punto de entrada
├── Dockerfile      # Configuración de Docker
└── docker-compose.dev.yml # Configuración de desarrollo
```

### Variables de entorno

| Variable | Descripción | Valor por defecto |
|----------|-------------|-------------------|
| `PORT` | Puerto del servidor | 8080 |
| `GIN_MODE` | Modo de Gin | debug |
| `MONGO_URI` | URI de conexión a MongoDB | mongodb://localhost:27017 |
| `MONGO_DATABASE` | Nombre de la base de datos | users_brm_dev |

### Comandos útiles

**Detener servicios**
```bash
docker compose -f docker-compose.dev.yml down
```

**Ver logs**
```bash
docker compose -f docker-compose.dev.yml logs -f
```

**Reconstruir sin cache**
```bash
docker compose -f docker-compose.dev.yml up --build --force-recreate
```

**Acceder al contenedor de la API**
```bash
docker exec -it api_users_brm_dev sh
```

## 📊 Monitoreo

### Health Check
```bash
curl http://localhost:8080/api/v1/health
```

Respuesta esperada:
```json
{
  "status": "ok",
  "message": "API Users BRM is running",
  "time": "2025-07-11T05:31:02Z"
}
```

### Logs en tiempo real
```bash
# Logs de la API
docker logs -f api_users_brm_dev

# Logs de MongoDB
docker logs -f mongodb_users_brm_dev

# Logs de Mongo Express
docker logs -f mongo_express_users_brm_dev
```

## 🔧 Troubleshooting

### Problema: Error de dependencias
```bash
go mod tidy
docker compose -f docker-compose.dev.yml up --build
```

### Problema: Puerto ocupado
```bash
# Verificar puertos en uso
lsof -i :8080
lsof -i :8081
lsof -i :27017

# Detener servicios y reiniciar
docker compose -f docker-compose.dev.yml down
docker compose -f docker-compose.dev.yml up --build
```

### Problema: MongoDB no conecta
```bash
# Verificar estado de MongoDB
docker logs mongodb_users_brm_dev

# Reiniciar solo MongoDB
docker restart mongodb_users_brm_dev
```

## 📚 Tecnologías utilizadas

- **Go 1.24.5** - Lenguaje de programación
- **Gin** - Framework web
- **MongoDB** - Base de datos NoSQL
- **Mongo Express** - Interfaz web para MongoDB
- **Swagger** - Documentación de API
- **Docker** - Containerización
- **Docker Compose** - Orquestación de contenedores

## 📄 Licencia

Este proyecto está bajo la Licencia Apache 2.0.
