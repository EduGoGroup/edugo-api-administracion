# 🟡 FASE 3: Value Objects y Validación

**Prioridad**: Media  
**Estimación**: 6 horas  
**Rama**: `refactor/fase-3-value-objects`

---

## Preparación Git

```bash
git checkout dev
git pull origin dev
git checkout -b refactor/fase-3-value-objects
```

---

## 3.1 Crear Value Object para MembershipRole

### Problema Actual
Roles hardcodeados en arrays dentro del servicio:
```go
validRoles := []string{"teacher", "student", "guardian", "coordinator", "admin", "assistant"}
```

### Ubicación del Problema
```
internal/application/service/unit_membership_service.go:75-85
```

### Tareas
1. Crear carpeta `internal/domain/valueobject/`
2. Crear archivo `membership_role.go`
3. Definir constantes y métodos de validación
4. Refactorizar servicios para usar el value object
5. Agregar tests

### Crear Archivo
```
internal/domain/valueobject/membership_role.go
```

### Código
```go
package valueobject

import "fmt"

// MembershipRole representa el rol de un usuario en una unidad académica
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

// IsValid verifica si el rol es válido
func (r MembershipRole) IsValid() bool {
    return validMembershipRoles[r]
}

// String retorna el rol como string
func (r MembershipRole) String() string {
    return string(r)
}

// ParseMembershipRole convierte un string a MembershipRole
func ParseMembershipRole(s string) (MembershipRole, error) {
    role := MembershipRole(s)
    if !role.IsValid() {
        return "", fmt.Errorf("invalid membership role: %s", s)
    }
    return role, nil
}

// AllMembershipRoles retorna todos los roles válidos
func AllMembershipRoles() []MembershipRole {
    return []MembershipRole{
        RoleTeacher,
        RoleStudent,
        RoleGuardian,
        RoleCoordinator,
        RoleAdmin,
        RoleAssistant,
        RoleDirector,
        RoleObserver,
    }
}

// AllMembershipRolesStrings retorna todos los roles como strings
func AllMembershipRolesStrings() []string {
    roles := AllMembershipRoles()
    result := make([]string, len(roles))
    for i, r := range roles {
        result[i] = string(r)
    }
    return result
}
```

### Uso en Servicio
```go
import "github.com/EduGoGroup/edugo-api-administracion/internal/domain/valueobject"

// Validar role usando value object
role, err := valueobject.ParseMembershipRole(req.Role)
if err != nil {
    return nil, errors.NewValidationError(err.Error())
}
```

### Esfuerzo
4 horas

---

## 3.2 Crear Value Object para UnitType

### Problema Actual
Tipos de unidad definidos en múltiples lugares (tags de validación, arrays).

### Ubicación del Problema
```
internal/application/dto/academic_unit_dto.go:13
```

### Código Actual
```go
Type string `json:"type" validate:"required,oneof=school grade section club department"`
```

### Tareas
1. Crear archivo `unit_type.go`
2. Definir constantes y helpers
3. Refactorizar DTOs y servicios

### Crear Archivo
```
internal/domain/valueobject/unit_type.go
```

### Código
```go
package valueobject

import "fmt"

// UnitType representa el tipo de una unidad académica
type UnitType string

const (
    UnitTypeSchool     UnitType = "school"
    UnitTypeGrade      UnitType = "grade"
    UnitTypeSection    UnitType = "section"
    UnitTypeClub       UnitType = "club"
    UnitTypeDepartment UnitType = "department"
)

var validUnitTypes = map[UnitType]bool{
    UnitTypeSchool:     true,
    UnitTypeGrade:      true,
    UnitTypeSection:    true,
    UnitTypeClub:       true,
    UnitTypeDepartment: true,
}

// IsValid verifica si el tipo es válido
func (t UnitType) IsValid() bool {
    return validUnitTypes[t]
}

// String retorna el tipo como string
func (t UnitType) String() string {
    return string(t)
}

// ParseUnitType convierte un string a UnitType
func ParseUnitType(s string) (UnitType, error) {
    unitType := UnitType(s)
    if !unitType.IsValid() {
        return "", fmt.Errorf("invalid unit type: %s", s)
    }
    return unitType, nil
}

// AllUnitTypes retorna todos los tipos válidos
func AllUnitTypes() []UnitType {
    return []UnitType{
        UnitTypeSchool,
        UnitTypeGrade,
        UnitTypeSection,
        UnitTypeClub,
        UnitTypeDepartment,
    }
}

// AllUnitTypesStrings retorna todos los tipos como strings
func AllUnitTypesStrings() []string {
    types := AllUnitTypes()
    result := make([]string, len(types))
    for i, t := range types {
        result[i] = string(t)
    }
    return result
}
```

### Esfuerzo
2 horas

---

## Tests

### Crear Archivo
```
internal/domain/valueobject/membership_role_test.go
internal/domain/valueobject/unit_type_test.go
```

### Ejemplo de Test
```go
package valueobject_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/EduGoGroup/edugo-api-administracion/internal/domain/valueobject"
)

func TestMembershipRole_IsValid(t *testing.T) {
    tests := []struct {
        role     valueobject.MembershipRole
        expected bool
    }{
        {valueobject.RoleTeacher, true},
        {valueobject.RoleStudent, true},
        {"invalid_role", false},
        {"", false},
    }

    for _, tt := range tests {
        t.Run(string(tt.role), func(t *testing.T) {
            assert.Equal(t, tt.expected, tt.role.IsValid())
        })
    }
}

func TestParseMembershipRole(t *testing.T) {
    role, err := valueobject.ParseMembershipRole("teacher")
    assert.NoError(t, err)
    assert.Equal(t, valueobject.RoleTeacher, role)

    _, err = valueobject.ParseMembershipRole("invalid")
    assert.Error(t, err)
}
```

---

## Documentación a Actualizar

Al completar esta fase, actualizar:

- `documents/improvements/REFACTORING.md` - Eliminar sección 2 y 5
- `documents/improvements/CODE_SMELLS.md` - Eliminar referencias a roles hardcodeados
- `documents/ARCHITECTURE.md` - Agregar sección sobre Value Objects en Domain Layer

---

## Finalización

```bash
git add .
git commit -m "refactor: crear value objects para MembershipRole y UnitType"
git push origin refactor/fase-3-value-objects
```

### Crear PR a dev con:
- Título: `refactor: crear value objects para MembershipRole y UnitType`
- Descripción: Fase 3 del plan de mejoras - Value Objects

---

## Checklist

- [ ] Carpeta `internal/domain/valueobject/` creada
- [ ] `membership_role.go` implementado
- [ ] `unit_type.go` implementado
- [ ] Tests para ambos value objects
- [ ] Servicios refactorizados para usar value objects
- [ ] `go build` exitoso
- [ ] Tests pasan
- [ ] Documentación actualizada
- [ ] PR creado a dev
