#!/bin/bash

# Script de ejemplo para probar la API Users BRM

BASE_URL="http://localhost:8080/api/v1"

echo "🧪 Probando API Users BRM..."

# Función para hacer requests HTTP
make_request() {
    local method=$1
    local endpoint=$2
    local data=$3
    
    if [ -n "$data" ]; then
        curl -s -X $method "$BASE_URL$endpoint" \
            -H "Content-Type: application/json" \
            -d "$data" | jq .
    else
        curl -s -X $method "$BASE_URL$endpoint" | jq .
    fi
}

# Health check
echo "🔍 Health Check..."
make_request "GET" "/health"

echo -e "\n"

# Crear usuarios de ejemplo
echo "👤 Creando usuarios de ejemplo..."

# Usuario 1
echo "Creando usuario 1..."
USER1_RESPONSE=$(make_request "POST" "/users" '{
    "name": "John Doe",
    "email": "john.doe@example.com",
    "age": 30,
    "phone": "+1234567890",
    "address": "123 Main St, City, Country"
}')
echo "$USER1_RESPONSE"

# Extraer ID del primer usuario
USER1_ID=$(echo "$USER1_RESPONSE" | jq -r '.data.id')
echo "ID del usuario 1: $USER1_ID"

echo -e "\n"

# Usuario 2
echo "Creando usuario 2..."
USER2_RESPONSE=$(make_request "POST" "/users" '{
    "name": "Jane Smith",
    "email": "jane.smith@example.com",
    "age": 25,
    "phone": "+0987654321",
    "address": "456 Oak Ave, Town, State"
}')
echo "$USER2_RESPONSE"

echo -e "\n"

# Obtener todos los usuarios
echo "📋 Obteniendo todos los usuarios..."
make_request "GET" "/users"

echo -e "\n"

# Obtener usuario específico
echo "🔍 Obteniendo usuario por ID..."
make_request "GET" "/users/$USER1_ID"

echo -e "\n"

# Actualizar usuario
echo "✏️  Actualizando usuario..."
make_request "PUT" "/users/$USER1_ID" '{
    "name": "John Updated Doe",
    "email": "john.updated@example.com",
    "age": 31
}'

echo -e "\n"

# Verificar usuario actualizado
echo "🔍 Verificando usuario actualizado..."
make_request "GET" "/users/$USER1_ID"

echo -e "\n"

# Eliminar usuario
echo "🗑️  Eliminando usuario..."
make_request "DELETE" "/users/$USER1_ID"

echo -e "\n"

# Verificar que el usuario fue eliminado
echo "🔍 Verificando que el usuario fue eliminado..."
make_request "GET" "/users/$USER1_ID"

echo -e "\n"

# Obtener usuarios restantes
echo "📋 Obteniendo usuarios restantes..."
make_request "GET" "/users"

echo -e "\n✅ Pruebas completadas!" 

