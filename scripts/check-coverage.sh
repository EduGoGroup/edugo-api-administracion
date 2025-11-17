#!/bin/bash
# check-coverage.sh - Verificar que la cobertura cumple el umbral mínimo
#
# Uso: ./scripts/check-coverage.sh <coverage-file.out> <threshold>
#
# Este script verifica que la cobertura de código sea mayor o igual al
# umbral especificado. Retorna código de salida 0 si pasa, 1 si falla.

set -euo pipefail

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Argumentos
COVERAGE_FILE="${1:-coverage.out}"
THRESHOLD="${2:-60}"

# Validar que existe el archivo de cobertura
if [[ ! -f "$COVERAGE_FILE" ]]; then
    echo -e "${RED}❌ Error: Archivo de cobertura no encontrado: $COVERAGE_FILE${NC}"
    echo "Ejecuta primero: go test -coverprofile=$COVERAGE_FILE ./..."
    exit 1
fi

# Validar que threshold es un número
if ! [[ "$THRESHOLD" =~ ^[0-9]+(\.[0-9]+)?$ ]]; then
    echo -e "${RED}❌ Error: Threshold debe ser un número: $THRESHOLD${NC}"
    exit 1
fi

echo -e "${BLUE}🎯 Verificando umbral de cobertura...${NC}"
echo "  Archivo:  $COVERAGE_FILE"
echo "  Umbral:   ${THRESHOLD}%"
echo ""

# Obtener cobertura total usando go tool cover
COVERAGE_OUTPUT=$(go tool cover -func="$COVERAGE_FILE" | grep "total:" || true)

if [[ -z "$COVERAGE_OUTPUT" ]]; then
    echo -e "${RED}❌ Error: No se pudo calcular la cobertura total${NC}"
    exit 1
fi

# Extraer porcentaje de cobertura
# Formato: "total:                                          (statements)             30.9%"
COVERAGE_PERCENT=$(echo "$COVERAGE_OUTPUT" | awk '{print $NF}' | sed 's/%//')

echo -e "${BLUE}📊 Resultado:${NC}"
echo "  Cobertura actual: ${COVERAGE_PERCENT}%"
echo "  Umbral mínimo:    ${THRESHOLD}%"
echo ""

# Comparar cobertura con umbral usando bc para decimales
COMPARISON=$(echo "$COVERAGE_PERCENT >= $THRESHOLD" | bc -l)

if [[ "$COMPARISON" -eq 1 ]]; then
    DIFF=$(echo "$COVERAGE_PERCENT - $THRESHOLD" | bc -l)
    echo -e "${GREEN}✅ ÉXITO: Cobertura cumple el umbral (+${DIFF}%)${NC}"
    echo ""
    echo -e "${GREEN}╔════════════════════════════════════════╗${NC}"
    echo -e "${GREEN}║   COBERTURA APROBADA: ${COVERAGE_PERCENT}% >= ${THRESHOLD}%   ║${NC}"
    echo -e "${GREEN}╚════════════════════════════════════════╝${NC}"
    exit 0
else
    DIFF=$(echo "$THRESHOLD - $COVERAGE_PERCENT" | bc -l)
    echo -e "${RED}❌ FALLO: Cobertura por debajo del umbral (-${DIFF}%)${NC}"
    echo ""
    echo -e "${RED}╔════════════════════════════════════════╗${NC}"
    echo -e "${RED}║   COBERTURA INSUFICIENTE: ${COVERAGE_PERCENT}% < ${THRESHOLD}%   ║${NC}"
    echo -e "${RED}╚════════════════════════════════════════╝${NC}"
    echo ""
    echo -e "${YELLOW}💡 Sugerencias:${NC}"
    echo "  1. Agrega tests a módulos sin cobertura"
    echo "  2. Revisa docs/TEST_ANALYSIS_REPORT.md para módulos críticos"
    echo "  3. Ejecuta: make test-coverage para ver reporte detallado"
    exit 1
fi
