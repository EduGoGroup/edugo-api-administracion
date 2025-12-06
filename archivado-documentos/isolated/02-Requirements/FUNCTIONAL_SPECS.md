# Especificaciones Funcionales - spec-03

## RF-001: Crear Escuela
**Prioridad:** MUST  
**Endpoint:** POST /v1/schools

### Criterios
- Validar nombre único
- Generar código automáticamente
- Solo admin puede crear

---

## RF-002: Crear Unidad Académica
**Prioridad:** MUST  
**Endpoint:** POST /v1/units

### Criterios
- Validar parent_id existe
- Prevenir ciclos (trigger SQL)
- Código único por escuela
- Tipos: grade, section, club

---

## RF-003: Obtener Árbol Jerárquico
**Prioridad:** MUST  
**Endpoint:** GET /v1/units/:id/tree

### Criterios
- Query recursiva (WITH RECURSIVE)
- Retornar JSON anidado
- Profundidad máxima: 5 niveles

---

## RF-004: Asignar Usuario a Unidad
**Prioridad:** MUST  
**Endpoint:** POST /v1/units/:id/members

### Criterios
- Validar usuario existe
- Validar unidad existe
- Roles: student, teacher, owner

---

**Total:** 12 RF (8 MUST, 4 SHOULD)
