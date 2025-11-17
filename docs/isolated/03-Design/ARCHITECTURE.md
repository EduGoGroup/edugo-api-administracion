# Arquitectura - spec-03
## Clean Architecture con Árbol Jerárquico
```
API REST (Gin) → Services → Repositories → PostgreSQL
                                          └→ Queries Recursivas (WITH RECURSIVE)
```
## Componentes
- Entities: School, AcademicUnit (con parent_id), UnitMembership
- Tree Traversal: Métodos para navegar jerarquía
- Repositories: Queries recursivas para árbol
