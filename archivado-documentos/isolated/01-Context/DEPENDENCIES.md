# DEPENDENCIES - API Admin

## Matriz de Dependencias

```
┌──────────────────────────────────────────────────┐
│                   API ADMIN                      │
├──────────────────────────────────────────────────┤
│ Dependencias Críticas                            │
│ ├─ SHARED v1.3.0+                               │
│ └─ PostgreSQL 15+ (con soporte recursivo)        │
│                                                  │
│ Dependencias Opcionales                          │
│ └─ API Mobile (para lectura de evaluaciones)     │
│                                                  │
│ Dependencias de Desarrollo/Ops                   │
│ ├─ Docker                                        │
│ └─ Go 1.21+                                      │
└──────────────────────────────────────────────────┘
```

---

## Dependencias Críticas

### 1. SHARED v1.3.0+

**Módulos requeridos:**

#### a) Logger Module
```go
import "github.com/EduGoGroup/edugo-shared/logger"

logger.Info("Escuela creada", map[string]interface{}{
    "school_id": school.ID,
})
```

#### b) Database Module (PostgreSQL)
```go
import "github.com/EduGoGroup/edugo-shared/database"

db := database.GetDB()

// Simple CRUD
var school School
db.First(&school, id)

// CTEs para queries recursivas
type HierarchyRow struct {
    ID    int64
    Name  string
    Depth int
}

var hierarchy []HierarchyRow
db.Raw(`WITH RECURSIVE tree AS (...)`, schoolID).Scan(&hierarchy)
```

#### c) Auth Module
```go
import "github.com/EduGoGroup/edugo-shared/auth"

middleware := auth.NewJWTValidator()
// Solo admin puede crear escuelas
```

---

### 2. PostgreSQL 15+

**Requisitos especiales:**
- Common Table Expressions (CTEs) con RECURSIVE
- Window functions
- JSON operators
- Full text search (opcional)

**Creación de BD:**
```bash
createdb -U postgres edugo_admin
```

**Habilitar extensiones:**
```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
```

**Schema básico:**
```sql
-- Escuelas
CREATE TABLE schools (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL UNIQUE,
  created_at TIMESTAMP DEFAULT NOW()
);

-- Unidades Académicas (árbol)
CREATE TABLE academic_units (
  id BIGSERIAL PRIMARY KEY,
  school_id BIGINT REFERENCES schools(id),
  parent_id BIGINT REFERENCES academic_units(id),
  type VARCHAR(50),
  name VARCHAR(255),
  created_at TIMESTAMP DEFAULT NOW()
);

-- Índices CRÍTICOS para performance de CTEs
CREATE INDEX idx_au_parent ON academic_units(parent_id);
CREATE INDEX idx_au_school_parent ON academic_units(school_id, parent_id);
```

**Query Recursiva Ejemplo:**
```sql
WITH RECURSIVE unit_tree AS (
  SELECT id, school_id, parent_id, name, 0 as depth
  FROM academic_units
  WHERE id = $1
  
  UNION ALL
  
  SELECT au.id, au.school_id, au.parent_id, au.name, ut.depth + 1
  FROM academic_units au
  INNER JOIN unit_tree ut ON au.parent_id = ut.id
)
SELECT * FROM unit_tree
ORDER BY depth;
```

---

### 3. Tablas Compartidas con API Mobile

```sql
-- Usuarios (creados por otro servicio)
CREATE TABLE users (
  id BIGSERIAL PRIMARY KEY,
  email VARCHAR(255) UNIQUE,
  first_name VARCHAR(100),
  last_name VARCHAR(100),
  school_id BIGINT REFERENCES schools(id),
  created_at TIMESTAMP DEFAULT NOW()
);

-- Docentes
CREATE TABLE teachers (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT UNIQUE REFERENCES users(id),
  school_id BIGINT REFERENCES schools(id),
  created_at TIMESTAMP DEFAULT NOW()
);

-- Estudiantes
CREATE TABLE students (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT UNIQUE REFERENCES users(id),
  school_id BIGINT REFERENCES schools(id),
  created_at TIMESTAMP DEFAULT NOW()
);

-- Membresías (usuarios en unidades académicas)
CREATE TABLE memberships (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT REFERENCES users(id),
  academic_unit_id BIGINT REFERENCES academic_units(id),
  role VARCHAR(50),
  UNIQUE(user_id, academic_unit_id, role)
);

-- Inscripciones (estudiantes en programas)
CREATE TABLE enrollments (
  id BIGSERIAL PRIMARY KEY,
  student_id BIGINT REFERENCES students(id),
  academic_unit_id BIGINT REFERENCES academic_units(id),
  UNIQUE(student_id, academic_unit_id)
);
```

---

## Dependencias Go principales

```bash
go get github.com/EduGoGroup/edugo-shared@v1.3.0
go get github.com/gin-gonic/gin@latest
go get gorm.io/gorm@latest
go get gorm.io/driver/postgres@latest
go get github.com/spf13/viper@latest
```

---

## Configuración Requerida

```bash
# PostgreSQL
DB_HOST=localhost
DB_PORT=5432
DB_USER=edugo_user
DB_PASSWORD=password
DB_NAME=edugo_admin

# API
API_PORT=8081
LOG_LEVEL=info

# Auth
JWT_SECRET=secret_key
```

---

## Checklist de Instalación

```markdown
- [ ] PostgreSQL 15+ instalado
- [ ] Base de datos "edugo_admin" creada
- [ ] Extensiones habilitadas (uuid, pgcrypto)
- [ ] Tablas creadas (schools, academic_units, etc)
- [ ] Índices creados (críticos para CTEs)
- [ ] Go 1.21+ instalado
- [ ] go mod download ejecutado
- [ ] SHARED v1.3.0+ en go.mod
- [ ] go build compila sin errores
- [ ] Tests pasan (go test ./...)
```
