# Build stage with Alpine linux image
FROM golang:1.24.5-alpine AS builder

# Instalar dependencias del sistema
RUN apk add --no-cache git ca-certificates tzdata

# Establecer directorio de trabajo
WORKDIR /app

# Copiar archivos de dependencias
COPY go.mod go.sum ./

# Descargar dependencias
RUN go mod download

# Copiar código fuente
COPY . .

# Generar documentación Swagger
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g main.go -o docs

# Construir la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

# Instalar ca-certificates para HTTPS
RUN apk --no-cache add ca-certificates

# Crear usuario no-root
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /root/

# Copiar el binario desde el stage de build
COPY --from=builder /app/main .

# Copiar documentación Swagger generada
COPY --from=builder /app/docs ./docs

# Copiar archivos de configuración
COPY --from=builder /app/env.* ./

# Cambiar propietario al usuario no-root
RUN chown -R appuser:appgroup /root/

# Cambiar al usuario no-root
USER appuser

# Exponer puerto
EXPOSE 8080

# Comando para ejecutar la aplicación
CMD ["./main"] 

