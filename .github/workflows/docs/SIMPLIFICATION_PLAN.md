# Plan de Simplificación de CI/CD

**Fecha**: 9 de noviembre de 2025  
**Objetivo**: Eliminar duplicación, simplificar workflows, mantener solo lo esencial

---

## 🔍 Análisis de Workflows Actuales

### Workflows Existentes (16 archivos)

| Archivo | Propósito | Estado | Acción Recomendada |
|---------|-----------|--------|-------------------|
| `pr-to-dev.yml` | Tests en PR a dev | ✅ Activo | **MANTENER** |
| `pr-to-main.yml` | Tests completos en PR a main | ✅ Activo | **MANTENER** |
| `manual-release.yml` | Release manual completo | ✅ Activo | **MANTENER** |
| `test.yml` | Tests manuales | ✅ Activo | **MANTENER** (simplificado) |
| `build-and-push.yml` | Build Docker automático | ⚠️ Duplicado | **ELIMINAR** |
| `ci.yml` | CI Pipeline genérico | ⚠️ Duplicado | **ELIMINAR** |
| `docker-only.yml` | Build Docker simple | ⚠️ Duplicado | **ELIMINAR** |
| `release.yml` | Release automático por tag | ⚠️ Duplicado | **ELIMINAR** |
| `sync-main-to-dev.yml` | Sync branches | ✅ Útil | **MANTENER** |
| `sync-main-to-dev-ff.yml` | Sync fast-forward | ⚠️ Duplicado | **ELIMINAR** |
| `test.yml.bak` | Backup | ❌ Basura | **ELIMINAR** |
| `integration-tests.yml.example` | Ejemplo | ❌ Basura | **ELIMINAR** |
| `README.md` | Doc vieja | ⚠️ Desactualizada | **MOVER a docs/** |
| `TESTING_STRATEGY.md` | Doc vieja | ⚠️ Desactualizada | **MOVER a docs/** |
| `CI_CD_STRATEGY.md` | Doc nueva | ✅ Activa | **MOVER a docs/** |
| `WORKFLOW_DIAGRAM.md` | Doc nueva | ✅ Activa | **MOVER a docs/** |

---

## ❌ Workflows a ELIMINAR (Duplicados)

### 1. `build-and-push.yml` - ELIMINAR

**Razón**: Duplica funcionalidad de `manual-release.yml`

**Problemas**:
- Se ejecuta automáticamente en push a main (no queremos esto)
- Ejecuta tests antes de build (ya se ejecutan en PR)
- Duplica lógica de Docker que ya está en `manual-release.yml`

**Qué hace**:
```yaml
on:
  push:
    branches: [main]  # ← PROBLEMA: Automático
  workflow_dispatch:
```

**Reemplazo**: Usar `manual-release.yml` que es más completo y controlado

---

### 2. `ci.yml` - ELIMINAR

**Razón**: Duplica funcionalidad de `pr-to-dev.yml` y `pr-to-main.yml`

**Problemas**:
- Se ejecuta en PRs a main y dev (ya cubierto por workflows específicos)
- Ejecuta tests que ya se ejecutan en otros workflows
- Lógica de tests menos optimizada que los nuevos workflows

**Qué hace**:
```yaml
on:
  pull_request:
    branches: [main, dev]  # ← Ya cubierto por pr-to-dev.yml y pr-to-main.yml
```

**Reemplazo**: Los workflows `pr-to-dev.yml` y `pr-to-main.yml` son más específicos y optimizados

---

### 3. `docker-only.yml` - ELIMINAR

**Razón**: Duplica funcionalidad de `manual-release.yml`

**Problemas**:
- Solo build Docker sin tests ni validación
- Funcionalidad ya incluida en `manual-release.yml`
- No agrega valor adicional

**Qué hace**:
```yaml
on:
  workflow_dispatch:  # Solo manual
```

**Reemplazo**: Usar `manual-release.yml` que incluye build Docker + tests + release

---

### 4. `release.yml` - ELIMINAR

**Razón**: Duplica funcionalidad de `manual-release.yml`

**Problemas**:
- Se ejecuta automáticamente en tags (no queremos esto ahora)
- Funcionalidad idéntica a `manual-release.yml`
- Menos control que el workflow manual

**Qué hace**:
```yaml
on:
  push:
    tags: ['v*']  # ← PROBLEMA: Automático
```

**Reemplazo**: Usar `manual-release.yml` que es más controlado

**Nota**: En el futuro, cuando vayamos a producción, podemos invocar `manual-release.yml` desde un workflow de merge a main

---

### 5. `sync-main-to-dev-ff.yml` - ELIMINAR

**Razón**: Duplica funcionalidad de `sync-main-to-dev.yml`

**Problemas**:
- Hace lo mismo que `sync-main-to-dev.yml` pero con fast-forward
- Tener dos workflows de sync confunde

**Reemplazo**: Mantener solo `sync-main-to-dev.yml`

---

### 6. `test.yml.bak` - ELIMINAR

**Razón**: Archivo de backup innecesario

**Acción**: Eliminar directamente

---

### 7. `integration-tests.yml.example` - ELIMINAR

**Razón**: Archivo de ejemplo que ya no se usa

**Acción**: Eliminar directamente

---

## ✅ Workflows a MANTENER

### 1. `pr-to-dev.yml` ✅

**Propósito**: Tests rápidos en PRs a dev

**Mantener porque**:
- Optimizado para velocidad (~2-3 min)
- Solo tests unitarios (suficiente para dev)
- Comentarios automáticos útiles

**Triggers**:
```yaml
on:
  pull_request:
    branches: [dev]
```

---

### 2. `pr-to-main.yml` ✅

**Propósito**: Tests completos en PRs a main

**Mantener porque**:
- Validación completa antes de producción
- Tests unitarios + integración + security
- Comentarios detallados

**Triggers**:
```yaml
on:
  pull_request:
    branches: [main]
```

---

### 3. `manual-release.yml` ✅

**Propósito**: Release completo manual

**Mantener porque**:
- Control total del proceso de release
- Incluye: tests + build Docker + GitHub release
- On-demand (no automático)
- Más completo que otros workflows

**Triggers**:
```yaml
on:
  workflow_dispatch:  # Solo manual
```

**Nota**: Este es el workflow maestro para releases

---

### 4. `test.yml` ✅

**Propósito**: Tests manuales

**Mantener porque**:
- Útil para ejecutar tests on-demand
- Permite elegir tipo de tests (unit/integration/all)
- No interfiere con otros workflows

**Triggers**:
```yaml
on:
  workflow_dispatch:  # Solo manual
```

---

### 5. `sync-main-to-dev.yml` ✅

**Propósito**: Sincronizar main → dev después de release

**Mantener porque**:
- Mantiene dev actualizado con main
- Se ejecuta automáticamente después de release
- Evita divergencia de branches

**Triggers**:
```yaml
on:
  push:
    branches: [main]
```

---

## 📊 Comparación: Antes vs Después

### ANTES (Sobrecargado)

```
Workflows Activos: 11
├─ PR a dev: ci.yml (duplicado)
├─ PR a main: ci.yml (duplicado)
├─ Build Docker: build-and-push.yml (automático)
├─ Build Docker: docker-only.yml (manual)
├─ Release: release.yml (automático)
├─ Release: manual-release.yml (manual)
├─ Sync: sync-main-to-dev.yml
├─ Sync: sync-main-to-dev-ff.yml (duplicado)
└─ Tests: test.yml

Problemas:
- 4 workflows duplicados
- 3 workflows automáticos no deseados
- Confusión sobre cuál usar
```

### DESPUÉS (Simplificado)

```
Workflows Activos: 5
├─ PR a dev: pr-to-dev.yml (optimizado)
├─ PR a main: pr-to-main.yml (completo)
├─ Release: manual-release.yml (on-demand)
├─ Sync: sync-main-to-dev.yml
└─ Tests: test.yml (manual)

Beneficios:
- Sin duplicación
- Todo on-demand excepto sync
- Claro qué usar en cada caso
- 54% menos workflows
```

---

## 🎯 Estrategia Simplificada

### Flujo de Desarrollo

```
feature/nueva-funcionalidad
  ↓ PR
dev ← pr-to-dev.yml (tests unitarios, ~2-3 min)
  ↓ PR
main ← pr-to-main.yml (tests completos, ~3-4 min)
  ↓ Merge
main → sync-main-to-dev.yml (automático)
```

### Flujo de Release (On-Demand)

```
main (listo para release)
  ↓ Manual
manual-release.yml
  ├─ Crear tag
  ├─ Build Docker
  ├─ Publicar imagen
  ├─ Crear GitHub release
  └─ Trigger sync-main-to-dev.yml
```

### Futuro: Release Automático en Merge

Cuando estemos listos para producción:

```yaml
# Nuevo workflow: auto-release-on-merge.yml
on:
  push:
    branches: [main]

jobs:
  trigger-release:
    runs-on: ubuntu-latest
    steps:
      - name: Trigger manual-release.yml
        uses: actions/github-script@v7
        with:
          script: |
            await github.rest.actions.createWorkflowDispatch({
              owner: context.repo.owner,
              repo: context.repo.repo,
              workflow_id: 'manual-release.yml',
              ref: 'main',
              inputs: {
                version: '...',  # Calcular automáticamente
                bump_type: 'patch'
              }
            });
```

---

## 📝 Plan de Acción

### Fase 1: Eliminar Duplicados (Ahora)

```bash
# Eliminar workflows duplicados
rm .github/workflows/build-and-push.yml
rm .github/workflows/ci.yml
rm .github/workflows/docker-only.yml
rm .github/workflows/release.yml
rm .github/workflows/sync-main-to-dev-ff.yml

# Eliminar basura
rm .github/workflows/test.yml.bak
rm .github/workflows/integration-tests.yml.example
```

### Fase 2: Organizar Documentación (Ahora)

```bash
# Mover docs a carpeta docs/
mv .github/workflows/README.md .github/workflows/docs/
mv .github/workflows/TESTING_STRATEGY.md .github/workflows/docs/
mv .github/workflows/CI_CD_STRATEGY.md .github/workflows/docs/
mv .github/workflows/WORKFLOW_DIAGRAM.md .github/workflows/docs/

# Crear índice unificado
# (ver WORKFLOWS_INDEX.md)
```

### Fase 3: Futuro (Cuando vayamos a producción)

- Crear `auto-release-on-merge.yml` que invoque `manual-release.yml`
- Configurar protección de branches
- Configurar required checks

---

## 🔒 Protección de Branches (Recomendado)

### Branch `main`

```yaml
Required checks:
  - Unit Tests (pr-to-main.yml)
  - Integration Tests (pr-to-main.yml)
  - Lint (pr-to-main.yml)
  - Security Scan (pr-to-main.yml)

Settings:
  - Require PR before merge: ✅
  - Require approvals: 1
  - Dismiss stale reviews: ✅
  - Require status checks: ✅
  - Require branches up to date: ✅
  - No direct pushes: ✅
```

### Branch `dev`

```yaml
Required checks:
  - Unit Tests (pr-to-dev.yml)
  - Lint (pr-to-dev.yml)

Settings:
  - Require PR before merge: ✅
  - Require approvals: 0 (opcional)
  - Require status checks: ✅
  - Allow direct pushes: ❌
```

---

## 📚 Documentación Unificada

Toda la documentación estará en `.github/workflows/docs/`:

1. **WORKFLOWS_INDEX.md** - Índice maestro de todos los workflows
2. **CI_CD_STRATEGY.md** - Estrategia general de CI/CD
3. **WORKFLOW_DIAGRAM.md** - Diagramas visuales
4. **SIMPLIFICATION_PLAN.md** - Este documento
5. **TESTING_STRATEGY.md** - Estrategia de testing (actualizada)
6. **TROUBLESHOOTING.md** - Guía de resolución de problemas

---

## ✅ Checklist de Implementación

- [ ] Eliminar 7 workflows duplicados/innecesarios
- [ ] Mover 4 documentos a `docs/`
- [ ] Crear `WORKFLOWS_INDEX.md`
- [ ] Actualizar `TESTING_STRATEGY.md`
- [ ] Crear `TROUBLESHOOTING.md`
- [ ] Verificar que workflows activos funcionan
- [ ] Actualizar README principal del proyecto
- [ ] Comunicar cambios al equipo

---

## 💡 Resumen

**Antes**: 11 workflows activos, 4 duplicados, confusión

**Después**: 5 workflows activos, 0 duplicados, claridad

**Ahorro**: 54% menos workflows, 100% menos confusión

**Beneficios**:
- ✅ Sin duplicación de código
- ✅ Claro qué workflow usar
- ✅ Todo on-demand excepto sync
- ✅ Preparado para futuro automático
- ✅ Documentación organizada

---

**Última actualización**: 9 de noviembre de 2025

