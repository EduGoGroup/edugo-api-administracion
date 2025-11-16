# Tareas Sprint 03
## TASK-03-001: UnitRepository con Query Recursiva
**Estimación:** 4h
```go
// GetTree obtiene árbol jerárquico con WITH RECURSIVE
func (r *UnitRepository) GetTree(ctx context.Context, unitID uuid.UUID) (*entities.AcademicUnit, error) {
    query := `
        WITH RECURSIVE unit_tree AS (
            SELECT * FROM academic_unit WHERE id = $1
            UNION ALL
            SELECT u.* FROM academic_unit u
            INNER JOIN unit_tree t ON u.parent_unit_id = t.id
        )
        SELECT * FROM unit_tree ORDER BY depth
    `
    // Construir árbol en memoria
}
```
**Tiempo:** 4h
