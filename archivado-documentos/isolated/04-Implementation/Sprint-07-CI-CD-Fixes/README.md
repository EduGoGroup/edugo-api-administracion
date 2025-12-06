# 📋 Sprint-07: Corrección Completa de CI/CD

**Fecha:** 17 de Noviembre, 2025  
**Versión:** 1.0.0  
**Duración Total:** 3 horas  
**Prioridad:** ALTA (workflows fallando en múltiples PRs)

---

## 🎯 Objetivo

Corregir completamente la infraestructura de CI/CD de `edugo-api-administracion` copiando archivos faltantes desde `edugo-api-mobile` (que funciona correctamente) y eliminando duplicaciones.

**Filosofía:** "Plomo al ampa" - Hacerlo bien desde el inicio, no parches temporales.

---

## 🐛 Problemas Identificados

### 1. **Scripts Faltantes** ❌ CRÍTICO

```bash
# Workflows hacen referencia a:
./scripts/check-coverage.sh          # ❌ NO EXISTE
./scripts/filter-coverage.sh         # ❌ NO EXISTE

# Resultado:
# - pr-to-main.yml falla en step "Verificar umbral de cobertura"
# - test.yml falla en step "Verificar umbral de cobertura"
```

**Impacto:** 
- PRs bloqueados
- No hay validación real de cobertura
- Falsos positivos en GitHub Actions

---

### 2. **Comandos Makefile Faltantes** ❌ CRÍTICO

```bash
# Workflows ejecutan:
make coverage-report                  # ❌ NO EXISTE

# Makefile actual solo tiene:
make test-coverage                    # ✅ Genera HTML pero diferente formato
```

**Impacto:**
- Steps de workflows fallan
- No se generan archivos esperados (`coverage/coverage-filtered.out`)

---

### 3. **Archivo .coverignore Faltante** ⚠️ MEDIO

```bash
# filter-coverage.sh busca:
.coverignore                          # ❌ NO EXISTE

# Resultado:
# - Cobertura incluye código generado (docs/docs.go)
# - Cobertura incluye DTOs/requests/responses
# - Métricas infladas o incorrectas
```

**Impacto:**
- Métricas de cobertura incorrectas
- Dificultad para alcanzar umbrales reales

---

### 4. **Workflows Duplicados** ⚠️ BAJO

```bash
.github/workflows/
├── ci.yml                    # ⚠️ DUPLICADO (contenido similar a pr-to-dev.yml)
├── docker-only.yml           # ⚠️ DUPLICADO (contenido en build-and-push.yml)
└── build-and-push.yml        # ✅ Funcional
```

**Impacto:**
- Confusión sobre cuál usar
- Mantenimiento duplicado
- Tiempo de CI desperdiciado

---

### 5. **Errores de Sintaxis YAML** ❌ CRÍTICO

**En `pr-to-main.yml:51-52`:**

```yaml
- name: 📊 Generar reporte de cobertura
  run: make coverage-report
  continue-on-error: true
  continue-on-error: true    # ❌ DUPLICADO - Error de sintaxis YAML
```

**Impacto:**
- Workflow puede fallar en parsing
- Comportamiento impredecible

---

### 6. **Versión Go Incorrecta** ⚠️ BAJO

```yaml
# En algunos workflows:
GO_VERSION: "1.25.3"          # ❌ NO EXISTE (última es 1.23.x)

# Debería ser:
GO_VERSION: "1.24"            # ✅ Como en test.yml
```

**Impacto:**
- Setup de Go falla en algunos runners
- Tests no se ejecutan

---

## 📊 Resumen de Correcciones

| Problema | Severidad | Solución | Origen |
|----------|-----------|----------|--------|
| Scripts faltantes | ❌ CRÍTICO | Copiar desde api-mobile | ✅ Existe |
| Comando Makefile | ❌ CRÍTICO | Copiar desde api-mobile | ✅ Existe |
| .coverignore | ⚠️ MEDIO | Copiar y adaptar de api-mobile | ✅ Existe |
| Workflows duplicados | ⚠️ BAJO | Eliminar ci.yml y docker-only.yml | - |
| Sintaxis YAML | ❌ CRÍTICO | Eliminar línea duplicada | - |
| Versión Go | ⚠️ BAJO | Estandarizar a 1.24 | - |

---

## 🚀 Plan de Ejecución

### TASK 1: Copiar Scripts desde api-mobile (30 min)

#### 1.1 Copiar `check-coverage.sh`

```bash
# Desde api-mobile a api-admin
cp /path/to/edugo-api-mobile/scripts/check-coverage.sh \
   /path/to/edugo-api-administracion/scripts/

# Dar permisos de ejecución
chmod +x scripts/check-coverage.sh
```

**Contenido esperado:**
- Lectura de archivo de cobertura
- Cálculo de porcentaje total
- Comparación con umbral
- Salida colorizada
- Exit code 0 (éxito) o 1 (falla)

**Validación:**
```bash
./scripts/check-coverage.sh coverage/coverage.out 33
# Debe mostrar porcentaje y comparar con 33%
```

---

#### 1.2 Copiar `filter-coverage.sh`

```bash
cp /path/to/edugo-api-mobile/scripts/filter-coverage.sh \
   /path/to/edugo-api-administracion/scripts/

chmod +x scripts/filter-coverage.sh
```

**Contenido esperado:**
- Lectura de `.coverignore`
- Filtrado de líneas de cobertura
- Generación de `coverage-filtered.out`
- Reporte de líneas filtradas

**Validación:**
```bash
./scripts/filter-coverage.sh coverage/coverage.out
# Debe generar: coverage/coverage-filtered.out
```

---

### TASK 2: Crear .coverignore (15 min)

#### 2.1 Copiar base desde api-mobile

```bash
cp /path/to/edugo-api-mobile/.coverignore \
   /path/to/edugo-api-administracion/.coverignore
```

#### 2.2 Adaptar patrones a api-admin

**Contenido de `.coverignore`:**

```gitignore
# ============================================
# .coverignore - API Administración
# ============================================
# Archivos excluidos del cálculo de cobertura

# Código generado
docs/docs.go
docs/swagger.json
docs/swagger.yaml

# DTOs (solo estructuras, sin lógica)
internal/application/dto/*

# HTTP Requests/Responses (solo mapeo)
internal/infrastructure/http/request/*
internal/infrastructure/http/response/*

# Entidades de dominio (solo datos)
internal/domain/entity/*_entity.go

# Value Objects (mayormente validación simple)
internal/domain/valueobject/*

# Configuración (solo structs)
internal/infrastructure/config/*

# Main (bootstrap, difícil de testear)
cmd/main.go

# Tests (no cubrir tests)
*_test.go
test/*
mock/*

# Migraciones SQL (no código Go)
scripts/postgresql/*

# Third-party o generado
vendor/*
```

**Validación:**
```bash
# Verificar que filter-coverage.sh lo lee
./scripts/filter-coverage.sh coverage/coverage.out | grep "Filtered"
# Debe mostrar: "Filtered X lines from coverage"
```

---

### TASK 3: Actualizar Makefile (30 min)

#### 3.1 Agregar comando `coverage-report`

Agregar al `Makefile` en la sección **Testing**:

```makefile
# ============================================
# Testing
# ============================================

coverage-report: ## Generar reporte de cobertura filtrado para CI
	@echo "$(YELLOW)📊 Generando reporte de cobertura (filtrado)...$(RESET)"
	@mkdir -p $(COVERAGE_DIR)
	@echo "$(BLUE)→ Ejecutando tests con cobertura...$(RESET)"
	@$(GOTEST) -v -race -coverprofile=$(COVERAGE_DIR)/coverage.out -covermode=atomic ./...
	@echo "$(BLUE)→ Filtrando archivos según .coverignore...$(RESET)"
	@./scripts/filter-coverage.sh $(COVERAGE_DIR)/coverage.out
	@echo "$(BLUE)→ Generando reporte HTML...$(RESET)"
	@$(GOCMD) tool cover -html=$(COVERAGE_DIR)/coverage-filtered.out -o $(COVERAGE_DIR)/coverage.html
	@echo "$(BLUE)→ Resumen de cobertura:$(RESET)"
	@$(GOCMD) tool cover -func=$(COVERAGE_DIR)/coverage-filtered.out | tail -1
	@echo "$(GREEN)✓ Reportes generados:$(RESET)"
	@echo "  - $(COVERAGE_DIR)/coverage.out (completo)"
	@echo "  - $(COVERAGE_DIR)/coverage-filtered.out (filtrado)"
	@echo "  - $(COVERAGE_DIR)/coverage.html"
	@echo "$(BLUE)💡 Abrir reporte: open $(COVERAGE_DIR)/coverage.html$(RESET)"

coverage-check: coverage-report ## Verificar umbral de cobertura
	@echo "$(YELLOW)✅ Verificando umbral de cobertura...$(RESET)"
	@./scripts/check-coverage.sh $(COVERAGE_DIR)/coverage-filtered.out 33
```

**Validación:**
```bash
make coverage-report
# Debe generar:
# - coverage/coverage.out
# - coverage/coverage-filtered.out
# - coverage/coverage.html

make coverage-check
# Debe verificar umbral (33%) y salir con código correcto
```

---

### TASK 4: Corregir Workflows (45 min)

#### 4.1 Eliminar workflows duplicados

```bash
# Eliminar archivos duplicados
git rm .github/workflows/ci.yml
git rm .github/workflows/docker-only.yml

# Razón:
# - ci.yml es redundante con pr-to-dev.yml
# - docker-only.yml está contenido en build-and-push.yml
```

---

#### 4.2 Corregir `pr-to-main.yml`

**Cambios:**

1. **Eliminar línea duplicada** (línea 52):

```yaml
# ANTES:
- name: 📊 Generar reporte de cobertura
  run: make coverage-report
  continue-on-error: true
  continue-on-error: true    # ❌ DUPLICADO

# DESPUÉS:
- name: 📊 Generar reporte de cobertura
  run: make coverage-report
  continue-on-error: true
```

2. **Estandarizar versión de Go** (si tiene 1.25.3):

```yaml
# ANTES:
env:
  GO_VERSION: "1.25.3"

# DESPUÉS:
env:
  GO_VERSION: "1.24"
```

**Validación:**
```bash
# Verificar sintaxis YAML
yamllint .github/workflows/pr-to-main.yml

# O usar:
cat .github/workflows/pr-to-main.yml | docker run --rm -i cytopia/yamllint
```

---

#### 4.3 Verificar otros workflows

Revisar y estandarizar `GO_VERSION` en:
- `test.yml` (ya tiene 1.24 ✅)
- `pr-to-dev.yml`
- `release.yml`
- `build-and-push.yml`

**Script de verificación:**

```bash
grep -r "GO_VERSION" .github/workflows/*.yml

# Resultado esperado:
# Todos deben mostrar: GO_VERSION: "1.24"
```

---

### TASK 5: Testing Local (30 min)

#### 5.1 Probar scripts localmente

```bash
# 1. Generar cobertura
make test-coverage

# 2. Verificar que coverage-report funciona
make coverage-report

# 3. Verificar que coverage-check funciona
make coverage-check

# 4. Probar scripts directamente
./scripts/filter-coverage.sh coverage/coverage.out
./scripts/check-coverage.sh coverage/coverage-filtered.out 33
```

**Validación exitosa:**
- ✅ Archivos generados en `coverage/`
- ✅ `coverage-filtered.out` es más pequeño que `coverage.out`
- ✅ `check-coverage.sh` muestra porcentaje correcto
- ✅ Exit codes correctos (0 si pasa umbral, 1 si no)

---

#### 5.2 Simular workflow localmente (opcional)

```bash
# Usar act (GitHub Actions local)
brew install act  # macOS
# o: curl https://raw.githubusercontent.com/nektos/act/master/install.sh | sudo bash

# Ejecutar workflow de tests
act -j unit-tests

# Ejecutar workflow completo de PR
act pull_request -W .github/workflows/pr-to-main.yml
```

---

### TASK 6: Commit y Push (15 min)

#### 6.1 Crear feature branch

```bash
git checkout dev
git pull origin dev
git checkout -b feature/fix-ci-cd-infrastructure
```

---

#### 6.2 Commit de cambios

```bash
# Agregar archivos nuevos
git add scripts/check-coverage.sh
git add scripts/filter-coverage.sh
git add .coverignore

# Agregar cambios en Makefile
git add Makefile

# Agregar correcciones de workflows
git add .github/workflows/pr-to-main.yml
# (otros workflows si se modificaron)

# Eliminar workflows duplicados
git rm .github/workflows/ci.yml
git rm .github/workflows/docker-only.yml

# Commit
git commit -m "fix(ci): corregir infraestructura completa de CI/CD

- Agregar scripts faltantes desde api-mobile
  - check-coverage.sh: Verificar umbral de cobertura
  - filter-coverage.sh: Filtrar archivos según .coverignore
  
- Crear .coverignore adaptado a api-admin
  - Excluir código generado (docs/docs.go)
  - Excluir DTOs, requests, responses
  - Excluir entidades de dominio
  
- Actualizar Makefile
  - Nuevo: coverage-report (genera reporte filtrado)
  - Nuevo: coverage-check (verifica umbral)
  
- Corregir workflows
  - pr-to-main.yml: Eliminar línea duplicada continue-on-error
  - Estandarizar GO_VERSION a 1.24 en todos los workflows
  
- Eliminar workflows duplicados
  - ci.yml (redundante con pr-to-dev.yml)
  - docker-only.yml (contenido en build-and-push.yml)

BREAKING CHANGE: Workflows ahora requieren scripts/ completo
Closes #XX"
```

---

#### 6.3 Push y crear PR

```bash
git push origin feature/fix-ci-cd-infrastructure

# Crear PR en GitHub:
# Base: dev
# Head: feature/fix-ci-cd-infrastructure
# Título: fix(ci): Corregir infraestructura completa de CI/CD
```

---

### TASK 7: Validación en GitHub Actions (15 min)

#### 7.1 Verificar que workflows pasan

Después de crear el PR, verificar en GitHub Actions:

```
✅ Unit Tests
  ├── ✅ Ejecutar tests unitarios
  ├── ✅ Generar reporte de cobertura (make coverage-report)
  ├── ✅ Verificar umbral de cobertura (scripts/check-coverage.sh)
  └── ✅ Subir reporte de cobertura

✅ Integration Tests
  └── ✅ Tests de integración pasan

✅ Lint & Format Check
  └── ✅ golangci-lint pasa

✅ Security Scan
  └── ✅ Gosec pasa
```

**Si alguno falla:**
1. Revisar logs del step que falla
2. Corregir localmente
3. Push al mismo branch (PR se actualiza automáticamente)

---

## ✅ Checklist de Validación

Antes de aprobar y mergear el PR, verificar:

### Scripts
- [ ] `scripts/check-coverage.sh` existe y tiene permisos +x
- [ ] `scripts/filter-coverage.sh` existe y tiene permisos +x
- [ ] Ambos scripts funcionan localmente

### Configuración
- [ ] `.coverignore` existe en la raíz
- [ ] `.coverignore` adaptado a api-admin (no es copia exacta de api-mobile)
- [ ] Patrones en `.coverignore` son correctos

### Makefile
- [ ] `make coverage-report` funciona
- [ ] `make coverage-check` funciona
- [ ] Genera archivos esperados en `coverage/`

### Workflows
- [ ] `ci.yml` eliminado
- [ ] `docker-only.yml` eliminado
- [ ] `pr-to-main.yml` corregido (sin línea duplicada)
- [ ] `GO_VERSION: "1.24"` en todos los workflows
- [ ] Sintaxis YAML válida en todos los workflows

### Testing
- [ ] `make coverage-report` genera archivos correctos localmente
- [ ] `./scripts/check-coverage.sh` retorna exit code correcto
- [ ] Workflows pasan en GitHub Actions
- [ ] PR muestra todos los checks en verde ✅

### Documentación
- [ ] Este README.md documenta todos los cambios
- [ ] Commit message sigue conventional commits
- [ ] PR description explica el problema y la solución

---

## 📈 Mejoras Esperadas

| Métrica | Antes | Después | Mejora |
|---------|-------|---------|--------|
| **Workflows fallando** | ~15 runs | 0 | ✅ 100% |
| **Scripts faltantes** | 2 | 0 | ✅ 100% |
| **Workflows duplicados** | 2 | 0 | ✅ 100% |
| **Cobertura calculada** | Incorrecta | Correcta | ✅ Precisa |
| **Tiempo de debug CI** | ~2h/semana | ~0h | ✅ 100% |
| **Confianza en CI** | Baja | Alta | ✅ Mejorada |

---

## 🎯 Criterios de Éxito

Sprint-07 está **COMPLETADO** cuando:

1. ✅ Scripts copiados y funcionando
2. ✅ `.coverignore` creado y funcional
3. ✅ Makefile actualizado con nuevos comandos
4. ✅ Workflows duplicados eliminados
5. ✅ Errores de sintaxis YAML corregidos
6. ✅ Versiones de Go estandarizadas
7. ✅ PR con todos los checks en verde
8. ✅ PR mergeado a `dev`

---

## 🔗 Referencias

- **api-mobile (referencia):** `/path/to/edugo-api-mobile/`
- **Workflows originales:** `.github/workflows/`
- **Scripts copiados:** `scripts/`
- **Configuración:** `.coverignore`

---

## 📝 Notas Adicionales

### ¿Por qué copiar desde api-mobile?

1. **Probado y funcionando:** api-mobile tiene CI/CD estable sin fallos
2. **Mismo ecosistema:** Misma estructura de proyecto
3. **Consistencia:** Mantener estándares uniformes entre APIs
4. **Ahorro de tiempo:** No reinventar la rueda

### ¿Qué NO copiar?

- ❌ Workflows completos (tienen configuraciones específicas de api-mobile)
- ❌ `.coverignore` sin adaptar (diferentes packages)
- ✅ SÍ copiar: Scripts shell (son genéricos)
- ✅ SÍ copiar: Comandos Makefile (son genéricos)

### Mantenimiento futuro

Cuando api-mobile actualice sus scripts:
1. Revisar cambios en api-mobile
2. Evaluar si aplican a api-admin
3. Copiar cambios relevantes
4. Testear localmente
5. Push y verificar en CI

---

**Documento creado:** 17 de Noviembre, 2025  
**Versión:** 1.0.0  
**Tiempo estimado:** 3 horas  
**Próximo paso:** Ejecutar TASK 1 (Copiar scripts)
