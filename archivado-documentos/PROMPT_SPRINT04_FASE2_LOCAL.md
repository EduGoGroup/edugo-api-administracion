# PROMPT SPRINT-04 FASE 2 - CLAUDE CODE LOCAL

**Proyecto:** edugo-api-administracion  
**Sprint:** Sprint-04 - Services/API  
**Ejecutor:** Claude Code Local  
**Duración estimada:** 1-2 horas  
**Branch:** `feature/sprint-04-services-api` (continuación)

---

## 🎯 TU OBJETIVO (Fase 2 - Con servidor local)

Completar el trabajo que Claude Code Web dejó en Fase 1, específicamente:
1. **Levantar servidor HTTP local**
2. **Descomentar y ejecutar tests E2E**
3. **Validar endpoints con Postman/curl**
4. **Crear PR a dev**
5. **Monitorear CI/CD y hacer merge**

---

## 📋 PREREQUISITOS

Antes de empezar, lee:
1. `HANDOFF_SPRINT04_FASE1_TO_FASE2.md` - Qué hizo Claude Web
2. `PROMPT_SPRINT04_FASE1_WEB.md` - Contexto de Fase 1

---

## 📋 TAREAS FASE 2

### TASK-01: Revisar Trabajo de Fase 1 (15min)

```bash
# 1. Checkout de la branch
git checkout feature/sprint-04-services-api
git pull origin feature/sprint-04-services-api

# 2. Verificar que compila
go build ./...

# 3. Leer handoff
cat HANDOFF_SPRINT04_FASE1_TO_FASE2.md

# 4. Ejecutar tests unitarios (deben pasar)
go test ./internal/application/service/... -v
```

**Validar que existe:**
- ✅ DTOs en `internal/infrastructure/http/dto/`
- ✅ HierarchyService en `internal/application/service/`
- ✅ Handlers en `internal/infrastructure/http/handler/`
- ✅ Router en `internal/infrastructure/http/router/`
- ✅ main.go en `cmd/api/`
- ✅ Tests E2E con stubs en `test/integration/http_api_test.go`

---

### TASK-02: Levantar Servidor HTTP (20min)

**Opción A: Con PostgreSQL local**

```bash
# 1. Verificar que PostgreSQL está corriendo
docker ps | grep postgres

# 2. Si no está, levantar
docker-compose up -d postgres

# 3. Ejecutar migraciones (incluye Sprint-03 ltree!)
migrate -path migrations -database "postgresql://edugo_user:edugo_pass@localhost:5432/edugo_admin?sslmode=disable" up

# 4. Verificar que extensión ltree está habilitada
psql -h localhost -U edugo_user -d edugo_admin -c "SELECT * FROM pg_extension WHERE extname = 'ltree';"

# 5. Configurar variables de entorno
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=edugo_user
export DB_PASSWORD=edugo_pass
export DB_NAME=edugo_admin
export PORT=8080

# 6. Levantar servidor
go run cmd/api/main.go

# Deberías ver:
# ✅ Connected to database
# 🚀 Server starting on port 8080
```

**Opción B: Con testcontainers (en tests)**

Los tests de integración levantarán su propio servidor usando testcontainers.

**✅ Validación:**

```bash
# En otra terminal
curl http://localhost:8080/health
# Debe retornar: {"status":"ok"}
```

---

### TASK-03: Tests Manuales con curl (30min)

**Test 1: CRUD de Escuelas**

```bash
# Crear escuela
SCHOOL=$(curl -s -X POST http://localhost:8080/api/v1/schools \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test School",
    "code": "TS001",
    "address": "123 Main St"
  }')

echo $SCHOOL | jq .
SCHOOL_ID=$(echo $SCHOOL | jq -r .id)

# Listar escuelas
curl http://localhost:8080/api/v1/schools | jq .

# Obtener escuela por ID
curl http://localhost:8080/api/v1/schools/$SCHOOL_ID | jq .

# Actualizar escuela
curl -X PUT http://localhost:8080/api/v1/schools/$SCHOOL_ID \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated School Name"
  }' | jq .
```

**Test 2: CRUD de Unidades Académicas**

```bash
# Crear grado (raíz)
GRADE=$(curl -s -X POST http://localhost:8080/api/v1/units \
  -H "Content-Type: application/json" \
  -d "{
    \"school_id\": \"$SCHOOL_ID\",
    \"type\": \"grade\",
    \"name\": \"Grade 1\",
    \"code\": \"G1\"
  }")

echo $GRADE | jq .
GRADE_ID=$(echo $GRADE | jq -r .id)

# Crear sección bajo el grado
SECTION=$(curl -s -X POST http://localhost:8080/api/v1/units \
  -H "Content-Type: application/json" \
  -d "{
    \"parent_unit_id\": \"$GRADE_ID\",
    \"school_id\": \"$SCHOOL_ID\",
    \"type\": \"section\",
    \"name\": \"Section A\",
    \"code\": \"G1-A\"
  }")

SECTION_ID=$(echo $SECTION | jq -r .id)

# Crear club bajo la sección
CLUB=$(curl -s -X POST http://localhost:8080/api/v1/units \
  -H "Content-Type: application/json" \
  -d "{
    \"parent_unit_id\": \"$SECTION_ID\",
    \"school_id\": \"$SCHOOL_ID\",
    \"type\": \"club\",
    \"name\": \"Math Club\",
    \"code\": \"G1-A-MC\"
  }")

CLUB_ID=$(echo $CLUB | jq -r .id)
```

**Test 3: Árbol Jerárquico (ltree!)**

```bash
# Obtener árbol completo del grado
curl http://localhost:8080/api/v1/units/$GRADE_ID/tree | jq .

# Debe mostrar:
# {
#   "id": "...",
#   "name": "Grade 1",
#   "code": "G1",
#   "type": "grade",
#   "depth": 0,
#   "children": [
#     {
#       "id": "...",
#       "name": "Section A",
#       "code": "G1-A",
#       "type": "section",
#       "depth": 1,
#       "children": [
#         {
#           "id": "...",
#           "name": "Math Club",
#           "code": "G1-A-MC",
#           "type": "club",
#           "depth": 2,
#           "children": []
#         }
#       ]
#     }
#   ]
# }
```

**Test 4: Filtrado por Profundidad (ltree!)**

```bash
# Listar solo nivel 1 (grados)
curl "http://localhost:8080/api/v1/units?school_id=$SCHOOL_ID&depth=1" | jq .

# Listar solo nivel 2 (secciones)
curl "http://localhost:8080/api/v1/units?school_id=$SCHOOL_ID&depth=2" | jq .

# Listar solo nivel 3 (clubs)
curl "http://localhost:8080/api/v1/units?school_id=$SCHOOL_ID&depth=3" | jq .
```

**Test 5: Mover Unidad (MoveSubtree con ltree)**

```bash
# Crear segundo grado
GRADE2=$(curl -s -X POST http://localhost:8080/api/v1/units \
  -H "Content-Type: application/json" \
  -d "{
    \"school_id\": \"$SCHOOL_ID\",
    \"type\": \"grade\",
    \"name\": \"Grade 2\",
    \"code\": \"G2\"
  }")

GRADE2_ID=$(echo $GRADE2 | jq -r .id)

# Mover Section A (con su club) de Grade 1 a Grade 2
curl -X PUT http://localhost:8080/api/v1/units/$SECTION_ID \
  -H "Content-Type: application/json" \
  -d "{
    \"parent_unit_id\": \"$GRADE2_ID\"
  }" | jq .

# Verificar que se movió correctamente
curl http://localhost:8080/api/v1/units/$GRADE2_ID/tree | jq .

# Árbol de Grade 2 ahora debe mostrar Section A y Math Club
```

**Test 6: Manejo de Errores**

```bash
# JSON inválido
curl -X POST http://localhost:8080/api/v1/schools \
  -H "Content-Type: application/json" \
  -d 'invalid json'
# Debe retornar 400

# ID inexistente
curl http://localhost:8080/api/v1/schools/00000000-0000-0000-0000-000000000000
# Debe retornar 404

# Código duplicado
curl -X POST http://localhost:8080/api/v1/schools \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Another School",
    "code": "TS001",
    "address": "456 Oak Ave"
  }'
# Debe retornar 400 (código ya existe)

# Intentar crear ciclo
curl -X PUT http://localhost:8080/api/v1/units/$GRADE_ID \
  -H "Content-Type: application/json" \
  -d "{
    \"parent_unit_id\": \"$CLUB_ID\"
  }"
# Debe retornar 400 (circular reference)
```

**Documentar resultados** en un archivo `MANUAL_TESTS_RESULTS.md`

---

### TASK-04: Implementar Helper para Tests E2E (20min)

**Archivo:** `test/integration/http_api_test.go`

Antes de descomentar los tests, agregar el helper:

```go
// setupTestServer levanta un servidor Gin de test con testcontainers
func setupTestServer(t *testing.T) (*httptest.Server, *sql.DB, func()) {
	// Setup PostgreSQL con testcontainers
	db, dbCleanup := setupTestDB(t)
	
	// Crear repositorios
	schoolRepo := repository.NewPostgresSchoolRepository(db)
	unitRepo := repository.NewPostgresAcademicUnitRepository(db)
	
	// Configurar router
	cfg := &router.Config{
		SchoolRepo: schoolRepo,
		UnitRepo:   unitRepo,
	}
	
	ginRouter := router.SetupRouter(cfg)
	
	// Crear test server
	server := httptest.NewServer(ginRouter)
	
	cleanup := func() {
		server.Close()
		dbCleanup()
	}
	
	return server, db, cleanup
}

// Helper para hacer requests HTTP
func doRequest(t *testing.T, server *httptest.Server, method, path string, body interface{}) *http.Response {
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, _ := json.Marshal(body)
		bodyReader = bytes.NewReader(bodyBytes)
	}
	
	req, err := http.NewRequest(method, server.URL+path, bodyReader)
	require.NoError(t, err)
	
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	
	return resp
}
```

---

### TASK-05: Descomentar y Ejecutar Tests E2E (45min)

**Para cada test en `test/integration/http_api_test.go`:**

1. **Quitar `t.Skip()`**
2. **Descomentar código**
3. **Implementar usando helpers**
4. **Ejecutar y validar**

**Ejemplo descomentado:**

```go
func TestSchoolAPI_CreateAndGet(t *testing.T) {
	server, db, cleanup := setupTestServer(t)
	defer cleanup()
	
	// 1. Crear escuela
	createReq := dto.CreateSchoolRequest{
		Name:    "Test School",
		Code:    "TS001",
		Address: "123 Main St",
	}
	
	resp := doRequest(t, server, "POST", "/api/v1/schools", createReq)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	
	var schoolResp dto.SchoolResponse
	json.NewDecoder(resp.Body).Decode(&schoolResp)
	resp.Body.Close()
	
	assert.NotEmpty(t, schoolResp.ID)
	assert.Equal(t, "Test School", schoolResp.Name)
	assert.Equal(t, "TS001", schoolResp.Code)
	
	// 2. Obtener escuela por ID
	resp = doRequest(t, server, "GET", "/api/v1/schools/"+schoolResp.ID, nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	
	var getResp dto.SchoolResponse
	json.NewDecoder(resp.Body).Decode(&getResp)
	resp.Body.Close()
	
	assert.Equal(t, schoolResp.ID, getResp.ID)
	assert.Equal(t, "Test School", getResp.Name)
}
```

**Ejecutar tests:**

```bash
go test -tags=integration ./test/integration/... -v -run TestAPI

# Ejecutar uno específico
go test -tags=integration ./test/integration/... -v -run TestSchoolAPI_CreateAndGet
```

**Criterio de éxito:**
- ✅ Todos los tests E2E pasan
- ✅ Sin `t.Skip()` en el código
- ✅ Coverage de handlers >= 60%

---

### TASK-06: Validación de Performance con ltree (15min)

**Benchmark informal:**

```bash
# Crear jerarquía grande (script en Go o bash loop)
for i in {1..10}; do
  # Crear 10 grados
  # Cada grado con 5 secciones
  # Cada sección con 3 clubs
  # Total: 10 + 50 + 150 = 210 unidades
done

# Medir tiempo del endpoint /tree
time curl http://localhost:8080/api/v1/units/$GRADE1_ID/tree > /dev/null

# Debería ser < 100ms gracias a ltree
```

**Validar en PostgreSQL:**

```sql
-- Ver que usa índice GIST
EXPLAIN ANALYZE
SELECT * FROM academic_units
WHERE path <@ (SELECT path FROM academic_units WHERE id = 'UNIT_ID');
```

Documentar resultados comparando con lo que tomaría sin ltree.

---

### TASK-07: Validación Completa (15min)

```bash
# Compilación
go build ./...

# Tests unitarios (deben seguir pasando)
go test ./internal/... -v

# Tests de integración (ahora DEBEN pasar)
go test -tags=integration ./test/integration/... -v

# Coverage
go test -tags=integration ./test/integration/... -coverprofile=coverage.out
go tool cover -html=coverage.out

# Lint
make lint
```

**Criterios:**
- ✅ Compilación OK
- ✅ Tests unitarios: 100% pasando
- ✅ Tests E2E: 100% pasando
- ✅ Coverage handlers >= 60%
- ✅ Coverage services >= 80%
- ✅ Lint sin errores críticos

---

### TASK-08: Crear PR a dev (10min)

**IMPORTANTE**: El PR debe ir hacia `dev`, NO hacia `main`. Workflow: feature → dev → main

**Título:** `feat(api): Sprint-04 - HTTP REST API with ltree support`

**Descripción:**

```markdown
## 🎯 Sprint-04: Services/API

Implementación completa de capa de aplicación y API REST con soporte ltree.

### ✅ Completado

**DTOs (Data Transfer Objects):**
- CreateSchoolRequest, UpdateSchoolRequest, SchoolResponse
- CreateUnitRequest, UpdateUnitRequest, UnitResponse
- UnitTreeNode para árbol jerárquico
- ErrorResponse, SuccessResponse

**Application Service:**
- HierarchyService con 4 métodos principales
- CreateUnit con validaciones (escuela existe, código único, etc.)
- GetUnitTree usando ltree FindDescendants
- MoveUnit usando ltree MoveSubtree
- ValidateNoCircularReference usando ltree

**HTTP Endpoints (10):**

Schools:
- POST   /api/v1/schools
- GET    /api/v1/schools
- GET    /api/v1/schools/:id
- PUT    /api/v1/schools/:id
- DELETE /api/v1/schools/:id

Academic Units:
- POST   /api/v1/units
- GET    /api/v1/units (con filtro ?depth= usando ltree nlevel)
- GET    /api/v1/units/:id
- GET    /api/v1/units/:id/tree (árbol completo con ltree!)
- PUT    /api/v1/units/:id (incluye mover unidad)
- DELETE /api/v1/units/:id

**Aprovechar Sprint-03 (ltree):**
- ✅ FindDescendants para árbol completo en 1 query
- ✅ FindBySchoolIDAndDepth para filtro por profundidad
- ✅ MoveSubtree para reorganizar jerarquías
- ✅ Paths se actualizan automáticamente

### 📊 Validación

- ✅ Tests unitarios: PASS
- ✅ Tests E2E: PASS (5/5)
- ✅ Tests manuales: PASS
- ✅ Coverage services: >= 80%
- ✅ Coverage handlers: >= 60%
- ✅ Servidor HTTP funcional
- ✅ Manejo de errores HTTP correcto

### 🔄 Proceso de Desarrollo

- **FASE 1 (Claude Web)**: DTOs, servicios, handlers, router, tests con stubs
- **FASE 2 (Claude Local)**: ✅ Servidor HTTP, tests E2E, validación completa

### 📁 Archivos Principales

```
internal/
├── application/service/
│   ├── hierarchy_service.go          (Servicio de aplicación)
│   └── hierarchy_service_test.go     (Tests unitarios)
├── infrastructure/http/
│   ├── dto/                           (DTOs con validaciones)
│   ├── handler/                       (Handlers HTTP)
│   └── router/                        (Configuración Gin)
└── cmd/api/main.go                    (Entry point)

test/integration/http_api_test.go      (Tests E2E)
```

### 🚀 Performance

Gracias a ltree (Sprint-03):
- Árbol de 200+ unidades: < 100ms
- Filtro por profundidad: Usa índice GIST
- Mover subárbol: Actualización en cascada eficiente

---

**Revisores**: @EduGoGroup
**Relacionado**: Sprint-04 - API REST con ltree
**Depende de**: Sprint-03 (ltree)
```

**Crear PR:**
```bash
# La tool de GitHub lo hará
```

---

### TASK-09: Monitorear CI/CD y Merge (Variable)

Una vez creado el PR:

1. **Esperar pipeline**
2. **Verificar jobs:**
   - Unit Tests
   - Integration Tests  
   - Lint
3. **Si falla:** analizar, corregir, push
4. **Cuando esté verde:** mergear con squash

---

## 📊 CHECKLIST FINAL FASE 2

Antes de mergear:

### Servidor
- [ ] Servidor levanta correctamente
- [ ] Health check funciona
- [ ] Conexión a PostgreSQL OK
- [ ] ltree extension habilitada

### Endpoints
- [ ] POST /schools funciona
- [ ] GET /schools/:id funciona
- [ ] GET /units/:id/tree retorna árbol correcto
- [ ] GET /units?depth=N usa ltree
- [ ] PUT /units/:id (move) funciona
- [ ] Errores retornan códigos HTTP correctos

### Tests
- [ ] Tests unitarios: ✅
- [ ] Tests E2E: ✅
- [ ] Coverage >= 60-80%
- [ ] Sin `t.Skip()` en código

### CI/CD
- [ ] Pipeline verde
- [ ] PR aprobado
- [ ] Listo para merge

---

## 🎯 COMANDOS RÁPIDOS

```bash
# Checkout
git checkout feature/sprint-04-services-api

# Levantar servidor
export DB_HOST=localhost DB_PORT=5432 DB_USER=edugo_user DB_PASSWORD=edugo_pass DB_NAME=edugo_admin
go run cmd/api/main.go

# Tests E2E
go test -tags=integration ./test/integration/... -v -run TestAPI

# Test manual
curl http://localhost:8080/health

# Coverage
go test -tags=integration ./test/integration/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

---

**¡Éxito en Fase 2!** 🚀
