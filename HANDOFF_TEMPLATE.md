# Template: Handoff Fase 1 → Fase 2

**Usar este template al finalizar Fase 1 (Claude Web)**

---

## ✅ Completado en Fase 1

### Archivos Creados
- [ ] Lista de archivos nuevos con descripción

### Código Implementado  
- [ ] Lista de funciones/métodos implementados
- [ ] Compilación: ✅ / ❌

### Tests Escritos
- [ ] Lista de tests creados
- [ ] Marcados con `t.Skip()`: Sí / No

---

## ⏸️ Stubs/Mocks Creados

### Stub #1: [Nombre]
**Ubicación:** `path/to/file.go:line`  
**Razón:** [Por qué es stub - ej: "Requiere PostgreSQL"]  
**Qué hacer en Fase 2:**
```
1. Paso 1
2. Paso 2
```

### Stub #2: [Nombre]
**Ubicación:** `path/to/file.go:line`  
**Razón:** [Por qué es stub]  
**Qué hacer en Fase 2:**
```
1. Paso 1
2. Paso 2
```

---

## 🔧 Pendiente para Fase 2

### Migraciones
- [ ] Ejecutar `migrations/XXX.up.sql`
- [ ] Validar que funcionan
- [ ] Validar triggers

### Tests de Integración
**Archivo:** `test/integration/XXX_test.go`

**Tests a descomentar:**
- [ ] TestFunctionName1
- [ ] TestFunctionName2
- [ ] TestFunctionName3

**Comando:**
```bash
go test -tags=integration ./test/integration/... -v
```

### Validaciones con DB Real
- [ ] Query 1 funciona
- [ ] Query 2 funciona
- [ ] Performance validada

---

## 📊 Métricas Baseline

**Pre-Fase 2:**
- Tests pasando: X/Y
- Tests skipeados: Z
- Coverage: ___%

**Objetivo Post-Fase 2:**
- Tests pasando: 100%
- Tests skipeados: 0
- Coverage: >= 80% (repository)

---

## 🚀 Comandos para Fase 2

### Setup
```bash
git checkout feature/sprint-03-repositorios-ltree
git pull origin feature/sprint-03-repositorios-ltree
```

### Desarrollo
```bash
# Quitar stubs
sed -i '' '/t.Skip/d' test/integration/file_test.go

# Ejecutar tests
go test -tags=integration ./test/integration/... -v

# Benchmark
go test -tags=integration ./test/integration/... -bench=.
```

### Finalización
```bash
# Validación completa
make test-unit
make test-integration  
make lint
make coverage-report

# Commit
git add -A
git commit -m "feat(infrastructure): complete ltree implementation (FASE 2)"
git push origin feature/sprint-03-repositorios-ltree
```

---

## ⚠️ Problemas Conocidos

### Problema #1: [Si hay alguno]
**Descripción:** ...  
**Workaround temporal:** ...  
**Solución en Fase 2:** ...

---

## 📝 Notas de Fase 1

[Claude Web: Agrega aquí cualquier nota relevante, decisiones tomadas, trade-offs, etc.]

---

## ✅ Checklist de Handoff

**Claude Web debe marcar antes de finalizar:**
- [ ] Código compila
- [ ] Lint pasa
- [ ] HANDOFF_FASE1_A_FASE2.md creado
- [ ] Stubs claramente documentados
- [ ] Branch pusheada
- [ ] Mensaje de finalización dejado

**Claude Local debe validar al inicio:**
- [ ] Handoff document leído
- [ ] Branch descargada
- [ ] Compilación verificada
- [ ] Plan de Fase 2 claro

---

**Última actualización:** [Fecha]  
**Creado por:** Claude Code Web  
**Para:** Claude Code Local
