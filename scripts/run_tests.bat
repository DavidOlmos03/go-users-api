@echo off
REM Script para ejecutar tests locales en Windows
REM Uso: scripts\run_tests.bat [opciones]

setlocal enabledelayedexpansion

REM Variables por defecto
set VERBOSE=false
set COVERAGE=false
set UNIT_ONLY=false
set INTEGRATION_ONLY=false
set FAIL_FAST=false

REM Parsear argumentos
:parse_args
if "%~1"=="" goto :end_parse
if "%~1"=="-v" set VERBOSE=true
if "%~1"=="--verbose" set VERBOSE=true
if "%~1"=="-c" set COVERAGE=true
if "%~1"=="--coverage" set COVERAGE=true
if "%~1"=="-u" set UNIT_ONLY=true
if "%~unit"=="--unit" set UNIT_ONLY=true
if "%~1"=="-i" set INTEGRATION_ONLY=true
if "%~1"=="--integration" set INTEGRATION_ONLY=true
if "%~1"=="-f" set FAIL_FAST=true
if "%~1"=="--fail-fast" set FAIL_FAST=true
if "%~1"=="-h" goto :show_help
if "%~1"=="--help" goto :show_help
shift
goto :parse_args

:show_help
echo Uso: %0 [opciones]
echo.
echo Opciones:
echo   -h, --help          Mostrar esta ayuda
echo   -v, --verbose       Ejecutar tests con output verbose
echo   -c, --coverage      Generar reporte de cobertura
echo   -u, --unit          Ejecutar solo tests unitarios
echo   -i, --integration   Ejecutar solo tests de integración
echo   -f, --fail-fast     Detener en el primer error
echo.
echo Ejemplos:
echo   %0                    # Ejecutar todos los tests
echo   %0 -v -c             # Tests verbose con cobertura
echo   %0 -u                # Solo tests unitarios
echo   %0 -i -f             # Tests de integración con fail-fast
exit /b 0

:end_parse

REM Verificar que estamos en el directorio correcto
if not exist "go.mod" (
    echo [ERROR] No se encontró go.mod. Ejecuta este script desde la raíz del proyecto.
    exit /b 1
)

echo [INFO] Iniciando ejecución de tests...

REM Construir argumentos para go test
set TEST_ARGS=

if "%VERBOSE%"=="true" set TEST_ARGS=%TEST_ARGS% -v
if "%FAIL_FAST%"=="true" set TEST_ARGS=%TEST_ARGS% -failfast
if "%COVERAGE%"=="true" (
    set TEST_ARGS=%TEST_ARGS% -cover -coverprofile=coverage.out
    echo [INFO] Generando reporte de cobertura...
)

REM Ejecutar tests
echo [INFO] Ejecutando tests...

if "%COVERAGE%"=="true" (
    go test %TEST_ARGS% ./tests -covermode=atomic
) else (
    go test %TEST_ARGS% ./tests
)

if %errorlevel% equ 0 (
    echo [SUCCESS] Tests completados exitosamente!
    
    REM Generar reporte de cobertura si se solicitó
    if "%COVERAGE%"=="true" (
        if exist "coverage.out" (
            echo [INFO] Generando reporte HTML de cobertura...
            go tool cover -html=coverage.out -o coverage.html
            
            echo [SUCCESS] Reporte de cobertura generado: coverage.html
            echo [INFO] Cobertura total:
            go tool cover -func=coverage.out | findstr "total:"
            
            REM Limpiar archivo temporal
            del coverage.out
        ) else (
            echo [WARNING] No se generó archivo de cobertura
        )
    )
    
    echo [SUCCESS] Todos los tests completados exitosamente!
) else (
    echo [ERROR] Tests fallaron
    exit /b 1
) 