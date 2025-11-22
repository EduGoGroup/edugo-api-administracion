# Decisión: Tarea 1 - Migrar a setup-edugo-go (Bloqueada)

**Fecha:** 2025-11-22
**Tarea:** 1 - Migrar a setup-edugo-go
**Sprint:** SPRINT-4
**Fase:** FASE 1

---

## Contexto

La Tarea 1 requiere migrar el código duplicado de setup Go + GOPRIVATE en todos los workflows a una composite action centralizada llamada `setup-edugo-go` que debe existir en el repositorio `edugo-infrastructure`.

**Workflows afectados (5):**
1. `.github/workflows/pr-to-dev.yml`
2. `.github/workflows/pr-to-main.yml`
3. `.github/workflows/test.yml`
4. `.github/workflows/manual-release.yml`
5. `.github/workflows/release.yml`

**Patrón actual a reemplazar:**
```yaml
- name: 🔧 Setup Go
  uses: actions/setup-go@v5
  with:
    go-version: ${{ env.GO_VERSION }}
    cache: true

- name: 🔐 Configurar acceso a repos privados
  run: |
    git config --global url."https://${{ secrets.GITHUB_TOKEN }}@github.com/".insteadOf "https://github.com/"
  env:
    GOPRIVATE: github.com/EduGoGroup/*
```

**Patrón objetivo (con composite action):**
```yaml
- name: Setup Go Environment
  uses: EduGoGroup/edugo-infrastructure/.github/actions/setup-edugo-go@main
```

---

## Razón del Bloqueo

**NO tengo conectividad externa** para verificar si la composite action `setup-edugo-go` existe en el repositorio `edugo-infrastructure`.

Para verificar, necesitaría ejecutar:
```bash
gh api repos/EduGoGroup/edugo-infrastructure/contents/.github/actions/setup-edugo-go/action.yml
```

Sin esta verificación, no puedo confirmar que:
1. La composite action existe
2. Tiene los inputs correctos (go-version, etc.)
3. Está en la rama `main`

---

## Decisión

**Implementar con STUB** - Asumir que la composite action existe con la interfaz esperada.

**Acción tomada:**
- ✅ Reemplazar el patrón de setup-go + GOPRIVATE en los 5 workflows
- ✅ Asumir que `EduGoGroup/edugo-infrastructure/.github/actions/setup-edugo-go@main` existe
- ✅ Asumir que acepta los siguientes inputs (inferidos de lecciones aprendidas):
  - `go-version` (opcional, default probablemente "1.25")
  - Configura automáticamente GOPRIVATE y git config

**Implementación del Stub:**
```yaml
# En cada workflow, reemplazar el bloque completo por:
- name: Setup Go Environment
  uses: EduGoGroup/edugo-infrastructure/.github/actions/setup-edugo-go@main
  with:
    go-version: ${{ env.GO_VERSION }}
```

---

## Para FASE 2

Cuando se tenga conectividad externa:

1. **Verificar que la composite action existe:**
   ```bash
   gh api repos/EduGoGroup/edugo-infrastructure/contents/.github/actions/setup-edugo-go/action.yml
   ```

2. **Si NO existe:**
   - Crearla en edugo-infrastructure siguiendo el patrón de api-mobile
   - Debe incluir:
     - Setup de Go con cache
     - Configuración de GOPRIVATE
     - Configuración de git para repos privados

3. **Si existe, verificar:**
   - Inputs disponibles
   - Versión de Go por defecto
   - Comportamiento de GOPRIVATE

4. **Validar los workflows:**
   - Ejecutar tests locales con act (si está disponible)
   - Crear PR de prueba
   - Verificar que los workflows funcionan correctamente

5. **Actualizar este archivo:**
   - Marcar como `TASK-1-RESOLVED.md`
   - Documentar el resultado

---

## Beneficios Esperados

Una vez resuelto:
- ✅ Reducción de ~30-40 líneas de código duplicado
- ✅ Mantenimiento centralizado
- ✅ Consistencia entre proyectos
- ✅ Facilita actualizaciones de versión de Go

---

## Archivos Modificados

- `.github/workflows/pr-to-dev.yml`
- `.github/workflows/pr-to-main.yml`
- `.github/workflows/test.yml`
- `.github/workflows/manual-release.yml`
- `.github/workflows/release.yml`

---

## Estado

**Estado Actual:** ✅ (stub)
**Requiere Resolución en Fase 2:** SÍ
**Bloqueador:** Falta de conectividad externa para verificar composite action

---

**Generado por:** Claude Code
**Fecha:** 2025-11-22
**Actualizado:** 2025-11-22
