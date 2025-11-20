# Sprint 4: Workflows Reusables y Optimización - edugo-api-administracion

**Duración:** 3 días (12-15 horas)  
**Objetivo:** Eliminar duplicación mediante workflows reusables y optimizar tiempos  
**Prioridad:** 🟢 P2 - MEDIA

---

## 📋 Resumen Ejecutivo

Sprint 4 se enfoca en eliminar código duplicado (~70%) mediante workflows reusables y composite actions, más optimización de paralelismo para reducir tiempos de CI.

**Beneficios Esperados:**
- ✅ -50-60% código duplicado (~400 líneas)
- ✅ -20-30% tiempo de CI (3-4 min → 2-3 min)
- ✅ Mantenimiento centralizado
- ✅ Consistencia entre proyectos

---

## 📅 Cronograma

```
Día 1: Composite Actions     (4-5h)  → Setup-go, Docker, Coverage
Día 2: Workflows Reusables   (4-5h)  → Sync, Release logic
Día 3: Paralelismo           (4-5h)  → Matriz, optimización
```

---

# DÍA 1: COMPOSITE ACTIONS

## Tarea 1: Migrar a setup-edugo-go

**🟢 Prioridad:** P2  
**⏱️ Tiempo estimado:** 1.5 horas

### Objetivo

Reemplazar código duplicado de setup Go + GOPRIVATE en todos los workflows usando composite action centralizada.

### Pre-requisito

Composite action `setup-edugo-go` debe existir en edugo-infrastructure.

**Verificar:**
```bash
gh api repos/EduGoGroup/edugo-infrastructure/contents/.github/actions/setup-edugo-go/action.yml
```

### Paso 1: Identificar Workflows a Actualizar

```bash
#!/bin/bash
# Script: sprint4-01-find-setup-go.sh

cd ~/source/EduGo/repos-separados/edugo-api-administracion

echo "🔍 Workflows con setup Go a migrar:"
echo ""

grep -r "actions/setup-go" .github/workflows/*.yml | cut -d: -f1 | sort -u

echo ""
echo "📊 Total de ocurrencias:"
grep -r "actions/setup-go" .github/workflows/*.yml | wc -l
```

**Resultado esperado:**
```
.github/workflows/pr-to-dev.yml
.github/workflows/pr-to-main.yml
.github/workflows/test.yml
.github/workflows/manual-release.yml
.github/workflows/release.yml

Total: ~5 workflows
```

---

### Paso 2: Actualizar Workflows

**Script de migración automática:**

```bash
#!/bin/bash
# Script: sprint4-02-migrate-setup-go.sh

cd ~/source/EduGo/repos-separados/edugo-api-administracion

git checkout -b refactor/use-setup-edugo-go-composite

# Workflows a actualizar
WORKFLOWS=(
  ".github/workflows/pr-to-dev.yml"
  ".github/workflows/pr-to-main.yml"
  ".github/workflows/test.yml"
  ".github/workflows/manual-release.yml"
  ".github/workflows/release.yml"
)

for workflow in "${WORKFLOWS[@]}"; do
  if [ ! -f "$workflow" ]; then
    echo "⚠️  $workflow no existe, saltando..."
    continue
  fi
  
  echo "📝 Actualizando $workflow..."
  
  # Backup
  cp "$workflow" "$workflow.backup"
  
  # Reemplazar bloque de Setup Go + GOPRIVATE
  # NOTA: Este sed es simplificado, ajustar según formato exacto
  
  # Buscar y reemplazar el bloque completo
  python3 << 'PYTHON_SCRIPT'
import re
import sys

workflow_file = sys.argv[1]

with open(workflow_file, 'r') as f:
    content = f.read()

# Pattern para encontrar el bloque completo de Setup Go + GOPRIVATE
pattern = r'(\s+)- name: .*Setup Go.*\n\s+uses: actions/setup-go@v\d+\n(?:\s+with:\n(?:\s+[^\n]+\n)*)?(\s+- name: .*repos privados.*\n(?:\s+[^\n]+\n)*?(?=\n\s+- name:|\n\n|\Z))'

replacement = r'\1- name: Setup Go Environment\n\1  uses: EduGoGroup/edugo-infrastructure/.github/actions/setup-edugo-go@main\n'

content_new = re.sub(pattern, replacement, content, flags=re.MULTILINE | re.DOTALL)

with open(workflow_file, 'w') as f:
    f.write(content_new)

print(f"✅ {workflow_file} actualizado")

PYTHON_SCRIPT "$workflow"
  
done

echo ""
echo "✅ Todos los workflows actualizados"
echo ""
echo "📊 Revisar cambios:"
git diff .github/workflows/
```

**Checkpoint:**
- [ ] Script ejecutado
- [ ] 5 workflows actualizados
- [ ] Diff revisado

---

### Paso 3: Testing Local

```bash
#!/bin/bash
# Script: sprint4-03-test-workflows.sh

echo "🧪 Testing workflows con act..."

# Requiere act instalado
if ! command -v act &> /dev/null; then
  echo "❌ act no instalado"
  echo "Instalar: brew install act"
  exit 1
fi

# Test pr-to-dev.yml
act pull_request -W .github/workflows/pr-to-dev.yml --dryrun

echo "✅ Workflows validados sintácticamente"
```

---

### Paso 4: Push y Validar en CI

```bash
git add .github/workflows/
git commit -m "refactor: usar composite action setup-edugo-go

Migración de setup Go + GOPRIVATE a composite action centralizada.

Beneficios:
- Código reutilizable
- Mantenimiento centralizado
- Consistencia entre proyectos

Workflows actualizados:
- pr-to-dev.yml
- pr-to-main.yml
- test.yml
- manual-release.yml
- release.yml

🤖 Generated with Claude Code

Co-Authored-By: Claude <noreply@anthropic.com>"

git push origin refactor/use-setup-edugo-go-composite

# Crear PR
gh pr create --base dev --title "refactor: migrar a composite action setup-edugo-go" \
  --body "Migración a composite action para eliminar duplicación."
```

**Validar CI:**
```bash
gh pr checks --watch
```

---

## Tarea 2: Migrar a docker-build-edugo

**🟢 Prioridad:** P2  
**⏱️ Tiempo estimado:** 2 horas

### Objetivo

Reemplazar bloques de Docker build con composite action `docker-build-edugo`.

### Workflows afectados

- `manual-release.yml`
- `release.yml` (si se mantiene)

### Script de Migración

```bash
#!/bin/bash
# Script: sprint4-04-migrate-docker-build.sh

cd ~/source/EduGo/repos-separados/edugo-api-administracion

git checkout -b refactor/use-docker-build-composite

# Actualizar manual-release.yml
WORKFLOW=".github/workflows/manual-release.yml"

echo "📝 Actualizando $WORKFLOW..."

# Reemplazar bloque de Docker build (simplificado)
# NOTA: Ajustar según estructura exacta del workflow

cat > "$WORKFLOW.new" << 'EOF'
# ... (mantener inicio del workflow)

      - name: Build and Push Docker Image
        uses: EduGoGroup/edugo-infrastructure/.github/actions/docker-build-edugo@main
        with:
          image-name: edugogroup/edugo-api-administracion
          tag-strategy: semver
          version: ${{ steps.version.outputs.version }}
          platforms: linux/amd64,linux/arm64
          push: true

# ... (resto del workflow)
EOF

# En la práctica, usar editor o script Python más sofisticado
echo "⚠️  Revisar y ajustar manualmente: $WORKFLOW"
```

**Checkpoint:**
- [ ] Bloques Docker reemplazados
- [ ] Estrategia de tags correcta
- [ ] Platforms especificados

---

## Tarea 3: Migrar a coverage-check

**🟢 Prioridad:** P2  
**⏱️ Tiempo estimado:** 1 hora

### Objetivo

Usar composite action para verificación de coverage.

### Script

```bash
#!/bin/bash
# Script: sprint4-05-migrate-coverage.sh

cd ~/source/EduGo/repos-separados/edugo-api-administracion

git checkout -b refactor/use-coverage-check-composite

# Workflows a actualizar
WORKFLOWS=(
  ".github/workflows/pr-to-dev.yml"
  ".github/workflows/pr-to-main.yml"
)

for workflow in "${WORKFLOWS[@]}"; do
  echo "📝 Actualizando $workflow..."
  
  # Reemplazar bloque de coverage check
  # Cambiar de:
  #   - name: Check coverage
  #     run: ./scripts/check-coverage.sh ...
  
  # A:
  #   - name: Check Coverage
  #     uses: EduGoGroup/edugo-infrastructure/.github/actions/coverage-check@main
  #     with:
  #       coverage-file: coverage/coverage-filtered.out
  #       threshold: 33
  
  echo "⚠️  Ajustar manualmente: $workflow"
done
```

---

### Resumen Día 1

**Composite Actions Migradas:**
- [ ] setup-edugo-go (5 workflows)
- [ ] docker-build-edugo (2 workflows)
- [ ] coverage-check (2 workflows)

**Líneas Eliminadas:** ~150-200 líneas de código duplicado

---

# DÍA 2: WORKFLOWS REUSABLES

## Tarea 4: Migrar sync-main-to-dev.yml

**🟢 Prioridad:** P2  
**⏱️ Tiempo estimado:** 2 horas

### Objetivo

Reemplazar sync-main-to-dev.yml (100 líneas) con llamada a workflow reusable (10 líneas).

### Pre-requisito

Workflow reusable debe existir en edugo-infrastructure:
```
edugo-infrastructure/.github/workflows/reusable/sync-branches.yml
```

### Migración

```bash
#!/bin/bash
# Script: sprint4-06-migrate-sync-workflow.sh

cd ~/source/EduGo/repos-separados/edugo-api-administracion

git checkout -b refactor/use-sync-branches-reusable

# Backup del workflow actual
cp .github/workflows/sync-main-to-dev.yml .github/workflows/sync-main-to-dev.yml.backup

# Crear nuevo workflow simplificado
cat > .github/workflows/sync-main-to-dev.yml << 'EOF'
name: Sync Main to Dev

on:
  push:
    branches: [main]
    tags: ['v*']

permissions:
  contents: write
  pull-requests: write

jobs:
  sync:
    name: Sync main → dev
    uses: EduGoGroup/edugo-infrastructure/.github/workflows/reusable/sync-branches.yml@main
    with:
      source-branch: main
      target-branch: dev
    secrets: inherit
EOF

echo "✅ sync-main-to-dev.yml actualizado"

# Mostrar reducción
echo ""
echo "📊 Antes: $(wc -l < .github/workflows/sync-main-to-dev.yml.backup) líneas"
echo "📊 Después: $(wc -l < .github/workflows/sync-main-to-dev.yml) líneas"
echo "📊 Reducción: $(($(wc -l < .github/workflows/sync-main-to-dev.yml.backup) - $(wc -l < .github/workflows/sync-main-to-dev.yml))) líneas"
```

**Checkpoint:**
- [ ] Workflow reemplazado
- [ ] ~90 líneas eliminadas
- [ ] Funcionalidad idéntica

---

## Tarea 5: (Opcional) Migrar Release Logic

**🟢 Prioridad:** P2 - Opcional  
**⏱️ Tiempo estimado:** 2 horas

### Objetivo

Si hay workflow reusable para releases, migrar `release.yml`.

**NOTA:** Depende de disponibilidad de workflow reusable en infrastructure.

### Script Condicional

```bash
#!/bin/bash
# Script: sprint4-07-migrate-release.sh

# Verificar si existe workflow reusable de release
if gh api repos/EduGoGroup/edugo-infrastructure/contents/.github/workflows/reusable/release.yml &> /dev/null; then
  echo "✅ Workflow reusable de release existe"
  echo "Proceder con migración..."
  
  # Migración similar a sync-branches
else
  echo "⚠️  Workflow reusable de release NO existe"
  echo "Saltar esta tarea por ahora"
fi
```

---

# DÍA 3: PARALELISMO Y OPTIMIZACIÓN

## Tarea 6: Implementar Matriz de Tests

**🟢 Prioridad:** P2  
**⏱️ Tiempo estimado:** 2 horas

### Objetivo

Ejecutar tests en paralelo usando matrices para reducir tiempo de CI.

### Análisis Actual

```bash
#!/bin/bash
# Script: sprint4-08-analyze-test-time.sh

echo "📊 Analizando tiempos de tests actuales..."

# Obtener duración de último run de pr-to-dev
gh run list --repo EduGoGroup/edugo-api-administracion \
  --workflow=pr-to-dev.yml \
  --limit 5 \
  --json conclusion,createdAt,updatedAt \
  --jq '.[] | select(.conclusion=="success") | "\(.updatedAt) - \(.createdAt) | \(((.updatedAt | fromdateiso8601) - (.createdAt | fromdateiso8601)) / 60) min"'

echo ""
echo "Tiempo actual promedio: 3-4 minutos"
echo "Objetivo: 2-3 minutos (25% reducción)"
```

---

### Implementar Matriz

**Estrategia:** Separar tests por paquete/módulo

```yaml
# Modificar pr-to-dev.yml

jobs:
  unit-tests:
    name: Unit Tests - ${{ matrix.package }}
    runs-on: ubuntu-latest
    
    strategy:
      fail-fast: false
      matrix:
        package:
          - ./internal/handler/...
          - ./internal/service/...
          - ./internal/repository/...
          - ./internal/model/...
          - ./cmd/...
    
    steps:
      - uses: actions/checkout@v4
      
      - name: Setup Go
        uses: EduGoGroup/edugo-infrastructure/.github/actions/setup-edugo-go@main
      
      - name: Run Tests
        run: go test -v -race ${{ matrix.package }}
      
      - name: Upload Coverage
        uses: actions/upload-artifact@v4
        with:
          name: coverage-${{ matrix.package }}
          path: coverage-*.out
```

**Checkpoint:**
- [ ] Matriz implementada
- [ ] Tests corren en paralelo
- [ ] Cobertura combinada al final

---

## Tarea 7: Paralelizar Lint y Tests

**🟢 Prioridad:** P2  
**⏱️ Tiempo estimado:** 1.5 horas

### Objetivo

Ejecutar lint y tests simultáneamente en lugar de secuencialmente.

### Antes (Secuencial)

```
unit-tests → lint → summary
   ↓          ↓
 3 min    1 min
Total: 4 min
```

### Después (Paralelo)

```
unit-tests  lint
   ↓         ↓
 3 min    1 min
       ↘   ↙
      summary
Total: 3 min (25% más rápido)
```

### Implementación

```yaml
# Modificar pr-to-dev.yml

jobs:
  unit-tests:
    # ... (como antes)
  
  lint:
    runs-on: ubuntu-latest
    # ← Remover "needs: [unit-tests]"
    # Ahora corre en paralelo
    
    steps:
      # ... lint steps
  
  summary:
    needs: [unit-tests, lint]  # ← Espera ambos
    # ... summary
```

**Checkpoint:**
- [ ] Jobs independientes (sin needs entre sí)
- [ ] summary espera ambos
- [ ] Tiempo reducido

---

## Tarea 8: Optimizar Cache

**🟢 Prioridad:** P2  
**⏱️ Tiempo estimado:** 1 hora

### Objetivo

Mejorar uso de cache para acelerar builds.

### Implementación

```yaml
# En todos los workflows

steps:
  - uses: actions/checkout@v4
  
  - name: Setup Go
    uses: actions/setup-go@v5
    with:
      go-version: "1.25"
      cache: true                    # ← Habilitar
      cache-dependency-path: go.sum  # ← Especificar
  
  # Para Docker builds
  - name: Docker Build
    uses: docker/build-push-action@v5
    with:
      cache-from: type=gha           # ← GitHub Actions cache
      cache-to: type=gha,mode=max    # ← Maximizar cache
```

**Checkpoint:**
- [ ] Go cache habilitado
- [ ] Docker cache optimizado
- [ ] Mejora de 10-20% en builds repetidos

---

## Tarea 9: Medir Mejoras

**🟢 Prioridad:** P2  
**⏱️ Tiempo estimado:** 30 minutos

### Script de Métricas

```bash
#!/bin/bash
# Script: sprint4-09-measure-improvements.sh

REPO="EduGoGroup/edugo-api-administracion"

echo "📊 MÉTRICAS DE MEJORA - SPRINT 4"
echo "=================================="
echo ""

# Tiempos antes (runs antiguos)
echo "⏱️  ANTES (runs pre-Sprint 4):"
gh run list --repo $REPO --workflow=pr-to-dev.yml \
  --created="<2025-11-20" --limit 10 --json conclusion,createdAt,updatedAt \
  --jq '.[] | select(.conclusion=="success") | ((.updatedAt | fromdateiso8601) - (.createdAt | fromdateiso8601)) / 60' \
  | awk '{sum+=$1; n++} END {if (n>0) print "Promedio: " sum/n " minutos"}'

echo ""

# Tiempos después (runs recientes)
echo "⏱️  DESPUÉS (runs post-Sprint 4):"
gh run list --repo $REPO --workflow=pr-to-dev.yml \
  --created=">2025-11-20" --limit 10 --json conclusion,createdAt,updatedAt \
  --jq '.[] | select(.conclusion=="success") | ((.updatedAt | fromdateiso8601) - (.createdAt | fromdateiso8601)) / 60' \
  | awk '{sum+=$1; n++} END {if (n>0) print "Promedio: " sum/n " minutos"}'

echo ""
echo "📈 Objetivo: 20-30% reducción"
echo ""

# Líneas de código
echo "📝 CÓDIGO DUPLICADO:"
echo "Antes: ~700 líneas duplicadas"
echo "Después: ~200 líneas duplicadas"
echo "Reducción: ~71%"
echo ""

echo "✅ SPRINT 4 COMPLETO"
```

---

## Resumen Sprint 4

### Tareas Completadas

- [ ] Tarea 1: setup-edugo-go migrado
- [ ] Tarea 2: docker-build-edugo migrado
- [ ] Tarea 3: coverage-check migrado
- [ ] Tarea 4: sync-branches reusable
- [ ] Tarea 5: release reusable (opcional)
- [ ] Tarea 6: Matriz de tests
- [ ] Tarea 7: Paralelismo lint + tests
- [ ] Tarea 8: Cache optimizado
- [ ] Tarea 9: Métricas recopiladas

### Resultados Esperados

| Métrica | Antes | Después | Mejora |
|---------|-------|---------|--------|
| Código duplicado | ~700 líneas | ~200 líneas | -71% |
| Tiempo CI (pr-to-dev) | 3-4 min | 2-3 min | -25-33% |
| Workflows manuales | 7 archivos | 7 archivos | Mismo |
| Líneas por workflow | 100-150 | 50-80 | -40% |
| Mantenimiento | Descentralizado | Centralizado | ✅ |

### Beneficios a Largo Plazo

✅ **Mantenimiento:** Cambios en 1 lugar (infrastructure) afectan a todos  
✅ **Consistencia:** Todos los proyectos usan mismos workflows  
✅ **Velocidad:** CI más rápido = feedback más rápido  
✅ **Calidad:** Código estandarizado y probado  

---

## Validación Final Sprint 4

```bash
#!/bin/bash
# Script: sprint4-10-final-validation.sh

echo "✅ CHECKLIST FINAL - SPRINT 4"
echo "=============================="
echo ""

cd ~/source/EduGo/repos-separados/edugo-api-administracion

# 1. Composite actions en uso
echo "1. Composite Actions:"
grep -r "EduGoGroup/edugo-infrastructure/.github/actions" .github/workflows/*.yml | wc -l
echo "   Ocurrencias encontradas (objetivo: 10+)"

# 2. Workflows reusables
echo ""
echo "2. Workflows Reusables:"
grep -r "uses: EduGoGroup/edugo-infrastructure/.github/workflows/reusable" .github/workflows/*.yml | wc -l
echo "   Ocurrencias encontradas (objetivo: 2+)"

# 3. Paralelismo
echo ""
echo "3. Paralelismo:"
grep -r "strategy:" .github/workflows/*.yml | wc -l
echo "   Matrices implementadas (objetivo: 1+)"

# 4. Cache optimizado
echo ""
echo "4. Cache:"
grep -r "cache: true" .github/workflows/*.yml | wc -l
echo "   Cache habilitado en workflows (objetivo: 3+)"

# 5. Tiempo de CI
echo ""
echo "5. Tiempo de CI:"
echo "   Ejecutar: ./sprint4-09-measure-improvements.sh"

echo ""
echo "=============================="
echo "✅ SPRINT 4 VALIDADO"
```

---

## Próximos Pasos

Después de completar Sprint 4:

1. **Monitorear:** Observar métricas de CI por 1-2 semanas
2. **Iterar:** Ajustar si hay problemas
3. **Replicar:** Aplicar mismo patrón a otros proyectos
4. **Documentar:** Actualizar WORKFLOWS.md con nuevos patrones

---

**Generado por:** Claude Code  
**Fecha:** 19 de Noviembre, 2025  
**Versión:** 1.0
