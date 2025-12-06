# 🔄 Refactorizaciones Pendientes

> Mejoras estructurales y de diseño identificadas en el código

---

## 1. Error Handling Repetitivo en Handlers (ALTA PRIORIDAD)

### Ubicación
```
internal/infrastructure/http/handler/*.go
```

### Problema
Cada handler repite el mismo patrón de manejo de errores (~15 líneas por método):

```go
// Este patrón se repite en TODOS los handlers
if err != nil {
    if appErr, ok := errors.GetAppError(err); ok {
        h.logger.Error("operation failed", "error", appErr.Message, "code", appErr.Code)
        c.JSON(appErr.StatusCode, ErrorResponse{Error: appErr.Message, Code: string(appErr.Code)})
        return
    }
    h.logger.Error("unexpected error", "error", err)
    c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error", Code: "INTERNAL_ERROR"})
    return
}
```

### Impacto
- **~300 líneas de código duplicado** entre todos los handlers
- Difícil mantener consistencia en respuestas de error
- Cambios requieren modificar múltiples archivos

### Solución Propuesta

**Opción A: Middleware de Error Handling**
```go
// internal/infrastructure/http/middleware/error_handler.go
package middleware

func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        
        if len(c.Errors) > 0 {
            err := c.Errors.Last().Err
            if appErr, ok := errors.GetAppError(err); ok {
                c.JSON(appErr.StatusCode, ErrorResponse{
                    Error: appErr.Message,
                    Code:  string(appErr.Code),
                })
                return
            }
            c.JSON(http.StatusInternalServerError, ErrorResponse{
                Error: "internal server error",
                Code:  "INTERNAL_ERROR",
            })
        }
    }
}

// En handlers, simplemente:
func (h *SchoolHandler) GetSchool(c *gin.Context) {
    school, err := h.schoolService.GetSchool(c.Request.Context(), id)
    if err != nil {
        _ = c.Error(err) // El middleware lo maneja
        return
    }
    c.JSON(http.StatusOK, school)
}
```

**Opción B: Helper Function**
```go
// internal/infrastructure/http/handler/helpers.go
func handleError(c *gin.Context, logger logger.Logger, err error, operation string) {
    if appErr, ok := errors.GetAppError(err); ok {
        logger.Error(operation+" failed", "error", appErr.Message, "code", appErr.Code)
        c.JSON(appErr.StatusCode, ErrorResponse{Error: appErr.Message, Code: string(appErr.Code)})
        return
    }
    logger.Error("unexpected error", "error", err, "operation", operation)
    c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error", Code: "INTERNAL_ERROR"})
}

// Uso:
func (h *SchoolHandler) GetSchool(c *gin.Context) {
    school, err := h.schoolService.GetSchool(c.Request.Context(), id)
    if err != nil {
        handleError(c, h.logger, err, "get school")
        return
    }
    c.JSON(http.StatusOK, school)
}
```

### Esfuerzo Estimado
- Implementar middleware/helper: 2 horas
- Refactorizar handlers existentes: 4 horas
- Tests: 2 horas
- **Total: ~8 horas**

---

## 2. Validación de Roles Hardcodeada (MEDIA PRIORIDAD)

### Ubicación
```
internal/application/service/unit_membership_service.go:75-85
```

### Código Actual
```go
// Validar role (lógica movida del value object)
validRoles := []string{"teacher", "student", "guardian", "coordinator", "admin", "assistant"}
isValid := false
for _, r := range validRoles {
    if req.Role == r {
        isValid = true
        break
    }
}
if !isValid {
    return nil, errors.NewValidationError("invalid membership role")
}
```

### Problemas
1. Roles hardcodeados en el código
2. Lógica de validación repetida en múltiples lugares
3. No usa los enums de `edugo-shared/common`
4. Difícil agregar/modificar roles

### Solución Propuesta

```go
// internal/domain/valueobject/membership_role.go
package valueobject

import "github.com/EduGoGroup/edugo-shared/common/types/enum"

type MembershipRole string

const (
    RoleTeacher     MembershipRole = "teacher"
    RoleStudent     MembershipRole = "student"
    RoleGuardian    MembershipRole = "guardian"
    RoleCoordinator MembershipRole = "coordinator"
    RoleAdmin       MembershipRole = "admin"
    RoleAssistant   MembershipRole = "assistant"
    RoleDirector    MembershipRole = "director"
    RoleObserver    MembershipRole = "observer"
)

var validMembershipRoles = map[MembershipRole]bool{
    RoleTeacher:     true,
    RoleStudent:     true,
    RoleGuardian:    true,
    RoleCoordinator: true,
    RoleAdmin:       true,
    RoleAssistant:   true,
    RoleDirector:    true,
    RoleObserver:    true,
}

func (r MembershipRole) IsValid() bool {
    return validMembershipRoles[r]
}

func (r MembershipRole) String() string {
    return string(r)
}

func ParseMembershipRole(s string) (MembershipRole, error) {
    role := MembershipRole(s)
    if !role.IsValid() {
        return "", fmt.Errorf("invalid membership role: %s", s)
    }
    return role, nil
}
```

### Uso en Service
```go
// Validar role usando value object
role, err := valueobject.ParseMembershipRole(req.Role)
if err != nil {
    return nil, errors.NewValidationError(err.Error())
}
```

### Esfuerzo Estimado
- Crear value object: 1 hora
- Refactorizar servicios: 2 horas
- Tests: 1 hora
- **Total: ~4 horas**

---

## 3. Parámetro `activeOnly` Sin Implementar (ALTA PRIORIDAD - BUG)

### Ubicación
```
internal/application/service/unit_membership_service.go:138-176
```

### Código Actual
```go
func (s *unitMembershipService) ListMembershipsByUnit(ctx context.Context, unitID string, activeOnly bool) ([]dto.MembershipResponse, error) {
    uid, err := uuid.Parse(unitID)
    if err != nil {
        return nil, errors.NewValidationError("invalid unit ID")
    }

    // ⚠️ activeOnly NO SE USA - siempre retorna todas
    memberships, err := s.membershipRepo.FindByUnit(ctx, uid)
    // ...
}
```

### Problema
- El parámetro `activeOnly` se recibe pero **nunca se usa**
- La API documenta que filtra por activos, pero no lo hace
- **Es un bug funcional**

### Solución Propuesta

**Opción A: Filtrar en el repositorio**
```go
// internal/domain/repository/unit_membership_repository.go
type UnitMembershipRepository interface {
    FindByUnit(ctx context.Context, unitID uuid.UUID, activeOnly bool) ([]*entities.Membership, error)
    // ...
}

// Implementación PostgreSQL
func (r *unitMembershipRepository) FindByUnit(ctx context.Context, unitID uuid.UUID, activeOnly bool) ([]*entities.Membership, error) {
    query := r.db.Where("academic_unit_id = ?", unitID)
    if activeOnly {
        query = query.Where("is_active = ?", true).Where("withdrawn_at IS NULL")
    }
    // ...
}
```

**Opción B: Filtrar en el servicio**
```go
func (s *unitMembershipService) ListMembershipsByUnit(ctx context.Context, unitID string, activeOnly bool) ([]dto.MembershipResponse, error) {
    memberships, err := s.membershipRepo.FindByUnit(ctx, uid)
    if err != nil {
        return nil, err
    }
    
    if activeOnly {
        filtered := make([]*entities.Membership, 0)
        for _, m := range memberships {
            if m.IsActive && m.WithdrawnAt == nil {
                filtered = append(filtered, m)
            }
        }
        memberships = filtered
    }
    // ...
}
```

### Esfuerzo Estimado
- Implementar filtro: 1 hora
- Actualizar tests: 1 hora
- **Total: ~2 horas**

---

## 4. `ListMembershipsByRole` No Filtra (ALTA PRIORIDAD - BUG)

### Ubicación
```
internal/application/service/unit_membership_service.go:174-177
```

### Código Actual
```go
func (s *unitMembershipService) ListMembershipsByRole(ctx context.Context, unitID string, role string, activeOnly bool) ([]dto.MembershipResponse, error) {
    // Implementación simplificada
    return s.ListMembershipsByUnit(ctx, unitID, activeOnly)
}
```

### Problema
- **IGNORA COMPLETAMENTE el parámetro `role`**
- Retorna todas las membresías de la unidad sin filtrar
- Comentario "Implementación simplificada" = incompleta
- **Bug crítico si clientes dependen de este filtro**

### Solución
```go
func (s *unitMembershipService) ListMembershipsByRole(ctx context.Context, unitID string, role string, activeOnly bool) ([]dto.MembershipResponse, error) {
    uid, err := uuid.Parse(unitID)
    if err != nil {
        return nil, errors.NewValidationError("invalid unit ID")
    }

    // Validar rol
    if _, err := valueobject.ParseMembershipRole(role); err != nil {
        return nil, errors.NewValidationError("invalid role")
    }

    memberships, err := s.membershipRepo.FindByUnitAndRole(ctx, uid, role, activeOnly)
    if err != nil {
        return nil, errors.NewDatabaseError("find memberships", err)
    }

    responses := make([]dto.MembershipResponse, len(memberships))
    for i, m := range memberships {
        responses[i] = dto.ToMembershipResponse(m)
    }
    return responses, nil
}
```

### Esfuerzo Estimado
- Implementar método en repositorio: 1 hora
- Implementar en servicio: 30 min
- Tests: 1 hora
- **Total: ~2.5 horas**

---

## 5. Centralizar Tipos de Unidad Académica (MEDIA PRIORIDAD)

### Ubicación
```
internal/application/dto/academic_unit_dto.go:13
internal/application/service/academic_unit_service.go (validación)
```

### Código Actual
```go
// En DTO (validación con tags)
Type string `json:"type" validate:"required,oneof=school grade section club department"`

// En algún servicio (si existe validación manual)
validTypes := []string{"school", "grade", "section", "club", "department"}
```

### Problema
- Tipos definidos en múltiples lugares (tags, arrays)
- No hay un lugar centralizado para agregar nuevos tipos
- Fácil que se desincronicen

### Solución
```go
// internal/domain/valueobject/unit_type.go
package valueobject

type UnitType string

const (
    UnitTypeSchool     UnitType = "school"
    UnitTypeGrade      UnitType = "grade"
    UnitTypeSection    UnitType = "section"
    UnitTypeClub       UnitType = "club"
    UnitTypeDepartment UnitType = "department"
)

func AllUnitTypes() []UnitType {
    return []UnitType{
        UnitTypeSchool,
        UnitTypeGrade,
        UnitTypeSection,
        UnitTypeClub,
        UnitTypeDepartment,
    }
}

func AllUnitTypesStrings() []string {
    types := AllUnitTypes()
    result := make([]string, len(types))
    for i, t := range types {
        result[i] = string(t)
    }
    return result
}
```

---

## 📊 Resumen de Refactorizaciones

| # | Refactorización | Prioridad | Esfuerzo | Impacto |
|---|-----------------|-----------|----------|---------|
| 1 | Error handling middleware | Alta | 8h | Alto - reduce ~300 líneas |
| 2 | Value object para roles | Media | 4h | Medio - mejor diseño |
| 3 | Implementar activeOnly | Alta | 2h | Alto - fix bug |
| 4 | Implementar ListByRole | Alta | 2.5h | Alto - fix bug |
| 5 | Value object unit types | Media | 2h | Medio - mejor diseño |

**Total estimado: ~18.5 horas de trabajo**

---

## ✅ Orden de Implementación Recomendado

```
1. [URGENTE] Fix ListMembershipsByRole (bug)     → 2.5h
2. [URGENTE] Implementar activeOnly (bug)        → 2h
3. [SPRINT] Error handling middleware            → 8h
4. [SPRINT] Value object roles                   → 4h
5. [BACKLOG] Value object unit types             → 2h
```
