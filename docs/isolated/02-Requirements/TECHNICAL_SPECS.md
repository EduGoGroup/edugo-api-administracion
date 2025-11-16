# Especificaciones Técnicas - spec-03

## Stack
- Go 1.21+, Gin, GORM, PostgreSQL 15+

## Performance
- GET /units/:id/tree: <500ms
- POST /units: <200ms

## Constraints
- Trigger SQL previene ciclos
- Unique index en (school_id, code)
- Foreign keys con ON DELETE CASCADE

## Escalabilidad
- Índices en parent_id para queries recursivas
- Cache de árboles frecuentes (Post-MVP)
