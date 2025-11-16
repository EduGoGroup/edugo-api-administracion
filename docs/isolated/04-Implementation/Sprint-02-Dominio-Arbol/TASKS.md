# Tareas Sprint 02
## TASK-02-001: Entity AcademicUnit con Tree Methods
**Estimación:** 4h

```go
// internal/domain/entities/academic_unit.go
type AcademicUnit struct {
    ID           uuid.UUID
    SchoolID     uuid.UUID
    ParentUnitID *uuid.UUID  // nil = raíz
    UnitType     string      // grade|section|club
    DisplayName  string
    Code         string
    Depth        int
    Children     []*AcademicUnit  // Para árbol en memoria
}

func (u *AcademicUnit) IsRoot() bool {
    return u.ParentUnitID == nil
}

func (u *AcademicUnit) AddChild(child *AcademicUnit) error {
    if child.ParentUnitID == nil || *child.ParentUnitID != u.ID {
        return errors.New("invalid parent")
    }
    u.Children = append(u.Children, child)
    return nil
}

func (u *AcademicUnit) GetAllDescendants() []*AcademicUnit {
    // Traversal recursivo
}
```

**Tiempo:** 4h
