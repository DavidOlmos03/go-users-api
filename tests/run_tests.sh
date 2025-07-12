#!/bin/bash

# Script para ejecutar tests del proyecto Go Users API
# Uso: ./tests/run_tests.sh [opciones]

set -e

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Función para imprimir mensajes con colores
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Función para mostrar ayuda
show_help() {
    echo "Uso: $0 [opciones]"
    echo ""
    echo "Opciones:"
    echo "  -h, --help          Mostrar esta ayuda"
    echo "  -v, --verbose       Ejecutar tests con output verbose"
    echo "  -c, --coverage      Generar reporte de cobertura"
    echo "  -u, --unit          Ejecutar solo tests unitarios"
    echo "  -i, --integration   Ejecutar solo tests de integración"
    echo "  -a, --all           Ejecutar todos los tests (default)"
    echo "  -f, --fail-fast     Detener en el primer error"
    echo ""
    echo "Ejemplos:"
    echo "  $0                    # Ejecutar todos los tests"
    echo "  $0 -v -c             # Tests verbose con cobertura"
    echo "  $0 -u                # Solo tests unitarios"
    echo "  $0 -i -f             # Tests de integración con fail-fast"
}

# Variables por defecto
VERBOSE=false
COVERAGE=false
UNIT_ONLY=false
INTEGRATION_ONLY=false
FAIL_FAST=false

# Parsear argumentos
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -c|--coverage)
            COVERAGE=true
            shift
            ;;
        -u|--unit)
            UNIT_ONLY=true
            shift
            ;;
        -i|--integration)
            INTEGRATION_ONLY=true
            shift
            ;;
        -a|--all)
            UNIT_ONLY=false
            INTEGRATION_ONLY=false
            shift
            ;;
        -f|--fail-fast)
            FAIL_FAST=true
            shift
            ;;
        *)
            print_error "Opción desconocida: $1"
            show_help
            exit 1
            ;;
    esac
done

# Verificar que estamos en el directorio correcto
if [[ ! -f "go.mod" ]]; then
    print_error "No se encontró go.mod. Ejecuta este script desde la raíz del proyecto."
    exit 1
fi

print_status "Iniciando ejecución de tests..."

# Construir argumentos para go test
TEST_ARGS=""

if [[ "$VERBOSE" == true ]]; then
    TEST_ARGS="$TEST_ARGS -v"
fi

if [[ "$FAIL_FAST" == true ]]; then
    TEST_ARGS="$TEST_ARGS -failfast"
fi

if [[ "$COVERAGE" == true ]]; then
    TEST_ARGS="$TEST_ARGS -cover -coverprofile=coverage.out"
    print_status "Generando reporte de cobertura..."
fi

# Función para ejecutar tests
run_tests() {
    local test_path=$1
    local test_name=$2
    
    print_status "Ejecutando $test_name..."
    
    if [[ "$COVERAGE" == true ]]; then
        go test $TEST_ARGS ./$test_path -covermode=atomic
    else
        go test $TEST_ARGS ./$test_path
    fi
    
    if [[ $? -eq 0 ]]; then
        print_success "$test_name completados exitosamente"
    else
        print_error "$test_name fallaron"
        exit 1
    fi
}

# Ejecutar tests según las opciones
if [[ "$UNIT_ONLY" == true ]]; then
    print_status "Ejecutando solo tests unitarios..."
    run_tests "tests" "Tests unitarios"
elif [[ "$INTEGRATION_ONLY" == true ]]; then
    print_status "Ejecutando solo tests de integración..."
    run_tests "tests" "Tests de integración"
else
    print_status "Ejecutando todos los tests..."
    run_tests "tests" "Tests completos"
fi

# Generar reporte de cobertura si se solicitó
if [[ "$COVERAGE" == true ]]; then
    if [[ -f "coverage.out" ]]; then
        print_status "Generando reporte HTML de cobertura..."
        go tool cover -html=coverage.out -o coverage.html
        
        print_success "Reporte de cobertura generado: coverage.html"
        print_status "Cobertura total:"
        go tool cover -func=coverage.out | tail -1
        
        # Limpiar archivo temporal
        rm coverage.out
    else
        print_warning "No se generó archivo de cobertura"
    fi
fi

print_success "Todos los tests completados exitosamente!" 