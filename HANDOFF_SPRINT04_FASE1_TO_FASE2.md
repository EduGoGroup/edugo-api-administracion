# Handoff: Sprint-04 Fase 1 (Web) → Fase 2 (Local)

**Sprint:** Sprint-04 - Services/API
**Ejecutor Fase 1:** Claude Code Web
**Fecha:** 2025-11-18
**Branch:** `claude/sprint-04-services-api-01HWh2zu7zjfyg6rWqNcsqeq`

---

## ✅ COMPLETADO EN FASE 1 (Claude Web)

### 1. DTOs Implementados

**Archivos creados:**
- `internal/infrastructure/http/dto/school_dto.go`
- `internal/infrastructure/http/dto/academic_unit_dto.go`
- `internal/infrastructure/http/dto/common_dto.go`

**DTOs disponibles:**
- `CreateSchoolRequest`, `UpdateSchoolRequest`, `SchoolResponse`
- `CreateUnitRequest`, `UpdateUnitRequest`, `UnitResponse`
- `UnitTreeNode` (para árbol jerárquico con ltree)
- `ErrorResponse`, `SuccessResponse`, `PaginationMeta`

**Validaciones:** Usando `binding` tags de Gin/validator

---

### 2. Application Service - HierarchyService

**Archivo:** `internal/application/service/hierarchy_service.go`

**Métodos implementados:**
- `CreateUnit()` - Crea unidad con validaciones (escuela existe, código único, etc.)
- `GetUnitTree()` - Obtiene árbol usando ltree (Sprint-03!)
- `MoveUnit()` - Mueve unidad usando MoveSubtree ltree
- `ValidateNoCircularReference()` - Previene ciclos usando FindDescendants ltree

**Tests unitarios:** `hierarchy_service_test.go` con mocks

**Aprovecha Sprint-03:**
- ✅ FindDescendants para obtener árbol completo
- ✅ MoveSubtree para reorganizar jerarquías
- ✅ FindBySchoolIDAndDepth para filtrado por nivel

---

### 3. HTTP Handlers (Ya existían)

**Archivos verificados:**
- `internal/infrastructure/http/handler/school_handler.go`
- `internal/infrastructure/http/handler/academic_unit_handler.go`

**Endpoints disponibles:**

**Schools (6):**
- POST   /api/v1/schools
- GET    /api/v1/schools
- GET    /api/v1/schools/:id
- GET    /api/v1/schools/code/:code
- PUT    /api/v1/schools/:id
- DELETE /api/v1/schools/:id

**Academic Units (9):**
- POST   /api/v1/schools/:schoolId/units
- GET    /api/v1/schools/:schoolId/units
- GET    /api/v1/schools/:schoolId/units/tree (árbol completo con ltree!)
- GET    /api/v1/schools/:schoolId/units/by-type
- GET    /api/v1/units/:id
- PUT    /api/v1/units/:id (incluye mover unidad)
- DELETE /api/v1/units/:id
- POST   /api/v1/units/:id/restore
- GET    /api/v1/units/:id/hierarchy-path (usa ltree!)

**Características:**
- Validación de DTOs con Gin binding
- Manejo de errores con códigos HTTP apropiados
- Conversión entity ↔ DTO
- Documentación Swagger con anotaciones

---

### 4. Router Configuration

**Archivo:** `internal/infrastructure/http/router/router.go`

**Configurado:**
- Gin router con middleware de recovery
- CORS básico
- Health check en /health
- Agrupación de rutas en /api/v1
- Inyección de dependencias (repositorios, servicios)

**Nota:** El archivo `cmd/main.go` ya existía y configura las rutas de forma similar usando el container pattern.

---

### 5. Main Entry Point

**Archivo:** `cmd/main.go` (YA EXISTÍA)

**Funcionalidad:**
- Conexión a PostgreSQL usando shared/bootstrap
- Inicialización de repositorios con container
- Configuración de router Gin
- Servidor HTTP con graceful shutdown
- Swagger en /swagger/*any
- Middleware de autenticación (pendiente implementar JWT)

---

### 6. Tests con STUBS

**Archivo:** `test/integration/http_api_test.go`

**Tests estructurados (todos con t.Skip):**
1. `TestSchoolAPI_CreateAndGet`
2. `TestUnitAPI_CreateTree`
3. `TestUnitAPI_MoveSubtree`
4. `TestUnitAPI_ListByDepth`
5. `TestAPI_ErrorHandling`
6. `TestUnitAPI_GetHierarchyPath`
7. `TestSchoolAPI_ListAll`
8. `TestSchoolAPI_UpdateAndDelete`
9. `TestUnitAPI_RestoreDeleted`

**Cada test tiene:**
- ⚠️ `t.Skip("STUB_FASE2: Requiere servidor HTTP")`
- Comentarios `TODO_FASE2` con pasos detallados

---

## ⏸️ PENDIENTE PARA FASE 2 (Claude Local)

### 1. Ejecutar Servidor HTTP ⚠️ CRÍTICO

**Razón:** Requiere levantar Gin server en local

**Tareas Fase 2:**
1. Verificar configuración en `.env` o variables de entorno
2. Asegurar que PostgreSQL está corriendo (docker-compose)
3. Ejecutar migraciones (incluye 013 de Sprint-03)
4. Ejecutar `go run cmd/main.go`
5. Verificar que servidor levanta en puerto 8081
6. Probar health check: `curl http://localhost:8081/health`

---

### 2. Descomentar y Ejecutar Tests E2E ⚠️ CRÍTICO

**Archivo:** `test/integration/http_api_test.go`

**Para cada test:**
1. Quitar `t.Skip()`
2. Descomentar código
3. Implementar helper para levantar servidor de test
4. Ejecutar requests HTTP (usar httptest o testcontainers)

**Ejemplo de helper necesario:**

```go
func setupTestServer(t *testing.T) (*gin.Engine, *sql.DB, func()) {
    // Cargar config
    cfg, _ := config.Load()

    // Inicializar recursos
    ctx := context.Background()
    resources, cleanup, _ := bootstrap.Initialize(ctx, cfg)

    // Crear container
    c := container.NewContainer(resources.PostgreSQL, resources.Logger)

    // Configurar router (similar a main.go)
    r := gin.Default()
    // ... configurar rutas ...

    return r, resources.PostgreSQL, cleanup
}
```

---

### 3. Validaciones Específicas

**Test:** `TestSchoolAPI_CreateAndGet`
- Crear escuela via POST
- Verificar response 201 con SchoolResponse válido
- GET por ID debe retornar misma escuela
- Validar que timestamps se generan

**Test:** `TestUnitAPI_CreateTree`
- Crear grado (raíz)
- Crear 2 secciones bajo el grado
- Crear club bajo sección
- GET /schools/:schoolId/units/tree del grado
- Verificar que árbol tiene estructura correcta
- **Validar que usa ltree** (verificar que hijos están ordenados por path)

**Test:** `TestUnitAPI_MoveSubtree`
- Crear Grade1 -> Section -> Club
- Crear Grade2 (vacío)
- PUT /units/:section_id con parent_unit_id = Grade2
- Verificar que Section se movió
- Verificar que Club sigue siendo hijo de Section
- GET /schools/:schoolId/units/tree debe mostrar Section y Club bajo Grade2

**Test:** `TestUnitAPI_GetHierarchyPath`
- Crear jerarquía: School -> Grade -> Section -> Club
- GET /units/:club_id/hierarchy-path
- Verificar que retorna el path completo desde la raíz
- Validar que el orden es correcto (de raíz a hoja)

**Test:** `TestAPI_ErrorHandling`
- POST con JSON inválido → 400
- POST con field faltante → 400 con detalles
- GET con UUID inválido → 400
- GET con ID inexistente → 404
- POST con código duplicado → 400 o 409
- PUT para crear ciclo → 400 con mensaje claro

---

### 4. Tests Manuales con Postman/curl

**Collection de Postman recomendada:**

```bash
# 1. Health check
curl http://localhost:8081/health

# 2. Crear escuela
curl -X POST http://localhost:8081/v1/schools \
  -H "Content-Type: application/json" \
  -d '{"name": "Test School", "code": "TS001", "address": "123 Main St"}'

# 3. Crear grado
curl -X POST http://localhost:8081/v1/schools/SCHOOL_ID/units \
  -H "Content-Type: application/json" \
  -d '{"type": "grade", "display_name": "Grade 1", "code": "G1"}'

# 4. Obtener árbol
curl http://localhost:8081/v1/schools/SCHOOL_ID/units/tree

# 5. Obtener path jerárquico
curl http://localhost:8081/v1/units/UNIT_ID/hierarchy-path
```

---

### 5. Validar Integración con ltree (Sprint-03)

**Endpoints que DEBEN usar ltree:**

| Endpoint | Método ltree usado | Validación |
|----------|-------------------|------------|
| GET /schools/:id/units/tree | FindDescendants | Árbol completo en 1 query |
| GET /units/:id/hierarchy-path | GetHierarchyPath | Obtiene path usando ltree |
| PUT /units/:id (move) | MoveSubtree | Actualiza paths en cascada |

**Cómo validar:**
1. Crear jerarquía de 100+ unidades
2. Medir tiempo de GET /schools/:id/units/tree
3. Verificar que es rápido (< 100ms)
4. Confirmar en logs de PostgreSQL que usa índice GIST
5. Revisar plan de query: `EXPLAIN ANALYZE SELECT ... WHERE path <@ ...`

---

### 6. Swagger UI (Opcional)

El main.go ya tiene configurado Swagger:

```bash
# Acceder a la documentación
http://localhost:8081/swagger/index.html
```

Si necesitas regenerar los docs:

```bash
go get -u github.com/swaggo/swag/cmd/swag
swag init -g cmd/main.go -o docs
```

---

## 📊 COBERTURA ESPERADA POST-FASE 2

### Código
- **HierarchyService:** >= 80% cobertura
- **Handlers:** Ya implementados y probados
- **Router:** 100% (es simple)

### Funcionalidad
- ✅ CRUD completo de escuelas
- ✅ CRUD completo de unidades
- ✅ Árbol jerárquico con ltree
- ✅ Mover unidades con validación de ciclos
- ✅ Path jerárquico con ltree
- ✅ Manejo de errores HTTP

---

## 🚀 COMANDOS PARA FASE 2

```bash
# Checkout
git checkout claude/sprint-04-services-api-01HWh2zu7zjfyg6rWqNcsqeq
git pull origin claude/sprint-04-services-api-01HWh2zu7zjfyg6rWqNcsqeq

# Levantar PostgreSQL (si no está corriendo)
docker-compose up -d postgres

# Ejecutar migraciones (incluye 013 de Sprint-03)
migrate -path migrations -database "postgresql://edugo_user:edugo_pass@localhost:5432/edugo_admin?sslmode=disable" up

# Levantar servidor
go run cmd/main.go

# En otra terminal, ejecutar tests E2E
go test -tags=integration ./test/integration/... -v

# Tests manuales
curl http://localhost:8081/health
curl http://localhost:8081/swagger/index.html
```

---

## 📝 ESTRUCTURA DE ARCHIVOS CREADOS/MODIFICADOS

```
internal/
├── application/
│   └── service/
│       ├── hierarchy_service.go          ← NUEVO (Fase 1)
│       └── hierarchy_service_test.go     ← NUEVO (Fase 1)
│
├── infrastructure/
│   └── http/
│       ├── dto/
│       │   ├── academic_unit_dto.go      ← NUEVO (Fase 1)
│       │   ├── common_dto.go             ← NUEVO (Fase 1)
│       │   └── school_dto.go             ← NUEVO (Fase 1)
│       │
│       ├── handler/
│       │   ├── academic_unit_handler.go  ← YA EXISTÍA
│       │   └── school_handler.go         ← YA EXISTÍA
│       │
│       └── router/
│           └── router.go                 ← NUEVO (Fase 1)
│
├── cmd/
│   └── main.go                           ← YA EXISTÍA
│
└── test/
    └── integration/
        └── http_api_test.go              ← NUEVO (Fase 1 - STUBS)
```

---

## 🔗 INTEGRACIÓN CON SPRINT-03

Este Sprint-04 **aprovecha completamente** el trabajo de Sprint-03 (ltree):

1. **GetUnitTree** usa `FindDescendants(ctx, unitID)` → Query ltree: `WHERE path <@ root_path`
2. **GetHierarchyPath** usa `GetHierarchyPath(ctx, unitID)` → Query ltree: `WHERE root_path @> path`
3. **MoveUnit** usa `MoveSubtree(ctx, unitID, newParentID)` → Update ltree en cascada

**Sin ltree, estos endpoints serían muy lentos** (N+1 queries o múltiples JOINs recursivos).

---

## ⚠️ NOTAS IMPORTANTES

1. **El main.go ya existía** - Usa el pattern de container y bootstrap de shared
2. **Los handlers ya existían** - Ya implementados con swagger annotations
3. **Los DTOs son nuevos** - Creados en `http/dto` para separación de concerns
4. **HierarchyService es nuevo** - Servicio de aplicación especializado en jerarquías
5. **Router es opcional** - Ya que main.go configura las rutas directamente

---

## ✅ CHECKLIST FINAL FASE 1

- [x] DTOs con validaciones Gin
- [x] HierarchyService implementado
- [x] Tests unitarios de servicio
- [x] Handlers verificados (ya existían)
- [x] Router configurado
- [x] main.go verificado (ya existía)
- [x] Tests E2E con t.Skip()
- [x] Código compila sin errores
- [x] Documentación de handoff completa

---

**¡Éxito en Fase 2!** 🚀

**Fin del documento de handoff**
