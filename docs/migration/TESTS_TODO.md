# Tests Pendientes - Post Migración DDD

**Fecha:** 2025-11-22  
**Estado:** ✅ Build y tests existentes pasan - Faltan crear tests para lógica de negocio

---

## ✅ Estado Actual

**Build:** ✅ Sin errores  
**Tests Existentes:** ✅ Todos pasan (handlers)  
**Tests Eliminados:** 5 archivos (tests DDD obsoletos)

---

## 📋 Tests que Fueron Eliminados

Los siguientes tests fueron eliminados porque validaban lógica DDD en entities:

1. **internal/domain/service/academic_unit_service_test.go** ❌ ELIMINADO
   - Testeaba métodos de domain service (SetParent, validaciones en entity)
   
2. **internal/domain/service/membership_service_test.go** ❌ ELIMINADO
   - Testeaba lógica de domain service para memberships

3. **internal/application/service/hierarchy_service_test.go** ❌ ELIMINADO
   - Testeaba ValidateNoCircularReference (ya no existe)
   - Usaba value objects obsoletos

4. **test/integration/academic_unit_ltree_test.go** ❌ ELIMINADO
   - Test de integración con jerarquías ltree
   - Usaba entities DDD

5. **test/integration/integration_flows_test.go** ❌ ELIMINADO
   - Flujos de integración end-to-end con entities DDD

---

## 📝 Tests Nuevos Requeridos

### PRIORIDAD ALTA: Lógica de Negocio en Services

La lógica que antes estaba en entities ahora está en application services. Se requieren tests para:

#### 1. UserService
**Lógica a testear:**
- ✅ Validación: no permitir crear usuarios admin
- ✅ Validación: role válido
- ✅ UpdateUser: validar cambio de role
- ✅ UpdateUser: no permitir promover a admin
- ✅ UpdateUser: validar estado (activar/desactivar)
- ✅ CreateUser: email único

**Archivo:** `internal/application/service/user_service_test.go` (CREAR)

#### 2. SchoolService
**Lógica a testear:**
- ✅ Validación: nombre mínimo 3 caracteres
- ✅ Validación: código mínimo 3 caracteres
- ✅ CreateSchool: código único
- ✅ UpdateSchool: validaciones de campos

**Archivo:** `internal/application/service/school_service_test.go` (CREAR)

#### 3. GuardianService
**Lógica a testear:**
- ✅ Validación: relationship_type válido
- ✅ CreateRelation: no duplicar relación activa
- ✅ Validación: guardian no puede ser el estudiante

**Archivo:** `internal/application/service/guardian_service_test.go` (CREAR)

#### 4. UnitMembershipService
**Lógica a testear:**
- ✅ Validación: role válido
- ✅ CreateMembership: no duplicar membresía activa
- ✅ ExpireMembership: establecer withdrawn_at

**Archivo:** `internal/application/service/unit_membership_service_test.go` (CREAR)

#### 5. AcademicUnitService
**Lógica a testear:**
- ✅ Validación: displayName mínimo 3 caracteres
- ✅ CreateUnit: código único por escuela
- ✅ CreateUnit: validar padre existe
- ✅ CreateUnit: padre en misma escuela
- ✅ UpdateUnit: unidad no puede ser su propio padre

**Archivo:** `internal/application/service/academic_unit_service_test.go` (CREAR)

---

### PRIORIDAD MEDIA: Tests de Integración

#### 1. Flujos End-to-End
- Crear escuela → crear unidad → asignar membresía
- Crear usuario → crear relación guardian
- Validar soft delete funciona

**Archivo:** `test/integration/flows_test.go` (CREAR)

#### 2. Jerarquías AcademicUnit
- Crear árbol de unidades
- Validar GetHierarchyPath
- Validar BuildUnitTree

**Archivo:** `test/integration/academic_unit_hierarchy_test.go` (CREAR)

---

## 📊 Cobertura Actual

**Handlers:** ✅ Testeados (todos pasan)  
**Application Services:** ❌ Sin tests (lógica de negocio SIN tests)  
**Repositorios:** ❌ Sin tests  
**DTOs:** ✅ Validación en producción

---

## 🎯 Plan de Tests Sugerido

### FASE A: Tests Unitarios de Services (CRÍTICO)
- UserService (~2 horas)
- SchoolService (~1 hora)
- GuardianService (~1 hora)
- UnitMembershipService (~1.5 horas)
- AcademicUnitService (~2 horas)

**Total:** ~7-8 horas

### FASE B: Tests de Integración (RECOMENDADO)
- Flows básicos (~2 horas)
- Jerarquías (~2 horas)

**Total:** ~4 horas

---

## ✅ Criterio de Éxito

Para considerar la migración completa y segura:

- ✅ Build sin errores
- ✅ Tests existentes pasan
- ⏳ Tests de lógica de negocio en services (PENDIENTE)
- ⏳ Tests de integración básicos (PENDIENTE)

---

## 📝 Recomendación

**Opción A (Conservadora):** Crear tests de lógica de negocio ANTES de hacer merge a dev  
**Opción B (Pragmática):** Hacer merge ahora, crear tests en siguiente sprint  
**Opción C (Híbrida):** Crear tests críticos (User, School) ahora, resto después  

**Recomendación:** Opción C - Crear tests de User y School (3 horas), luego merge.
