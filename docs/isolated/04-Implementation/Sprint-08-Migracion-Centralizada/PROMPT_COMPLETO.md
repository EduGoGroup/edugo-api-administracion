# SPRINT ESPECIAL: MIGRACIÓN CENTRALIZADA A INFRASTRUCTURE

**Proyecto:** edugo-api-administracion  
**Ejecutor:** Claude Code (Web o Local - decide según capacidades)  
**Duración estimada:** 30-45 minutos  
**Prioridad:** CRÍTICA  

---

## 🎯 OBJETIVO

Migrar **edugo-api-administracion** para que use **100% infrastructure centralizado** en lugar de manejar migraciones locales.

### Estado Actual (INCORRECTO)

- ❌ **go.mod NO tiene** dependencia de `github.com/EduGoGroup/edugo-infrastructure`
- ❌ Tests usan comentario "Las migraciones vienen de infrastructure" pero **NO es cierto**
- ❌ No hay integración real con infrastructure v0.7.1

### Estado Deseado (CORRECTO)

- ✅ **go.mod incluye** `github.com/EduGoGroup/edugo-infrastructure/postgres@v0.7.1`
- ✅ Tests de integración **usan migraciones de infrastructure** vía `InitScripts`
- ✅ **CERO migraciones locales** en api-administracion

---

## 📋 PLAN DE EJECUCIÓN

### TASK 1: Agregar Dependencia de Infrastructure (5 min)

**Objetivo:** Agregar infrastructure v0.7.1 a go.mod

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-administracion

# Crear branch
git checkout dev
git pull origin dev
git checkout -b feature/sprint-08-migracion-centralizada

# Agregar dependencia
go get github.com/EduGoGroup/edugo-infrastructure/postgres@v0.7.1

# Limpiar
go mod tidy

# Validar
cat go.mod | grep "edugo-infrastructure"
# Debe mostrar: github.com/EduGoGroup/edugo-infrastructure/postgres v0.7.1
```

---

### TASK 2: Actualizar setup.go para Usar Migraciones de Infrastructure (15 min)

**Archivo:** `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-administracion/test/integration/setup.go`

**Cambios necesarios:**

```go
//go:build integration

package integration

import (
	"context"
	"database/sql"
	"testing"

	"github.com/EduGoGroup/edugo-shared/testing/containers"
	infraPostgres "github.com/EduGoGroup/edugo-infrastructure/postgres" // NUEVO IMPORT
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
)

// setupTestDB crea una instancia de PostgreSQL para tests usando shared/testing
func setupTestDB(t *testing.T) (*sql.DB, func()) {
	ctx := context.Background()

	// NUEVO: Obtener migraciones de infrastructure
	migrationScripts := infraPostgres.GetMigrationScripts() // Función que devuelve []string con paths

	// Configurar PostgreSQL CON scripts de infrastructure
	cfg := containers.NewConfig().
		WithPostgreSQL(&containers.PostgresConfig{
			Database:    "edugo_test",
			Username:    "edugo_user",
			Password:    "edugo_pass",
			InitScripts: migrationScripts, // NUEVO: Usar migraciones de infrastructure
		}).
		Build()

	manager, err := containers.GetManager(t, cfg)
	if err != nil {
		t.Fatalf("Failed to get manager: %v", err)
	}

	pg := manager.PostgreSQL()
	if pg == nil {
		t.Fatal("Failed to get PostgreSQL container")
	}

	connString, err := pg.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("Failed to get Postgres connection string: %v", err)
	}

	db, err := sql.Open("postgres", connString)
	if err != nil {
		t.Fatalf("Failed to connect to Postgres: %v", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		t.Fatalf("Failed to ping Postgres: %v", err)
	}

	cleanup := func() {
		db.Close()
	}

	return db, cleanup
}

// setupTestMongoDB crea una instancia de MongoDB para tests usando shared/testing
func setupTestMongoDB(t *testing.T) (*mongo.Database, func()) {
	ctx := context.Background()

	cfg := containers.NewConfig().
		WithMongoDB(&containers.MongoConfig{
			Database: "edugo_test",
			Username: "edugo_admin",
			Password: "edugo_pass",
		}).
		Build()

	manager, err := containers.GetManager(t, cfg)
	if err != nil {
		t.Fatalf("Failed to get manager: %v", err)
	}

	mongoDB := manager.MongoDB()
	if mongoDB == nil {
		t.Fatal("Failed to get MongoDB container")
	}

	db := mongoDB.Database()

	cleanup := func() {
		mongoDB.DropAllCollections(ctx)
	}

	return db, cleanup
}

// setupTestDBWithMigrations crea una BD con las migraciones aplicadas
// Ahora realmente usa las migraciones de infrastructure
func setupTestDBWithMigrations(t *testing.T) (*sql.DB, func()) {
	return setupTestDB(t)
}
```

**NOTA CRÍTICA:** Necesitas verificar si `edugo-infrastructure/postgres` exporta una función para obtener las migraciones. Si NO la tiene, crea una alternativa:

**Opción A (Si infrastructure tiene embed):**
```go
// En infrastructure/postgres/migrations.go debería existir:
//go:embed migrations/*.sql
var migrationsFS embed.FS

func GetMigrationScripts() []string {
    // Retorna paths a archivos .up.sql en orden
}
```

**Opción B (Si infrastructure NO tiene función helper):**
```go
// En test/integration/setup.go
import (
	_ "embed"
)

//go:embed ../../vendor/github.com/EduGoGroup/edugo-infrastructure/postgres/migrations/*.up.sql
var migrationsFS embed.FS

func getMigrationScripts() []string {
	// Leer archivos del embed y retornar paths en orden
	// O usar paths absolutos a vendor/
}
```

**Opción C (Más simple - usar paths relativos):**
```go
cfg := containers.NewConfig().
    WithPostgreSQL(&containers.PostgresConfig{
        Database:    "edugo_test",
        Username:    "edugo_user",
        Password:    "edugo_pass",
        InitScripts: []string{
            "../../vendor/github.com/EduGoGroup/edugo-infrastructure/postgres/migrations/001_initial.up.sql",
            "../../vendor/github.com/EduGoGroup/edugo-infrastructure/postgres/migrations/002_create_schools.up.sql",
            "../../vendor/github.com/EduGoGroup/edugo-infrastructure/postgres/migrations/003_create_academic_units.up.sql",
            "../../vendor/github.com/EduGoGroup/edugo-infrastructure/postgres/migrations/004_create_memberships.up.sql",
            // ... resto de migraciones hasta 011
        },
    }).
    Build()
```

**DECISIÓN:** Usa la **Opción C** por simplicidad. Listar manualmente las migraciones necesarias.

---

### TASK 3: Validar Que Infrastructure Tiene Las Migraciones (5 min)

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-infrastructure

# Verificar que existen las migraciones v0.7.1
ls -la postgres/migrations/*.up.sql

# Debe mostrar:
# 001_initial.up.sql
# 002_create_schools.up.sql (CON metadata JSONB)
# 003_create_academic_units.up.sql (CON parent_unit_id, metadata, función anti-ciclos, vista)
# 004_create_memberships.up.sql (CON metadata, roles extendidos)
# ... hasta 011

# Verificar contenido de 003 (debe tener jerarquía)
grep -q "parent_unit_id" postgres/migrations/003_create_academic_units.up.sql && echo "✅ Jerarquía presente"
grep -q "prevent_academic_unit_cycles" postgres/migrations/003_create_academic_units.up.sql && echo "✅ Función anti-ciclos presente"
grep -q "v_academic_unit_tree" postgres/migrations/003_create_academic_units.up.sql && echo "✅ Vista recursiva presente"
```

**Si las migraciones NO tienen jerarquía:**
- DETENER y reportar: "infrastructure v0.7.1 NO tiene las migraciones correctas"
- Verificar que estás usando el tag correcto

---

### TASK 4: Actualizar setup.go con Paths Correctos (10 min)

**Archivo:** `test/integration/setup.go`

**Implementación final:**

```go
//go:build integration

package integration

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/EduGoGroup/edugo-shared/testing/containers"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
)

// getMigrationScripts retorna paths a migraciones de infrastructure
func getMigrationScripts() []string {
	// Obtener path al módulo de infrastructure en vendor
	_, filename, _, _ := runtime.Caller(0)
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(filename)))
	
	infrastructurePath := filepath.Join(projectRoot, "vendor", "github.com", "EduGoGroup", "edugo-infrastructure", "postgres", "migrations")

	// Listar migraciones en orden (001 a 011, solo .up.sql)
	migrations := []string{
		"001_initial.up.sql",
		"002_create_schools.up.sql",
		"003_create_academic_units.up.sql",
		"004_create_memberships.up.sql",
		// Agregar resto si existen (005-011)
	}

	var fullPaths []string
	for _, migration := range migrations {
		fullPaths = append(fullPaths, filepath.Join(infrastructurePath, migration))
	}

	return fullPaths
}

// setupTestDB crea una instancia de PostgreSQL para tests usando shared/testing
func setupTestDB(t *testing.T) (*sql.DB, func()) {
	ctx := context.Background()

	// Obtener migraciones de infrastructure
	migrationScripts := getMigrationScripts()

	// Log para debug
	t.Logf("Using %d migration scripts from infrastructure", len(migrationScripts))
	for i, script := range migrationScripts {
		t.Logf("  [%d] %s", i+1, filepath.Base(script))
	}

	// Configurar PostgreSQL CON scripts de infrastructure
	cfg := containers.NewConfig().
		WithPostgreSQL(&containers.PostgresConfig{
			Database:    "edugo_test",
			Username:    "edugo_user",
			Password:    "edugo_pass",
			InitScripts: migrationScripts,
		}).
		Build()

	manager, err := containers.GetManager(t, cfg)
	if err != nil {
		t.Fatalf("Failed to get manager: %v", err)
	}

	pg := manager.PostgreSQL()
	if pg == nil {
		t.Fatal("Failed to get PostgreSQL container")
	}

	connString, err := pg.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("Failed to get Postgres connection string: %v", err)
	}

	db, err := sql.Open("postgres", connString)
	if err != nil {
		t.Fatalf("Failed to connect to Postgres: %v", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		t.Fatalf("Failed to ping Postgres: %v", err)
	}

	cleanup := func() {
		db.Close()
	}

	return db, cleanup
}

// setupTestMongoDB - SIN CAMBIOS
func setupTestMongoDB(t *testing.T) (*mongo.Database, func()) {
	ctx := context.Background()

	cfg := containers.NewConfig().
		WithMongoDB(&containers.MongoConfig{
			Database: "edugo_test",
			Username: "edugo_admin",
			Password: "edugo_pass",
		}).
		Build()

	manager, err := containers.GetManager(t, cfg)
	if err != nil {
		t.Fatalf("Failed to get manager: %v", err)
	}

	mongoDB := manager.MongoDB()
	if mongoDB == nil {
		t.Fatal("Failed to get MongoDB container")
	}

	db := mongoDB.Database()

	cleanup := func() {
		mongoDB.DropAllCollections(ctx)
	}

	return db, cleanup
}

// setupTestDBWithMigrations - Ahora realmente usa infrastructure
func setupTestDBWithMigrations(t *testing.T) (*sql.DB, func()) {
	return setupTestDB(t)
}
```

---

### TASK 5: Ejecutar Tests de Integración (10 min)

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-administracion

# Verificar que Docker está corriendo
docker ps

# Ejecutar tests de integración
go test -v -tags=integration ./test/integration/... -count=1

# Debe pasar TODOS los tests
# Si falla:
# 1. Revisar logs (paths de migraciones incorrectos?)
# 2. Verificar que vendor/ tiene infrastructure v0.7.1
# 3. Verificar que migraciones existen en vendor/
```

**Validaciones esperadas:**
- ✅ Container PostgreSQL levanta
- ✅ Migraciones se aplican correctamente
- ✅ Tablas schools, academic_units, unit_memberships existen
- ✅ Columnas nuevas existen (parent_unit_id, metadata, description)
- ✅ Tests de repositorios pasan

---

### TASK 6: Commit y Push (5 min)

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-administracion

git add go.mod go.sum test/integration/setup.go

git commit -m "feat: migrar a infrastructure centralizado v0.7.1

## Cambios

- Add: Dependencia github.com/EduGoGroup/edugo-infrastructure/postgres@v0.7.1
- Update: test/integration/setup.go usa migraciones de infrastructure
- Remove: Dependencia de migraciones locales

## Arquitectura

Ahora las migraciones están 100% centralizadas en edugo-infrastructure:
- ✅ Schema unificado entre api-mobile, api-admin, worker
- ✅ Una sola fuente de verdad
- ✅ Mantenimiento simplificado

## Testing

- ✅ Tests integración: PASSING
- ✅ Migraciones de infrastructure aplicadas correctamente
- ✅ Jerarquía académica disponible (parent_unit_id, vista recursiva)

BREAKING CHANGE: Ahora requiere infrastructure v0.7.1

Co-authored-by: infrastructure team <noreply@edugo.com>
Co-Authored-By: Claude Code <noreply@anthropic.com>"

git push origin feature/sprint-08-migracion-centralizada
```

---

### TASK 7: Crear Pull Request (5 min)

```bash
gh pr create \
  --title "feat: Migración a Infrastructure Centralizado v0.7.1" \
  --body "## Sprint Especial: Migración Centralizada

### 🎯 Objetivo

Eliminar dependencia de migraciones locales y usar **100% infrastructure centralizado**.

### ✅ Cambios

- **go.mod**: Agregado \`edugo-infrastructure/postgres@v0.7.1\`
- **test/integration/setup.go**: Actualizado para usar migraciones de infrastructure vía vendor/
- **Arquitectura**: Schema unificado entre todos los servicios

### 🏗️ Beneficios

- ✅ **Una sola fuente de verdad**: infrastructure tiene el schema completo
- ✅ **Sincronización automática**: Todos los servicios usan mismo schema
- ✅ **Mantenimiento simplificado**: Cambios en un solo lugar
- ✅ **Jerarquía académica disponible**: parent_unit_id, metadata, vista recursiva

### 🧪 Testing

- ✅ Tests integración: **PASSING**
- ✅ Migraciones aplicadas: **SUCCESS**
- ✅ Schema validado: **OK**

### 🔗 Dependencias

- Requiere: \`edugo-infrastructure@v0.7.1\`
- Compatible con: \`api-mobile v0.2.0\`

---
**Sprint:** Especial (Migración Centralizada)  
**Duración:** ~45 minutos  
**Prioridad:** CRÍTICA" \
  --base dev \
  --head feature/sprint-08-migracion-centralizada
```

---

### TASK 8: Monitorear CI/CD y Merge (10 min)

```bash
# Obtener número de PR
PR_NUMBER=$(gh pr view --json number -q '.number')

# Monitorear (max 5 minutos)
MAX_WAIT=300
START_TIME=$(date +%s)

while true; do
  ELAPSED=$(($(date +%s) - START_TIME))
  
  if [ $ELAPSED -gt $MAX_WAIT ]; then
    echo "❌ TIMEOUT: CI/CD no completó en 5 minutos"
    gh pr checks
    exit 1
  fi

  STATUS=$(gh pr checks --json state,conclusion | jq '[.[] | select(.state=="IN_PROGRESS")] | length')
  
  if [ "$STATUS" -eq 0 ]; then
    FAILED=$(gh pr checks --json conclusion | jq '[.[] | select(.conclusion=="FAILURE")] | length')
    
    if [ "$FAILED" -eq 0 ]; then
      echo "✅ CI/CD completado exitosamente"
      break
    else
      echo "❌ CI/CD falló"
      gh pr checks
      exit 1
    fi
  fi
  
  echo "⏳ CI/CD en progreso... ($ELAPSED/$MAX_WAIT seg)"
  sleep 60
done

# Merge si todo OK
gh pr merge --squash --delete-branch

echo "✅ Sprint Especial completado y mergeado"
```

---

## ✅ CHECKLIST COMPLETO

- [ ] Agregar dependencia infrastructure v0.7.1 a go.mod
- [ ] Actualizar test/integration/setup.go con getMigrationScripts()
- [ ] Validar que vendor/ tiene infrastructure con migraciones
- [ ] Ejecutar tests integración (deben pasar)
- [ ] Commit y push
- [ ] Crear PR
- [ ] Monitorear CI/CD (max 5 min)
- [ ] Merge a dev

---

## 🎯 RESULTADO ESPERADO

Al finalizar:
- ✅ go.mod tiene `github.com/EduGoGroup/edugo-infrastructure/postgres@v0.7.1`
- ✅ Tests usan migraciones de infrastructure (vendor/)
- ✅ CERO dependencia de migraciones locales
- ✅ Arquitectura centralizada
- ✅ PR mergeado a dev

---

## 📝 NOTAS IMPORTANTES

1. **Infrastructure v0.7.1** ya existe y tiene:
   - Migraciones 002, 003, 004 con jerarquía académica
   - Función anti-ciclos
   - Vista recursiva
   - Metadata JSONB

2. **NO necesitas modificar infrastructure**, solo consumirlo

3. **Vendor approach**: Usamos `vendor/` para acceder a migraciones. Go mod vendor las descargará automáticamente.

4. **Si tests fallan**: Verificar paths a migraciones en vendor/

---

## 🚀 ¿LISTO PARA EJECUTAR?

Este sprint es **simple y directo**:
1. Agregar dependencia
2. Actualizar setup.go
3. Validar con tests
4. Merge

**Tiempo estimado:** 30-45 minutos

¿Estás listo para iniciar en modo desatendido?
