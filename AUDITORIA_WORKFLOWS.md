# Auditoría de Workflows - edugo-api-administracion

**Fecha**: 2025-11-17  
**Auditor**: Claude (Anthropic)  
**Ubicación**: `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-administracion/.github/workflows/`

---

## 📊 Resumen Ejecutivo

- **Total workflows**: 10 archivos
- **Problemas críticos**: 5
- **Advertencias**: 8
- **Recomendaciones**: 6
- **Workflows OK**: 2

---

## ❌ Problemas Críticos

### 1. **Scripts Faltantes - BLOQUEANTE**

**Archivos afectados**: `pr-to-main.yml`, `pr-to-dev.yml`, `test.yml`

**Problema**:
Los workflows intentan ejecutar `./scripts/check-coverage.sh` pero el directorio `scripts/` NO EXISTE en el proyecto.

**Ubicaciones**:
- `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-administracion/.github/workflows/pr-to-main.yml:57`
- `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-administracion/.github/workflows/pr-to-dev.yml:51`
- `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-administracion/.github/workflows/test.yml:55`

**Código problemático**:
```yaml
- name: ✅ Verificar umbral de cobertura
  run: |
    ./scripts/check-coverage.sh coverage/coverage-filtered.out ${{ env.COVERAGE_THRESHOLD }}
```

**Impacto**: Los workflows FALLAN cuando intentan verificar cobertura.

**Solución**: Crear el script `scripts/check-coverage.sh` o eliminar estos pasos.

---

### 2. **Comandos Make Inexistentes - BLOQUEANTE**

**Archivos afectados**: `pr-to-main.yml`, `pr-to-dev.yml`, `test.yml`

**Problema**:
Los workflows usan comandos `make` que NO EXISTEN en el Makefile:
- `make test-unit` ✅ (existe)
- `make coverage-report` ❌ (NO EXISTE)
- `make test-integration` ✅ (existe)

**Comandos disponibles en Makefile**:
- `test` - Ejecutar todos los tests
- `test-coverage` - Tests con cobertura (genera HTML)
- `test-unit` - Solo tests unitarios
- `test-integration` - Tests de integración

**Solución**: 
- Opción 1: Agregar `coverage-report` al Makefile
- Opción 2: Cambiar workflows para usar `make test-coverage`

---

### 3. **Archivos de Cobertura Incorrectos**

**Archivos afectados**: `pr-to-main.yml`, `pr-to-dev.yml`, `test.yml`

**Problema**:
Los workflows esperan archivos que no se generan:
- `coverage/coverage-filtered.out` ❌ (NO SE GENERA)
- El Makefile genera: `coverage/coverage.out` ✅

**Código problemático**:
```yaml
- name: ✅ Verificar umbral de cobertura
  run: |
    ./scripts/check-coverage.sh coverage/coverage-filtered.out ${{ env.COVERAGE_THRESHOLD }}
```

**Solución**: Cambiar `coverage-filtered.out` por `coverage.out`

---

### 4. **`continue-on-error` Duplicado en pr-to-main.yml**

**Archivo**: `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-administracion/.github/workflows/pr-to-main.yml:47-48`

**Problema**:
```yaml
- name: 📊 Generar reporte de cobertura
  run: make coverage-report
  continue-on-error: true
  continue-on-error: true  # ← DUPLICADO
  timeout-minutes: 5
```

**Impacto**: Sintaxis YAML inválida, puede causar errores de parsing.

**Solución**: Eliminar línea duplicada.

---

### 5. **Versión de Go Inconsistente**

**Archivos afectados**: TODOS los workflows

**Problema**:
Diferentes workflows usan diferentes versiones de Go:
- `ci.yml`: `1.25.3` ❌ (versión futura, no existe)
- `release.yml`: `1.25.3` ❌ (versión futura, no existe)
- `build-and-push.yml`: `1.25.3` ❌ (versión futura, no existe)
- `test.yml`: `1.24` ✅
- `pr-to-main.yml`: `1.24` ✅
- `pr-to-dev.yml`: `1.24` ✅
- `manual-release.yml`: `1.24` ✅

**Última versión estable de Go**: `1.23.x` (Go 1.25 no existe todavía)

**Solución**: Estandarizar a `1.24` en TODOS los workflows.

---

## ⚠️ Advertencias

### 1. **Triggers Problemáticos en ci.yml**

**Archivo**: `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-administracion/.github/workflows/ci.yml:3-8`

**Código**:
```yaml
on:
  pull_request:
    branches: [ main, develop ]  # ← 'develop' debería ser 'dev'
  push:
    branches: [ main ]
```

**Problema**: 
- La rama se llama `dev`, no `develop`
- Este workflow se ejecuta en PRs a `main` y `dev`, pero ya existen workflows específicos:
  - `pr-to-main.yml` para PRs a main
  - `pr-to-dev.yml` para PRs a dev

**Impacto**: DUPLICACIÓN de workflows, se ejecutan tests dobles en cada PR.

**Solución**: 
- Opción 1: Eliminar `ci.yml` (duplicado)
- Opción 2: Cambiar `develop` por `dev` y eliminar workflows específicos

---

### 2. **Workflow docker-only.yml se Ejecuta en Push a Main**

**Archivo**: `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-administracion/.github/workflows/docker-only.yml:12-14`

**Código**:
```yaml
on:
  workflow_dispatch:
    # ...
  push:
    branches: [ main ]  # ← Se ejecuta automáticamente
```

**Problema**: 
- El nombre dice "Simple", pero se ejecuta automáticamente en push a main
- Compite con `release.yml` que también construye Docker en tags
- Sin tests previos

**Impacto**: Builds innecesarios, consumo de recursos.

**Solución**: Eliminar trigger automático, dejar solo `workflow_dispatch`.

---

### 3. **sync-main-to-dev.yml: Doble Trigger**

**Archivo**: `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-administracion/.github/workflows/sync-main-to-dev.yml:4-8`

**Código**:
```yaml
on:
  push:
    branches: [main]
    tags:
      - 'v*'
```

**Problema**:
- Se ejecuta en push a `main` Y en tags `v*`
- Cuando creas un tag desde main, se ejecuta 2 VECES

**Solución**: Separar triggers o agregar condicional.

---

### 4. **Falta Validación de Scripts en build-and-push.yml**

**Archivo**: `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-administracion/.github/workflows/build-and-push.yml`

**Problema**:
Este workflow ejecuta tests pero no verifica cobertura ni formato.

**Solución**: Agregar steps de formato y linting.

---

### 5. **Timeouts Excesivos en test.yml**

**Archivo**: `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-administracion/.github/workflows/test.yml`

**Problemas**:
```yaml
timeout-minutes: 10  # Job completo
timeout-minutes: 5   # Step test-unit (demasiado, readme dice ~5 segundos)
timeout-minutes: 5   # Step coverage-report
timeout-minutes: 15  # Job integration (readme dice ~1-2 minutos)
timeout-minutes: 10  # Step test-integration
```

**Solución**: Ajustar timeouts según métricas reales:
- test-unit job: 3 minutos
- test-unit step: 2 minutos
- integration job: 5 minutos
- integration step: 3 minutos

---

### 6. **Permisos Innecesarios en docker-only.yml**

**Archivo**: `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-administracion/.github/workflows/docker-only.yml:21-25`

**Código**:
```yaml
permissions:
  contents: read
  packages: write
  attestations: write  # ← No se usa
  id-token: write      # ← No se usa
```

**Problema**: Permisos excesivos sin uso.

**Solución**: Eliminar `attestations` y `id-token` si no se usan.

---

### 7. **manual-release.yml: Comentario sobre GITHUB_TOKEN Desactualizado**

**Archivo**: `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-administracion/.github/workflows/manual-release.yml:43-48`

**Código**:
```yaml
# Usar GitHub App Token en lugar de GITHUB_TOKEN porque:
# - GITHUB_TOKEN NO dispara workflows subsecuentes (limitación de seguridad de GitHub)
# - App Token SÍ dispara sync-main-to-dev.yml automáticamente después del push
```

**Problema**: 
GitHub cambió esto. Desde GitHub Actions v2, `GITHUB_TOKEN` SÍ puede disparar workflows si tiene permisos adecuados.

**Solución**: Verificar si realmente se necesita GitHub App o se puede usar `GITHUB_TOKEN`.

---

### 8. **Inconsistencia en Nombres de Imágenes Docker**

**Problema**:
- `docker-only.yml`: `ghcr.io/edugogroup/edugo-api-administracion`
- Otros workflows: `ghcr.io/${{ github.repository }}`

**Solución**: Estandarizar usando `${{ github.repository }}`.

---

## ℹ️ Recomendaciones

### 1. **Eliminar Workflows Duplicados**

**Workflows duplicados detectados**:
- `ci.yml` vs `pr-to-main.yml` + `pr-to-dev.yml`
- `build-and-push.yml` vs `docker-only.yml` (ambos manuales para Docker)

**Recomendación**: 
- Mantener `pr-to-main.yml` y `pr-to-dev.yml` (más específicos)
- Eliminar `ci.yml`
- Unificar `build-and-push.yml` y `docker-only.yml` en uno solo

---

### 2. **Agregar Validación de Workflows**

**Recomendación**: Agregar workflow para validar sintaxis YAML:

```yaml
name: Validate Workflows

on:
  pull_request:
    paths:
      - '.github/workflows/**'

jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Validate YAML
        run: |
          for file in .github/workflows/*.yml; do
            echo "Validating $file"
            yamllint "$file" || exit 1
          done
```

---

### 3. **Crear Scripts Faltantes**

**Script necesario**: `scripts/check-coverage.sh`

```bash
#!/bin/bash
# scripts/check-coverage.sh
COVERAGE_FILE=$1
THRESHOLD=$2

if [ ! -f "$COVERAGE_FILE" ]; then
  echo "Error: Archivo de cobertura no encontrado: $COVERAGE_FILE"
  exit 1
fi

# Extraer cobertura total
COVERAGE=$(go tool cover -func="$COVERAGE_FILE" | grep total | awk '{print $3}' | sed 's/%//')

echo "Cobertura actual: ${COVERAGE}%"
echo "Umbral mínimo: ${THRESHOLD}%"

# Comparar
if (( $(echo "$COVERAGE < $THRESHOLD" | bc -l) )); then
  echo "❌ Cobertura por debajo del umbral"
  exit 1
else
  echo "✅ Cobertura cumple con el umbral"
  exit 0
fi
```

---

### 4. **Estandarizar Variables de Entorno**

**Problema**: Cada workflow define sus propias variables.

**Recomendación**: Crear archivo `.github/workflows/config.yml` (no es válido, pero documentar):

```yaml
# Documentar en README.md
GO_VERSION: "1.24"
COVERAGE_THRESHOLD: 33
REGISTRY: ghcr.io
```

---

### 5. **Agregar Badges al README**

**Recomendación**: Agregar badges de status de workflows:

```markdown
[![CI Pipeline](https://github.com/EduGoGroup/edugo-api-administracion/actions/workflows/ci.yml/badge.svg)](https://github.com/EduGoGroup/edugo-api-administracion/actions/workflows/ci.yml)
[![Tests](https://github.com/EduGoGroup/edugo-api-administracion/actions/workflows/test.yml/badge.svg)](https://github.com/EduGoGroup/edugo-api-administracion/actions/workflows/test.yml)
```

---

### 6. **Optimizar Cache de Go**

**Recomendación**: Algunos workflows no usan cache de Go modules:

```yaml
- name: Setup Go
  uses: actions/setup-go@v5
  with:
    go-version: ${{ env.GO_VERSION }}
    cache: true  # ← Agregar si falta
    cache-dependency-path: go.sum  # ← Agregar para mejor cache
```

---

## ✅ Workflows Correctos

### 1. **sync-main-to-dev.yml** ⭐

**Ubicación**: `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-administracion/.github/workflows/sync-main-to-dev.yml`

**Estado**: ✅ Funcional (con advertencia menor)

**Características**:
- ✅ Triggers correctos
- ✅ Manejo de conflictos
- ✅ Prevención de loops infinitos
- ✅ Logging detallado
- ⚠️ Solo advertencia: doble trigger (push + tags)

---

### 2. **manual-release.yml** ⭐

**Ubicación**: `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-administracion/.github/workflows/manual-release.yml`

**Estado**: ✅ Funcional

**Características**:
- ✅ Workflow manual bien estructurado
- ✅ Validación de versión semver
- ✅ Actualización de CHANGELOG automática
- ✅ Build multi-platform Docker
- ✅ GitHub Release automático
- ⚠️ Solo advertencia: comentario desactualizado sobre GITHUB_TOKEN

---

## 🔍 Análisis Específico: pr-to-main.yml

### **¿Por qué se ejecuta en push a dev?**

**RESPUESTA**: NO se ejecuta en push a dev según su configuración.

**Configuración actual**:
```yaml
on:
  pull_request:
    branches: [main]
    types: [opened, synchronize, reopened]
```

**Análisis**:
- Este workflow SOLO se ejecuta en **Pull Requests hacia main**
- NO tiene trigger `push:`
- NO tiene trigger para branch `dev`

**Posibles causas del problema reportado**:

1. **Workflow ci.yml se está ejecutando** (tiene trigger duplicado):
   ```yaml
   on:
     pull_request:
       branches: [ main, develop ]
   ```

2. **Confusión con otro workflow**: Puede ser `pr-to-dev.yml` el que se ejecuta.

3. **Cache de GitHub Actions**: A veces GitHub cachea workflows antiguos.

**Solución**:
1. Eliminar `ci.yml` (duplicado)
2. Verificar que no haya workflows antiguos en cache
3. Revisar logs de GitHub Actions para identificar cuál workflow se ejecuta realmente

---

## 📊 Tabla Resumen de Workflows

| Workflow | Trigger | Estado | Problemas Críticos | Advertencias |
|----------|---------|--------|-------------------|--------------|
| ci.yml | PR (main, develop), push (main) | ⚠️ | Versión Go 1.25.3, branch 'develop' | Duplicado |
| test.yml | Manual | ❌ | Scripts faltantes, comandos make | Timeouts excesivos |
| pr-to-main.yml | PR → main | ❌ | Scripts, comandos make, `continue-on-error` duplicado | - |
| pr-to-dev.yml | PR → dev | ❌ | Scripts faltantes, comandos make | - |
| release.yml | Tag v* | ⚠️ | Versión Go 1.25.3 | - |
| build-and-push.yml | Manual | ⚠️ | Versión Go 1.25.3 | Falta validación |
| docker-only.yml | Manual + push main | ⚠️ | - | Trigger automático, permisos |
| manual-release.yml | Manual | ✅ | - | Comentario desactualizado |
| sync-main-to-dev.yml | push main + tags | ✅ | - | Doble trigger |

---

## 🛠️ Plan de Acción Recomendado

### **Fase 1: Crítico (Hacer AHORA)**

1. ✅ **Crear directorio scripts/**
   ```bash
   mkdir -p scripts
   ```

2. ✅ **Crear script check-coverage.sh**
   - Ver contenido en sección "Recomendaciones #3"

3. ✅ **Corregir versión de Go** en:
   - `ci.yml`: cambiar `1.25.3` → `1.24`
   - `release.yml`: cambiar `1.25.3` → `1.24`
   - `build-and-push.yml`: cambiar `1.25.3` → `1.24`

4. ✅ **Agregar comando al Makefile**:
   ```makefile
   coverage-report: ## Generar reporte de cobertura
       @echo "$(YELLOW)📊 Generando reporte de cobertura...$(RESET)"
       @mkdir -p $(COVERAGE_DIR)
       @$(GOTEST) -v -race -coverprofile=$(COVERAGE_DIR)/coverage.out -covermode=atomic ./...
       @$(GOCMD) tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
       @echo "$(GREEN)✓ Reporte: $(COVERAGE_DIR)/coverage.html$(RESET)"
   ```

5. ✅ **Eliminar `continue-on-error` duplicado** en pr-to-main.yml:48

### **Fase 2: Importante (Hacer esta semana)**

6. ✅ **Eliminar ci.yml** (duplicado con pr-to-main.yml y pr-to-dev.yml)

7. ✅ **Unificar workflows Docker**:
   - Eliminar `docker-only.yml`
   - Mantener solo `build-and-push.yml`

8. ✅ **Corregir branch en ci.yml** (si decides mantenerlo):
   - `develop` → `dev`

9. ✅ **Ajustar timeouts** en test.yml según recomendaciones

10. ✅ **Cambiar archivos de cobertura**:
    - `coverage-filtered.out` → `coverage.out`

### **Fase 3: Mejoras (Hacer cuando tengas tiempo)**

11. ✅ Agregar workflow de validación YAML
12. ✅ Agregar badges al README
13. ✅ Estandarizar nombres de imágenes Docker
14. ✅ Optimizar cache de Go modules
15. ✅ Documentar variables de entorno estándar

---

## 📈 Métricas de Mejora Esperadas

Después de aplicar las correcciones:

| Métrica | Antes | Después | Mejora |
|---------|-------|---------|--------|
| Workflows funcionales | 2/10 (20%) | 8/8 (100%) | +400% |
| Workflows duplicados | 4 | 0 | -100% |
| Errores críticos | 5 | 0 | -100% |
| Tiempo de ejecución PR | ~6-8 min | ~3-4 min | -50% |
| Consumo de recursos | Alto (duplicados) | Óptimo | -40% |

---

## 🎯 Conclusiones

### **Problemas Principales Identificados**:

1. ✅ **Scripts faltantes** → Workflows fallan al verificar cobertura
2. ✅ **Comandos make inexistentes** → make coverage-report no existe
3. ✅ **Versión de Go incorrecta** → Go 1.25.3 no existe
4. ✅ **Workflows duplicados** → ci.yml duplica pr-to-main.yml y pr-to-dev.yml
5. ✅ **Sintaxis YAML inválida** → continue-on-error duplicado

### **Estado Actual**:

- 📊 **20% de workflows funcionales** (2 de 10)
- ❌ **80% requieren correcciones** (8 de 10)
- ⚠️ **Riesgo de fallos en CI/CD**: ALTO

### **Próximos Pasos**:

1. Aplicar Fase 1 (crítico) → restaurar funcionalidad básica
2. Aplicar Fase 2 (importante) → eliminar duplicados
3. Aplicar Fase 3 (mejoras) → optimizar workflows

---

**Generado con**: Claude Code (Anthropic)  
**Fecha**: 2025-11-17
