# EXECUTION PLAN - API Admin

## Información del Proyecto

**Proyecto:** EduGo API Administración  
**Objetivo:** API REST para jerarquía académica  
**Duración:** 6 Sprints (12 semanas)  
**Equipo:** Backend engineers + DevOps  
**Repositorio:** https://github.com/EduGoGroup/edugo-api-admin

---

## Fase 1: Setup + Escuelas CRUD (Sprint 1)

### 1.1 Configuración Inicial
- [ ] Clonar repo
- [ ] go mod download && go mod tidy
- [ ] Conectar a PostgreSQL
- [ ] Crear base de datos edugo_admin
- [ ] Crear tabla schools

### 1.2 Endpoints de Escuelas
- [ ] POST /api/v1/schools
- [ ] GET /api/v1/schools
- [ ] GET /api/v1/schools/:id
- [ ] PUT /api/v1/schools/:id
- [ ] DELETE /api/v1/schools/:id

### 1.3 Validaciones
- [ ] Nombre único
- [ ] Campos requeridos
- [ ] Soft delete

### 1.4 Tests
- [ ] Tests unitarios
- [ ] Tests de integración
- [ ] Cobertura >= 80%

---

## Fase 2: Unidades Académicas (sin recursión) (Sprint 2)

### 2.1 Tabla academic_units
- [ ] Crear tabla con parent_id
- [ ] Crear índices (críticos)
- [ ] Relación con schools

### 2.2 Endpoints CRUD
- [ ] POST /api/v1/academic-units
- [ ] GET /api/v1/academic-units
- [ ] GET /api/v1/academic-units/:id
- [ ] PUT /api/v1/academic-units/:id
- [ ] DELETE /api/v1/academic-units/:id

### 2.3 Validaciones
- [ ] Tipo válido
- [ ] Parent existe
- [ ] School existe
- [ ] No crear ciclos (simple check)

---

## Fase 3: Queries Recursivas + Árbol (Sprint 3)

### 3.1 Queries con CTE
- [ ] Query obtener árbol completo
- [ ] Query obtener descendientes
- [ ] Query obtener ascendientes (breadcrumb)

### 3.2 Endpoints de Árbol
- [ ] GET /api/v1/schools/:id/hierarchy
- [ ] GET /api/v1/academic-units/:id/tree
- [ ] GET /api/v1/academic-units/:id/ancestors

### 3.3 Optimizaciones
- [ ] Índices de performance
- [ ] Caché si es necesario
- [ ] Límites de profundidad

### 3.4 Tests de Recursión
- [ ] Tests de árbol simple
- [ ] Tests de árbol profundo
- [ ] Tests de ciclos

---

## Fase 4: Docentes + Estudiantes CRUD (Sprint 4)

### 4.1 Tablas
- [ ] Teachers
- [ ] Students

### 4.2 Endpoints Docentes
- [ ] CRUD completo
- [ ] Asignar a unidad académica
- [ ] Listar por escuela

### 4.3 Endpoints Estudiantes
- [ ] CRUD completo
- [ ] Listar por escuela

---

## Fase 5: Membresías + Inscripciones (Sprint 5)

### 5.1 Tablas
- [ ] Memberships
- [ ] Enrollments

### 5.2 Endpoints
- [ ] POST /api/v1/memberships
- [ ] GET /api/v1/users/:id/memberships
- [ ] POST /api/v1/students/:id/enroll
- [ ] GET /api/v1/academic-units/:id/students

---

## Fase 6: Reportes + Optimizaciones (Sprint 6)

### 6.1 Reportes
- [ ] GET /api/v1/schools/:id/stats
- [ ] GET /api/v1/academic-units/:id/enrollment-report

### 6.2 Optimizaciones
- [ ] Benchmarks
- [ ] Índices adicionales si es necesario
- [ ] Caché de queries frecuentes

### 6.3 Documentación
- [ ] Swagger 100%
- [ ] Tests >= 80%
- [ ] Runbook para producción

---

## Post-Sprints

### Performance Testing
- [ ] Load test 1000 unidades académicas
- [ ] Load test queries recursivas
- [ ] Optimizar si es necesario

### Producción
- [ ] Docker image
- [ ] Deployment a staging
- [ ] Smoke tests
- [ ] Deploy a producción

---

**Próxima revisión:** Después de Sprint 1  
**Última actualización:** 15 de Noviembre, 2025
