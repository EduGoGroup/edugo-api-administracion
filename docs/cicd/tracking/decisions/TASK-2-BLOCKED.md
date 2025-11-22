# Decisión: Tarea 2 - Migrar a docker-build-edugo (Bloqueada)

**Fecha:** 2025-11-22
**Tarea:** 2 - Migrar a docker-build-edugo
**Sprint:** SPRINT-4
**Fase:** FASE 1

---

## Contexto

La Tarea 2 requiere migrar los bloques de Docker build en los workflows de release a una composite action centralizada llamada `docker-build-edugo` que debe existir en el repositorio `edugo-infrastructure`.

**Workflows afectados (2):**
1. `.github/workflows/manual-release.yml` - Job: `build-docker-image`
2. `.github/workflows/release.yml` - Job: `build-and-push-docker`

**Bloque actual en manual-release.yml:**
```yaml
- name: Setup Docker Buildx
  uses: docker/setup-buildx-action@v3

- name: Login a GitHub Container Registry
  uses: docker/login-action@v3
  with:
    registry: ${{ env.REGISTRY }}
    username: ${{ github.actor }}
    password: ${{ secrets.GITHUB_TOKEN }}

- name: Extraer metadata para Docker
  id: meta
  uses: docker/metadata-action@v5
  with:
    images: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}
    tags: |
      type=raw,value=v${{ inputs.version }}
      type=raw,value=${{ inputs.version }}
      type=raw,value=latest

- name: Build and push Docker image
  uses: docker/build-push-action@v5
  with:
    context: .
    platforms: linux/amd64,linux/arm64
    push: true
    tags: ${{ steps.meta.outputs.tags }}
    labels: ${{ steps.meta.outputs.labels }}
    cache-from: type=gha
    cache-to: type=gha,mode=max
    build-args: |
      VERSION=v${{ inputs.version }}
      BUILD_DATE=${{ github.event.repository.updated_at }}
```

**Bloque en release.yml es similar pero con tags diferentes (semver).**

---

## Razón del Bloqueo

**NO tengo conectividad externa** para verificar si la composite action `docker-build-edugo` existe en el repositorio `edugo-infrastructure`.

Para verificar, necesitaría ejecutar:
```bash
gh api repos/EduGoGroup/edugo-infrastructure/contents/.github/actions/docker-build-edugo/action.yml
```

Además, el bloque de Docker build es **muy complejo** y requiere:
1. **Setup Buildx** - Configuración del builder multi-plataforma
2. **Login a GHCR** - Autenticación con secrets
3. **Metadata** - Generación de tags y labels dinámicos
4. **Build & Push** - Build multi-plataforma con cache

Una composite action que reemplace todo esto debería:
- Aceptar múltiples inputs (image-name, version, platforms, tags, build-args, etc.)
- Manejar diferentes estrategias de tags (semver, raw, sha)
- Configurar cache correctamente
- Tener outputs para usar en pasos posteriores

---

## Análisis: ¿Implementar Stub o SKIP?

### Opción A: SKIP - No implementar (RECOMENDADA)

**Razones:**
1. ✅ **Complejidad alta** - El bloque requiere 4 pasos coordinados
2. ✅ **No crítico** - Solo afecta workflows de release (manual y automático)
3. ✅ **Funcionando actualmente** - Los workflows de release funcionan bien
4. ✅ **Riesgo medio** - Error en stub podría romper releases
5. ✅ **Opcional en plan** - SPRINT-4-TASKS.md lo marca como posible

**Impacto de SKIP:**
- ❌ No se elimina duplicación en estos 2 workflows (~30 líneas por workflow)
- ❌ Mantenimiento de Docker build sigue descentralizado
- ✅ Releases siguen funcionando sin riesgo

### Opción B: Stub - Implementar con suposiciones

**Stub propuesto:**
```yaml
- name: Build and Push Docker Image
  uses: EduGoGroup/edugo-infrastructure/.github/actions/docker-build-edugo@main
  with:
    image-name: ${{ env.IMAGE_NAME }}
    registry: ${{ env.REGISTRY }}
    version: ${{ inputs.version }}  # o steps.tag.outputs.version
    platforms: linux/amd64,linux/arm64
    push: true
    tag-strategy: manual  # o 'semver' para release.yml
```

**Riesgos del stub:**
- ⚠️ **Alto** - Si la composite action no existe o tiene interfaz diferente, releases fallarán
- ⚠️ **Crítico** - Los releases son flujo importante de producción
- ⚠️ **Sin validación** - No puedo probar localmente

---

## Decisión

**SKIP** - No implementar esta tarea en FASE 1.

**Justificación:**
- El riesgo de romper releases con un stub es demasiado alto
- La ganancia (reducción de ~60 líneas) no justifica el riesgo
- Los workflows de release funcionan correctamente actualmente
- Se puede considerar en futuro cuando haya conectividad para validar

**Acción tomada:**
- ⏭️ Marcar tarea como SKIP
- ⏭️ Documentar razón en este archivo
- ⏭️ Continuar con Tarea 3 (coverage-check) que es de menor riesgo

---

## Para Futuro (Post-FASE 1)

Si en el futuro se decide implementar:

1. **Verificar que la composite action existe:**
   ```bash
   gh api repos/EduGoGroup/edugo-infrastructure/contents/.github/actions/docker-build-edugo/action.yml
   ```

2. **Si NO existe, crearla con inputs:**
   - `image-name` (required)
   - `registry` (default: ghcr.io)
   - `version` (required)
   - `platforms` (default: linux/amd64)
   - `push` (default: true)
   - `tag-strategy` (manual | semver | custom)
   - `custom-tags` (opcional)
   - `build-args` (opcional)

3. **Si existe, validar:**
   - Inputs disponibles
   - Outputs disponibles (tags, digest, etc.)
   - Estrategias de tags soportadas

4. **Testing:**
   - Probar en ambiente de desarrollo primero
   - Validar con release de prueba
   - NO probar directo en producción

---

## Alternativa: Composite Action Simplificada

Una alternativa más segura sería crear una composite action que solo maneje lo básico:
- Setup Buildx
- Login a registry
- Metadata básica

Y dejar el `build-push-action` inline en cada workflow con su configuración específica.

Esto reduciría ~15 líneas por workflow con menor riesgo.

---

## Estado

**Estado Actual:** ⏭️ SKIP
**Razón:** Complejidad alta + Riesgo en releases + Funcionando actualmente
**Requiere Resolución:** NO - Opcional
**Bloqueador:** Riesgo de romper releases con stub sin validación

---

## Beneficios Perdidos al SKIP

- ❌ No se elimina ~60 líneas de código duplicado
- ❌ Mantenimiento de Docker builds sigue descentralizado

## Beneficios de SKIP

- ✅ Releases siguen funcionando sin riesgo
- ✅ No introducir posibles errores en flujo crítico
- ✅ Tiempo ahorrado para otras tareas de menor riesgo

---

**Generado por:** Claude Code
**Fecha:** 2025-11-22
**Actualizado:** 2025-11-22
