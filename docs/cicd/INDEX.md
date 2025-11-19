# Índice - Plan de Implementación edugo-api-administracion

**🎯 Punto de Entrada Principal**

---

## 🗺️ Navegación Rápida

### Para Empezar
1. **[INDEX.md](./INDEX.md)** ⭐ - Estás aquí (5 min)
2. **[README.md](./README.md)** ⭐⭐ - Contexto completo del proyecto (15-20 min)

### Para Implementar
3. **[SPRINT-2-TASKS.md](./SPRINT-2-TASKS.md)** ⭐⭐⭐ - Plan detallado Sprint 2 (P0 + P1)
4. **[SPRINT-4-TASKS.md](./SPRINT-4-TASKS.md)** - Plan detallado Sprint 4 (P2)

---

## 📊 Resumen Ultra-Rápido

```
Proyecto: edugo-api-administracion
Tipo: A (Aplicación Desplegable - API Administrativa)
Puerto: 8081
Estado Actual: ⚠️ CRÍTICO - Success Rate 40%

Plan Completo: 4 archivos principales
├── Sprint 2: CRÍTICO Y ALTA PRIORIDAD 🔴🟡
│   ├── Duración: 5 días / 18-22 horas
│   ├── P0: Resolver fallos (CRÍTICO)
│   ├── P0: Eliminar Docker duplicado (CRÍTICO)
│   ├── P1: Agregar pr-to-main.yml
│   ├── P1: Migrar a Go 1.25
│   └── ~35 scripts bash ready-to-use
│
└── Sprint 4: WORKFLOWS REUSABLES 🟢
    ├── Duración: 3 días / 12-15 horas
    ├── P2: Paralelismo y optimización
    └── P2: Migrar a workflows reusables

Total Estimado: 30-37 horas de implementación
```

---

## 🚨 CONTEXTO CRÍTICO

### Problema Principal
```
Success Rate: 40% (4/10 últimos runs)
Último fallo: Run 19485500426 (release.yml)
Fecha: 2025-11-19T00:38:48Z
```

### Problemas Identificados

🔴 **P0 - CRÍTICO (Resolver primero):**
1. **Workflow release.yml fallando** - Bloqueando releases
2. **Duplicación workflow Docker** - build-and-push.yml Y release.yml
3. **Falta pr-to-main.yml** - No hay validación pre-merge a main

🟡 **P1 - ALTA PRIORIDAD:**
4. **Go 1.24** - Necesita migrar a 1.25 (ya validado en api-mobile)
5. **No tiene tests de integración** - Solo unitarios
6. **Coverage threshold sin bypass** - No tiene label skip-coverage

🟢 **P2 - MEDIA PRIORIDAD:**
7. **Sin paralelismo** - Tests corren secuencialmente
8. **Código duplicado** - 70% código repetido con otros proyectos
9. **Sin workflows reusables** - Mantenimiento difícil

---

## 🎯 Quick Actions

### Acción 1: Ver Estado Actual
```bash
cd ~/source/EduGo/repos-separados/edugo-api-administracion
git status
git log --oneline -5
gh run list --limit 10
```

### Acción 2: Comenzar Sprint 2 AHORA
```bash
open SPRINT-2-TASKS.md
# Ir a Tarea 1.1: Investigar fallos en release.yml
# Seguir paso a paso
```

### Acción 3: Modo Lectura (Entender sin Ejecutar)
```bash
open README.md
# Leer contexto completo
# Entender arquitectura
# Revisar roadmap
```

---

## 📁 Estructura de Archivos

```
04-api-administracion/
├── INDEX.md                    ← Estás aquí
├── README.md                   ← Contexto del proyecto (~400 líneas)
├── SPRINT-2-TASKS.md          ← ⭐ Sprint 2 completo (~2,500 líneas)
└── SPRINT-4-TASKS.md          ← Sprint 4 parcial (~800 líneas)

Total: ~3,700+ líneas de documentación
```

---

## 🎯 Por Rol

### Soy el Implementador
→ Lee: **README.md** → **SPRINT-2-TASKS.md**  
→ Ejecuta: Tareas P0 primero, luego P1  
→ Tiempo: 18-22 horas Sprint 2

### Soy el DevOps/SRE
→ Lee: **README.md** (sección Workflows Actuales)  
→ Foco: Resolver fallos en release.yml  
→ Tiempo: 2-4 horas investigación + fix

### Soy el Tech Lead
→ Lee: **README.md** + **INDEX.md**  
→ Revisa: Priorización y estimaciones  
→ Tiempo: 30-45 minutos

### Quiero Ver Solo los Problemas
→ Sección: **Problemas Críticos Detallados** (abajo)  
→ Tiempo: 10 minutos

---

## 🔥 Top 5 Tareas Críticas (Sprint 2)

Si solo tienes tiempo limitado, ejecuta estas:

### 1. 🔴 Investigar y Resolver Fallo en release.yml (2-4h)
```bash
# Ver logs del fallo
gh run view 19485500426 --repo EduGoGroup/edugo-api-administracion --log-failed

# Identificar step exacto que falla
# Reproducir localmente
# Aplicar fix
```

### 2. 🔴 Eliminar Workflow Docker Duplicado (1h)
```bash
# Eliminar build-and-push.yml
# Consolidar en manual-release.yml
# Agregar control por variable ENABLE_AUTO_RELEASE
```

### 3. 🟡 Crear pr-to-main.yml (1.5h)
```bash
# Copiar de api-mobile
# Adaptar para api-administracion
# Agregar tests de integración (placeholder)
```

### 4. 🟡 Migrar a Go 1.25 (45 min)
```bash
# Script ya validado en api-mobile
# Actualizar go.mod, workflows, Dockerfile
# Ejecutar tests
```

### 5. 🟡 Configurar Pre-commit Hooks (1h)
```bash
# Agregar .githooks/pre-commit
# Configurar formato + lint + tests
# Actualizar Makefile
```

**Total Quick Wins:** ~6-7 horas (en lugar de 18-22h completas)

---

## 📊 Workflows Actuales

### Lista de Workflows (7 archivos)

| Workflow | Trigger | Estado | Problema |
|----------|---------|--------|----------|
| `pr-to-dev.yml` | PR → dev | ✅ OK | Ninguno |
| `pr-to-main.yml` | PR → main | ❌ NO EXISTE | **FALTANTE** |
| `test.yml` | Manual | ✅ OK | Ninguno |
| `manual-release.yml` | Manual | ✅ OK | Sin GitHub App token |
| `build-and-push.yml` | Manual/Push | ⚠️ Duplicado | **ELIMINAR** |
| `release.yml` | Tag v* | ❌ FALLA | **RESOLVER** |
| `sync-main-to-dev.yml` | Push main | ✅ OK | Duplicado (código) |

---

## 🔍 Problemas Críticos Detallados

### Problema 1: release.yml Fallando

**Evidencia:**
```
Run ID: 19485500426
Workflow: Release CI/CD (release.yml)
Conclusion: failure
Fecha: 2025-11-19T00:38:48Z
Trigger: Tag push (v*)
```

**Hipótesis de Causa:**
1. Docker build fallando
2. Tests fallando antes de build
3. Problema con permisos de GHCR
4. Archivo version.txt o CHANGELOG faltante
5. Dependencias no resueltas

**Acción Requerida:**
```bash
# 1. Ver logs completos
gh run view 19485500426 --repo EduGoGroup/edugo-api-administracion --log-failed

# 2. Buscar línea exacta de fallo
# 3. Reproducir localmente
cd ~/source/EduGo/repos-separados/edugo-api-administracion
git checkout <commit-del-fallo>
# Ejecutar step que falla

# 4. Aplicar fix y crear PR
```

---

### Problema 2: Duplicación de Workflows Docker

**Situación:**
- `build-and-push.yml` - Trigger: Manual + opcional push
- `release.yml` - Trigger: Tag push (v*)

**Ambos construyen imágenes Docker** → Desperdicio de recursos + confusión

**Estrategia de Tags:**

`build-and-push.yml`:
```yaml
tags: |
  type=raw,value=${{ inputs.environment }}     # development, staging, production
  type=raw,value=latest,enable=${{ inputs.push_latest }}
  type=sha,prefix=${{ inputs.environment }}-
```

`release.yml`:
```yaml
tags: |
  type=semver,pattern={{version}}              # 1.0.0
  type=semver,pattern={{major}}.{{minor}}      # 1.0
  type=semver,pattern={{major}}                # 1
  type=raw,value=latest
  type=raw,value=production
  type=sha,prefix=${{ tag }}-
```

**Problema:** Si se hace manual build Y tag el mismo día:
- `latest` se sobreescribe entre workflows
- Múltiples SHA tags: `staging-abc123`, `production-abc123`, `v1.0.0-abc123`

**Solución Propuesta:**
1. **Eliminar** `build-and-push.yml`
2. **Mantener** `manual-release.yml` (consolidado)
3. **Opcional:** Habilitar `release.yml` solo cuando se use auto-release

---

### Problema 3: Falta pr-to-main.yml

**Consecuencia:**
- No hay validación de tests antes de merge a main
- Errores pueden llegar a main sin detectarse
- No hay tests de integración en gate de calidad

**Comparación con api-mobile:**
```
api-mobile tiene:
✅ pr-to-dev.yml  - Tests unitarios + lint
✅ pr-to-main.yml - Tests unitarios + INTEGRACIÓN + lint + security

api-administracion tiene:
✅ pr-to-dev.yml  - Tests unitarios + lint
❌ pr-to-main.yml - NO EXISTE
```

**Solución:**
1. Copiar `pr-to-main.yml` de api-mobile
2. Adaptar para api-administracion
3. Agregar placeholder para tests de integración (implementar después)

---

## 📈 Estadísticas del Proyecto

### Estado de Salud CI/CD

```
Success Rate: 40% (4 success / 10 runs)
Failure Rate: 60% (6 failures / 10 runs)
```

**Comparación con otros proyectos:**
```
api-mobile:        90% ✅ (excelente)
api-administracion: 40% ⚠️ (crítico)
worker:            70% ⚠️ (necesita mejora)
shared:           100% ✅ (perfecto)
infrastructure:    20% 🔴 (emergencia)
```

### Workflows por Categoría

```
✅ Funcionales:        4/7 (57%)
⚠️ Con problemas:     2/7 (29%)
❌ Faltantes:          1/7 (14%)
```

### Código Duplicado

```
Duplicación estimada: ~70%
Líneas duplicadas:    ~800 líneas
Oportunidad ahorro:   ~600 líneas (con reusables)
```

---

## 🚀 Roadmap de Implementación

### Sprint 2: Resolver Críticos + Alta Prioridad (5 días)

**Día 1: Investigación y Análisis** (4-5h)
- [ ] Tarea 1.1: Investigar fallos release.yml (2-4h)
- [ ] Tarea 1.2: Analizar logs y reproducir (1-2h)

**Día 2: Resolución de Fallos** (4-5h)
- [ ] Tarea 2.1: Aplicar fix a release.yml (2-3h)
- [ ] Tarea 2.2: Eliminar workflow duplicado (1h)
- [ ] Tarea 2.3: Testing y validación (1h)

**Día 3: Agregar pr-to-main.yml** (4-5h)
- [ ] Tarea 3.1: Crear pr-to-main.yml (1.5h)
- [ ] Tarea 3.2: Configurar tests integración placeholder (1h)
- [ ] Tarea 3.3: Testing workflow (1h)
- [ ] Tarea 3.4: Documentar (30min)

**Día 4: Migrar a Go 1.25** (3-4h)
- [ ] Tarea 4.1: Ejecutar script migración (45min)
- [ ] Tarea 4.2: Tests completos (1h)
- [ ] Tarea 4.3: Actualizar docs (30min)
- [ ] Tarea 4.4: Crear PR y merge (1h)

**Día 5: Mejoras Adicionales** (3-4h)
- [ ] Tarea 5.1: Configurar pre-commit hooks (1h)
- [ ] Tarea 5.2: Agregar label skip-coverage (30min)
- [ ] Tarea 5.3: Documentación final (1h)
- [ ] Tarea 5.4: Revisión y cierre sprint (30min)

**Total Sprint 2:** 18-22 horas

---

### Sprint 4: Workflows Reusables (3 días)

**Día 1: Migrar a Composite Actions** (4-5h)
- [ ] Usar setup-edugo-go
- [ ] Usar docker-build-edugo
- [ ] Usar coverage-check

**Día 2: Migrar a Workflows Reusables** (4-5h)
- [ ] Migrar sync-main-to-dev.yml
- [ ] Migrar release logic

**Día 3: Paralelismo y Optimización** (4-5h)
- [ ] Implementar matriz de tests
- [ ] Paralelizar lint + tests
- [ ] Optimizar tiempos de CI

**Total Sprint 4:** 12-15 horas

---

## 💾 Backup y Seguridad

### Antes de Comenzar

```bash
# Crear backup del estado actual
cd ~/source/EduGo/repos-separados/edugo-api-administracion
git checkout dev
git pull origin dev
git checkout -b backup/pre-sprint2-$(date +%Y%m%d)
git push origin backup/pre-sprint2-$(date +%Y%m%d)
```

### Puntos de Restore

```bash
# Si algo sale mal, restaurar desde:
git checkout backup/pre-sprint2-YYYYMMDD
```

---

## 🆘 Ayuda Rápida

### ¿Por dónde empiezo?
**Respuesta:** README.md (contexto) → SPRINT-2-TASKS.md Tarea 1.1

### ¿Cuánto tiempo necesito?
**Respuesta:** Sprint 2 completo = 18-22h. Modo rápido (solo P0) = 4-6h.

### ¿Puedo saltar tareas?
**Respuesta:** NO saltes tareas P0 (críticas). P1 y P2 son opcionales.

### ¿Los scripts funcionan?
**Respuesta:** Sí, diseñados para copiar/pegar y ejecutar. Testeados en api-mobile.

### ¿Qué hago si release.yml sigue fallando?
**Respuesta:** 
1. Verificar logs completos
2. Consultar sección Troubleshooting en SPRINT-2-TASKS.md
3. Considerar deshabilitar release.yml temporalmente

### ¿Debo seguir el orden exacto?
**Respuesta:** Sí para P0. P1 y P2 pueden reordenarse.

---

## 📞 Referencias

### Documentación Base
- [Análisis Estado Actual](../../01-ANALISIS-ESTADO-ACTUAL.md)
- [Propuestas de Mejora](../../02-PROPUESTAS-MEJORA.md)
- [Duplicidades Detalladas](../../03-DUPLICIDADES-DETALLADAS.md)
- [Quick Wins](../../05-QUICK-WINS.md)
- [Resultado Pruebas Go 1.25](../../08-RESULTADO-PRUEBAS-GO-1.25.md)

### Repositorio
- **URL:** https://github.com/EduGoGroup/edugo-api-administracion
- **Ruta Local:** `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-administracion`
- **Puerto:** 8081
- **Tipo:** A (Aplicación Desplegable)

### Plan de api-mobile (Referencia)
- **Ruta:** `../03-api-mobile/`
- **Uso:** Como ejemplo de implementación exitosa

---

## ✅ Checklist Pre-Lectura

Antes de comenzar:
- [x] Directorio correcto
- [x] Tienes tiempo para leer (mínimo 30 min)
- [ ] Editor markdown disponible
- [ ] Listo para tomar notas
- [ ] Acceso al repositorio local
- [ ] gh CLI configurado

---

## 🎯 Próxima Acción

```bash
# Opción A: Comenzar implementación inmediata
open SPRINT-2-TASKS.md

# Opción B: Entender contexto primero
open README.md

# Opción C: Ver solo problemas críticos
# Buscar sección "Problemas Críticos" en README.md

# Opción D: Quick win (resolver fallo)
gh run view 19485500426 --repo EduGoGroup/edugo-api-administracion --log-failed
```

---

## 📊 Métricas del Plan

| Métrica | Valor |
|---------|-------|
| Archivos totales | 4 markdown |
| Líneas totales (est.) | ~3,700 |
| Scripts incluidos | ~35 bash scripts |
| Tareas P0 | 3 tareas |
| Tareas P1 | 4 tareas |
| Tareas P2 | 5 tareas |
| Tiempo P0 | 4-6 horas |
| Tiempo total Sprint 2 | 18-22 horas |
| Tiempo total Sprint 4 | 12-15 horas |
| Nivel de detalle | Ultra-alto |

---

## 🎉 ¡Listo para Comenzar!

Has llegado al final del índice. Tienes una visión completa del proyecto.

**Siguiente paso recomendado:**

```bash
# Para entender el contexto
open README.md

# Para empezar a trabajar
open SPRINT-2-TASKS.md
```

---

**Generado por:** Claude Code  
**Fecha:** 19 de Noviembre, 2025  
**Versión:** 1.0  
**Basado en:** Plan de api-mobile + análisis específico de api-administracion
