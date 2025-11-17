#!/bin/bash
# filter-coverage.sh - Filtrar reporte de cobertura según exclusiones configuradas
#
# Uso: ./scripts/filter-coverage.sh <coverage-input.out> [coverage-output.out]
#
# Este script lee los patrones de exclusión desde .coverignore y filtra
# el reporte de cobertura para eliminar archivos que no deben contarse.

set -euo pipefail

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Argumentos
COVERAGE_INPUT="${1:-coverage.out}"
COVERAGE_OUTPUT="${2:-coverage-filtered.out}"
COVERIGNORE_FILE=".coverignore"

# Validar que existe el archivo de entrada
if [[ ! -f "$COVERAGE_INPUT" ]]; then
    echo -e "${RED}❌ Error: Archivo de cobertura no encontrado: $COVERAGE_INPUT${NC}"
    echo "Ejecuta primero: go test -coverprofile=$COVERAGE_INPUT ./..."
    exit 1
fi

# Validar que existe .coverignore
if [[ ! -f "$COVERIGNORE_FILE" ]]; then
    echo -e "${RED}❌ Error: Archivo .coverignore no encontrado${NC}"
    echo "Crea el archivo .coverignore con los patrones de exclusión"
    exit 1
fi

echo -e "${BLUE}🔍 Filtrando reporte de cobertura...${NC}"
echo "  Input:  $COVERAGE_INPUT"
echo "  Output: $COVERAGE_OUTPUT"
echo ""

# Leer patrones de .coverignore (ignorar comentarios y líneas vacías)
PATTERNS=$(grep -v '^#' "$COVERIGNORE_FILE" | grep -v '^$' || true)

if [[ -z "$PATTERNS" ]]; then
    echo -e "${YELLOW}⚠️  Warning: .coverignore está vacío, copiando archivo sin filtrar${NC}"
    cp "$COVERAGE_INPUT" "$COVERAGE_OUTPUT"
    exit 0
fi

# Convertir patrones a regex de grep
# Escapar caracteres especiales y convertir wildcards de shell a regex
GREP_PATTERN=""
while IFS= read -r pattern; do
    # Convertir patrón de shell a regex de grep
    # *_mock.go -> .*_mock\.go
    # internal/domain/entity/ -> internal/domain/entity/
    regex_pattern=$(echo "$pattern" | sed 's/\./\\./g' | sed 's/\*/\.\*/g')

    if [[ -z "$GREP_PATTERN" ]]; then
        GREP_PATTERN="$regex_pattern"
    else
        GREP_PATTERN="$GREP_PATTERN|$regex_pattern"
    fi
done <<< "$PATTERNS"

echo -e "${BLUE}📋 Patrones de exclusión:${NC}"
echo "$PATTERNS" | sed 's/^/  - /'
echo ""

# Copiar primera línea (mode: ...)
head -n 1 "$COVERAGE_INPUT" > "$COVERAGE_OUTPUT"

# Filtrar líneas que NO coincidan con los patrones
# -v: invertir match (mostrar líneas que NO coinciden)
# -E: regex extendido
tail -n +2 "$COVERAGE_INPUT" | grep -vE "$GREP_PATTERN" >> "$COVERAGE_OUTPUT" || true

# Contar líneas antes y después
LINES_BEFORE=$(wc -l < "$COVERAGE_INPUT" | tr -d ' ')
LINES_AFTER=$(wc -l < "$COVERAGE_OUTPUT" | tr -d ' ')
LINES_FILTERED=$((LINES_BEFORE - LINES_AFTER))

echo -e "${GREEN}✅ Filtrado completado${NC}"
echo "  Líneas antes:    $LINES_BEFORE"
echo "  Líneas después:  $LINES_AFTER"
echo "  Líneas filtradas: $LINES_FILTERED"
echo ""

# Calcular cobertura aproximada
if command -v go &> /dev/null; then
    echo -e "${BLUE}📊 Análisis de cobertura:${NC}"
    echo ""
    go tool cover -func="$COVERAGE_OUTPUT" | tail -20
fi

echo ""
echo -e "${GREEN}✓ Reporte filtrado guardado en: $COVERAGE_OUTPUT${NC}"
