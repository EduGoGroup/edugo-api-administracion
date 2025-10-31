# 🔄 Workflows de CI/CD - edugo-api-administracion

## 📋 Workflows Configurados

### 1️⃣ **release.yml** - Release Automático con Docker (TAGS)

**Trigger:** Solo cuando creas un tag `v*` (ej: `v1.0.0`)

**Ejecuta:**
- ✅ Verificación de formato
- ✅ Análisis estático (go vet)
- ✅ Tests con race detection
- ✅ Cobertura de código
- ✅ Build de binarios para producción
- ✅ **Build y push de imagen Docker a GHCR**
- ✅ Creación automática de GitHub Release con binarios

**Cuándo se ejecuta:**
```bash
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0  # ← AQUÍ se ejecuta automáticamente
```

**Docker Tags creados:**
- `ghcr.io/edugogroup/edugo-api-administracion:1.0.0` (versión específica)
- `ghcr.io/edugogroup/edugo-api-administracion:1.0` (major.minor)
- `ghcr.io/edugogroup/edugo-api-administracion:1` (major)
- `ghcr.io/edugogroup/edugo-api-administracion:latest`
- `ghcr.io/edugogroup/edugo-api-administracion:production`

**Duración estimada:** 5-7 minutos

---

### 2️⃣ **ci.yml** - Pipeline de Integración Continua

**Trigger:**
- ✅ Pull Requests a `main` o `develop`
- ✅ Push directo a `main` (red de seguridad)

**Ejecuta:**
- ✅ Verificación de formato
- ✅ Verificación de go.mod/go.sum
- ✅ Análisis estático (go vet)
- ✅ Tests con race detection
- ✅ Build verification
- ✅ Linter (opcional, no bloquea)
- ✅ Docker build test (sin push)

**Cuándo se ejecuta:**
```bash
# Cuando creas un PR
gh pr create --title "..." --body "..."  # ← AQUÍ se ejecuta

# O cuando alguien hace push directo a main (no recomendado)
git push origin main  # ← AQUÍ se ejecuta
```

**Duración estimada:** 3-4 minutos

---

### 3️⃣ **test.yml** - Tests con Cobertura (MANUAL/PR)

**Trigger:**
- ✅ Manual (workflow_dispatch desde GitHub UI)
- ✅ Pull Requests a `main` o `develop`

**Ejecuta:**
- ✅ Tests con cobertura detallada
- ✅ Generación de reporte HTML
- ✅ Upload de reportes a Codecov
- ✅ Artifacts con reportes de cobertura (30 días)

**Cuándo se ejecuta:**
```bash
# Manual desde GitHub UI:
# Actions → Tests with Coverage → Run workflow

# O automáticamente en PRs
gh pr create  # ← AQUÍ se ejecuta junto con ci.yml
```

**Duración estimada:** 2-3 minutos

---

### 4️⃣ **build-and-push.yml** - Build Manual On-Demand

**Trigger:**
- ✅ Solo Manual (workflow_dispatch)

**Ejecuta:**
- ✅ Tests completos
- ✅ Build de imagen Docker
- ✅ Push a GHCR con tags custom

**Parámetros:**
- `environment`: development | staging | production
- `push_latest`: ¿Tagear como latest?

**Cuándo se ejecuta:**
```bash
# Manual desde GitHub UI:
# Actions → Build and Push Docker Image → Run workflow
# Seleccionar: environment = "staging"
```

**Docker Tags creados (ejemplo: staging):**
- `ghcr.io/edugogroup/edugo-api-administracion:staging`
- `ghcr.io/edugogroup/edugo-api-administracion:staging-abc1234` (SHA)
- `ghcr.io/edugogroup/edugo-api-administracion:latest` (si push_latest=true)

**Duración estimada:** 4-5 minutos

---

## 🎯 Estrategia de CI/CD Optimizada

### **Flujo Normal de Desarrollo:**

```
┌─────────────────────────────────────────────────────────────┐
│  1. Desarrollo Local                                        │
│     - Hacer cambios en código                               │
│     - go test ./... (local)                                 │
│     - git commit                                            │
│     ✅ NO GASTA MINUTOS DE GITHUB                           │
└─────────────────────────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│  2. Crear Pull Request                                      │
│     - gh pr create                                          │
│     - CI automático (ci.yml + test.yml)                     │
│     - Revisar resultados en GitHub                          │
│     ✅ VALIDA ANTES DE MERGE                                │
└─────────────────────────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│  3. Merge a Main                                            │
│     - gh pr merge                                           │
│     - CI de seguridad (ci.yml) si se hace push directo     │
│     ✅ CÓDIGO VALIDADO EN MAIN                              │
└─────────────────────────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│  4. Deploy Manual (Opcional)                                │
│     - GitHub UI → build-and-push.yml                        │
│     - Seleccionar environment (dev/staging/prod)            │
│     ✅ DEPLOY ON-DEMAND                                     │
└─────────────────────────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│  5. Crear Release (IMPORTANTE)                              │
│     - git tag -a v1.0.0 -m "Release v1.0.0"                 │
│     - git push origin v1.0.0                                │
│     - Release automático (release.yml)                      │
│     - ✅ Docker images creadas AUTOMÁTICAMENTE              │
│     - ✅ GitHub Release con binarios                        │
│     ✅ RELEASE CON VALIDACIÓN COMPLETA                      │
└─────────────────────────────────────────────────────────────┘
```

---

## 🐳 Gestión de Imágenes Docker

### **Releases (Automático con Tags)**

Cuando creas un tag, **release.yml** construye automáticamente las imágenes Docker:

```bash
# 1. Crear tag
git tag -a v1.2.0 -m "Release v1.2.0"
git push origin v1.2.0

# 2. GitHub Actions construye AUTOMÁTICAMENTE estas imágenes:
# - ghcr.io/edugogroup/edugo-api-administracion:1.2.0
# - ghcr.io/edugogroup/edugo-api-administracion:1.2
# - ghcr.io/edugogroup/edugo-api-administracion:1
# - ghcr.io/edugogroup/edugo-api-administracion:latest
# - ghcr.io/edugogroup/edugo-api-administracion:production

# 3. Descargar la imagen:
docker pull ghcr.io/edugogroup/edugo-api-administracion:1.2.0
```

### **Deploys On-Demand (Manual)**

Para despliegues a ambientes específicos sin crear release:

```bash
# 1. GitHub UI → Actions → "Build and Push Docker Image"
# 2. Run workflow
# 3. Seleccionar: environment = "staging"
# 4. Opcional: push_latest = true

# Resultado:
# - ghcr.io/edugogroup/edugo-api-administracion:staging
# - ghcr.io/edugogroup/edugo-api-administracion:staging-abc1234
# - ghcr.io/edugogroup/edugo-api-administracion:latest (si activaste push_latest)
```

---

## 💰 Ahorro de Minutos de GitHub Actions

### **Antes (sin optimización):**
```
Push a main → 3 workflows × 5 min = 15 minutos
10 pushes al día = 150 minutos/día
Mes = 4,500 minutos (¡casi 100% del plan gratuito!)
```

### **Después (optimizado):**
```
Push a main → 1 workflow × 3 min = 3 minutos
PR → 2 workflows × 6 min = 12 minutos
Tag/Release → 1 workflow × 6 min = 6 minutos
Manual deploy → 1 workflow × 5 min = 5 minutos (solo cuando necesitas)

Mes típico:
- 5 PRs = 60 minutos
- 2 releases = 12 minutos
- 3 deploys manuales = 15 minutos
- 5 pushes directos = 15 minutos
Total = 102 minutos/mes (✅ Solo 4-5% del plan gratuito)
```

**Ahorro:** ~95% de minutos 🎉

---

## 🚀 Guía Rápida

### **Para desarrollo normal:**
```bash
# 1. Desarrollar localmente
vim internal/application/service/user_service.go

# 2. Probar localmente (NO usa GitHub)
go test ./...

# 3. Commit
git commit -m "feat: nueva funcionalidad de usuarios"

# 4. Push a tu rama
git push origin feature/nueva-funcionalidad

# 5. Crear PR (ejecuta CI automáticamente)
gh pr create --title "Nueva funcionalidad" --body "..."

# 6. Esperar aprobación y merge
```

### **Para crear una release con Docker:**
```bash
# 1. Probar todo localmente
go test ./...
go build ./...

# 2. Actualizar CHANGELOG.md (opcional pero recomendado)
vim CHANGELOG.md

# 3. Commit cambios finales
git add .
git commit -m "chore: preparar release v1.2.0"
git push origin main

# 4. Crear y push tag (ejecuta release.yml AUTOMÁTICAMENTE)
git tag -a v1.2.0 -m "Release v1.2.0"
git push origin v1.2.0

# 5. GitHub Actions automáticamente:
#    - Ejecuta tests
#    - Construye binarios
#    - Construye imagen Docker
#    - Publica a GHCR
#    - Crea GitHub Release

# 6. Descargar y usar la imagen:
docker pull ghcr.io/edugogroup/edugo-api-administracion:1.2.0
docker run -d -p 8081:8081 ghcr.io/edugogroup/edugo-api-administracion:1.2.0
```

### **Para deploy manual a staging/production:**
```bash
# 1. GitHub UI → Actions → "Build and Push Docker Image (On-Demand)"
# 2. Click "Run workflow"
# 3. Seleccionar:
#    - Branch: main
#    - Environment: staging (o production)
#    - Tag as latest: false (o true si quieres)
# 4. Click "Run workflow"

# 5. Esperar 4-5 minutos

# 6. Descargar la imagen:
docker pull ghcr.io/edugogroup/edugo-api-administracion:staging
```

---

## 🔑 Autenticación con GitHub Container Registry

Para descargar imágenes privadas:

```bash
# 1. Crear un Personal Access Token (PAT) en GitHub
# Settings → Developer settings → Personal access tokens → Tokens (classic)
# Permisos necesarios: read:packages

# 2. Login a GHCR
echo $GITHUB_TOKEN | docker login ghcr.io -u TU_USERNAME --password-stdin

# 3. Pull de la imagen
docker pull ghcr.io/edugogroup/edugo-api-administracion:latest

# 4. Run
docker run -d -p 8081:8081 \
  -e DB_HOST=postgres \
  -e DB_PORT=5432 \
  ghcr.io/edugogroup/edugo-api-administracion:latest
```

---

## 📊 Comparación de Configuraciones

| Escenario | Workflow | Triggers | Docker Build | Duración |
|-----------|----------|----------|--------------|----------|
| Pull Request | ci.yml + test.yml | PR creado | Test only | 6 min |
| Push a main | ci.yml | Push directo | Test only | 3 min |
| Release | release.yml | Tag v* | ✅ **Si + Push** | 6 min |
| Deploy manual | build-and-push.yml | Manual | ✅ **Si + Push** | 5 min |

---

## 🛡️ Branch Protection (Recomendado)

Para forzar el uso de PRs, configura protección de rama:

1. GitHub → Settings → Branches → Add rule
2. Branch name pattern: `main`
3. Configurar:
   - ✅ Require pull request before merging
   - ✅ Require status checks to pass before merging
   - ✅ Status checks: "Test and Build", "Test Coverage"
   - ✅ Require branches to be up to date before merging

Esto previene push directo a `main` y garantiza que todo pase por PR + CI.

---

## 🔍 Ver Estado de Workflows

```bash
# Ver últimos workflows ejecutados
gh run list --limit 10

# Ver detalles de un workflow específico
gh run view <run-id>

# Ver logs de un workflow
gh run view <run-id> --log

# Re-ejecutar un workflow fallido
gh run rerun <run-id>

# Ver workflows de release.yml
gh run list --workflow=release.yml
```

---

## 📝 Notas Importantes

### **¿Cuándo se construyen imágenes Docker?**

✅ **SÍ se construye y publica:**
- Cuando creas un tag (`git push origin v1.0.0`) → **release.yml**
- Cuando ejecutas manualmente build-and-push.yml → **build-and-push.yml**

❌ **NO se construye/publica:**
- En Pull Requests → solo test de build
- En push a main → solo test de build

### **Recomendaciones:**

1. **Para producción**: Usa tags (`v1.0.0`) → release automático
2. **Para staging**: Usa build manual on-demand
3. **Para development**: Usa build manual on-demand o local

---

## 📚 Recursos

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [GitHub Container Registry](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry)
- [Dockerfile](../../Dockerfile)

---

**Última actualización:** 2025-10-31
**Mantenedor:** Equipo EduGo
**Proyecto:** edugo-api-administracion
