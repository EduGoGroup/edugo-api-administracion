# 🦨 Code Smells

> Malas prácticas y patrones problemáticos identificados en el código

---

## 1. Valores Hardcodeados en Services (MEDIA PRIORIDAD)

### Ubicación
```
internal/application/service/school_service.go:68-85
```

### Código Problemático
```go
school := &entities.School{
    ID:               uuid.New(),
    Name:             req.Name,
    Code:             req.Code,
    Address:          addr,
    City:             nil,                // TODO: agregar cuando se agregue al DTO
    Country:          "CO",               // ⚠️ HARDCODED - valor por defecto
    Phone:            phone,
    Email:            email,
    Metadata:         metadataJSON,
    IsActive:         true,
    SubscriptionTier: "free",             // ⚠️ HARDCODED - valor por defecto
    MaxTeachers:      50,                 // ⚠️ HARDCODED - valor por defecto
    MaxStudents:      500,                // ⚠️ HARDCODED - valor por defecto
    CreatedAt:        now,
    UpdatedAt:        now,
    DeletedAt:        nil,
}
```

### Problemas
1. **Country "CO"** - No todos los clientes son de Colombia
2. **SubscriptionTier "free"** - Debería venir del request o configuración
3. **MaxTeachers/MaxStudents** - Límites fijos sin configuración
4. **City nil** - TODO pendiente

### Solución

**Paso 1: Agregar campos al DTO**
```go
type CreateSchoolRequest struct {
    Name             string                 `json:"name" validate:"required,min=3"`
    Code             string                 `json:"code" validate:"required,min=3"`
    Address          string                 `json:"address"`
    City             string                 `json:"city"`              // NUEVO
    Country          string                 `json:"country"`           // NUEVO - default en config
    ContactEmail     string                 `json:"contact_email"`
    ContactPhone     string                 `json:"contact_phone"`
    SubscriptionTier string                 `json:"subscription_tier"` // NUEVO - opcional
    Metadata         map[string]interface{} `json:"metadata"`
}
```

**Paso 2: Configuración de defaults**
```yaml
# config/config.yaml
defaults:
  school:
    country: "CO"
    subscription_tier: "free"
    max_teachers: 50
    max_students: 500
```

**Paso 3: Usar configuración en service**
```go
func NewSchoolService(repo SchoolRepository, logger Logger, config SchoolConfig) SchoolService {
    return &schoolService{
        repo:   repo,
        logger: logger,
        config: config,
    }
}

// En CreateSchool:
country := req.Country
if country == "" {
    country = s.config.DefaultCountry
}
```

---

## 2. Metadata Vacío Hardcodeado (BAJA PRIORIDAD)

### Ubicación
```
internal/application/service/unit_membership_service.go:100
```

### Código Problemático
```go
membership := &entities.Membership{
    // ...
    Metadata: []byte("{}"),  // ⚠️ String literal para JSON vacío
    // ...
}
```

### Problema
- String literal `"{}"` en lugar de nil o constante
- Inconsistente con otros lugares que usan `json.Marshal`

### Solución
```go
// Constante o helper
var emptyJSONObject = []byte("{}")

// O mejor, permitir nil
Metadata: nil, // GORM maneja nil como NULL o default
```

---

## 3. Validación Duplicada (MEDIA PRIORIDAD)

### Ubicación
```
internal/application/service/school_service.go:49-54
internal/application/dto/school_dto.go (tags validate)
```

### Código Problemático
```go
// En DTO (validación con tags)
type CreateSchoolRequest struct {
    Name string `json:"name" validate:"required,min=3"`
    Code string `json:"code" validate:"required,min=3"`
    // ...
}

// En Service (validación manual DUPLICADA)
func (s *schoolService) CreateSchool(...) {
    // Validaciones (lógica movida del entity)
    if req.Name == "" || len(req.Name) < 3 {
        return nil, errors.NewValidationError("name must be at least 3 characters")
    }
    if req.Code == "" || len(req.Code) < 3 {
        return nil, errors.NewValidationError("code must be at least 3 characters")
    }
    // ...
}
```

### Problema
- **Validación duplicada** en DTO tags y service
- Si se cambia en un lugar, puede olvidarse en el otro
- El handler debería validar el DTO, el service no debería re-validar

### Solución
```go
// Opción 1: Solo validar en handler con validator
func (h *SchoolHandler) CreateSchool(c *gin.Context) {
    var req dto.CreateSchoolRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        // Binding ya valida los tags
        c.JSON(http.StatusBadRequest, ErrorResponse{...})
        return
    }
    // Service confía en que DTO ya está validado
}

// Opción 2: Método Validate() en DTO
func (r *CreateSchoolRequest) Validate() error {
    v := validator.New()
    v.Required(r.Name, "name")
    v.MinLength(r.Name, 3, "name")
    // Única fuente de verdad
    return v.GetError()
}
```

---

## 4. Logging Inconsistente (BAJA PRIORIDAD)

### Ubicación
Múltiples archivos en handlers y services

### Código Problemático
```go
// A veces con campos estructurados
h.logger.Error("create school failed", "error", appErr.Message, "code", appErr.Code)

// A veces sin contexto suficiente
h.logger.Error("unexpected error", "error", err)

// A veces con información útil
h.logger.Info("school created", "school_id", school.ID, "name", school.Name)

// A veces sin información
s.logger.Info("school updated", "id", id)  // Falta "name" como en create
```

### Problema
- Logs inconsistentes dificultan debugging
- Algunos logs tienen más contexto que otros
- No hay estándar definido

### Solución
Crear guía de logging:

```go
// Estándar propuesto:
// Nivel INFO: Operaciones exitosas con IDs relevantes
logger.Info("entity created", 
    "entity_type", "school",
    "entity_id", school.ID,
    "name", school.Name,
    "user_id", userIDFromContext,
)

// Nivel ERROR: Siempre incluir operación, error, y contexto
logger.Error("operation failed",
    "operation", "create_school",
    "error", err.Error(),
    "code", appErr.Code,
    "school_name", req.Name,
    "user_id", userIDFromContext,
)

// Nivel WARN: Situaciones no ideales pero manejadas
logger.Warn("validation failed",
    "field", "email",
    "value", req.Email,
    "reason", "invalid format",
)
```

---

## 5. Error Nil Check Inconsistente (BAJA PRIORIDAD)

### Ubicación
Múltiples services

### Código Problemático
```go
// Patrón 1: Separado
if err != nil {
    return nil, err
}
if school == nil {
    return nil, errors.NewNotFoundError("school")
}

// Patrón 2: Combinado
if err != nil || school == nil {
    return nil, errors.NewNotFoundError("school")
}
```

### Problema
- El patrón 2 **oculta errores de base de datos**
- Si hay un error de conexión, se retorna "not found" en lugar del error real

### Código Correcto
```go
// SIEMPRE manejar err primero
school, err := s.schoolRepo.FindByID(ctx, schoolID)
if err != nil {
    s.logger.Error("database error", "error", err)
    return nil, errors.NewDatabaseError("find school", err)
}
if school == nil {
    return nil, errors.NewNotFoundError("school")
}
```

### Archivos a Revisar
```bash
grep -rn "err != nil || .* == nil" --include="*.go" internal/
```

---

## 6. Context No Propagado Consistentemente (MEDIA PRIORIDAD)

### Ubicación
Algunos handlers y services

### Código Problemático
```go
// ✅ Correcto - usa c.Request.Context()
school, err := h.schoolService.GetSchool(c.Request.Context(), id)

// ⚠️ Potencial problema si se crea context nuevo
ctx := context.Background()  // Pierde información del request
```

### Verificación
```bash
# Buscar usos de context.Background() en código de producción
grep -rn "context.Background()" --include="*.go" internal/
```

---

## 📊 Resumen de Code Smells

| # | Code Smell | Severidad | Esfuerzo | Archivos Afectados |
|---|------------|-----------|----------|-------------------|
| 1 | Valores hardcodeados | Media | 3h | 2 |
| 2 | Metadata hardcodeado | Baja | 30min | 3 |
| 3 | Validación duplicada | Media | 2h | 10+ |
| 4 | Logging inconsistente | Baja | 2h | 15+ |
| 5 | Error nil check | Baja | 1h | 5 |
| 6 | Context propagation | Media | 1h | Verificar |

**Total estimado: ~9.5 horas**

---

## ✅ Priorización

### Hacer Ahora (Sprint Actual)
- [ ] Fix error nil check pattern
- [ ] Agregar campos country/city al DTO

### Próximo Sprint
- [ ] Configuración de defaults
- [ ] Eliminar validación duplicada

### Backlog
- [ ] Estandarizar logging
- [ ] Verificar context propagation
