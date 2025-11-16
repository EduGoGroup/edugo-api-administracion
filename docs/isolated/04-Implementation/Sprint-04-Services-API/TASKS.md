# Tareas Sprint 04
## TASK-04-001: HierarchyService
**Estimación:** 5h
```go
func (s *HierarchyService) CreateUnit(ctx, parentID, schoolID, unitType, name, code) (*Unit, error)
func (s *HierarchyService) GetUnitTree(ctx, unitID) (*UnitTreeDTO, error)
func (s *HierarchyService) ValidateNoCircularReference(ctx, unitID, newParentID) error
```
## TASK-04-002: CRUD Endpoints (10 endpoints)
POST/GET/PUT/DELETE /schools, /units, /units/:id/members
**Tiempo:** 10h
