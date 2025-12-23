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

---

### 3. ~~legacy_handlers.go - Endpoints a remover~~ ✅ RESUELTO

**Ubicación**: `cmd/legacy_handlers.go:9-17`

**Estado**: ✅ Completado en Fase 2
**Fecha**: 2025-12-23

El archivo `legacy_handlers.go` no existe (fue eliminado previamente).
Los tipos de respuesta han sido centralizados en `internal/infrastructure/http/dto/response.go`.

---

## Resumen por Archivo

| Archivo | TODOs | Prioridad General |
|---------|-------|-------------------|
| `school_service.go` | 5 | Media |
| ~~`legacy_handlers.go`~~ | ~~1~~ ✅ | ~~Alta~~ Resuelto |
| **Total** | **5** | - |

---

## TODOs por Prioridad

### 🔴 Alta Prioridad
| TODO | Archivo | Línea | Descripción |
|------|---------|-------|-------------|
| ~~Código deprecated~~ ✅ | ~~`legacy_handlers.go`~~ | ~~9~~ | ~~Eliminar en v0.6.0~~ Completado Fase 2 |

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
- [x] ~~Eliminar `legacy_handlers.go`~~ ✅ Fase 2 (2025-12-23)
- [x] ~~Implementar `ListMembershipsByRole` correctamente~~ ✅ Fase 1 (2025-12-22)

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
| 2025-12-23 | Código deprecated | `legacy_handlers.go` | Pendiente | Centralización de response types en dto/response.go |
| 2025-12-22 | ListMembershipsByRole no filtra | `unit_membership_service.go` | #57 | Implementado FindByUnitAndRole en repositorio |
| 2025-11-20 | Auth centralizado | `container.go` | #45 | Migrado a shared/auth |
| 2025-11-15 | Mock repositories | `factory.go` | #42 | Implementado factory pattern |
