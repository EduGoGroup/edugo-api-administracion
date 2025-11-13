# Troubleshooting - CI/CD Workflows

**Última actualización**: 9 de noviembre de 2025

---

## 🔍 Problemas Comunes

### 1. Tests Unitarios Fallan en CI pero Pasan Localmente

**Síntomas**:
```
❌ pr-to-dev.yml failed
✅ Local: make test-unit passes
```

**Causas Posibles**:

#### A. Dependencias no sincronizadas

**Solución**:
```bash
# Sincronizar go.mod y go.sum
go mod tidy
git add go.mod go.sum
git commit -m "chore: sync dependencies"
git push
```

#### B. Tests con race conditions

**Solución**:
```bash
# Ejecutar localmente con race detection
go test -race ./...

# Si falla, agregar t.Parallel() solo a tests independientes
func TestSomething(t *testing.T) {
    t.Parallel()  // Solo si el test no comparte estado
    // ...
}
```

#### C. Variables de entorno faltantes

**Solución**:
```go
// En el test, usar valores por defecto
dbHost := os.Getenv("DB_HOST")
if dbHost == "" {
    dbHost = "localhost"  // Default para tests
}
```

---

### 2. Tests de Integración Timeout

**Síntomas**:
```
❌ TestProgressFlow_UpsertProgress (62.04s) - FAIL
Error: context deadline exceeded
```

**Causas Posibles**:

#### A. RabbitMQ tarda en iniciar

**Solución**: Ya implementada en contenedores compartidos

**Verificar**:
```bash
# Localmente
make test-integration

# Debe completar en ~1-2 minutos
```

#### B. Docker no disponible en CI

**Solución**:
```yaml
# En el workflow, agregar:
- name: Setup Docker Buildx
  uses: docker/setup-buildx-action@v3

- name: Verificar Docker
  run: docker --version
```

---

### 3. Cobertura Por Debajo del Umbral

**Síntomas**:
```
❌ Coverage check failed
Current: 32.5%
Required: 33%
```

**Soluciones**:

#### A. Agregar tests para código nuevo

```bash
# Ver qué no está cubierto
make coverage-html
open coverage/coverage.html
```

#### B. Usar label temporal

```bash
# En GitHub PR, agregar label:
skip-coverage
```

#### C. Ajustar umbral (temporal)

```yaml
# En el workflow
env:
  COVERAGE_THRESHOLD: 30  # Reducir temporalmente
```

---

### 4. Manual Release Falla

**Síntomas**:
```
❌ manual-release.yml failed
Error: Tag v0.1.0 already exists
```

**Causas Posibles**:

#### A. Tag ya existe

**Solución**:
```bash
# Eliminar tag local y remoto
git tag -d v0.1.0
git push origin :refs/tags/v0.1.0

# Volver a ejecutar workflow
```

#### B. Versión inválida

**Solución**:
```
# Usar formato correcto
✅ Correcto: 0.1.0
❌ Incorrecto: v0.1.0 (sin 'v')
❌ Incorrecto: 0.1 (falta patch)
```

#### C. GitHub App secrets faltantes

**Solución**:
```
1. Verificar secrets en Settings → Secrets:
   - APP_ID
   - APP_PRIVATE_KEY

2. Si faltan, contactar admin del repo
```

---

### 5. Docker Build Falla

**Síntomas**:
```
❌ Build and push Docker image failed
Error: failed to solve: process "/bin/sh -c go build" did not complete successfully
```

**Causas Posibles**:

#### A. Dependencias privadas no accesibles

**Solución**:
```dockerfile
# En Dockerfile, agregar:
ARG GITHUB_TOKEN
RUN git config --global url."https://${GITHUB_TOKEN}@github.com/".insteadOf "https://github.com/"
```

#### B. Build context muy grande

**Solución**:
```bash
# Agregar a .dockerignore
coverage/
test/
*.md
.git/
```

---

### 6. Sync Main to Dev Falla

**Síntomas**:
```
❌ sync-main-to-dev.yml failed
Error: merge conflict
```

**Solución**:
```bash
# Resolver manualmente
git checkout dev
git pull origin dev
git merge main

# Resolver conflictos
git add .
git commit -m "chore: sync main to dev"
git push origin dev
```

---

### 7. PR Bloqueado por Checks

**Síntomas**:
```
⚠️ Some checks haven't completed yet
⏳ pr-to-dev.yml / Unit Tests - In progress
```

**Soluciones**:

#### A. Esperar a que termine

```
Tiempo esperado:
- pr-to-dev.yml: 2-3 minutos
- pr-to-main.yml: 3-4 minutos
```

#### B. Si tarda más de 10 minutos

```
1. Cancelar workflow en Actions
2. Hacer push vacío para re-trigger:
   git commit --allow-empty -m "chore: re-trigger CI"
   git push
```

---

### 8. Workflow No Se Ejecuta

**Síntomas**:
```
PR creado pero no hay checks ejecutándose
```

**Causas Posibles**:

#### A. Branch incorrecto

**Solución**:
```
pr-to-dev.yml solo se ejecuta en PRs a 'dev'
pr-to-main.yml solo se ejecuta en PRs a 'main'

Verificar que el PR apunte al branch correcto
```

#### B. Workflow deshabilitado

**Solución**:
```
1. Ir a Actions
2. Seleccionar el workflow
3. Verificar que no esté deshabilitado
4. Click "Enable workflow" si es necesario
```

---

### 9. Tests Pasan pero PR No Se Puede Mergear

**Síntomas**:
```
✅ All checks passed
⚠️ Merge blocked
```

**Causas Posibles**:

#### A. Branch protection rules

**Solución**:
```
1. Verificar Settings → Branches → Branch protection rules
2. Verificar que todos los required checks pasen
3. Verificar que tenga las aprobaciones necesarias
```

#### B. Branch desactualizado

**Solución**:
```bash
# Actualizar branch
git checkout feature/mi-feature
git pull origin dev
git push
```

---

### 10. Contenedores Residuales en Tests de Integración

**Síntomas**:
```
⚠️ Warning: 2 containers not cleaned up
- rabbitmq:3.12-alpine (abc123)
- postgres:15-alpine (def456)
```

**Solución**:
```bash
# Limpiar manualmente
docker ps -a | grep -E "postgres|mongo|rabbitmq" | awk '{print $1}' | xargs docker rm -f

# Verificar
docker ps -a
```

**Prevención**: Ya implementada en workflows con cleanup automático

---

## 🔧 Comandos Útiles

### Ejecutar Tests Localmente

```bash
# Tests unitarios
make test-unit

# Tests de integración
make test-integration

# Todos los tests
make test-all

# Con cobertura
make coverage-report

# Ver cobertura en navegador
make coverage-html
open coverage/coverage.html
```

### Verificar Workflows

```bash
# Validar sintaxis de workflow
act -l  # Requiere 'act' instalado

# O usar GitHub CLI
gh workflow list
gh workflow view pr-to-dev.yml
```

### Debug de Docker

```bash
# Ver logs de contenedor
docker logs <container_id>

# Inspeccionar contenedor
docker inspect <container_id>

# Ejecutar comando en contenedor
docker exec -it <container_id> /bin/sh
```

---

## 📞 Escalación

### Nivel 1: Auto-resolución
- Revisar este documento
- Verificar logs en GitHub Actions
- Ejecutar tests localmente

### Nivel 2: Equipo
- Preguntar en Slack #dev-help
- Revisar PRs similares recientes
- Consultar documentación en `docs/`

### Nivel 3: DevOps
- Crear issue en GitHub
- Tag: `ci/cd`, `bug`
- Incluir: logs, screenshots, pasos para reproducir

---

## 📚 Recursos Adicionales

- **[Workflows Index](WORKFLOWS_INDEX.md)** - Índice de todos los workflows
- **[CI/CD Strategy](CI_CD_STRATEGY.md)** - Estrategia general
- **[Testing Strategy](TESTING_STRATEGY.md)** - Estrategia de testing
- **[GitHub Actions Docs](https://docs.github.com/en/actions)** - Documentación oficial

---

**Última actualización**: 9 de noviembre de 2025

