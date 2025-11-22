# Campos Pendientes en edugo-infrastructure

> **Fecha:** 2025-11-22  
> **Proyecto:** edugo-api-administracion  
> **Propósito:** Documentar campos que necesita admin pero que infrastructure aún no tiene

---

## Resumen

Este documento lista los campos que el proyecto `edugo-api-administracion` necesita pero que actualmente NO están disponibles en las entidades de `edugo-infrastructure`.

Para cada campo faltante se debe:
1. Crear una tarea en el proyecto `edugo-infrastructure` para agregar el campo
2. Crear un DTO temporal en `admin` con el campo
3. En el mapper, usar un valor fijo + TODO hasta que infrastructure libere la versión con el campo

---

## Estado Actual de Infrastructure

- **Versión en main:** v0.3.0
- **Última versión postgres:** postgres/v0.10.0
- **Última versión mongodb:** mongodb/v0.10.0

---

## Campos Faltantes por Entidad

### 1. User (users table)

#### Campos que admin NO necesita (existen en infra):
- ✅ `PasswordHash` - No se maneja en admin, solo en autenticación
- ✅ `EmailVerified` - No se usa en admin actualmente
- ✅ `DeletedAt` - Soft delete (admin usa `IsActive`)

#### Campos que admin necesita pero infra NO tiene:
Ninguno. **User está completo en infrastructure.**

#### Acción:
- ✅ Usar directamente `infrastructure.User`
- No se requiere DTO temporal

---

### 2. School (schools table)

#### Campos que admin NO necesita (existen en infra):
- ✅ `City` - No se usa en admin actualmente
- ✅ `Country` - No se usa en admin actualmente
- ✅ `SubscriptionTier` - No se maneja en admin
- ✅ `MaxTeachers` - No se maneja en admin
- ✅ `MaxStudents` - No se maneja en admin
- ✅ `IsActive` - Admin no tiene este campo actualmente
- ✅ `DeletedAt` - Soft delete

#### Campos que admin necesita pero infra NO tiene:
- ❌ `ContactEmail` como campo separado de `Email`
- ❌ `ContactPhone` como campo separado de `Phone`

**Nota:** Infrastructure tiene `Email` y `Phone`, pero admin los llamaba `ContactEmail` y `ContactPhone`. 

#### Acción:
- ✅ Usar `Email` y `Phone` de infrastructure directamente
- No se requiere DTO temporal
- Renombrar referencias en el código de admin

---

### 3. Subject (subjects table)

#### Campos que admin NO necesita (existen en infra):
Ninguno.

#### Campos que admin necesita pero infra NO tiene:
Ninguno. **Subject está completo en infrastructure.**

#### Acción:
- ✅ Usar directamente `infrastructure.Subject`
- No se requiere DTO temporal

---

### 4. Unit (units table)

#### Campos que admin NO necesita (existen en infra):
Ninguno.

#### Campos que admin necesita pero infra NO tiene:
Ninguno. **Unit está completo en infrastructure.**

#### Acción:
- ✅ Usar directamente `infrastructure.Unit`
- No se requiere DTO temporal

---

### 5. GuardianRelation (guardian_relations table)

#### Campos que admin NO necesita (existen en infra):
Ninguno.

#### Campos que admin necesita pero infra NO tiene:
Ninguno. **GuardianRelation está completo en infrastructure.**

#### Acción:
- ✅ Usar directamente `infrastructure.GuardianRelation`
- No se requiere DTO temporal

---

### 6. UnitMembership → Membership (memberships table)

**Nota:** Infrastructure la llama `Membership`, admin la llamaba `UnitMembership`.

#### Campos que admin necesita pero infra NO tiene:
- ❌ `ValidFrom` - Admin lo usaba para fecha de inicio
- ❌ `ValidUntil` - Admin lo usaba para fecha de fin

**Infrastructure tiene:**
- ✅ `EnrolledAt` - Similar a ValidFrom
- ✅ `WithdrawnAt` - Similar a ValidUntil

#### Acción:
- ✅ Usar `EnrolledAt` en lugar de `ValidFrom`
- ✅ Usar `WithdrawnAt` en lugar de `ValidUntil`
- No se requiere DTO temporal
- Solo renombrar en el código

---

### 7. AcademicUnit (academic_units table)

#### Campos que admin NO necesita (existen en infra):
- ✅ `Level` - No se usa en admin
- ✅ `AcademicYear` - No se usa en admin
- ✅ `IsActive` - Admin no tiene este campo

#### Campos que admin necesita pero infra NO tiene:
- ❌ `DisplayName` como campo separado de `Name`

**Infrastructure tiene:**
- ✅ `Name` - Que sería equivalente a `DisplayName`
- ✅ `Code`
- ✅ `Type`

#### Acción:
- ✅ Usar `Name` de infrastructure en lugar de `DisplayName`
- No se requiere DTO temporal
- Renombrar referencias en el código

---

## Campos Especiales: Value Objects

Admin usa **value objects** para IDs y tipos enumerados. Infrastructure usa **tipos primitivos**.

### IDs (UUIDs):

| Admin Value Object | Infrastructure Type | Acción |
|-------------------|---------------------|--------|
| `valueobject.UserID` | `uuid.UUID` | ✅ Eliminar value object, usar uuid.UUID |
| `valueobject.SchoolID` | `uuid.UUID` | ✅ Eliminar value object, usar uuid.UUID |
| `valueobject.SubjectID` | `uuid.UUID` | ✅ Eliminar value object, usar uuid.UUID |
| `valueobject.UnitID` | `uuid.UUID` | ✅ Eliminar value object, usar uuid.UUID |
| `valueobject.GuardianID` | `uuid.UUID` | ✅ Eliminar value object, usar uuid.UUID |
| `valueobject.StudentID` | `uuid.UUID` | ✅ Eliminar value object, usar uuid.UUID |
| `valueobject.MembershipID` | `uuid.UUID` | ✅ Eliminar value object, usar uuid.UUID |

### Enums:

| Admin Value Object | Infrastructure Type | Acción |
|-------------------|---------------------|--------|
| `enum.SystemRole` | `string` | ✅ Mantener enum de shared, validar en services |
| `valueobject.UnitType` | `string` | ✅ Convertir a constantes o enum, validar en services |
| `valueobject.RelationshipType` | `string` | ✅ Convertir a constantes o enum, validar en services |
| `valueobject.MembershipRole` | `string` | ✅ Convertir a constantes o enum, validar en services |

### Email:

| Admin Value Object | Infrastructure Type | Acción |
|-------------------|---------------------|--------|
| `valueobject.Email` | `string` | ✅ Eliminar value object, usar string, validar en services |

---

## Resumen de Acciones

### ✅ NO SE REQUIEREN DTOs TEMPORALES

Todos los campos necesarios están disponibles en infrastructure. Solo se requieren cambios de nombres:

1. **School:**
   - `ContactEmail` → `Email`
   - `ContactPhone` → `Phone`

2. **AcademicUnit:**
   - `DisplayName` → `Name`

3. **UnitMembership → Membership:**
   - `ValidFrom` → `EnrolledAt`
   - `ValidUntil` → `WithdrawnAt`

### ✅ NO SE REQUIEREN TAREAS EN INFRASTRUCTURE

Infrastructure tiene todos los campos necesarios. La migración es directa.

---

## Plan de Migración

### FASE 1: Preparación
1. Agregar dependencia `edugo-infrastructure` v0.3.0 al go.mod
2. Verificar compatibilidad con shared v0.7.0

### FASE 2: Eliminación de Value Objects
1. Reemplazar todos los value objects de IDs por `uuid.UUID`
2. Reemplazar `valueobject.Email` por `string`
3. Convertir enums de value objects a constantes o mantener los de shared

### FASE 3: Reemplazo de Entidades
1. Reemplazar imports de `entity` por imports de `infrastructure.entities`
2. Actualizar nombres de campos según tabla arriba
3. Actualizar repositorios para usar entidades de infrastructure

### FASE 4: Migración de Lógica de Negocio
1. Mover métodos de validación de entities a services/validators
2. Mover métodos de negocio de entities a domain services
3. Eliminar constructores `New*()` y `Reconstruct*()`

### FASE 5: Testing
1. Ejecutar tests unitarios
2. Ejecutar tests de integración
3. Verificar build completo

---

## Notas Adicionales

- ⚠️ **Soft Delete:** Infrastructure usa `DeletedAt`, admin usaba `IsActive`. Decidir estrategia.
- ⚠️ **Metadata:** Infrastructure usa `[]byte` (JSONB), admin usaba `map[string]interface{}`. Requiere helpers de serialización.
- ⚠️ **Campos Null:** Infrastructure usa punteros para NULL (`*string`, `*time.Time`), admin usaba valores directos. Actualizar manejo de nulls.

---

## Siguiente Paso

✅ Proceder con la migración directa SIN necesidad de DTOs temporales.
