# Estrategia de CI/CD - EduGo API Mobile

**Última actualización**: 9 de noviembre de 2025  
**Versión**: 2.0

---

## 📋 Resumen

Este documento describe la estrategia de CI/CD implementada para el proyecto, optimizada para balance entre velocidad y confiabilidad.

### Principios

1. **Velocidad en desarrollo**: Tests rápidos en PRs a `dev`
2. **Confiabilidad en producción**: Tests completos en PRs a `main`
3. **Feedback rápido**: Resultados en minutos, no horas
4. **Eficiencia de recursos**: Ejecutar solo lo necesario

---

## 🔄 Workflows por Tipo de PR

### 1. PR a `dev` (Feature → Dev)

**Archivo**: `.github/workflows/pr-to-dev.yml`

**Cuándo se ejecuta**: Al abrir/actualizar PR de cualquier rama hacia `dev`

**Tests ejecutados**:
- ✅ Tests Unitarios (77 tests, ~5 segundos)
- ✅ Lint & Format Check
- ✅ Verificación de cobertura (umbral: 33%)

**NO se ejecutan**:
- ❌ Tests de Integración (se ejecutan solo en PR a main)
- ❌ Security Scan (se ejecuta solo en PR a main)

**Tiempo total**: ~2-3 minutos

**Justificación**:
- Los PRs a `dev` son frecuentes (múltiples por día)
- Los tests unitarios son suficientes para validar lógica de negocio
- Los tests de integración son más lentos y se ejecutan antes de producción

**Ejemplo de flujo**:
```
feature/nueva-funcionalidad → dev
  ├─ ✅ Tests Unitarios (5s)
  ├─ ✅ Lint (30s)
  └─ ✅ Cobertura (5s)
  
Total: ~2 minutos
```

---

### 2. PR a `main` (Dev → Main)

**Archivo**: `.github/workflows/pr-to-main.yml`

**Cuándo se ejecuta**: Al abrir/actualizar PR de `dev` hacia `main`

**Tests ejecutados**:
- ✅ Tests Unitarios (77 tests, ~5 segundos)
- ✅ Tests de Integración (18 tests, ~1-2 minutos)
- ✅ Lint & Format Check
- ✅ Security Scan
- ✅ Verificación de cobertura (umbral: 33%)

**Tiempo total**: ~3-4 minutos

**Justificación**:
- Los PRs a `main` son menos frecuentes (1-2 por semana)
- Requieren validación completa antes de producción
- Los tests de integración validan que todo funciona end-to-end

**Ejemplo de flujo**:
```
dev → main
  ├─ ✅ Tests Unitarios (5s)
  ├─ ✅ Tests Integración (1-2 min) ← NUEVO
  ├─ ✅ Lint (30s)
  ├─ ✅ Security Scan (30s) ← NUEVO
  └─ ✅ Cobertura (5s)
  
Total: ~3-4 minutos
```

---

### 3. Ejecución Manual

**Archivo**: `.github/workflows/test.yml`

**Cuándo se ejecuta**: Manualmente desde GitHub Actions UI

**Opciones**:
- `unit`: Solo tests unitarios
- `integration`: Solo tests de integración
- `all`: Todos los tests

**Uso**:
1. Ir a Actions → Tests with Coverage (Manual)
2. Click en "Run workflow"
3. Seleccionar tipo de tests
4. Configurar umbral de cobertura (opcional)

---

## 📊 Comparación de Tiempos

### Antes de la Optimización

| Workflow | Tests | Tiempo |
|----------|-------|--------|
| PR a dev | Unit + Integration | ~8-10 min |
| PR a main | Unit + Integration | ~8-10 min |

**Problema**: Todos los PRs tardaban lo mismo, ralentizando desarrollo

### Después de la Optimización

| Workflow | Tests | Tiempo | Mejora |
|----------|-------|--------|--------|
| PR a dev | Unit only | ~2-3 min | **-70%** 🚀 |
| PR a main | Unit + Integration | ~3-4 min | **-60%** 🚀 |

**Beneficio**: PRs a dev son 3x más rápidos, manteniendo calidad en main

---

## 🎯 Métricas de Performance

### Tests Unitarios

- **Cantidad**: 77 tests
- **Tiempo**: ~5 segundos
- **Cobertura**: 33.6%
- **Ejecutados en**: Todos los PRs

### Tests de Integración

- **Cantidad**: 18 tests
- **Tiempo**: ~1-2 minutos (con contenedores compartidos)
- **Mejora vs antes**: 81.5% más rápido (de 7:18 a 1:21)
- **Ejecutados en**: Solo PRs a main

### Optimizaciones Implementadas

1. **Contenedores Compartidos**: Reutilización entre tests (-81.5% tiempo)
2. **RabbitMQ Ligero**: Sin management plugin (-3s por setup)
3. **Retry Logic**: Manejo de errores TCP temporales
4. **Cleanup Optimizado**: TRUNCATE en lugar de DROP (-0.5s por test)

---

## 🔒 Checks Obligatorios

### PR a Dev

| Check | Obligatorio | Puede Fallar CI |
|-------|-------------|-----------------|
| Tests Unitarios | ✅ Sí | ✅ Sí |
| Cobertura >= 33% | ✅ Sí | ✅ Sí |
| Lint | ✅ Sí | ✅ Sí |

**Excepciones**:
- Label `skip-coverage`: Permite merge sin cumplir umbral de cobertura

### PR a Main

| Check | Obligatorio | Puede Fallar CI |
|-------|-------------|-----------------|
| Tests Unitarios | ✅ Sí | ✅ Sí |
| Tests Integración | ✅ Sí | ✅ Sí |
| Cobertura >= 33% | ✅ Sí | ✅ Sí |
| Lint | ✅ Sí | ✅ Sí |
| Security Scan | ⚠️ Recomendado | ❌ No |

**Nota**: Security scan no bloquea merge pero genera alertas

---

## 📝 Comentarios Automáticos en PRs

### PR a Dev

```markdown
## 📊 Cobertura de Tests Unitarios

**Cobertura Total**: 33.6%
**Umbral Mínimo**: 33%

✅ Cobertura cumple con el umbral

📄 Ver reporte completo
```

### PR a Main

```markdown
## 🚀 Resumen de Checks - PR a Main

| Check | Estado | Descripción |
|-------|--------|-------------|
| Tests Unitarios | ✅ success | 77 tests, ~5 segundos |
| Tests Integración | ✅ success | 18 tests, ~1-2 minutos |
| Lint & Format | ✅ success | Calidad de código |
| Security Scan | ✅ success | Análisis de seguridad |

---

✅ **Todos los checks pasaron** - PR listo para merge a main 🎉

### 📊 Métricas de Performance

- **Tiempo total estimado**: ~3-4 minutos
- **Tests unitarios**: ~5 segundos (77 tests)
- **Tests integración**: ~1-2 minutos (18 tests con contenedores compartidos)
- **Mejora vs anterior**: 81.5% más rápido 🚀
```

---

## 🚀 Flujo de Trabajo Recomendado

### Desarrollo de Feature

```bash
# 1. Crear rama desde dev
git checkout dev
git pull origin dev
git checkout -b feature/nueva-funcionalidad

# 2. Desarrollar y hacer commits
git add .
git commit -m "feat: nueva funcionalidad"

# 3. Ejecutar tests localmente (opcional pero recomendado)
make test-unit

# 4. Push y crear PR a dev
git push origin feature/nueva-funcionalidad
# Crear PR en GitHub: feature/nueva-funcionalidad → dev

# 5. Esperar checks (~2-3 min)
# - Tests Unitarios ✅
# - Lint ✅
# - Cobertura ✅

# 6. Merge a dev después de review
```

### Release a Producción

```bash
# 1. Crear PR de dev a main
# En GitHub: dev → main

# 2. Esperar checks completos (~3-4 min)
# - Tests Unitarios ✅
# - Tests Integración ✅
# - Lint ✅
# - Security Scan ✅
# - Cobertura ✅

# 3. Review y aprobación

# 4. Merge a main
# → Deploy automático a producción
```

---

## 🔧 Configuración de Variables

### Variables de Entorno

| Variable | Valor | Descripción |
|----------|-------|-------------|
| `GO_VERSION` | 1.25.3 | Versión de Go |
| `COVERAGE_THRESHOLD` | 33 | Umbral mínimo de cobertura |
| `RUN_INTEGRATION_TESTS` | true | Habilitar tests de integración |

### Secrets Requeridos

| Secret | Descripción | Usado en |
|--------|-------------|----------|
| `GITHUB_TOKEN` | Token automático de GitHub | Todos los workflows |
| `CODECOV_TOKEN` | Token de Codecov (opcional) | PR a main |

---

## 📈 Monitoreo y Métricas

### Dashboards

1. **GitHub Actions**: Ver historial de ejecuciones
2. **Codecov**: Ver evolución de cobertura
3. **GitHub Insights**: Ver tiempo promedio de PRs

### Métricas a Rastrear

- **Tiempo promedio de PR a dev**: Objetivo < 3 minutos
- **Tiempo promedio de PR a main**: Objetivo < 5 minutos
- **Tasa de éxito de tests**: Objetivo > 95%
- **Cobertura de código**: Objetivo >= 33% (incrementar gradualmente)

---

## 🐛 Troubleshooting

### Tests de Integración Fallan en CI

**Problema**: Tests pasan localmente pero fallan en CI

**Soluciones**:
1. Verificar que Docker esté disponible
2. Verificar timeouts (pueden ser más lentos en CI)
3. Revisar logs de contenedores
4. Verificar que `RUN_INTEGRATION_TESTS=true` esté configurado

### Tests Unitarios Lentos

**Problema**: Tests unitarios tardan más de 10 segundos

**Soluciones**:
1. Verificar que no haya tests de integración mezclados
2. Usar mocks en lugar de dependencias reales
3. Paralelizar tests con `t.Parallel()`

### Cobertura Baja

**Problema**: PR rechazado por cobertura < 33%

**Soluciones**:
1. Agregar tests para código nuevo
2. Usar label `skip-coverage` si es temporal
3. Revisar qué código no está cubierto: `make coverage-html`

---

## 📚 Referencias

- **Documentación de Tests**: `docs/TESTING_GUIDE.md`
- **Análisis de Performance**: `docs/TEST_PERFORMANCE_ANALYSIS.md`
- **Resultados de Optimización**: `docs/TEST_PERFORMANCE_RESULTS.md`
- **Plan de Cobertura**: `docs/TEST_COVERAGE_PLAN.md`

---

## 🔄 Changelog

### v2.0 (2025-11-09)

- ✅ Separación de workflows por tipo de PR
- ✅ Tests de integración solo en PR a main
- ✅ Optimización de tests de integración (81.5% más rápido)
- ✅ Comentarios automáticos en PRs
- ✅ Security scan en PR a main

### v1.0 (Anterior)

- Tests unitarios + integración en todos los PRs
- Tiempo: ~8-10 minutos por PR
- Sin optimización de contenedores

---

**Mantenido por**: Equipo de DevOps  
**Última revisión**: 9 de noviembre de 2025

