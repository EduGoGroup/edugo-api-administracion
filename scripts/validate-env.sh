#!/bin/bash
# scripts/validate-env.sh
# Script de validación de variables de entorno para auth centralizada

set -e

# Colores
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo "============================================"
echo "Validación de Variables de Entorno"
echo "Proyecto: edugo-api-administracion"
echo "============================================"
echo ""

# Cargar .env si existe
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | grep -v '^$' | xargs)
    echo -e "${GREEN}✓ Archivo .env cargado${NC}"
else
    echo -e "${RED}✗ Archivo .env no encontrado${NC}"
    echo "  Ejecuta: cp .env.example .env"
    exit 1
fi

ERRORS=0

# Función para validar variable requerida
validate_required() {
    local var_name=$1
    local var_value="${!var_name}"
    
    if [ -z "$var_value" ]; then
        echo -e "${RED}✗ $var_name: NO DEFINIDA${NC}"
        ((ERRORS++))
    else
        # Ocultar valores sensibles
        if [[ $var_name == *"SECRET"* ]] || [[ $var_name == *"PASSWORD"* ]] || [[ $var_name == *"KEY"* ]]; then
            echo -e "${GREEN}✓ $var_name: ***oculto***${NC}"
        else
            echo -e "${GREEN}✓ $var_name: $var_value${NC}"
        fi
    fi
}

# Función para validar longitud mínima
validate_min_length() {
    local var_name=$1
    local min_length=$2
    local var_value="${!var_name}"
    
    if [ ${#var_value} -lt $min_length ]; then
        echo -e "${RED}✗ $var_name: Debe tener al menos $min_length caracteres (tiene ${#var_value})${NC}"
        ((ERRORS++))
    fi
}

# Función para validar valor esperado
validate_value() {
    local var_name=$1
    local expected=$2
    local var_value="${!var_name}"
    
    if [ "$var_value" != "$expected" ]; then
        echo -e "${YELLOW}⚠ $var_name: '$var_value' (esperado: '$expected')${NC}"
    fi
}

echo ""
echo "--- Variables de Ambiente ---"
validate_required "APP_ENV"

echo ""
echo "--- Variables de Base de Datos ---"
validate_required "POSTGRES_PASSWORD"
validate_required "MONGODB_URI"

echo ""
echo "--- Variables de JWT ---"
validate_required "AUTH_JWT_SECRET"
validate_min_length "AUTH_JWT_SECRET" 32
validate_required "AUTH_JWT_ISSUER"
validate_value "AUTH_JWT_ISSUER" "edugo-central"

echo ""
echo "--- Variables de Rate Limiting ---"
validate_required "AUTH_RATE_LIMIT_LOGIN_ATTEMPTS"
validate_required "AUTH_RATE_LIMIT_LOGIN_WINDOW"

echo ""
echo "--- Variables de Servicios Internos ---"
validate_required "AUTH_INTERNAL_SERVICES_API_KEYS"
validate_required "AUTH_INTERNAL_SERVICES_IP_RANGES"

echo ""
echo "--- Variables de Cache ---"
validate_required "AUTH_CACHE_TOKEN_VALIDATION_TTL"

echo ""
echo "--- Variables de Redis ---"
validate_required "REDIS_HOST"
validate_required "REDIS_PORT"

echo ""
echo "============================================"
if [ $ERRORS -eq 0 ]; then
    echo -e "${GREEN}✓ VALIDACIÓN EXITOSA - Todas las variables configuradas${NC}"
    exit 0
else
    echo -e "${RED}✗ VALIDACIÓN FALLIDA - $ERRORS errores encontrados${NC}"
    exit 1
fi
