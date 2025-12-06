# PRD - spec-03: API Administración - Jerarquía Académica

**Versión:** 1.0.0  
**Fecha:** 14 de Noviembre, 2025  
**Proyecto:** edugo-api-administracion  
**Repositorio:** /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-administracion

---

## 1. RESUMEN EJECUTIVO

### Visión
API administrativa que permite gestionar la jerarquía académica de instituciones educativas (escuelas → grados → secciones → clubes) y asignar usuarios (estudiantes, docentes) a unidades específicas.

### Problema
Sin jerarquía académica, el sistema no puede:
- Organizar estudiantes por grado/sección
- Asignar materiales a grupos específicos
- Gestionar permisos jerárquicos
- Generar reportes por unidad

### Solución
API REST con endpoints CRUD para:
1. **Schools** (Escuelas)
2. **Academic Units** (Unidades con árbol jerárquico: parent_id)
3. **Unit Memberships** (Asignación usuarios a unidades)

---

## 2. OBJETIVOS

### OBJ-1: Jerarquía Funcional
**Métrica:** CRUD completo de 3 entidades  
**Objetivo:** Schools + Units + Memberships operativos  
**Prioridad:** P0

### OBJ-2: Árbol Jerárquico
**Métrica:** Consulta recursiva de árbol <500ms  
**Objetivo:** Obtener jerarquía completa eficientemente  
**Prioridad:** P0

### OBJ-3: Prevención de Ciclos
**Métrica:** 0 jerarquías circulares en BD  
**Objetivo:** Trigger SQL previene parent_id inválidos  
**Prioridad:** P0

---

## 3. REQUERIMIENTOS FUNCIONALES

### RF-001: CRUD de Schools
Endpoints: POST /schools, GET /schools, GET /schools/:id, PUT, DELETE

### RF-002: CRUD de Academic Units
Endpoints con árbol jerárquico (parent_id)

### RF-003: CRUD de Memberships
Asignar estudiantes/docentes a unidades

### RF-004: Consulta de Árbol Jerárquico
GET /units/:id/tree (recursivo)

### RF-005: Validación de Permisos
Solo owners pueden modificar unidad

---

## 4. STACK TECNOLÓGICO

| Componente | Tecnología |
|------------|------------|
| Lenguaje | Go 1.21+ |
| Framework | Gin |
| ORM | GORM |
| BD | PostgreSQL 15+ |
| Auth | JWT (reutilizar shared) |

---

**Generado con:** Claude Code
