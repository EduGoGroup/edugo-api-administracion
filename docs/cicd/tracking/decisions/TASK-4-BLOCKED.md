# Decisión: Tarea 4 - Migrar sync-main-to-dev.yml (Bloqueada)

**Fecha:** 2025-11-22
**Tarea:** 4 - Migrar sync-main-to-dev.yml
**Sprint:** SPRINT-4
**Fase:** FASE 1

---

## Contexto

La Tarea 4 requiere reemplazar el workflow completo `sync-main-to-dev.yml` con una llamada a un workflow reusable que debe existir en el repositorio `edugo-infrastructure`.

**Archivo afectado:**
- `.github/workflows/sync-main-to-dev.yml` (~50 líneas actualmente)

**Patrón objetivo:**
Reemplazar el workflow completo con una llamada a workflow reusable:
```yaml
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
```

---

## Razón del Bloqueo

**NO tengo conectividad externa** para verificar si el workflow reusable `sync-branches.yml` existe en el repositorio `edugo-infrastructure`.

Para verificar, necesitaría ejecutar:
```bash
gh api repos/EduGoGroup/edugo-infrastructure/contents/.github/workflows/reusable/sync-branches.yml
```

---

## Análisis: ¿Implementar Stub o SKIP?

### Opción A: SKIP - No implementar (RECOMENDADA)

**Razones:**
1. ✅ **Workflow completo** - No es una action sino un workflow entero
2. ✅ **Funcionando correctamente** - sync-main-to-dev.yml funciona bien actualmente
3. ✅ **Riesgo medio** - Error podría afectar sincronización entre branches
4. ✅ **Ganancia menor** - Solo ~40 líneas, no crítico
5. ✅ **Requiere validación** - Workflows reusables tienen limitaciones específicas

**Impacto de SKIP:**
- ❌ No se elimina ~40 líneas de código
- ❌ Sincronización de branches sigue descentralizada
- ✅ Sincronización main→dev continúa funcionando sin riesgo

### Opción B: Stub - Implementar con suposiciones

**Riesgos del stub:**
- ⚠️ **Medio** - Si el workflow reusable no existe, la sincronización fallará
- ⚠️ **Sin validación** - No puedo probar localmente workflows reusables
- ⚠️ **Complejidad** - Workflows reusables tienen restricciones de permisos y secrets

---

## Decisión

**SKIP** - No implementar esta tarea en FASE 1.

**Justificación:**
- El workflow actual funciona correctamente
- Riesgo de afectar sincronización de branches
- Ganancia marginal (~40 líneas)
- Requiere validación de permisos y secrets en workflow reusable

**Acción tomada:**
- ⏭️ Marcar tarea como SKIP
- ⏭️ Documentar razón en este archivo
- ⏭️ Continuar con Tarea 5

---

## Para Futuro (Post-FASE 1)

Si en el futuro se decide implementar:

1. **Verificar que el workflow reusable existe:**
   ```bash
   gh api repos/EduGoGroup/edugo-infrastructure/contents/.github/workflows/reusable/sync-branches.yml
   ```

2. **Si existe, validar:**
   - Inputs disponibles (source-branch, target-branch)
   - Permisos requeridos
   - Manejo de secrets
   - Comportamiento con tags

3. **Testing:**
   - Probar en repository de prueba primero
   - Validar que sincronización funciona igual que antes
   - Verificar que tags también sincronizan correctamente

---

## Estado

**Estado Actual:** ⏭️ SKIP
**Razón:** Workflow funcionando + Riesgo innecesario + Ganancia marginal
**Requiere Resolución:** NO - Opcional
**Bloqueador:** Riesgo de afectar sincronización de branches sin validación

---

## Beneficios Perdidos al SKIP

- ❌ No se elimina ~40 líneas de código duplicado
- ❌ Sincronización sigue descentralizada

## Beneficios de SKIP

- ✅ Sincronización main→dev sigue funcionando sin riesgo
- ✅ No introducir posibles errores en flujo de sincronización
- ✅ Tiempo ahorrado para tareas de mayor valor

---

**Generado por:** Claude Code
**Fecha:** 2025-11-22
**Actualizado:** 2025-11-22
