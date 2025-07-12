@echo off
REM Script para ejecutar tests con Docker en Windows
REM Uso: scripts\run_docker_tests.bat

echo 🧪 Starting Docker tests for Go Users API...
echo ==============================================

REM Verificar si Docker está ejecutándose
docker info >nul 2>&1
if %errorlevel% neq 0 (
    echo [ERROR] Docker is not running. Please start Docker and try again.
    exit /b 1
)

REM Detener contenedores existentes
echo [INFO] Stopping existing containers...
docker compose -f docker-compose.test.yml down --remove-orphans 2>nul

REM Limpiar imágenes existentes
echo [INFO] Cleaning up existing images...
for /f "tokens=*" %%i in ('docker images -q go-users-api 2^>nul') do docker rmi %%i 2>nul

REM Iniciar MongoDB para tests
echo [INFO] Starting MongoDB for tests...
docker compose -f docker-compose.test.yml up -d mongodb

REM Esperar a que MongoDB esté listo
echo [INFO] Waiting for MongoDB to be ready...
timeout /t 10 /nobreak >nul

REM Ejecutar tests
echo [INFO] Running tests...
docker compose -f docker-compose.test.yml run --rm test
if %errorlevel% equ 0 (
    echo [SUCCESS] ✅ All tests passed!
    
    REM Construir la aplicación
    echo [INFO] Building application...
    docker compose -f docker-compose.test.yml build api
    if %errorlevel% equ 0 (
        echo [SUCCESS] ✅ Application built successfully!
        
        REM Iniciar la aplicación
        echo [INFO] Starting application...
        docker compose -f docker-compose.test.yml up -d api mongo-express
        
        echo [SUCCESS] 🎉 Application is now running!
        echo.
        echo 📋 Service URLs:
        echo    - API: http://localhost:8080
        echo    - Swagger Docs: http://localhost:8080/swagger/index.html
        echo    - MongoDB Express: http://localhost:8081
        echo.
        echo 📊 Test data has been loaded with 10 users in MongoDB
        echo 🔧 To view logs: docker compose -f docker-compose.test.yml logs -f
        echo 🛑 To stop: docker compose -f docker-compose.test.yml down
        
    ) else (
        echo [ERROR] ❌ Failed to build application
        exit /b 1
    )
) else (
    echo [ERROR] ❌ Tests failed! Application will not be built.
    echo [INFO] Stopping test containers...
    docker compose -f docker-compose.test.yml down
    exit /b 1
) 