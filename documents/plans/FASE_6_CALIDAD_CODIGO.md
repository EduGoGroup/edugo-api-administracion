# 🟢 FASE 6: Mejoras de Calidad de Código

**Prioridad**: Baja  
**Estimación**: 6 horas  
**Rama**: `refactor/fase-6-calidad-codigo`

---

## Preparación Git

```bash
git checkout dev
git pull origin dev
git checkout -b refactor/fase-6-calidad-codigo
```

---

## 6.1 Estandarizar Logging

### Problema
Logs inconsistentes con diferentes niveles de contexto en handlers y services.

### Ejemplos de Inconsistencia
```go
// A veces con campos estructurados
h.logger.Error("create school failed", "error", appErr.Message, "code", appErr.Code)

// A veces sin contexto suficiente
h.logger.Error("unexpected error", "error", err)

// A veces con información útil
s.logger.Info("school created", "school_id", school.ID, "name", school.Name)

// A veces sin información
s.logger.Info("school updated", "id", id)  // Falta "name"
```

### Estándar a Implementar

```go
// INFO: Operaciones exitosas
logger.Info("entity created",
    "entity_type", "school",
    "entity_id", school.ID,
    "name", school.Name,
)

// ERROR: Siempre incluir operación, error y contexto
logger.Error("operation failed",
    "operation", "create_school",
    "error", err.Error(),
    "school_name", req.Name,
)

// WARN: Situaciones no ideales pero manejadas
logger.Warn("validation failed",
    "field", "email",
    "value", req.Email,
    "reason", "invalid format",
)
```

### Tareas
1. Definir guía de logging
2. Refactorizar logs existentes siguiendo el estándar
3. Asegurar que todos los logs incluyan: operación, entity_id (cuando aplique)

### Esfuerzo
2 horas

---

## 6.2 Corregir patrón de Error Nil Check

### Problema
Algunos lugares combinan `err != nil || entity == nil` ocultando errores de base de datos.

### Código Problemático
```go
// MAL - Oculta errores de DB
if err != nil || school == nil {
    return nil, errors.NewNotFoundError("school")
}
```

### Código Correcto
```go
// BIEN - Maneja cada caso
if err != nil {
    s.logger.Error("database error", "error", err)
    return nil, errors.NewDatabaseError("find school", err)
}
if school == nil {
    return nil, errors.NewNotFoundError("school")
}
```

### Tareas
1. Buscar patrones problemáticos: `grep -rn "err != nil || .* == nil" --include="*.go" internal/`
2. Separar verificación de error y nil
3. Asegurar que errores de DB se propaguen correctamente

### Archivos a Revisar
- Services en `internal/application/service/`
- Auth services en `internal/auth/service/`

### Esfuerzo
1 hora

---

## 6.3 Verificar propagación de Context

### Problema
Posibles usos de `context.Background()` en código de producción que pierden información del request.

### Código Problemático
```go
// MAL - Pierde información del request
ctx := context.Background()
school, err := s.schoolRepo.FindByID(ctx, id)
```

### Código Correcto
```go
// BIEN - Propaga context del request
school, err := s.schoolRepo.FindByID(c.Request.Context(), id)
```

### Tareas
1. Auditar uso de context: `grep -rn "context.Background()" --include="*.go" internal/`
2. Asegurar que se propague `c.Request.Context()` en handlers
3. Documentar patrones correctos

### Esfuerzo
1 hora

---

## 6.4 Eliminar Validación Duplicada

### Problema
Validación tanto en tags de DTO como manual en service.

### Código Actual
```go
// En DTO (validación con tags)
type CreateSchoolRequest struct {
    Name string `json:"name" validate:"required,min=3"`
}

// En Service (validación manual DUPLICADA)
if req.Name == "" || len(req.Name) < 3 {
    return nil, errors.NewValidationError("name must be at least 3 characters")
}
```

### Decisión
Mantener validación solo en DTOs con tags de `validator`. El handler valida con binding.

### Código Correcto
```go
// Solo en DTO
type CreateSchoolRequest struct {
    Name string `json:"name" validate:"required,min=3" binding:"required,min=3"`
}

// Handler confía en validación del binding
func (h *SchoolHandler) CreateSchool(c *gin.Context) {
    var req dto.CreateSchoolRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        _ = c.Error(errors.NewValidationError(err.Error()))
        return
    }
    // Service NO re-valida
}
```

### Tareas
1. Asegurar que DTOs tengan tags `binding` además de `validate`
2. Eliminar validaciones manuales redundantes en services
3. Verificar que errores de validación sean descriptivos

### Esfuerzo
2 horas

---

## Documentación a Actualizar

Al completar esta fase, actualizar:

- `documents/improvements/CODE_SMELLS.md` - Eliminar secciones resueltas (3, 4, 5, 6)
- `documents/improvements/README.md` - Actualizar métricas finales y marcar como completado

---

## Finalización

```bash
git add .
git commit -m "refactor: mejorar calidad de código (logging, nil checks, validación)"
git push origin refactor/fase-6-calidad-codigo
```

### Crear PR a dev con:
- Título: `refactor: mejorar calidad de código (logging, nil checks, validación)`
- Descripción: Fase 6 del plan de mejoras - Calidad de código

---

## Checklist

- [ ] Logs estandarizados en handlers
- [ ] Logs estandarizados en services
- [ ] Patrones `err != nil || x == nil` corregidos
- [ ] Uso de `context.Background()` auditado y corregido
- [ ] Validaciones duplicadas eliminadas
- [ ] DTOs tienen tags `binding`
- [ ] Todos los tests pasan
- [ ] Documentación actualizada
- [ ] PR creado a dev
