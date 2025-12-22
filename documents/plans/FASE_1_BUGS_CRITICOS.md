# 🔴 FASE 1: Corrección de Bugs Críticos

**Prioridad**: Alta  
**Estimación**: 4.5 horas  
**Rama**: `fix/fase-1-bugs-criticos`

---

## Preparación Git

```bash
git checkout dev
git pull origin dev
git checkout -b fix/fase-1-bugs-criticos
```

---

## 1.1 Implementar `ListMembershipsByRole` correctamente

### Ubicación
```
internal/application/service/unit_membership_service.go:174-177
```

### Problema
La función ignora completamente el parámetro `role` y retorna todas las membresías sin filtrar.

### Código Actual (Problemático)
```go
func (s *unitMembershipService) ListMembershipsByRole(ctx context.Context, unitID string, role string, activeOnly bool) ([]dto.MembershipResponse, error) {
    // Implementación simplificada - IGNORA role
    return s.ListMembershipsByUnit(ctx, unitID, activeOnly)
}
```

### Tareas
1. Agregar método `FindByUnitAndRole(ctx, unitID, role, activeOnly)` al repositorio
2. Implementar la lógica de filtrado en el servicio
3. Agregar tests unitarios

### Código Esperado
```go
func (s *unitMembershipService) ListMembershipsByRole(ctx context.Context, unitID string, role string, activeOnly bool) ([]dto.MembershipResponse, error) {
    uid, err := uuid.Parse(unitID)
    if err != nil {
        return nil, errors.NewValidationError("invalid unit ID")
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

### Esfuerzo
2.5 horas

---

## 1.2 Implementar parámetro `activeOnly` en queries

### Ubicación
```
internal/application/service/unit_membership_service.go:138-176
```

### Problema
El parámetro `activeOnly` se recibe pero nunca se usa en las queries.

### Código Actual (Problemático)
```go
func (s *unitMembershipService) ListMembershipsByUnit(ctx context.Context, unitID string, activeOnly bool) ([]dto.MembershipResponse, error) {
    // ...
    // activeOnly NO SE USA - siempre retorna todas
    memberships, err := s.membershipRepo.FindByUnit(ctx, uid)
    // ...
}
```

### Tareas
1. Actualizar interfaz `FindByUnit` en el repositorio para aceptar `activeOnly`
2. Agregar condición WHERE para filtrar membresías activas
3. Actualizar tests

### Código Esperado en Repositorio
```go
func (r *unitMembershipRepository) FindByUnit(ctx context.Context, unitID uuid.UUID, activeOnly bool) ([]*entities.Membership, error) {
    query := r.db.Where("academic_unit_id = ?", unitID)
    if activeOnly {
        query = query.Where("is_active = ?", true).Where("withdrawn_at IS NULL")
    }
    
    var memberships []*entities.Membership
    if err := query.Find(&memberships).Error; err != nil {
        return nil, err
    }
    return memberships, nil
}
```

### Esfuerzo
2 horas

---

## Documentación a Actualizar

Al completar esta fase, actualizar:

- `documents/improvements/REFACTORING.md` - Eliminar secciones 3 y 4 (bugs corregidos)
- `documents/improvements/TODO_LIST.md` - Eliminar TODOs resueltos
- `documents/improvements/README.md` - Actualizar estado de items

---

## Finalización

```bash
git add .
git commit -m "fix: implementar filtros role y activeOnly en membresías"
git push origin fix/fase-1-bugs-criticos
```

### Crear PR a dev con:
- Título: `fix: implementar filtros role y activeOnly en membresías`
- Descripción: Fase 1 del plan de mejoras - Corrección de bugs críticos

---

## Checklist

- [ ] `FindByUnitAndRole` implementado en repositorio
- [ ] `FindByUnit` actualizado con parámetro `activeOnly`
- [ ] `ListMembershipsByRole` usa el nuevo método del repositorio
- [ ] `ListMembershipsByUnit` pasa `activeOnly` al repositorio
- [ ] Tests unitarios agregados/actualizados
- [ ] Documentación actualizada
- [ ] PR creado a dev
