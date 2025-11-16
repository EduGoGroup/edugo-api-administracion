# Preguntas Sprint 01
## Q001: ¿Cómo prevenir ciclos en jerarquía?
**Decisión:** **Trigger SQL + Validación aplicación**
- CHECK (id != parent_id)
- Validar en app antes de INSERT

## Q002: ¿Materialized path o solo parent_id?
**Decisión:** **Ambos** - parent_id (principal) + path (cache)
- path facilita queries "todos los descendientes"
