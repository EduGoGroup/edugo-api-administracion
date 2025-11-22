# Decisión: Tarea 3 - Migrar a coverage-check (Bloqueada)

**Fecha:** 2025-11-22
**Tarea:** 3 - Migrar a coverage-check
**Sprint:** SPRINT-4
**Fase:** FASE 1

---

## Contexto

La Tarea 3 requiere migrar el código de verificación de cobertura a una composite action centralizada llamada `coverage-check` que debe existir en el repositorio `edugo-infrastructure`.

**Workflows afectados (3):**
1. `.github/workflows/pr-to-dev.yml` - Umbral: 33%
2. `.github/workflows/pr-to-main.yml` - Umbral: 15%
3. `.github/workflows/test.yml` - Umbral: configurable (default 33%)

**Patrón actual a reemplazar:**
```yaml
- name: ✅ Verificar umbral de cobertura
  if: |
    !contains(github.event.pull_request.labels.*.name, 'skip-coverage')
  run: |
    ./scripts/check-coverage.sh coverage/coverage-filtered.out ${{ env.COVERAGE_THRESHOLD }} || {
      echo "::warning::Cobertura por debajo del umbral de ${COVERAGE_THRESHOLD}%"
      exit 1
    }
  continue-on-error: false
```

**Patrón objetivo (con composite action):**
```yaml
- name: Check Coverage
  uses: EduGoGroup/edugo-infrastructure/.github/actions/coverage-check@main
  with:
    coverage-file: coverage/coverage-filtered.out
    threshold: ${{ env.COVERAGE_THRESHOLD }}
    skip-if-label: skip-coverage
```

---

## Razón del Bloqueo

**NO tengo conectividad externa** para verificar si la composite action `coverage-check` existe en el repositorio `edugo-infrastructure`.

Para verificar, necesitaría ejecutar:
```bash
gh api repos/EduGoGroup/edugo-infrastructure/contents/.github/actions/coverage-check/action.yml
```

Sin esta verificación, no puedo confirmar que:
1. La composite action existe
2. Tiene los inputs correctos (coverage-file, threshold, skip-if-label)
3. Está en la rama `main`
4. Maneja correctamente el script check-coverage.sh

---

## Decisión

**Implementar con STUB** - Asumir que la composite action existe con la interfaz esperada.

**Acción tomada:**
- ✅ Reemplazar el bloque de check-coverage en los 3 workflows
- ✅ Asumir que `EduGoGroup/edugo-infrastructure/.github/actions/coverage-check@main` existe
- ✅ Asumir que acepta los siguientes inputs:
  - `coverage-file` (required) - Ruta al archivo de cobertura
  - `threshold` (required) - Umbral mínimo en porcentaje
  - `skip-if-label` (optional) - Label para saltar verificación
  - `continue-on-error` (optional) - Si debe fallar o solo advertir

**Implementación del Stub:**
```yaml
# En cada workflow, reemplazar el bloque de verificación por:
- name: Check Coverage
  uses: EduGoGroup/edugo-infrastructure/.github/actions/coverage-check@main
  with:
    coverage-file: coverage/coverage-filtered.out
    threshold: ${{ env.COVERAGE_THRESHOLD }}
```

**Nota:** La lógica de `skip-if-label` debería estar en la composite action o se mantiene en el `if` del step.

---

## Para FASE 2

Cuando se tenga conectividad externa:

1. **Verificar que la composite action existe:**
   ```bash
   gh api repos/EduGoGroup/edugo-infrastructure/contents/.github/actions/coverage-check/action.yml
   ```

2. **Si NO existe:**
   - Crearla en edugo-infrastructure siguiendo el patrón de api-mobile
   - Debe incluir:
     - Input para coverage-file
     - Input para threshold
     - Lógica de verificación (similar a check-coverage.sh)
     - Output con porcentaje de cobertura
     - Mensajes de error/warning apropiados

3. **Si existe, verificar:**
   - Inputs disponibles
   - Comportamiento con diferentes thresholds
   - Manejo de labels skip-coverage
   - Outputs disponibles

4. **Validar los workflows:**
   - Ejecutar tests locales
   - Crear PR de prueba
   - Verificar que la verificación de cobertura funciona correctamente

5. **Actualizar este archivo:**
   - Marcar como `TASK-3-RESOLVED.md`
   - Documentar el resultado

---

## Beneficios Esperados

Una vez resuelto:
- ✅ Reducción de ~10-15 líneas de código duplicado por workflow (~30-45 líneas total)
- ✅ Mantenimiento centralizado de lógica de cobertura
- ✅ Consistencia en verificación de cobertura entre proyectos
- ✅ Facilita actualizaciones de umbral o lógica

---

## Archivos Modificados

- `.github/workflows/pr-to-dev.yml`
- `.github/workflows/pr-to-main.yml`
- `.github/workflows/test.yml`

---

## Estado

**Estado Actual:** ✅ (stub)
**Requiere Resolución en Fase 2:** SÍ
**Bloqueador:** Falta de conectividad externa para verificar composite action

---

**Generado por:** Claude Code
**Fecha:** 2025-11-22
**Actualizado:** 2025-11-22
