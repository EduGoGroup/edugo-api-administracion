# PROJECT OVERVIEW - API Admin

## Información General

**Proyecto:** EduGo API Administración  
**Tipo:** API REST Microservicio  
**Puerto:** 8081  
**Lenguaje:** Go 1.21+  
**Marco de Trabajo:** Gin Framework  
**Especificación de Origen:** spec-03-api-administracion  
**Estado:** En Desarrollo (Sprint 1/6)

---

## Propósito del Proyecto

API REST especializada en la gestión de la **jerarquía académica** de EduGo. Proporciona endpoints para administrar escuelas, unidades académicas, docentes, estudiantes e inscripciones.

### Responsabilidades Principales
- Gestión de escuelas (multi-tenant)
- Jerarquía académica (facultades, departamentos, programas)
- Gestión de docentes y sus asignaciones
- Gestión de estudiantes y sus inscripciones
- Relaciones de membresía (users en academias)
- Queries recursivas y árbol de unidades
- Reportes administrativos

---

## Característica Principal: Jerarquía Académica

### Estructura de Árbol
```
SCHOOL (Escuela)
  ├─ FACULTY (Facultad)
  │  ├─ DEPARTMENT (Departamento)
  │  │  ├─ PROGRAM (Programa/Carrera)
  │  │  │  ├─ COURSE (Asignatura)
  │  │  │  │  ├─ CLASS (Clase)
  │  │  │  │  │  └─ STUDENTS (Estudiantes)
  │  │  │  │  └─ TEACHERS (Docentes)
  │  │  │  └─ STUDENTS (Inscritos en programa)
  │  │  └─ TEACHERS (Docentes de depto)
  │  └─ TEACHERS (Docentes de facultad)
  └─ STUDENTS (Estudiantes de escuela)

Nota: Algunos niveles pueden ser opcionales según configuración
```

### Ejemplo Concreto
```
Instituto EduGo
├─ Facultad de Ingeniería
│  ├─ Depto. Ingeniería de Sistemas
│  │  ├─ Carrera: Ingeniero en Sistemas
│  │  │  ├─ Asignatura: Programación I
│  │  │  │  ├─ Clase: T-01 (Turno mañana)
│  │  │  │  └─ Clase: T-02 (Turno tarde)
│  │  │  └─ Asignatura: Base de Datos
│  │  └─ Depto. Ingeniería Industrial
└─ Facultad de Ciencias
```

---

## Arquitectura del Proyecto

### Stack Tecnológico
```
┌─────────────────────────────────────────────────────┐
│ API Admin (Gin Framework - Go)                      │
├─────────────────────────────────────────────────────┤
│ Layers:                                             │
│ ├─ HTTP Handlers (routes + middleware)             │
│ ├─ Service Layer (lógica de negocio)               │
│ ├─ Repository Layer (acceso a datos + queries)     │
│ └─ Domain Models (estructuras de datos)             │
├─────────────────────────────────────────────────────┤
│ Características Especiales:                         │
│ ├─ Queries recursivas (CTEs, WITH RECURSIVE)      │
│ ├─ Multi-tenancy (por school_id)                  │
│ ├─ Árbol académico completo                        │
│ └─ Membresías y relaciones complejas               │
├─────────────────────────────────────────────────────┤
│ Dependencias Externas:                              │
│ ├─ PostgreSQL 15+ (con soporte recursivo)         │
│ ├─ shared v1.3.0+ (librerías compartidas)         │
│ └─ API Mobile (referencia a evaluaciones)         │
└─────────────────────────────────────────────────────┘
```

### Estructura de Carpetas
```
api-admin/
├── cmd/
│   └── api-admin/
│       └── main.go              # Punto de entrada
├── internal/
│   ├── handlers/                # HTTP handlers
│   │   ├── school_handler.go
│   │   ├── academic_unit_handler.go
│   │   ├── teacher_handler.go
│   │   └── student_handler.go
│   ├── services/                # Lógica de negocio
│   │   ├── hierarchy_service.go
│   │   └── membership_service.go
│   ├── repositories/            # Acceso a datos
│   │   ├── school_repository.go
│   │   ├── academic_unit_repository.go
│   │   └── queries/
│   │       └── hierarchy_queries.sql  # Queries recursivas
│   ├── models/                  # Estructuras de datos
│   ├── middleware/              # Middleware
│   └── config/                  # Configuración
├── migrations/                  # Migraciones GORM
├── docker/
│   └── Dockerfile
├── go.mod
├── go.sum
└── docker-compose.yml
```

---

## Entidades de Base de Datos

### PostgreSQL (Datos Relacionales)

```sql
-- 1. Escuelas
CREATE TABLE schools (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL UNIQUE,
  description TEXT,
  country VARCHAR(100),
  city VARCHAR(100),
  address VARCHAR(255),
  email VARCHAR(255),
  phone VARCHAR(20),
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 2. Unidades Académicas (árbol)
CREATE TABLE academic_units (
  id BIGSERIAL PRIMARY KEY,
  school_id BIGINT NOT NULL REFERENCES schools(id),
  parent_id BIGINT REFERENCES academic_units(id),
  type VARCHAR(50) NOT NULL, -- 'faculty', 'department', 'program', 'course'
  name VARCHAR(255) NOT NULL,
  code VARCHAR(50),
  description TEXT,
  level INT, -- 0=faculty, 1=department, 2=program, 3=course
  position INT,
  is_active BOOLEAN DEFAULT TRUE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  UNIQUE(school_id, parent_id, name)
);

-- 3. Docentes
CREATE TABLE teachers (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL, -- FK a tabla users (shared)
  school_id BIGINT NOT NULL REFERENCES schools(id),
  first_name VARCHAR(100),
  last_name VARCHAR(100),
  email VARCHAR(255),
  phone VARCHAR(20),
  specialization VARCHAR(255),
  is_active BOOLEAN DEFAULT TRUE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 4. Estudiantes
CREATE TABLE students (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL, -- FK a tabla users (shared)
  school_id BIGINT NOT NULL REFERENCES schools(id),
  first_name VARCHAR(100),
  last_name VARCHAR(100),
  email VARCHAR(255),
  phone VARCHAR(20),
  id_number VARCHAR(50) UNIQUE,
  admission_date DATE,
  is_active BOOLEAN DEFAULT TRUE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 5. Membresías (relación N:M entre users y academic units)
CREATE TABLE memberships (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL,
  academic_unit_id BIGINT NOT NULL REFERENCES academic_units(id),
  role VARCHAR(50) NOT NULL, -- 'teacher', 'student', 'coordinator'
  joined_date DATE DEFAULT NOW(),
  left_date DATE,
  is_active BOOLEAN DEFAULT TRUE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  UNIQUE(user_id, academic_unit_id, role)
);

-- 6. Inscripciones de Estudiantes
CREATE TABLE enrollments (
  id BIGSERIAL PRIMARY KEY,
  student_id BIGINT NOT NULL REFERENCES students(id),
  academic_unit_id BIGINT NOT NULL REFERENCES academic_units(id),
  enrollment_date DATE DEFAULT NOW(),
  completion_date DATE,
  status VARCHAR(50), -- 'active', 'completed', 'suspended'
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  UNIQUE(student_id, academic_unit_id)
);

-- Índices para performance
CREATE INDEX idx_academic_units_school_parent ON academic_units(school_id, parent_id);
CREATE INDEX idx_academic_units_type ON academic_units(school_id, type);
CREATE INDEX idx_memberships_user ON memberships(user_id, is_active);
CREATE INDEX idx_memberships_academic_unit ON memberships(academic_unit_id, is_active);
CREATE INDEX idx_enrollments_student ON enrollments(student_id);
CREATE INDEX idx_enrollments_unit ON enrollments(academic_unit_id);
```

---

## Queries Recursivas (CTEs)

### Obtener Árbol Completo de Unidad Académica

```sql
WITH RECURSIVE unit_tree AS (
  -- Base: la unidad solicitada
  SELECT
    id, school_id, parent_id, type, name, level, 0 as depth
  FROM academic_units
  WHERE id = $1

  UNION ALL

  -- Recursión: todos los descendientes
  SELECT
    au.id, au.school_id, au.parent_id, au.type, au.name, au.level, ut.depth + 1
  FROM academic_units au
  INNER JOIN unit_tree ut ON au.parent_id = ut.id
)
SELECT * FROM unit_tree
ORDER BY depth, name;
```

### Obtener Todos los Estudiantes de una Facultad (incluyendo sub-unidades)

```sql
WITH RECURSIVE faculty_units AS (
  -- Base: la facultad
  SELECT id, parent_id FROM academic_units
  WHERE id = $1 AND type = 'faculty'
  
  UNION ALL
  
  -- Recursión: todos los descendientes
  SELECT au.id, au.parent_id FROM academic_units au
  INNER JOIN faculty_units fu ON au.parent_id = fu.id
)
SELECT DISTINCT s.*
FROM students s
INNER JOIN enrollments e ON s.id = e.student_id
INNER JOIN faculty_units fu ON e.academic_unit_id = fu.id
WHERE s.school_id = $2
ORDER BY s.last_name, s.first_name;
```

### Obtener Ruta Completa de una Unidad (breadcrumb)

```sql
WITH RECURSIVE unit_path AS (
  -- Base: la unidad solicitada
  SELECT id, parent_id, name, 1 as level
  FROM academic_units
  WHERE id = $1
  
  UNION ALL
  
  -- Recursión: hacia arriba hasta la raíz
  SELECT au.id, au.parent_id, au.name, up.level + 1
  FROM academic_units au
  INNER JOIN unit_path up ON au.id = up.parent_id
)
SELECT * FROM unit_path
ORDER BY level DESC;
```

---

## Responsabilidades por Módulo

### 1. Escuelas
- **Crear escuela** → POST /api/v1/schools
- **Listar escuelas** → GET /api/v1/schools
- **Obtener detalle** → GET /api/v1/schools/:id
- **Editar escuela** → PUT /api/v1/schools/:id
- **Eliminar escuela** → DELETE /api/v1/schools/:id

### 2. Unidades Académicas
- **Crear unidad** → POST /api/v1/schools/:id/academic-units
- **Listar árbol** → GET /api/v1/schools/:id/hierarchy
- **Obtener unidad** → GET /api/v1/academic-units/:id
- **Editar unidad** → PUT /api/v1/academic-units/:id
- **Eliminar unidad** → DELETE /api/v1/academic-units/:id
- **Reordenar** → POST /api/v1/academic-units/:id/reorder

### 3. Docentes
- **Crear docente** → POST /api/v1/teachers
- **Listar docentes** → GET /api/v1/schools/:id/teachers
- **Obtener detalle** → GET /api/v1/teachers/:id
- **Editar docente** → PUT /api/v1/teachers/:id
- **Asignar a unidad** → POST /api/v1/teachers/:id/assign
- **Eliminar docente** → DELETE /api/v1/teachers/:id

### 4. Estudiantes
- **Crear estudiante** → POST /api/v1/students
- **Listar estudiantes** → GET /api/v1/schools/:id/students
- **Obtener detalle** → GET /api/v1/students/:id
- **Editar estudiante** → PUT /api/v1/students/:id
- **Inscribir en programa** → POST /api/v1/students/:id/enroll
- **Eliminar estudiante** → DELETE /api/v1/students/:id

### 5. Membresías
- **Crear membresía** → POST /api/v1/memberships
- **Listar membresías** → GET /api/v1/users/:id/memberships
- **Terminar membresía** → DELETE /api/v1/memberships/:id

---

## Flujos Principales

### Flujo 1: Crear Jerarquía Académica Completa
```
1. Crear Escuela
   └─ POST /api/v1/schools

2. Crear Facultades
   └─ POST /api/v1/schools/:id/academic-units (type: "faculty")

3. Crear Departamentos
   └─ POST /api/v1/schools/:id/academic-units (parent: faculty_id, type: "department")

4. Crear Programas
   └─ POST /api/v1/schools/:id/academic-units (parent: department_id, type: "program")

5. Crear Asignaturas
   └─ POST /api/v1/schools/:id/academic-units (parent: program_id, type: "course")

6. Asignar Docentes
   └─ POST /api/v1/teachers/:id/assign (academic_unit_id)

7. Inscribir Estudiantes
   └─ POST /api/v1/students/:id/enroll (academic_unit_id)
```

---

## Configuración Requerida

### Variables de Entorno
```bash
# Base de datos PostgreSQL
DB_HOST=localhost
DB_PORT=5432
DB_USER=edugo_user
DB_PASSWORD=secure_password
DB_NAME=edugo_admin

# API Configuration
API_PORT=8081
API_ENV=development
API_TIMEOUT=30s

# Authentication
JWT_SECRET=your_secret_key
JWT_EXPIRY=24h

# Shared
SHARED_LOG_LEVEL=info
SHARED_CONTEXT_TIMEOUT=30s
```

---

## Compilación y Despliegue

### Compilación Local
```bash
go mod download
go mod tidy
go build -o api-admin ./cmd/api-admin
./api-admin
```

### Docker
```bash
docker build -t edugo/api-admin:latest -f docker/Dockerfile .
docker run -p 8081:8081 \
  -e DB_HOST=postgres \
  edugo/api-admin:latest
```

---

## Testing

### Pruebas Unitarias
```bash
go test ./...
```

### Pruebas de Queries Recursivas
```bash
# Tests especiales para CTEs
go test -v ./internal/repositories -run TestHierarchy
```

---

## Sprint Planning (6 Sprints)

| Sprint | Funcionalidad | Duración |
|--------|---------------|----------|
| 1 | Setup + Escuelas CRUD | 2 semanas |
| 2 | Unidades Académicas (sin recursión) | 2 semanas |
| 3 | Queries recursivas + árbol | 2 semanas |
| 4 | Docentes + Estudiantes CRUD | 2 semanas |
| 5 | Membresías + Inscripciones | 2 semanas |
| 6 | Reportes + Optimizaciones | 2 semanas |

---

## Contacto y Referencias

- **Repositorio GitHub:** https://github.com/EduGoGroup/edugo-api-admin
- **Especificación Completa:** docs/ESTADO_PROYECTO.md (repo análisis)
- **Documentación Técnica:** Este directorio (01-Context/)
