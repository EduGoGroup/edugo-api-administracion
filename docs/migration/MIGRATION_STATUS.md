# Estado de Migración: DDD → Infrastructure Entities

**Fecha:** 2025-11-22  
**Objetivo:** Eliminar DDD y usar entidades de infrastructure directamente  
**Estado:** ✅ FASE 1 COMPLETADA - User migrado exitosamente

---

## Resumen Ejecutivo

Se ha completado exitosamente la migración del modelo **User** desde entidades de dominio DDD a entidades de infrastructure. El build compila correctamente y el patrón está establecido para migrar las entidades restantes.

---

## ✅ Completado

### 1. Análisis y Documentación
- ✅ Análisis completo de entidades actuales vs infrastructure
- ✅ Documento de campos pendientes (no se requieren DTOs temporales)
- ✅ Identificación de dependencias

### 2. Setup de Infrastructure
- ✅ Agregada dependencia `github.com/EduGoGroup/edugo-infrastructure/postgres@v0.10.0`
- ✅ Configurado replace local para desarrollo: `replace github.com/EduGoGroup/edugo-infrastructure/postgres => ../edugo-infrastructure/postgres`
- ✅ go.mod actualizado y limpio

### 3. Migración de User (COMPLETA)

#### Repositorio (`user_repository`)
**Antes:**
- Usaba `entity.User` (dominio DDD)
- Usaba `valueobject.UserID`, `valueobject.Email`
- Método `scanToEntity()` con lógica de conversión

**Después:**
- Usa `entities.User` (infrastructure)
- Usa `uuid.UUID`, `string` directamente
- Scan directo a la entidad de infrastructure
- Soporta soft delete con `deleted_at`
- Incluye campos `PasswordHash` y `EmailVerified`

**Archivo:** `internal/infrastructure/persistence/postgres/repository/user_repository_impl.go`

#### Interfaz de Repositorio
**Antes:**
```go
FindByID(ctx context.Context, id valueobject.UserID) (*entity.User, error)
FindByEmail(ctx context.Context, email valueobject.Email) (*entity.User, error)
```

**Después:**
```go
FindByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
FindByEmail(ctx context.Context, email string) (*entities.User, error)
```

**Archivo:** `internal/domain/repository/user_repository.go`

#### DTO
**Antes:**
```go
func ToUserResponse(user *entity.User) *UserResponse {
    return &UserResponse{
        ID: user.ID().String(),
        Email: user.Email().String(),
        FullName: user.FullName(), // Método de la entidad
        ...
    }
}
```

**Después:**
```go
func ToUserResponse(user *entities.User) *UserResponse {
    return &UserResponse{
        ID: user.ID.String(),
        Email: user.Email,
        FullName: user.FirstName + " " + user.LastName, // Lógica en DTO
        ...
    }
}
```

**Archivo:** `internal/application/dto/user_dto.go`

#### Service
**Lógica de Negocio Migrada del Entity al Service:**

1. **Validación de Role** (antes en `entity.NewUser`)
   ```go
   // Antes: entity.NewUser validaba internamente
   // Después: service valida explícitamente
   role := enum.SystemRole(req.Role)
   if !role.IsValid() {
       return nil, errors.NewValidationError("invalid role")
   }
   if role == enum.SystemRoleAdmin {
       return nil, errors.NewBusinessRuleError("cannot create admin users")
   }
   ```

2. **UpdateName** (antes `user.UpdateName()`)
   ```go
   // Antes: user.UpdateName(*req.FirstName, *req.LastName)
   // Después: validación y asignación directa en service
   if *req.FirstName == "" || *req.LastName == "" {
       return nil, errors.NewValidationError("names required")
   }
   user.FirstName = *req.FirstName
   user.LastName = *req.LastName
   ```

3. **ChangeRole** (antes `user.ChangeRole()`)
   ```go
   // Validaciones y reglas de negocio ahora en service
   if newRole == enum.SystemRoleAdmin {
       return nil, errors.NewBusinessRuleError("cannot promote to admin")
   }
   ```

4. **Activate/Deactivate** (antes `user.Activate()`, `user.Deactivate()`)
   ```go
   // Validación de estado y cambio directo en service
   if *req.IsActive && user.IsActive {
       return nil, errors.NewBusinessRuleError("already active")
   }
   user.IsActive = *req.IsActive
   ```

**Archivo:** `internal/application/service/user_service.go`

### 4. Verificación
- ✅ Build completo exitoso: `go build ./...`
- ✅ Binario principal compila: `go build -o /tmp/edugo-admin ./cmd`
- ✅ No hay errores de compilación relacionados con User

---

## 📋 Pendiente

### Entidades a Migrar

Faltan migrar las siguientes entidades siguiendo el mismo patrón de User:

1. **School** (`school_repository`, `school_dto`, services)
2. **Subject** (`subject_repository`)
3. **Unit** (`unit_repository`)
4. **GuardianRelation** (`guardian_repository`)
5. **UnitMembership → Membership** (`unit_membership_repository`)
6. **AcademicUnit** (`academic_unit_repository`, `academic_unit_dto`, `academic_unit_service`)

### Archivos Identificados que Requieren Migración

```
internal/infrastructure/persistence/postgres/repository/
├── school_repository_impl.go
├── subject_repository_impl.go
├── unit_repository_impl.go
├── guardian_repository_impl.go
├── unit_membership_repository_impl.go
└── academic_unit_repository_impl.go

internal/domain/repository/
├── school_repository.go
├── subject_repository.go
├── unit_repository.go
├── guardian_repository.go
├── unit_membership_repository.go
└── academic_unit_repository.go

internal/application/service/
├── academic_unit_service.go
└── hierarchy_service.go

internal/domain/service/
├── membership_service.go
└── academic_unit_service.go

internal/infrastructure/http/dto/
├── school_dto.go
└── academic_unit_dto.go

tests/
├── internal/domain/service/*_test.go
├── internal/application/service/*_test.go
└── test/integration/*_test.go
```

---

## 🎯 Patrón de Migración Establecido

### Para cada entidad, seguir estos pasos:

#### 1. Actualizar Interfaz de Repositorio
```go
// Antes
Create(ctx context.Context, entity *entity.X) error
FindByID(ctx context.Context, id valueobject.XID) (*entity.X, error)

// Después
Create(ctx context.Context, entity *entities.X) error
FindByID(ctx context.Context, id uuid.UUID) (*entities.X, error)
```

#### 2. Actualizar Implementación de Repositorio
- Cambiar import: `entity` → `entities` (infrastructure)
- Eliminar conversiones de value objects
- Usar campos públicos de la entidad: `user.ID` en lugar de `user.ID()`
- Scan directo a la entidad de infrastructure

#### 3. Actualizar DTOs
- Cambiar import de entidad
- Acceder campos directamente (públicos)
- Mover lógica simple (como FullName) al DTO

#### 4. Migrar Lógica de Negocio al Service
- Identificar métodos de la entidad DDD
- Mover validaciones al service
- Mover reglas de negocio al service
- Asignar campos directamente a la entidad

#### 5. Actualizar Tests
- Cambiar imports
- Actualizar mocks si es necesario
- Ajustar aserciones para campos públicos

---

## 📊 Estimación de Trabajo Restante

| Entidad | Complejidad | Archivos Afectados | Estimación |
|---------|-------------|-------------------|------------|
| School | Baja | 3 archivos | 20 min |
| Subject | Baja | 2 archivos | 15 min |
| Unit | Media | 3 archivos | 20 min |
| GuardianRelation | Baja | 2 archivos | 15 min |
| UnitMembership | Media | 3 archivos | 25 min |
| AcademicUnit | Alta | 6+ archivos | 40 min |
| Domain Services | Alta | 2 archivos | 30 min |
| Tests | Media | 5+ archivos | 45 min |
| **TOTAL** | - | **~30 archivos** | **~3-4 horas** |

---

## 🔧 Comandos Útiles

### Verificar build
```bash
go build ./...
```

### Compilar binario
```bash
go build -o /tmp/edugo-admin ./cmd
```

### Buscar archivos que usan entidades antiguas
```bash
grep -r "internal/domain/entity" --include="*.go" .
```

### Ejecutar tests
```bash
go test ./...
```

---

## ⚠️ Notas Importantes

### Cambios en el Schema de BD
- Infrastructure usa `deleted_at` para soft delete
- Admin anteriormente usaba `is_active` (boolean)
- **Decisión pendiente:** ¿Migrar a `deleted_at` o mantener `is_active`?

### Metadata
- Infrastructure usa `[]byte` (JSONB)
- Admin usaba `map[string]interface{}`
- Se requiere serialización/deserialización

### Value Objects
- **Eliminados:** UserID, SchoolID, SubjectID, UnitID, Email, etc.
- **Reemplazados por:** uuid.UUID, string
- **Enums:** Se mantienen de edugo-shared (SystemRole, etc.)

---

## 🚀 Próximos Pasos

1. **Migrar School**
   - Actualizar repositorio e interfaz
   - Actualizar DTO
   - Nota: `ContactEmail` → `Email`, `ContactPhone` → `Phone`

2. **Migrar Subject**
   - Repositorio simple, sin lógica compleja

3. **Migrar Unit**
   - Similar a School

4. **Migrar GuardianRelation**
   - Repositorio simple

5. **Migrar UnitMembership → Membership**
   - Nota: Cambio de nombre
   - `ValidFrom` → `EnrolledAt`
   - `ValidUntil` → `WithdrawnAt`

6. **Migrar AcademicUnit**
   - La más compleja
   - Tiene domain services
   - Tiene jerarquías (árbol)
   - Requiere migrar `academic_unit_service`

7. **Eliminar Entidades Antiguas**
   - Eliminar carpeta `internal/domain/entity/`
   - Eliminar carpeta `internal/domain/valueobject/`

8. **Ejecutar Tests Completos**
   - Corregir tests unitarios
   - Corregir tests de integración

9. **Commit**
   - Hacer commit atómico de la migración completa

---

## ✅ Criterios de Éxito

- [ ] Todas las entidades usan `entities` de infrastructure
- [ ] No hay imports de `internal/domain/entity`
- [ ] No hay imports de `internal/domain/valueobject` (excepto los que no son IDs)
- [ ] Build completo exitoso: `go build ./...`
- [ ] Tests pasan: `go test ./...`
- [ ] Carpetas `entity/` y `valueobject/` eliminadas
- [ ] Commit creado con mensaje descriptivo

---

## 📝 Conclusión

La migración de User establece el patrón correcto para eliminar DDD del proyecto. El enfoque es:

1. **Entidades simples** sin lógica de negocio (anémicas)
2. **Lógica en Services** (validaciones, reglas de negocio)
3. **Types primitivos** (uuid.UUID, string) en lugar de value objects
4. **Entidades de infrastructure** como fuente de verdad

El próximo paso es replicar este patrón en las entidades restantes.
