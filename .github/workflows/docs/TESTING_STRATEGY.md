# 🧪 Estrategia de Testing en CI/CD - Con Controles ON/OFF

## 🎯 Filosofía

**"Tests automáticos en CI, pero con control total del desarrollador"**

- ✅ Tests automáticos en PRs (calidad garantizada)
- ✅ Controles para deshabilitar cuando sea necesario
- ✅ Tests rápidos primero, lentos opcionales
- ✅ Armonizado con workflows existentes
- ✅ No bloquea desarrollo en features

---

## 🔧 Controles ON/OFF Implementados

### **1. Variable de Ambiente en Workflow**

Cada workflow de testing tiene una variable para habilitarlo/deshabilitarlo:

```yaml
env:
  ENABLE_UNIT_TESTS: true        # ← Cambiar a false para deshabilitar
  ENABLE_INTEGRATION_TESTS: false # ← true para habilitar
  ENABLE_COVERAGE_CHECK: true     # ← false para no fallar por cobertura
  COVERAGE_THRESHOLD: 60          # ← Ajustable según necesidad
```

### **2. Labels en Pull Requests**

Puedes controlar qué tests ejecutar usando labels:

```yaml
# En el PR, agregar labels:
skip-tests          # ← Salta TODOS los tests (emergencias)
skip-integration    # ← Solo tests unitarios
skip-coverage       # ← No verifica umbral de cobertura
run-full-suite      # ← Ejecuta TODO (unitarios + integración + coverage)
```

**Uso:**
```bash
# Crear PR sin tests automáticos
gh pr create --label "skip-tests"

# PR con todos los tests
gh pr create --label "run-full-suite"

# Agregar label a PR existente
gh pr edit 123 --add-label "skip-integration"
```

### **3. Archivo de Configuración**

```yaml
# .github/testing-config.yml
testing:
  # Global ON/OFF
  enabled: true
  
  # Por tipo de test
  unit_tests:
    enabled: true
    timeout: 5m
    
  integration_tests:
    enabled: false  # ← Deshabilitado por defecto en desarrollo
    timeout: 15m
    require_docker: true
    
  coverage:
    enabled: true
    threshold: 60
    fail_on_decrease: true  # ← Fallar si cobertura baja
    
  # Por branch
  branches:
    main:
      require_all_tests: true
      require_coverage: true
    dev:
      require_all_tests: false
      require_coverage: false
    feature:
      require_all_tests: false
      require_coverage: false
```

### **4. Manual Dispatch (Ejecutar cuando quieras)**

Todos los workflows tienen `workflow_dispatch` para ejecución manual:

```bash
# Ejecutar tests manualmente cuando quieras
gh workflow run test-unit.yml

# Con parámetros
gh workflow run test-coverage.yml \
  -f enable_integration=true \
  -f coverage_threshold=70
```

---

## 📋 Workflows de Testing Propuestos

### **Workflow 1: test-unit.yml** (Rápido, Siempre)

```yaml
name: Unit Tests

on:
  pull_request:
    branches: [main, dev]
  push:
    branches: [main]
  workflow_dispatch:  # Manual
    inputs:
      skip_tests:
        description: 'Saltar tests (emergencias)'
        type: boolean
        default: false

env:
  ENABLE_TESTS: true  # ← Control ON/OFF global

jobs:
  unit-tests:
    name: Tests Unitarios
    runs-on: ubuntu-latest
    
    # Skip si tiene label o input manual
    if: |
      (env.ENABLE_TESTS == 'true') &&
      (!contains(github.event.pull_request.labels.*.name, 'skip-tests')) &&
      (github.event.inputs.skip_tests != 'true')
    
    steps:
      - uses: actions/checkout@v4
      
      - uses: actions/setup-go@v5
        with:
          go-version: '1.25.3'
          cache: true
      
      - name: Configurar acceso a repos privados
        run: |
          git config --global url."https://${{ secrets.GITHUB_TOKEN }}@github.com/".insteadOf "https://github.com/"
        env:
          GOPRIVATE: github.com/EduGoGroup/*
      
      - name: Descargar dependencias
        run: go mod download
      
      - name: Ejecutar tests unitarios
        run: make test-unit
        timeout-minutes: 5
      
      - name: Resumen
        run: echo "✅ Tests unitarios pasaron correctamente"
```

**Características**:
- ⚡ Rápido (< 5 min)
- 🎯 Solo tests unitarios (sin Docker)
- 🔧 Control con label `skip-tests`
- 📍 Se ejecuta en PRs a main/dev

### **Workflow 2: test-coverage.yml** (Con Cobertura)

```yaml
name: Coverage Check

on:
  pull_request:
    branches: [main, dev]
  workflow_dispatch:
    inputs:
      threshold:
        description: 'Umbral de cobertura (%)'
        type: number
        default: 60
      fail_on_low_coverage:
        description: 'Fallar si cobertura < umbral'
        type: boolean
        default: true

env:
  ENABLE_COVERAGE_CHECK: true  # ← Control ON/OFF
  COVERAGE_THRESHOLD: 60

jobs:
  coverage:
    name: Verificar Cobertura
    runs-on: ubuntu-latest
    
    if: |
      (env.ENABLE_COVERAGE_CHECK == 'true') &&
      (!contains(github.event.pull_request.labels.*.name, 'skip-coverage'))
    
    steps:
      - uses: actions/checkout@v4
      
      - uses: actions/setup-go@v5
        with:
          go-version: '1.25.3'
          cache: true
      
      - name: Configurar acceso a repos privados
        run: |
          git config --global url."https://${{ secrets.GITHUB_TOKEN }}@github.com/".insteadOf "https://github.com/"
        env:
          GOPRIVATE: github.com/EduGoGroup/*
      
      - name: Generar reporte de cobertura
        run: make coverage-report
      
      - name: Verificar umbral
        if: github.event.inputs.fail_on_low_coverage != 'false'
        run: |
          THRESHOLD=${{ github.event.inputs.threshold || env.COVERAGE_THRESHOLD }}
          make coverage-check THRESHOLD=$THRESHOLD || {
            echo "⚠️ Cobertura por debajo de ${THRESHOLD}%"
            echo "💡 Agrega label 'skip-coverage' al PR si es temporal"
            exit 1
          }
      
      - name: Upload reporte
        uses: actions/upload-artifact@v4
        with:
          name: coverage-report
          path: coverage/
      
      - name: Comentar en PR
        if: github.event_name == 'pull_request'
        uses: actions/github-script@v7
        with:
          script: |
            const fs = require('fs');
            const { execSync } = require('child_process');
            
            const coverage = execSync('go tool cover -func=coverage/coverage-filtered.out | grep total').toString();
            const match = coverage.match(/(\d+\.\d+)%/);
            const percentage = match ? match[1] : 'N/A';
            
            const body = `## 📊 Reporte de Cobertura
            
**Cobertura total:** ${percentage}%
**Umbral mínimo:** ${{ env.COVERAGE_THRESHOLD }}%
            
${percentage >= ${{ env.COVERAGE_THRESHOLD }} ? '✅ Cobertura aprobada' : '⚠️ Cobertura por debajo del umbral'}

[Ver reporte completo](https://github.com/${{ github.repository }}/actions/runs/${{ github.run_id }})`;
            
            github.rest.issues.createComment({
              owner: context.repo.owner,
              repo: context.repo.repo,
              issue_number: context.issue.number,
              body: body
            });
```

**Características**:
- 📊 Reporte de cobertura con filtrado
- 🎯 Umbral configurable (default 60%)
- 💬 Comentario automático en PR
- 🔧 Control con label `skip-coverage`
- 🛑 Puede fallar o solo advertir

### **Workflow 3: test-integration.yml** (Opcional, Pesado)

```yaml
name: Integration Tests

on:
  workflow_dispatch:  # ← SOLO MANUAL (por defecto)
    inputs:
      enable_tests:
        description: 'Ejecutar tests de integración'
        type: boolean
        default: true
  # Opcional: Descomentar para ejecutar en PRs específicos
  # pull_request:
  #   branches: [main]
  #   types: [labeled]

env:
  ENABLE_INTEGRATION_TESTS: false  # ← Deshabilitado por defecto

jobs:
  integration-tests:
    name: Tests de Integración
    runs-on: ubuntu-latest
    
    # Solo si:
    # - Se ejecuta manual Y enable_tests=true
    # - O el PR tiene label 'run-integration-tests'
    if: |
      (github.event_name == 'workflow_dispatch' && github.event.inputs.enable_tests == 'true') ||
      (contains(github.event.pull_request.labels.*.name, 'run-integration-tests'))
    
    steps:
      - uses: actions/checkout@v4
      
      - uses: actions/setup-go@v5
        with:
          go-version: '1.25.3'
          cache: true
      
      - name: Configurar acceso a repos privados
        run: |
          git config --global url."https://${{ secrets.GITHUB_TOKEN }}@github.com/".insteadOf "https://github.com/"
        env:
          GOPRIVATE: github.com/EduGoGroup/*
      
      - name: Ejecutar tests de integración
        run: RUN_INTEGRATION_TESTS=true make test-integration
        timeout-minutes: 15
      
      - name: Resumen
        run: echo "✅ Tests de integración completados"
```

**Características**:
- 🐳 Usa testcontainers (requiere Docker en runner)
- 🐢 Lento (~5-15 min)
- 🎯 SOLO manual o con label específico
- 📍 **No bloquea desarrollo normal**

---

## 🎮 Guía de Uso - Controles ON/OFF

### **Escenario 1: Desarrollo Normal**

```bash
# Tu PR automáticamente ejecuta:
- test-unit.yml ✅ (rápido, siempre)
- test-coverage.yml ✅ (cobertura)

# NO ejecuta:
- test-integration.yml ❌ (solo manual)
```

### **Escenario 2: Work in Progress (WIP)**

```bash
# Crear PR con label para saltar tests
gh pr create --draft --label "skip-tests"

# O agregar label después
gh pr edit <num> --add-label "skip-tests"

# Resultado: NO ejecuta tests automáticos
# Cuando esté listo, quitar label y convertir a ready
gh pr edit <num> --remove-label "skip-tests"
gh pr ready <num>
```

### **Escenario 3: Debugging de Tests Fallidos**

```bash
# 1. PR falla por tests
# 2. Agregar label temporal
gh pr edit <num> --add-label "skip-coverage"

# 3. Trabajar en fixes localmente
go test ./...
make coverage-check

# 4. Push de fixes
git push

# 5. Quitar label cuando esté listo
gh pr edit <num> --remove-label "skip-coverage"
```

### **Escenario 4: Validación Completa Antes de Release**

```bash
# Ejecutar TODO manualmente antes de mergear a main
gh pr edit <num> --add-label "run-full-suite"

# O ejecutar manual:
gh workflow run test-integration.yml -f enable_tests=true
```

### **Escenario 5: Deshabilitar TODOS los Tests Temporalmente**

```yaml
# Editar .github/workflows/test-unit.yml
env:
  ENABLE_TESTS: false  # ← Cambiar a false

# Commit y push
git add .github/workflows/test-unit.yml
git commit -m "ci: deshabilitar tests temporalmente"
git push
```

---

## 📊 Integración con Workflows Existentes

### **NO se duplica, se complementa:**

```
Workflows EXISTENTES (se mantienen):
├── ci.yml                    ← Validación general (formato, vet, build)
├── test.yml                  ← Cobertura manual
├── build-and-push.yml        ← Docker build
├── release.yml               ← Release tags
├── manual-release.yml        ← Release manual
└── sync-main-to-dev-ff.yml   ← Sincronización

Workflows NUEVOS (testing específico):
├── test-unit-quick.yml       ← Tests unitarios rápidos (NUEVO)
├── test-coverage-check.yml   ← Verificación de cobertura (NUEVO)
└── test-integration-manual.yml ← Integración solo manual (NUEVO)
```

### **Estrategia de Integración:**

#### **Opción A: Mejorar workflows existentes** ⭐ RECOMENDADO

```yaml
# Actualizar ci.yml para usar nuevos comandos make
- name: Ejecutar tests
  run: make test-unit  # ← Usa el nuevo comando

# Actualizar test.yml para usar filtrado
- name: Coverage
  run: make coverage-report  # ← Usa scripts de filtrado
```

#### **Opción B: Workflows separados**

```yaml
# Crear test-unit-quick.yml (no reemplaza ci.yml)
# Se ejecuta ADEMÁS de ci.yml pero más rápido
```

---

## 🎛️ Panel de Control Centralizado

Crear archivo `.github/testing-config.yml`:

```yaml
# Panel de control central para testing en CI/CD
# Editar este archivo para habilitar/deshabilitar tests

testing:
  # 🌐 Control Global
  enabled: true  # ← false para deshabilitar TODO
  
  # 🧪 Tests Unitarios
  unit_tests:
    enabled: true
    timeout_minutes: 5
    fail_on_error: true
    
  # 🐳 Tests de Integración
  integration_tests:
    enabled: false  # ← false por defecto (solo manual)
    timeout_minutes: 15
    fail_on_error: true
    require_label: "run-integration-tests"  # ← Requiere label en PR
    
  # 📊 Cobertura
  coverage:
    enabled: true
    threshold: 60
    fail_below_threshold: true  # ← false para solo advertir
    upload_to_codecov: true
    comment_on_pr: true
    
  # 🏷️ Control por Branch
  branch_rules:
    main:
      require_unit_tests: true
      require_coverage_check: true
      min_coverage: 60
      
    dev:
      require_unit_tests: true
      require_coverage_check: false  # ← Más permisivo
      min_coverage: 50
      
    feature:
      require_unit_tests: false  # ← No bloquea features
      require_coverage_check: false
      
  # 🏃 Performance
  optimization:
    cache_dependencies: true
    parallel_tests: true
    fail_fast: false  # ← Ejecutar todos aunque uno falle
```

**Uso en workflows:**
```yaml
- name: Cargar configuración
  id: config
  run: |
    CONFIG=$(cat .github/testing-config.yml)
    ENABLED=$(echo "$CONFIG" | yq '.testing.unit_tests.enabled')
    echo "enabled=$ENABLED" >> $GITHUB_OUTPUT

- name: Tests
  if: steps.config.outputs.enabled == 'true'
  run: make test-unit
```

---

## 🚦 Flujo Completo con Controles

### **Desarrollo en Feature Branch:**

```
1. git checkout -b feature/nueva-funcionalidad
2. Desarrollo local + tests locales
3. git push origin feature/nueva-funcionalidad
   
   ✅ NO ejecuta workflows (ahorra minutos)
   
4. gh pr create --base dev
   
   ✅ test-unit-quick.yml ejecuta (3-5 min)
   ✅ test-coverage-check.yml ejecuta (4-6 min)
   ❌ test-integration NO ejecuta (deshabilitado)
   
   Si necesitas integración:
   gh pr edit --add-label "run-integration-tests"
   
5. Merge a dev
   ✅ ci.yml ejecuta validación
   ✅ sync actualiza main si corresponde
```

### **Release a Main:**

```
1. gh pr create --base main --head dev --title "Release v0.1.7"
   
   ✅ test-unit-quick.yml ejecuta
   ✅ test-coverage-check.yml ejecuta (threshold 60%)
   ✅ test-integration ejecuta SI tiene label
   
2. Aprobar y mergear
   
3. Crear tag (dispara release.yml)
   git tag v0.1.7
   git push origin v0.1.7
   
   ✅ release.yml build Docker + GitHub Release
   ✅ sync-main-to-dev-ff.yml sincroniza
```

### **Emergencia: Deshabilitar Tests Temporalmente:**

```bash
# Opción 1: Label en PR (recomendado)
gh pr edit <num> --add-label "skip-tests"

# Opción 2: Editar config (afecta todos los PRs)
vim .github/testing-config.yml
# Cambiar testing.enabled: false
git commit -m "ci: deshabilitar tests temporalmente"
git push

# Opción 3: Editar workflow directamente
vim .github/workflows/test-unit-quick.yml
# Cambiar ENABLE_TESTS: false
```

---

## 🏅 Badges con Control

### **Badge Inteligente** (muestra estado real):

```markdown
<!-- README.md -->

# EduGo API Mobile

<!-- Badge de tests (solo si están habilitados) -->
![Tests](https://github.com/EduGoGroup/edugo-api-mobile/workflows/Unit%20Tests/badge.svg)

<!-- Badge de cobertura con Codecov -->
[![codecov](https://codecov.io/gh/EduGoGroup/edugo-api-mobile/branch/main/graph/badge.svg)](https://codecov.io/gh/EduGoGroup/edugo-api-mobile)

<!-- Badge de Go version -->
![Go Version](https://img.shields.io/github/go-mod/go-version/EduGoGroup/edugo-api-mobile)

<!-- Badge de último release -->
![Release](https://img.shields.io/github/v/release/EduGoGroup/edugo-api-mobile)
```

**Cómo se ven:**

```
✓ Tests Passing    Coverage 35.3%    Go 1.25    v0.1.6
```

**Qué es Codecov Badge:**
- Muestra % de cobertura en tiempo real
- Color verde (>70%), amarillo (40-70%), rojo (<40%)
- Click lleva a reporte detallado
- Se actualiza automáticamente con cada push

**Setup Codecov (1 vez):**
```bash
# 1. Ir a https://codecov.io/
# 2. Login con GitHub
# 3. Agregar repositorio edugo-api-mobile
# 4. Copiar token (o usar sin token si repo es público)
# 5. Agregar a secrets: CODECOV_TOKEN (opcional)
```

---

## 💡 Recomendaciones

### **Para Desarrollo (Actual):**

```yaml
ENABLE_UNIT_TESTS: true          # ← Siempre
ENABLE_INTEGRATION_TESTS: false  # ← Solo manual
ENABLE_COVERAGE_CHECK: true      # ← Sí, pero sin bloquear
COVERAGE_THRESHOLD: 60           # ← Meta a largo plazo
```

### **Para Producción (Futuro):**

```yaml
ENABLE_UNIT_TESTS: true
ENABLE_INTEGRATION_TESTS: true   # ← Habilitar cuando estén estables
ENABLE_COVERAGE_CHECK: true
COVERAGE_THRESHOLD: 70           # ← Más estricto
```

### **Migración Gradual:**

```
Mes 1 (AHORA):
  ✅ Tests unitarios automáticos
  ⚠️ Cobertura informativa (no bloquea)
  ❌ Integración solo manual

Mes 2:
  ✅ Tests unitarios automáticos
  ⚠️ Cobertura bloquea si baja mucho (< 50%)
  🔄 Integración en PRs a main (con label)

Mes 3+:
  ✅ Tests unitarios automáticos
  ❌ Cobertura bloquea si < 60%
  ✅ Integración automática en PRs a main
```

---

**¿Te gusta este diseño?** Tiene:
- ✅ Control total con labels y variables
- ✅ No rompe workflows existentes
- ✅ Migración gradual
- ✅ Emergencias cubiertas

¿Procedo a implementarlo?
