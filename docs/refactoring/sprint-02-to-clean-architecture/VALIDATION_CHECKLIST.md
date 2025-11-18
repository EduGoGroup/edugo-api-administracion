# Checklist de Validación - Refactor Clean Architecture

**Proyecto:** edugo-api-administracion  
**Fecha:** 2025-11-17

---

## ✅ Pre-Refactor (Baseline)

### Código
- [ ] Branch limpia desde main
- [ ] Todos los tests pasando
- [ ] Build exitoso
- [ ] Lint sin errores

### Métricas Baseline
- [ ] Coverage total: _____%
- [ ] Coverage domain/entity: _____%
- [ ] Número de archivos en domain/: _____
- [ ] LOC en academic_unit.go: _____
- [ ] LOC en unit_membership.go: _____

---

## 🎯 FASE 1: Domain Services

### Implementación
- [ ] Directorio `internal/domain/service/` creado
- [ ] `academic_unit_service.go` implementado
- [ ] `membership_service.go` implementado
- [ ] Tests básicos creados

### Validación
- [ ] Compila sin errores
- [ ] Tests de service pasando
- [ ] Coverage service >= 80%
- [ ] Lint sin errores

### Code Review
- [ ] Todas las validaciones migradas
- [ ] Sin lógica duplicada
- [ ] Documentación adecuada
- [ ] Tests cubren edge cases

---

## 🎯 FASE 2: Entities Anemic

### Implementación
- [ ] `academic_unit.go` simplificado
- [ ] `unit_membership.go` simplificado  
- [ ] Getters/setters agregados
- [ ] Métodos deprecated marcados

### Validación
- [ ] Compila sin errores
- [ ] LOC academic_unit.go <= 200
- [ ] LOC unit_membership.go <= 150
- [ ] Tests legacy siguen pasando

### Code Review
- [ ] Solo datos + getters/setters
- [ ] Sin lógica de negocio
- [ ] Deprecated correctamente marcados

---

## 🎯 FASE 3: Tests

### Implementación
- [ ] Tests migrados a service_test.go
- [ ] Tests de entity reducidos
- [ ] Coverage validada

### Validación
- [ ] Todos los tests pasando
- [ ] Coverage service >= 85%
- [ ] Coverage entity >= 90% (solo getters/setters)
- [ ] No hay tests duplicados

### Code Review
- [ ] Misma cobertura que antes
- [ ] Tests claros y mantenibles
- [ ] Usa table-driven tests donde aplique

---

## 🎯 FASE 4: Application Layer

### Implementación
- [ ] Application services actualizados
- [ ] Domain services inyectados
- [ ] Dependency container configurado
- [ ] Llamadas migradas de entity a service

### Validación
- [ ] Compila sin errores
- [ ] Integration tests pasando
- [ ] No hay imports cíclicos
- [ ] Dependency injection funciona

### Code Review
- [ ] Services correctamente inyectados
- [ ] No hay llamadas directas a entity methods
- [ ] Error handling apropiado

---

## 🎯 FASE 5: Validación Final

### Tests Completos
- [ ] `make test-unit` ✅
- [ ] `make test-integration` ✅
- [ ] `make coverage-report` ✅
- [ ] Coverage total >= 35%

### Build & Lint
- [ ] `make build` ✅
- [ ] `make lint` ✅
- [ ] `go vet ./...` ✅
- [ ] No warnings

### Limpieza
- [ ] Deprecated methods eliminados
- [ ] Imports no usados eliminados
- [ ] Comentarios obsoletos removidos
- [ ] `.coverignore` actualizado

### Documentación
- [ ] README.md actualizado
- [ ] ARCHITECTURE.md actualizado
- [ ] Ejemplos de código actualizados
- [ ] CHANGELOG.md actualizado

---

## 📊 Validación de Métricas

### Post-Refactor vs Pre-Refactor

| Métrica | Pre | Post | ✅/❌ |
|---------|-----|------|-------|
| Coverage Total | 13.2% | ___% | |
| Coverage Service | 0% | ___% | |
| Coverage Entity | 48.2% | ___% | |
| LOC academic_unit.go | 400 | ___  | |
| LOC academic_unit_service.go | 0 | ___ | |
| Tests Pasando | X/Y | X/Y | |

---

## 🚀 Pre-PR Checklist

### Código
- [ ] Todos los commits con mensajes convencionales
- [ ] Sin console.logs o debug code
- [ ] Sin TODOs sin issue asociado
- [ ] Code formateado (gofmt)

### Tests
- [ ] 100% tests pasando
- [ ] Coverage >= objetivo
- [ ] Integration tests validados
- [ ] Performance similar (±5%)

### Documentación
- [ ] README actualizado
- [ ] Comentarios de código claros
- [ ] Ejemplos funcionando
- [ ] Migration guide si es necesario

### PR
- [ ] PR description completa
- [ ] Screenshots/demos si aplica
- [ ] Breaking changes documentados
- [ ] Reviewers asignados

---

## ⚠️ Red Flags (Detener si ocurre)

- 🚫 Coverage cae más del 5%
- 🚫 Tests fallando sin explicación
- 🚫 Performance degrada >10%
- 🚫 Build time aumenta significativamente
- 🚫 Imports cíclicos
- 🚫 Más de 3 niveles de abstracción

---

## ✅ Criterios de Aceptación Final

### Must Have
- ✅ Todos los tests pasando
- ✅ Coverage >= 35%
- ✅ Build exitoso
- ✅ Lint sin errores
- ✅ Documentación actualizada

### Nice to Have
- 🎯 Coverage >= 40%
- 🎯 Performance igual o mejor
- 🎯 Code review aprobado sin cambios
- 🎯 CI/CD pipeline verde

---

**Aprobación Final:**
- [ ] Validado por: _________________
- [ ] Fecha: _________________
- [ ] ¿Listo para merge? ☐ Sí ☐ No ☐ Con cambios
