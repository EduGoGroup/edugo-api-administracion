package repository

import (
	"context"
	"sync"
	"time"

	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/repository"
	mockData "github.com/EduGoGroup/edugo-api-administracion/internal/infrastructure/persistence/mock/data"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-shared/common/errors"
	"github.com/google/uuid"
)

// MockAcademicUnitRepository implementa repository.AcademicUnitRepository para testing
type MockAcademicUnitRepository struct {
	mu            sync.RWMutex
	academicUnits map[uuid.UUID]*entities.AcademicUnit
}

// NewMockAcademicUnitRepository crea una nueva instancia de MockAcademicUnitRepository
func NewMockAcademicUnitRepository() repository.AcademicUnitRepository {
	repo := &MockAcademicUnitRepository{
		academicUnits: make(map[uuid.UUID]*entities.AcademicUnit),
	}

	// Pre-cargar datos desde mockData
	for _, unit := range mockData.GetAcademicUnits() {
		unitCopy := *unit
		repo.academicUnits[unit.ID] = &unitCopy
	}

	return repo
}

// Create crea una nueva unidad académica
func (r *MockAcademicUnitRepository) Create(ctx context.Context, unit *entities.AcademicUnit) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Validar duplicado por SchoolID + Code
	for _, u := range r.academicUnits {
		if u.DeletedAt != nil {
			continue
		}
		if u.SchoolID == unit.SchoolID && u.Code == unit.Code {
			return errors.NewConflictError("academic unit with this code already exists in the school")
		}
	}

	// Validar que el padre existe y está en la misma escuela
	if unit.ParentUnitID != nil {
		parent, exists := r.academicUnits[*unit.ParentUnitID]
		if !exists || parent.DeletedAt != nil {
			return errors.NewNotFoundError("parent unit not found")
		}
		if parent.SchoolID != unit.SchoolID {
			return errors.NewValidationError("parent unit must be in the same school")
		}

		// Validar que no se crea un ciclo
		if r.wouldCreateCycle(unit.ID, *unit.ParentUnitID) {
			return errors.NewValidationError("creating this parent relationship would create a cycle in the hierarchy")
		}
	}

	// Generar ID si no existe
	if unit.ID == uuid.Nil {
		unit.ID = uuid.New()
	}

	// Establecer timestamps
	now := time.Now()
	unit.CreatedAt = now
	unit.UpdatedAt = now

	// Almacenar copia
	unitCopy := *unit
	r.academicUnits[unit.ID] = &unitCopy

	return nil
}

// FindByID busca una unidad por ID
func (r *MockAcademicUnitRepository) FindByID(ctx context.Context, id uuid.UUID, includeDeleted bool) (*entities.AcademicUnit, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	unit, exists := r.academicUnits[id]
	if !exists {
		return nil, errors.NewNotFoundError("academic unit not found")
	}

	if !includeDeleted && unit.DeletedAt != nil {
		return nil, errors.NewNotFoundError("academic unit not found")
	}

	// Retornar copia
	unitCopy := *unit
	return &unitCopy, nil
}

// FindBySchoolIDAndCode busca una unidad por escuela y código
func (r *MockAcademicUnitRepository) FindBySchoolIDAndCode(ctx context.Context, schoolID uuid.UUID, code string) (*entities.AcademicUnit, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, unit := range r.academicUnits {
		if unit.DeletedAt != nil {
			continue
		}
		if unit.SchoolID == schoolID && unit.Code == code {
			// Retornar copia
			unitCopy := *unit
			return &unitCopy, nil
		}
	}

	return nil, errors.NewNotFoundError("academic unit not found")
}

// FindBySchoolID lista todas las unidades de una escuela
func (r *MockAcademicUnitRepository) FindBySchoolID(ctx context.Context, schoolID uuid.UUID, includeDeleted bool) ([]*entities.AcademicUnit, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*entities.AcademicUnit

	for _, unit := range r.academicUnits {
		if unit.SchoolID != schoolID {
			continue
		}

		if !includeDeleted && unit.DeletedAt != nil {
			continue
		}

		// Agregar copia
		unitCopy := *unit
		result = append(result, &unitCopy)
	}

	return result, nil
}

// FindByParentID lista unidades hijas de una unidad padre
func (r *MockAcademicUnitRepository) FindByParentID(ctx context.Context, parentID uuid.UUID, includeDeleted bool) ([]*entities.AcademicUnit, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*entities.AcademicUnit

	for _, unit := range r.academicUnits {
		if unit.ParentUnitID == nil || *unit.ParentUnitID != parentID {
			continue
		}

		if !includeDeleted && unit.DeletedAt != nil {
			continue
		}

		// Agregar copia
		unitCopy := *unit
		result = append(result, &unitCopy)
	}

	return result, nil
}

// FindRootUnits lista unidades raíz (sin padre) de una escuela
func (r *MockAcademicUnitRepository) FindRootUnits(ctx context.Context, schoolID uuid.UUID) ([]*entities.AcademicUnit, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*entities.AcademicUnit

	for _, unit := range r.academicUnits {
		if unit.SchoolID != schoolID {
			continue
		}

		if unit.DeletedAt != nil {
			continue
		}

		// Solo incluir unidades sin padre
		if unit.ParentUnitID == nil {
			// Agregar copia
			unitCopy := *unit
			result = append(result, &unitCopy)
		}
	}

	return result, nil
}

// FindByType lista unidades por tipo
func (r *MockAcademicUnitRepository) FindByType(ctx context.Context, schoolID uuid.UUID, unitType string, includeDeleted bool) ([]*entities.AcademicUnit, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*entities.AcademicUnit

	for _, unit := range r.academicUnits {
		if unit.SchoolID != schoolID {
			continue
		}

		if unit.Type != unitType {
			continue
		}

		if !includeDeleted && unit.DeletedAt != nil {
			continue
		}

		// Agregar copia
		unitCopy := *unit
		result = append(result, &unitCopy)
	}

	return result, nil
}

// Update actualiza una unidad existente
func (r *MockAcademicUnitRepository) Update(ctx context.Context, unit *entities.AcademicUnit) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Verificar que la unidad existe
	existing, exists := r.academicUnits[unit.ID]
	if !exists || existing.DeletedAt != nil {
		return errors.NewNotFoundError("academic unit not found")
	}

	// Validar duplicado por SchoolID + Code (excepto el mismo registro)
	for _, u := range r.academicUnits {
		if u.DeletedAt != nil {
			continue
		}
		if u.ID != unit.ID && u.SchoolID == unit.SchoolID && u.Code == unit.Code {
			return errors.NewConflictError("academic unit with this code already exists in the school")
		}
	}

	// Validar que el padre existe y está en la misma escuela
	if unit.ParentUnitID != nil {
		parent, exists := r.academicUnits[*unit.ParentUnitID]
		if !exists || parent.DeletedAt != nil {
			return errors.NewNotFoundError("parent unit not found")
		}
		if parent.SchoolID != unit.SchoolID {
			return errors.NewValidationError("parent unit must be in the same school")
		}

		// Validar que no se crea un ciclo (solo si el padre cambió)
		if existing.ParentUnitID == nil || *existing.ParentUnitID != *unit.ParentUnitID {
			if r.wouldCreateCycle(unit.ID, *unit.ParentUnitID) {
				return errors.NewValidationError("updating this parent relationship would create a cycle in the hierarchy")
			}
		}
	}

	// Actualizar timestamp
	unit.UpdatedAt = time.Now()

	// Preservar CreatedAt original
	unit.CreatedAt = existing.CreatedAt

	// Almacenar copia actualizada
	unitCopy := *unit
	r.academicUnits[unit.ID] = &unitCopy

	return nil
}

// SoftDelete marca una unidad como eliminada
func (r *MockAcademicUnitRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	unit, exists := r.academicUnits[id]
	if !exists || unit.DeletedAt != nil {
		return errors.NewNotFoundError("academic unit not found")
	}

	// Soft delete
	now := time.Now()
	unit.DeletedAt = &now
	unit.UpdatedAt = now

	return nil
}

// Restore restaura una unidad eliminada
func (r *MockAcademicUnitRepository) Restore(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	unit, exists := r.academicUnits[id]
	if !exists {
		return errors.NewNotFoundError("academic unit not found")
	}

	if unit.DeletedAt == nil {
		return errors.NewValidationError("academic unit is not deleted")
	}

	// Restore
	unit.DeletedAt = nil
	unit.UpdatedAt = time.Now()

	return nil
}

// HardDelete elimina permanentemente una unidad
func (r *MockAcademicUnitRepository) HardDelete(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.academicUnits[id]; !exists {
		return errors.NewNotFoundError("academic unit not found")
	}

	// Hard delete - eliminar del mapa
	delete(r.academicUnits, id)

	return nil
}

// GetHierarchyPath obtiene el path jerárquico desde raíz hasta la unidad
func (r *MockAcademicUnitRepository) GetHierarchyPath(ctx context.Context, id uuid.UUID) ([]*entities.AcademicUnit, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var path []*entities.AcademicUnit

	currentUnit, exists := r.academicUnits[id]
	if !exists {
		return nil, errors.NewNotFoundError("academic unit not found")
	}

	// Protección contra ciclos infinitos
	visited := make(map[uuid.UUID]bool)

	// Construir el path desde la unidad actual hacia arriba
	for currentUnit != nil {
		// Detectar ciclo
		if visited[currentUnit.ID] {
			return nil, errors.NewValidationError("cycle detected in hierarchy")
		}
		visited[currentUnit.ID] = true

		// Insertar al inicio para mantener el orden de raíz a hoja
		unitCopy := *currentUnit
		path = append([]*entities.AcademicUnit{&unitCopy}, path...)

		// Buscar el padre
		if currentUnit.ParentUnitID != nil {
			parent, exists := r.academicUnits[*currentUnit.ParentUnitID]
			if !exists {
				// Si el padre no existe, romper el ciclo
				break
			}
			currentUnit = parent
		} else {
			break
		}
	}

	return path, nil
}

// ExistsBySchoolIDAndCode verifica si existe una unidad con ese código en la escuela
func (r *MockAcademicUnitRepository) ExistsBySchoolIDAndCode(ctx context.Context, schoolID uuid.UUID, code string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, unit := range r.academicUnits {
		if unit.DeletedAt != nil {
			continue
		}
		if unit.SchoolID == schoolID && unit.Code == code {
			return true, nil
		}
	}

	return false, nil
}

// wouldCreateCycle verifica si establecer newParentID como padre de unitID crearía un ciclo
// Nota: Esta función debe ser llamada con el mutex ya adquirido
func (r *MockAcademicUnitRepository) wouldCreateCycle(unitID, newParentID uuid.UUID) bool {
	// Si el nuevo padre es el mismo que la unidad, es un ciclo directo
	if unitID == newParentID {
		return true
	}

	// Recorrer la jerarquía del nuevo padre hacia arriba
	visited := make(map[uuid.UUID]bool)
	currentID := newParentID

	for {
		// Detectar ciclo infinito en la validación
		if visited[currentID] {
			return true
		}
		visited[currentID] = true

		// Si encontramos la unidad original en la jerarquía del padre, es un ciclo
		if currentID == unitID {
			return true
		}

		// Buscar el padre del actual
		current, exists := r.academicUnits[currentID]
		if !exists || current.ParentUnitID == nil {
			break
		}

		currentID = *current.ParentUnitID
	}

	return false
}

// Reset limpia todos los datos y recarga los datos mock (útil para testing)
func (r *MockAcademicUnitRepository) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Limpiar mapa
	r.academicUnits = make(map[uuid.UUID]*entities.AcademicUnit)

	// Recargar datos desde mockData
	for _, unit := range mockData.GetAcademicUnits() {
		unitCopy := *unit
		r.academicUnits[unit.ID] = &unitCopy
	}
}
