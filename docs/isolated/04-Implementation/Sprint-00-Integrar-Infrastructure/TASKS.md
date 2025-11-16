# TASKS Sprint-00: Integrar con infrastructure

## TASK-001: Actualizar go.mod

**Descripción:** Agregar dependencia de infrastructure y actualizar shared a v0.7.0

```bash
# 1. Agregar infrastructure
go get github.com/EduGoGroup/edugo-infrastructure/database@v0.2.0

# 2. Actualizar shared a v0.7.0
go get github.com/EduGoGroup/edugo-shared/auth@v0.7.0
go get github.com/EduGoGroup/edugo-shared/logger@v0.7.0
go get github.com/EduGoGroup/edugo-shared/config@v0.7.0
go get github.com/EduGoGroup/edugo-shared/bootstrap@v0.7.0
go get github.com/EduGoGroup/edugo-shared/lifecycle@v0.7.0
go get github.com/EduGoGroup/edugo-shared/database/postgres@v0.7.0

# 3. Limpiar
go mod tidy
```

**Validación:**
```bash
cat go.mod | grep infrastructure
# Debe mostrar: github.com/EduGoGroup/edugo-infrastructure/database v0.2.0

cat go.mod | grep shared
# Todas las versiones deben ser v0.7.0
```

---

## TASK-002: Comparar migraciones locales vs infrastructure

**Descripción:** Identificar diferencias entre migraciones locales y centralizadas

```bash
# Ver migraciones locales
cat scripts/postgresql/01_academic_hierarchy.sql

# Clonar infrastructure (si no está)
cd ..
git clone https://github.com/EduGoGroup/edugo-infrastructure

# Comparar con infrastructure
cat edugo-infrastructure/database/migrations/postgres/003_create_academic_units.up.sql

# Documentar diferencias
```

**Crear archivo:** `docs/MIGRATION_COMPARISON.md`

```markdown
# Comparación de Migraciones

## Tabla: academic_units

### En scripts/postgresql/01_academic_hierarchy.sql (LOCAL):
- Campo X con tipo Y
- Índice A

### En infrastructure/003_create_academic_units.up.sql (CENTRALIZADO):
- Campo X con tipo Z (DIFERENTE)
- Índice B (DIFERENTE)

## Decisión: Usar infrastructure (tiene la última palabra)

## Acciones:
- Eliminar scripts/postgresql/
- Usar infrastructure/database en go.mod
```

---

## TASK-003: Eliminar migraciones locales

**Descripción:** Eliminar scripts/postgresql/ (infrastructure es fuente de verdad)

```bash
# 1. Backup (por si acaso)
cp -r scripts/postgresql scripts/postgresql.backup

# 2. Eliminar
rm -rf scripts/postgresql/

# 3. Actualizar .gitignore
echo "scripts/postgresql.backup/" >> .gitignore
```

**Validación:**
```bash
ls scripts/
# NO debe mostrar postgresql/
```

---

## TASK-004: Actualizar README.md del proyecto

**Descripción:** Documentar que ahora usa infrastructure

Agregar sección en README.md:

```markdown
## Migraciones de Base de Datos

**IMPORTANTE:** Este proyecto usa migraciones centralizadas de `edugo-infrastructure`.

### Setup de Base de Datos

bash
# Opción 1: Usar infrastructure
cd /path/to/edugo-infrastructure
make dev-setup

# Opción 2: Ejecutar migraciones manualmente
cd /path/to/edugo-infrastructure/database
go run migrate.go up


### Tablas Usadas por este Proyecto

**Owner (crea y mantiene):**
- users
- schools  
- academic_units
- unit_membership

Ver: https://github.com/EduGoGroup/edugo-infrastructure/blob/dev/database/TABLE_OWNERSHIP.md
```

---

## TASK-005: Actualizar documentación de desarrollo

**Descripción:** Actualizar guías de desarrollo

En `docs/DEVELOPMENT.md` o similar, actualizar:

```markdown
## Prerequisitos

- Go 1.24+
- PostgreSQL 15+ (ejecutar migraciones de infrastructure)
- edugo-infrastructure clonado localmente

## Setup

1. Clonar infrastructure:
   bash
   git clone https://github.com/EduGoGroup/edugo-infrastructure
   cd edugo-infrastructure
   make dev-setup
   

2. Clonar este proyecto:
   bash
   git clone https://github.com/EduGoGroup/edugo-api-administracion
   cd edugo-api-administracion
   

3. Instalar dependencias:
   bash
   go mod download
   

4. Ejecutar API:
   bash
   go run cmd/api/main.go
   
```

---

## TASK-006: Actualizar CI/CD

**Descripción:** Actualizar GitHub Actions para usar infrastructure

En `.github/workflows/ci.yml`:

```yaml
jobs:
  test:
    runs-on: ubuntu-latest
    
    services:
      postgres:
        image: postgres:15-alpine
        env:
          POSTGRES_DB: edugo_test
          POSTGRES_USER: edugo
          POSTGRES_PASSWORD: test
        ports:
          - 5432:5432
    
    steps:
      - uses: actions/checkout@v4
      
      - uses: actions/setup-go@v4
        with:
          go-version: '1.24'
      
      # NUEVO: Clonar infrastructure para migraciones
      - name: Clone infrastructure
        run: |
          git clone https://github.com/EduGoGroup/edugo-infrastructure /tmp/infrastructure
      
      # NUEVO: Ejecutar migraciones
      - name: Run migrations
        env:
          DATABASE_URL: postgres://edugo:test@localhost:5432/edugo_test?sslmode=disable
        run: |
          cd /tmp/infrastructure/database
          go run migrate.go up
      
      - name: Run tests
        run: go test ./... -v
```

---

## ✅ Checklist de Completitud

- [ ] go.mod actualizado (infrastructure v0.2.0, shared v0.7.0)
- [ ] MIGRATION_COMPARISON.md creado
- [ ] scripts/postgresql/ eliminado
- [ ] README.md actualizado
- [ ] DEVELOPMENT.md actualizado
- [ ] CI/CD actualizado (.github/workflows/)
- [ ] Tests pasan con migraciones de infrastructure
- [ ] Documentación de integración completa
