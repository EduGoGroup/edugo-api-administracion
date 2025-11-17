# ECOSYSTEM CONTEXT - API Admin

## Posición en EduGo

**Rol:** Microservicio Backend - Gestión de jerarquía académica  
**Interacción:** Fuente de datos de escuelas, usuarios, unidades académicas

---

## Mapa de Ecosistema

```
┌──────────────────────────┐
│   Aplicación Web Admin   │
│  (Panel administrativo)  │
└────────────┬─────────────┘
             │ HTTP REST
             ▼
        ┌────────────────────────┐
        │   API ADMIN (8081)     │ ◄─── ESTE PROYECTO
        │  - Escuelas            │
        │  - Jerarquía           │
        │  - Docentes            │
        │  - Estudiantes         │
        └────────┬───────────────┘
                 │
        ┌────────┴──────────────────────┐
        │                               │
        ▼                               ▼
    ┌─────────────┐            ┌──────────────┐
    │ PostgreSQL  │            │  API Mobile  │
    │             │            │  (consultas) │
    │             │            └──────────────┘
    └─────────────┘

┌─ Referencia a Evaluaciones
│  └─ Docentes pueden ver quizzes
│     que sus estudiantes respondieron
```

---

## Interacciones con Otros Servicios

### 1. Integración con SHARED (v1.3.0+)

**Módulos utilizados:**

#### a) Logger
```go
logger.Info("Escuela creada", map[string]interface{}{
    "school_id": school.ID,
    "name": school.Name,
})
```

#### b) Database (PostgreSQL)
```go
db := database.GetDB()

// Queries relacionales
var school School
db.First(&school, id)

// Queries recursivas (CTE)
var units []AcademicUnit
db.Raw(hierarchyQuery, schoolID).Scan(&units)
```

#### c) Auth (JWT)
```go
middleware := auth.NewJWTValidator()
// Validar tokens de usuarios administrativos
```

---

### 2. Compartir PostgreSQL con API Mobile

**Tablas compartidas:**

```sql
-- Ambos APIs acceden (lectura/escritura)
├─ users                    # Usuarios del sistema
├─ schools                  # Escuelas
├─ academic_units           # Jerarquía
├─ teachers                 # Docentes
├─ students                 # Estudiantes
├─ memberships              # Relaciones user-unit
└─ enrollments              # Inscripciones

-- API Admin escribe, API Mobile lee
├─ schools
├─ academic_units
├─ teachers
├─ students
└─ enrollments
```

**Coordinación:**
```
API Admin:
├─ POST /api/v1/schools            (crea escuela)
└─ POST /api/v1/academic-units     (crea unidades)
       │
       ▼ (ambos APIs leen)
API Mobile:
└─ GET /api/v1/evaluaciones        (usa school_id del contexto)
```

---

### 3. Referencia a API Mobile (evaluaciones)

**API Mobile proporciona datos que API Admin consulta:**

```
API Admin puede hacer queries tipo:
- GET /api/v1/evaluaciones?school_id=1
  → Obtener evaluaciones de una escuela
  
- GET /api/v1/evaluaciones/:id/results
  → Ver resultados de una evaluación
```

**No hay integración de mensajería**, solo consultas HTTP si es necesario.

---

## Dependencias Directas

| Servicio | Versión | Tipo | Críticidad |
|----------|---------|------|-----------|
| SHARED | v1.3.0+ | Librería Go | CRÍTICA |
| PostgreSQL | 15+ | Base datos | CRÍTICA |
| API Mobile | Latest | Microservicio | MEDIA |

---

## Ciclo de Vida de Datos

```
Datos de Escuela:
1. API Admin crea (POST /schools)
2. PostgreSQL persiste
3. API Mobile accede vía shared DB
4. Usa school_id para contexto de evaluaciones

Datos de Jerarquía:
1. API Admin crea/edita unidades
2. PostgreSQL con CTEs para validar árbol
3. API Mobile referencia al asignar evaluaciones
4. Worker usa para contexto de procesamiento
```

---

## Checklist de Integración

- [ ] SHARED v1.3.0+ importado
- [ ] Conexión PostgreSQL funcionando
- [ ] Datos de schools creados
- [ ] Queries recursivas funcionan
- [ ] Middleware de auth validando
- [ ] Roles y permisos configurados
- [ ] Tests de integridad de árbol
