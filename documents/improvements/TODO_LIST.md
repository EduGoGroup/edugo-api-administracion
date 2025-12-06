# 📝 Lista de TODOs

> Comentarios TODO encontrados en el código que requieren acción

---

## Búsqueda de TODOs

```bash
# Comando para encontrar todos los TODOs
grep -rn "TODO\|FIXME\|HACK\|XXX" --include="*.go" .
```

---

## TODOs Encontrados

### 1. school_service.go - Campos faltantes en DTO

**Ubicación**: `internal/application/service/school_service.go:73-74`

```go
school := &entities.School{
    // ...
    City:             nil,  // TODO: agregar cuando se agregue al DTO
    Country:          "CO", // TODO: valor por defecto, agregar al DTO
    // ...
}
```

**Acción Requerida**:
1. Agregar campo `City` a `CreateSchoolRequest`
2. Agregar campo `Country` a `CreateSchoolRequest` con valor default
3. Actualizar documentación Swagger

**Prioridad**: Media
**Esfuerzo**: 1 hora

---

### 2. school_service.go - Valores por defecto de suscripción

**Ubicación**: `internal/application/service/school_service.go:79-81`

```go
school := &entities.School{
    // ...
    SubscriptionTier: "free", // TODO: valor por defecto
    MaxTeachers:      50,     // TODO: valor por defecto
    MaxStudents:      500,    // TODO: valor por defecto
    // ...
}
```

**Acción Requerida**:
1. Mover valores a configuración (`config.yaml`)
2. Opcionalmente permitir override en el DTO
3. Considerar diferentes tiers de suscripción

**Prioridad**: Media
**Esfuerzo**: 2 horas

---

### 3. unit_membership_service.go - Implementación simplificada

**Ubicación**: `internal/application/service/unit_membership_service.go:175-176`

```go
func (s *unitMembershipService) ListMembershipsByRole(...) ([]dto.MembershipResponse, error) {
    // Implementación simplificada
    return s.ListMembershipsByUnit(ctx, unitID, activeOnly)
}
```

**Acción Requerida**:
1. Implementar filtro por rol real
2. Agregar método `FindByUnitAndRole` al repositorio
3. Agregar tests

**Prioridad**: Alta (es un bug)
**Esfuerzo**: 2 horas

---

### 4. legacy_handlers.go - Endpoints a remover

**Ubicación**: `cmd/legacy_handlers.go:9-17`

```go
// ==================== LEGACY HANDLERS ====================
//
// DEPRECATED: Estos endpoints están deprecated y serán removidos en v0.6.0
// NO implementan lógica real, solo retornan datos mock para compatibilidad.
//
// Si necesitas estas funcionalidades, deberás:
// 1. Implementar los handlers reales en internal/interface/http/handler/
// 2. Crear los services correspondientes en internal/application/service/
// 3. Actualizar la documentación Swagger
```

**Acción Requerida**:
1. Eliminar archivo completo en v0.6.0
2. Verificar que no hay clientes usando estos endpoints
3. Actualizar Swagger

**Prioridad**: Alta (deadline v0.6.0)
**Esfuerzo**: 30 minutos

---

## Resumen por Archivo

| Archivo | TODOs | Prioridad General |
|---------|-------|-------------------|
| `school_service.go` | 5 | Media |
| `unit_membership_service.go` | 1 | Alta |
| `legacy_handlers.go` | 1 | Alta |
| **Total** | **7** | - |

---

## TODOs por Prioridad

### 🔴 Alta Prioridad
| TODO | Archivo | Línea | Descripción |
|------|---------|-------|-------------|
| Implementación incompleta | `unit_membership_service.go` | 175 | ListMembershipsByRole no filtra |
| Código deprecated | `legacy_handlers.go` | 9 | Eliminar en v0.6.0 |

### 🟡 Media Prioridad
| TODO | Archivo | Línea | Descripción |
|------|---------|-------|-------------|
| Campo City | `school_service.go` | 73 | Agregar al DTO |
| Campo Country | `school_service.go` | 74 | Agregar al DTO con default |
| SubscriptionTier | `school_service.go` | 79 | Mover a config |
| MaxTeachers | `school_service.go` | 80 | Mover a config |
| MaxStudents | `school_service.go` | 81 | Mover a config |

### 🟢 Baja Prioridad
(Ninguno identificado actualmente)

---

## Proceso para Resolver TODOs

```
1. Crear issue/ticket con referencia a este documento
2. Asignar a sprint según prioridad
3. Implementar solución con tests
4. Actualizar documentación
5. Remover TODO del código
6. Actualizar este documento
```

---

## Checklist de Resolución

### Sprint Actual (v0.6.0)
- [ ] Eliminar `legacy_handlers.go`
- [ ] Implementar `ListMembershipsByRole` correctamente

### Próximo Sprint
- [ ] Agregar campos City y Country al DTO de School
- [ ] Mover valores de suscripción a configuración

### Backlog
- [ ] Definir tiers de suscripción formalmente
- [ ] Crear documentación de límites por tier

---

## Scripts de Mantenimiento

### Buscar nuevos TODOs
```bash
#!/bin/bash
# scripts/find-todos.sh

echo "=== TODOs en el código ==="
grep -rn "TODO" --include="*.go" . | grep -v "_test.go" | grep -v "vendor/"

echo ""
echo "=== FIXMEs en el código ==="
grep -rn "FIXME" --include="*.go" . | grep -v "_test.go" | grep -v "vendor/"

echo ""
echo "=== HACKs en el código ==="
grep -rn "HACK" --include="*.go" . | grep -v "_test.go" | grep -v "vendor/"
```

### Contar TODOs
```bash
#!/bin/bash
# scripts/count-todos.sh

echo "TODOs: $(grep -rn 'TODO' --include='*.go' . | grep -v '_test.go' | wc -l)"
echo "FIXMEs: $(grep -rn 'FIXME' --include='*.go' . | grep -v '_test.go' | wc -l)"
echo "HACKs: $(grep -rn 'HACK' --include='*.go' . | grep -v '_test.go' | wc -l)"
```

---

## Historial de TODOs Resueltos

| Fecha | TODO | Archivo | PR | Notas |
|-------|------|---------|-----|-------|
| 2025-11-20 | Auth centralizado | `container.go` | #45 | Migrado a shared/auth |
| 2025-11-15 | Mock repositories | `factory.go` | #42 | Implementado factory pattern |
