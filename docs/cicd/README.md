# Plan de Implementación: edugo-api-administracion

**Proyecto:** API Administrativa EduGo  
**Tipo:** A (Aplicación Desplegable)  
**Puerto:** 8081  
**Estado:** ⚠️ CRÍTICO - Success Rate 40%  
**Fecha:** 19 de Noviembre, 2025

---

## 📋 Resumen Ejecutivo

Este plan detalla la implementación de mejoras CI/CD para edugo-api-administracion, proyecto con **tasa de éxito crítica de 40%** y múltiples problemas que requieren atención inmediata.

### Estado Actual vs Objetivo

| Aspecto | Actual | Objetivo | Mejora |
|---------|--------|----------|--------|
| Success Rate | 40% 🔴 | 90%+ ✅ | +125% |
| Workflows Docker | 2 (duplicados) | 1 consolidado | -50% |
| Workflows totales | 7 (1 faltante) | 7 completos | +14% |
| Go Version | 1.24 | 1.25 | Latest |
| Tests Integración | ❌ No | ✅ Sí | Nuevo |
| Código duplicado | ~70% | ~20% | -71% |
| Pre-commit hooks | ❌ No | ✅ Sí | Nuevo |

### Inversión vs Retorno

```
Tiempo Total: 30-37 horas (2 sprints)
├── Sprint 2: 18-22h (P0 + P1)
└── Sprint 4: 12-15h (P2)

ROI Esperado:
✅ +50% tasa de éxito (40% → 90%+)
✅ -3-4h mantenimiento mensual (código reusable)
✅ -1h por release (workflow consolidado)
✅ Prevención bugs en main (pr-to-main.yml)
```

---

## 🎯 Contexto del Proyecto

### ¿Qué es edugo-api-administracion?

API REST administrativa para el ecosistema EduGo. Proporciona endpoints para:
- Gestión de usuarios y permisos
- Administración de escuelas
- Gestión de materiales educativos
- Reportes y estadísticas
- Configuración del sistema

### Stack Tecnológico

```yaml
Lenguaje: Go 1.24 (migrar a 1.25)
Framework: Gin
ORM: GORM
Bases de Datos:
  - PostgreSQL 15 (principal)
  - MongoDB 7.0 (opcional)
Mensajería: RabbitMQ 3.12
Autenticación: JWT via edugo-shared/auth
Logger: edugo-shared/logger
Puerto: 8081
```

### Arquitectura

```
┌─────────────────────────────────────────────┐
│   edugo-api-administracion (Puerto 8081)    │
├─────────────────────────────────────────────┤
│                                             │
│  internal/                                  │
│  ├── handler/        ← HTTP handlers        │
│  ├── service/        ← Business logic       │
│  ├── repository/     ← Data access          │
│  └── model/          ← Domain models        │
│                                             │
│  cmd/                                       │
│  └── server/         ← Entry point          │
│                                             │
└─────────────────────────────────────────────┘
         │              │              │
         ▼              ▼              ▼
    PostgreSQL     RabbitMQ      edugo-shared
      (datos)      (eventos)      (librería)
```

### Dependencias Clave

```go
// go.mod
module github.com/EduGoGroup/edugo-api-administracion

require (
    github.com/EduGoGroup/edugo-shared/logger v0.7.0
    github.com/EduGoGroup/edugo-shared/auth v0.7.0
    github.com/EduGoGroup/edugo-shared/database/postgres v0.7.0
    github.com/EduGoGroup/edugo-shared/messaging/rabbit v0.7.0
    github.com/gin-gonic/gin v1.10.0
    gorm.io/gorm v1.25.12
)
```

---

## 🚨 Problemas Críticos Identificados

### Problema #1: release.yml Fallando 🔴

**Severidad:** CRÍTICA  
**Impacto:** Bloqueando releases de producción  
**Evidencia:**
```
Run ID: 19485500426
Workflow: Release CI/CD (release.yml)
Conclusion: failure
Fecha: 2025-11-19T00:38:48Z
Trigger: Tag v*
```

**Síntomas:**
- Último release exitoso: [fecha desconocida]
- Últimos 3 intentos: TODOS fallidos
- Workflow no previene merge (trigger post-merge)

**Hipótesis de Causas:**
1. **Docker build fallando**
   - Dependencias no resueltas
   - Multi-platform build issue
   - Permisos GHCR

2. **Tests fallando pre-build**
   - Coverage threshold no alcanzado
   - Tests unitarios con errores
   - Lint fallando

3. **Archivos faltantes**
   - `.github/version.txt` no existe
   - `CHANGELOG.md` mal formateado

4. **Permisos**
   - GITHUB_TOKEN sin permisos write:packages
   - Registry GHCR no accesible

**Acción Requerida:** Investigación urgente + fix (Tarea 1.1 y 1.2)

---

### Problema #2: Workflows Docker Duplicados 🔴

**Severidad:** CRÍTICA  
**Impacto:** Confusión, recursos desperdiciados, tags conflictivos

**Situación:**

| Workflow | Trigger | Propósito | Estado |
|----------|---------|-----------|--------|
| `build-and-push.yml` | Manual + opcional push | Build on-demand | ⚠️ Duplicado |
| `release.yml` | Tag push (v*) | Build en releases | ⚠️ Duplicado + Falla |

**Ambos construyen imágenes Docker → DUPLICACIÓN**

**Análisis de Tags:**

**Escenario A: Manual build (staging)**
```yaml
# build-and-push.yml
inputs:
  environment: staging

Genera:
- ghcr.io/edugogroup/edugo-api-administracion:staging
- ghcr.io/edugogroup/edugo-api-administracion:staging-abc1234
```

**Escenario B: Release v1.5.0**
```yaml
# release.yml
tag: v1.5.0

Genera:
- ghcr.io/edugogroup/edugo-api-administracion:1.5.0
- ghcr.io/edugogroup/edugo-api-administracion:1.5
- ghcr.io/edugogroup/edugo-api-administracion:1
- ghcr.io/edugogroup/edugo-api-administracion:latest
- ghcr.io/edugogroup/edugo-api-administracion:production
- ghcr.io/edugogroup/edugo-api-administracion:v1.5.0-abc1234
```

**Problema:**
- Si ambos corren el mismo día: `latest` se sobreescribe
- Múltiples tags SHA duplicados
- Confusión sobre cuál imagen usar
- Desperdicio de espacio en registry

**Solución Propuesta:**
1. **Eliminar** `build-and-push.yml` completamente
2. **Consolidar** todo en `manual-release.yml`
3. **Opcional:** Habilitar `release.yml` solo si se usa

---

### Problema #3: Falta pr-to-main.yml 🔴

**Severidad:** ALTA  
**Impacto:** Código no validado puede llegar a main

**Comparación:**

```
api-mobile (✅ TIENE):
├── pr-to-dev.yml     → Tests unitarios + lint
└── pr-to-main.yml    → Tests unitarios + INTEGRACIÓN + lint + security

api-administracion (❌ FALTA):
├── pr-to-dev.yml     → Tests unitarios + lint
└── pr-to-main.yml    → ❌ NO EXISTE
```

**Consecuencias:**
- PRs a main NO tienen gate de calidad adicional
- Tests de integración NO corren antes de merge
- Errores pueden llegar a main sin detectarse
- No hay validación de security issues

**Solución:** Crear pr-to-main.yml basado en api-mobile (Tarea 3.1)

---

### Problema #4: Go 1.24 (Migrar a 1.25) 🟡

**Severidad:** MEDIA-ALTA  
**Impacto:** Incompatibilidades futuras, sin mejoras de 1.25

**Estado Actual del Ecosistema:**
```
api-mobile:        1.24 → Migrar
api-administracion: 1.24 → Migrar ✅ (este proyecto)
worker:            1.25 ✅ (ya migrado)
shared:            1.25 ✅ (ya migrado)
infrastructure:    1.24 → Migrar
```

**Razón de Migración:**
- Go 1.25 ya validado en api-mobile (tests exitosos)
- Mejoras de performance y seguridad
- Compatibilidad con shared v0.7.0
- Alineación con resto del ecosistema

**Validación Realizada:**
```
✅ Build con golang:1.25-alpine: OK
✅ Tests unitarios: OK
✅ golangci-lint compatible: OK
✅ Dependencias: Todas compatibles
```

**Solución:** Ejecutar script de migración (Tarea 4.1)

---

### Problema #5: Sin Tests de Integración 🟡

**Severidad:** MEDIA  
**Impacto:** Bugs en integración no detectados en CI

**Comparación:**

```
api-mobile:
✅ Tests unitarios: 156 tests
✅ Tests integración: 23 tests (Testcontainers)
✅ Coverage: 39.2%

api-administracion:
✅ Tests unitarios: ~100 tests (estimado)
❌ Tests integración: NO IMPLEMENTADOS
✅ Coverage: 33%+
```

**Por qué es importante:**
- Tests unitarios NO cubren:
  - Interacción con PostgreSQL real
  - Queries GORM complejos
  - Transacciones y rollbacks
  - RabbitMQ messaging

**Solución Gradual:**
1. **Sprint 2:** Agregar placeholder en pr-to-main.yml
2. **Sprint 3:** Implementar tests básicos con Testcontainers
3. **Sprint 4:** Expandir coverage de integración

---

## 📊 Inventario de Workflows

### Workflows Existentes (7 archivos)

#### 1. pr-to-dev.yml ✅

**Propósito:** Validar PRs antes de merge a dev  
**Trigger:** `pull_request` a branch `dev`  
**Estado:** FUNCIONAL

**Jobs:**
```yaml
1. unit-tests:
   - Setup Go 1.24
   - go test ./...
   - Coverage check (33% threshold)
   - Comentar resultado en PR

2. lint:
   - golangci-lint v1.64.7
   - Verificar formato

3. summary:
   - Resumen de checks
   - Comentar en PR
```

**Métricas:**
- Duración promedio: 3-4 minutos
- Tasa de éxito: ~85%

---

#### 2. pr-to-main.yml ❌

**Estado:** **NO EXISTE - FALTANTE**

**Debería tener:**
```yaml
1. unit-tests:
   - Tests unitarios completos
   - Coverage check strict

2. integration-tests:
   - Tests con Testcontainers
   - PostgreSQL + RabbitMQ

3. lint:
   - golangci-lint strict

4. security:
   - gosec scan
   - nancy (dependency check)

5. summary:
   - Resumen completo
```

**Acción:** Crear en Sprint 2 Día 3

---

#### 3. test.yml ✅

**Propósito:** Tests on-demand con coverage detallado  
**Trigger:** `workflow_dispatch` (manual)  
**Estado:** FUNCIONAL

**Features:**
```yaml
- Coverage detallado por paquete
- Upload de artifact con reporte
- Comentario opcional en PR
- Sin threshold enforcement
```

**Uso:**
```bash
gh workflow run test.yml --repo EduGoGroup/edugo-api-administracion
```

---

#### 4. manual-release.yml ✅

**Propósito:** Release manual controlado  
**Trigger:** `workflow_dispatch` con inputs  
**Estado:** FUNCIONAL (pero sin GitHub App token)

**Inputs:**
```yaml
version: (required) - ej: 1.5.0
environment: (required) - development, staging, production
push_latest: (optional) - boolean
```

**Jobs:**
1. Validar version format
2. Actualizar version.txt
3. Build Docker multi-platform
4. Push a GHCR
5. Create GitHub release
6. Commit version.txt

**Problema Detectado:**
- Usa `GITHUB_TOKEN` en lugar de GitHub App token
- Consecuencia: No dispara `sync-main-to-dev.yml` automáticamente

**Solución:** Agregar GitHub App token (Sprint 2)

---

#### 5. build-and-push.yml ⚠️

**Estado:** **DUPLICADO - ELIMINAR**

**Propósito:** Build Docker on-demand  
**Trigger:** `workflow_dispatch` + opcional `push`

**Por qué eliminar:**
- Funcionalidad duplicada con manual-release.yml
- Genera tags conflictivos
- Confusión sobre cuál usar
- Mantenimiento duplicado

**Acción:** Eliminar en Sprint 2 Día 2

---

#### 6. release.yml ❌

**Estado:** **FALLANDO - RESOLVER**

**Propósito:** Release automático al crear tag  
**Trigger:** `push` de tag `v*`

**Problema:** Últimos runs TODOS fallidos

**Debería hacer:**
1. Extraer versión del tag
2. Run tests completos
3. Build Docker multi-platform
4. Push con múltiples tags (semver)
5. Create GitHub release con changelog

**Decisión Pendiente:**
- ¿Reparar y mantener?
- ¿O eliminar y usar solo manual-release.yml?

**Recomendación:** Resolver en Sprint 2, decidir si mantener o eliminar

---

#### 7. sync-main-to-dev.yml ✅

**Propósito:** Sincronizar main → dev automáticamente  
**Trigger:** `push` a `main` o tag `v*`  
**Estado:** FUNCIONAL (pero código duplicado)

**Lógica:**
```yaml
1. Check si dev existe (crear si no)
2. Verificar diferencias main vs dev
3. Merge main → dev (auto)
4. Si hay conflictos → fallar y notificar
```

**Problema:**
- Código 96% idéntico en 6 repositorios
- 100 líneas duplicadas

**Solución:** Migrar a workflow reusable (Sprint 4)

---

## 🎯 Plan de Sprints

### Sprint 2: Resolver Críticos + Alta Prioridad

**Objetivo:** Estabilizar CI/CD y resolver problemas críticos  
**Duración:** 5 días / 18-22 horas  
**Prioridad:** 🔴 P0 + 🟡 P1

#### Día 1: Investigación (4-5h)

**Tareas:**
- [Tarea 1.1] Investigar fallos en release.yml (2-4h)
- [Tarea 1.2] Analizar logs y reproducir localmente (1-2h)

**Entregables:**
- Documento de análisis de fallo
- Reproducción local del error
- Plan de corrección

---

#### Día 2: Resolución de Fallos (4-5h)

**Tareas:**
- [Tarea 2.1] Aplicar fix a release.yml (2-3h)
- [Tarea 2.2] Eliminar build-and-push.yml (1h)
- [Tarea 2.3] Testing y validación (1h)

**Entregables:**
- PR con fix de release.yml
- build-and-push.yml eliminado
- Tests CI/CD pasando

---

#### Día 3: Agregar pr-to-main.yml (4-5h)

**Tareas:**
- [Tarea 3.1] Crear pr-to-main.yml (1.5h)
- [Tarea 3.2] Configurar tests integración placeholder (1h)
- [Tarea 3.3] Testing workflow completo (1h)
- [Tarea 3.4] Documentar workflow (30min)

**Entregables:**
- pr-to-main.yml funcional
- Placeholder integración tests
- Documentación actualizada

---

#### Día 4: Migrar a Go 1.25 (3-4h)

**Tareas:**
- [Tarea 4.1] Ejecutar script de migración (45min)
- [Tarea 4.2] Tests completos (build + unit + lint) (1h)
- [Tarea 4.3] Actualizar documentación (30min)
- [Tarea 4.4] Crear PR y merge (1h)

**Entregables:**
- Go 1.25 en todos los archivos
- Tests pasando en Go 1.25
- PR merged a dev

---

#### Día 5: Mejoras Adicionales (3-4h)

**Tareas:**
- [Tarea 5.1] Configurar pre-commit hooks (1h)
- [Tarea 5.2] Agregar label skip-coverage (30min)
- [Tarea 5.3] Agregar GitHub App token (30min)
- [Tarea 5.4] Documentación final (1h)

**Entregables:**
- Pre-commit hooks activos
- Label skip-coverage disponible
- GitHub App token configurado
- README.md actualizado

---

### Sprint 4: Workflows Reusables

**Objetivo:** Eliminar duplicación, optimizar tiempos  
**Duración:** 3 días / 12-15 horas  
**Prioridad:** 🟢 P2

#### Día 1: Migrar a Composite Actions (4-5h)

**Tareas:**
- Usar setup-edugo-go
- Usar docker-build-edugo
- Usar coverage-check
- Testing

---

#### Día 2: Workflows Reusables (4-5h)

**Tareas:**
- Migrar sync-main-to-dev.yml
- Migrar release logic (si aplica)
- Testing

---

#### Día 3: Paralelismo (4-5h)

**Tareas:**
- Implementar matriz de tests
- Paralelizar lint + tests + build
- Medir mejoras de tiempo

**Objetivo de Performance:**
- Actual: ~3-4 minutos
- Objetivo: ~2-3 minutos (20-30% mejora)

---

## 🛠️ Herramientas y Comandos

### GitHub CLI (gh)

```bash
# Ver workflows
gh workflow list --repo EduGoGroup/edugo-api-administracion

# Ver runs recientes
gh run list --repo EduGoGroup/edugo-api-administracion --limit 10

# Ver logs de run fallido
gh run view 19485500426 --repo EduGoGroup/edugo-api-administracion --log-failed

# Ejecutar workflow manual
gh workflow run manual-release.yml \
  --repo EduGoGroup/edugo-api-administracion \
  --field version=1.5.0 \
  --field environment=staging

# Crear PR
gh pr create --base dev --title "fix: resolver fallo en release.yml" --body "..."
```

### Testing Local

```bash
cd ~/source/EduGo/repos-separados/edugo-api-administracion

# Tests unitarios
go test ./... -v

# Tests con coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out

# Lint
golangci-lint run

# Build
go build ./cmd/server

# Docker build (local)
docker build -t edugo-api-admin:local .
docker run -p 8081:8081 edugo-api-admin:local
```

### act (GitHub Actions Localmente)

```bash
# Instalar act
brew install act

# Listar workflows
act -l

# Correr pr-to-dev.yml localmente
act pull_request -W .github/workflows/pr-to-dev.yml

# Correr con secrets
act -s GITHUB_TOKEN=ghp_xxx
```

---

## 📚 Documentación de Referencia

### Workflows de api-mobile (Ejemplo a Seguir)

```
../03-api-mobile/
├── SPRINT-2-TASKS.md    ← Cómo implementar pr-to-main.yml
├── README.md            ← Arquitectura similar
└── workflows/
    └── pr-to-main.yml   ← Copiar y adaptar
```

### Análisis Previo

```
../../
├── 01-ANALISIS-ESTADO-ACTUAL.md        ← Estado inicial
├── 03-DUPLICIDADES-DETALLADAS.md       ← Código duplicado
├── 05-QUICK-WINS.md                    ← Mejoras rápidas
└── 08-RESULTADO-PRUEBAS-GO-1.25.md     ← Validación Go 1.25
```

---

## ✅ Criterios de Éxito

### Sprint 2 (Completado cuando...)

- [ ] release.yml pasa o está deshabilitado con justificación
- [ ] build-and-push.yml eliminado
- [ ] pr-to-main.yml existe y funciona
- [ ] Go 1.25 en todos los archivos (go.mod, workflows, Dockerfile)
- [ ] Tests pasan con Go 1.25
- [ ] Pre-commit hooks activos
- [ ] Success rate > 80%
- [ ] Documentación actualizada

### Sprint 4 (Completado cuando...)

- [ ] 3+ composite actions en uso
- [ ] sync-main-to-dev.yml usa workflow reusable
- [ ] Tests corren en paralelo
- [ ] Tiempo CI reducido 20%+
- [ ] Código duplicado < 30%

---

## 🚀 Quick Start

```bash
# 1. Clone o actualiza repo
cd ~/source/EduGo/repos-separados/edugo-api-administracion
git checkout dev
git pull origin dev

# 2. Crear backup
git checkout -b backup/pre-sprint2-$(date +%Y%m%d)
git push origin backup/pre-sprint2-$(date +%Y%m%d)
git checkout dev

# 3. Revisar estado actual
gh run list --limit 10
git log --oneline -5

# 4. Abrir plan de tareas
open ../Analisys/00-Projects-Isolated/cicd-analysis/implementation-plans/04-api-administracion/SPRINT-2-TASKS.md

# 5. Comenzar con Tarea 1.1
```

---

**Generado por:** Claude Code  
**Fecha:** 19 de Noviembre, 2025  
**Versión:** 1.0  
**Basado en:** Análisis CI/CD completo + plan de api-mobile
