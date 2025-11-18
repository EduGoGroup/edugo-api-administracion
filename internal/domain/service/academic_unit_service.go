package service

import (
	"time"

	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/entity"
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/valueobject"
	"github.com/EduGoGroup/edugo-shared/common/errors"
)

// AcademicUnitDomainService maneja la lógica de negocio de unidades académicas
// Esta capa contiene las reglas de negocio y validaciones que antes estaban en la entity
type AcademicUnitDomainService struct{}

// NewAcademicUnitDomainService crea una nueva instancia del servicio
func NewAcademicUnitDomainService() *AcademicUnitDomainService {
	return &AcademicUnitDomainService{}
}

// SetParent establece la unidad padre en la jerarquía con validaciones completas
func (s *AcademicUnitDomainService) SetParent(
	unit *entity.AcademicUnit,
	parentID valueobject.UnitID,
	parentType valueobject.UnitType,
) error {
	// No puede ser su propio padre (validar primero)
	if unit.ID().Equals(parentID) {
		return errors.NewBusinessRuleError("unit cannot be its own parent")
	}

	// Validar que el tipo padre puede tener hijos
	if !parentType.CanHaveChildren() {
		return errors.NewBusinessRuleError("parent unit type cannot have children: " + parentType.String())
	}

	// Validar que el tipo de hijo está permitido
	allowedTypes := parentType.AllowedChildTypes()
	isAllowed := false
	for _, allowed := range allowedTypes {
		if unit.UnitType() == allowed {
			isAllowed = true
			break
		}
	}

	if !isAllowed {
		return errors.NewBusinessRuleError(
			"unit type " + unit.UnitType().String() + " cannot be child of " + parentType.String(),
		)
	}

	// Aplicar cambio
	unit.SetParentID(parentID)
	unit.SetUpdatedAt(time.Now())
	return nil
}

// RemoveParent remueve la unidad padre (convierte en unidad raíz)
func (s *AcademicUnitDomainService) RemoveParent(unit *entity.AcademicUnit) {
	unit.RemoveParentID()
	unit.SetUpdatedAt(time.Now())
}

// AddChild agrega un hijo a una unidad con todas las validaciones
func (s *AcademicUnitDomainService) AddChild(
	parent *entity.AcademicUnit,
	child *entity.AcademicUnit,
) error {
	if child == nil {
		return errors.NewValidationError("child cannot be nil")
	}

	// Validar que el padre puede tener hijos
	if !s.CanHaveChildren(parent) {
		return errors.NewBusinessRuleError("this unit type cannot have children: " + parent.UnitType().String())
	}

	// Validar que el hijo no es el mismo padre
	if parent.ID().Equals(child.ID()) {
		return errors.NewBusinessRuleError("unit cannot be its own child")
	}

	// Validar que el tipo de hijo está permitido (antes de validar parent_id)
	allowedTypes := parent.UnitType().AllowedChildTypes()
	isAllowed := false
	for _, allowed := range allowedTypes {
		if child.UnitType() == allowed {
			isAllowed = true
			break
		}
	}

	if !isAllowed {
		return errors.NewBusinessRuleError(
			"unit type " + child.UnitType().String() + " cannot be child of " + parent.UnitType().String(),
		)
	}

	// Validar que el hijo apunta a esta unidad como padre
	if child.ParentUnitID() == nil {
		return errors.NewBusinessRuleError("child must have a parent_id")
	}

	if !child.ParentUnitID().Equals(parent.ID()) {
		return errors.NewBusinessRuleError("child's parent_id does not match this unit's id")
	}

	// Validar que el hijo no está ya agregado
	for _, existingChild := range parent.Children() {
		if existingChild.ID().Equals(child.ID()) {
			return errors.NewBusinessRuleError("child is already added")
		}
	}

	// Agregar el hijo
	parent.AddChildToSlice(child)
	parent.SetUpdatedAt(time.Now())
	return nil
}

// RemoveChild remueve un hijo de una unidad
func (s *AcademicUnitDomainService) RemoveChild(
	parent *entity.AcademicUnit,
	childID valueobject.UnitID,
) error {
	if childID.IsZero() {
		return errors.NewValidationError("child_id is required")
	}

	// Buscar el hijo
	children := parent.Children()
	indexToRemove := -1
	for i, child := range children {
		if child.ID().Equals(childID) {
			indexToRemove = i
			break
		}
	}

	if indexToRemove == -1 {
		return errors.NewBusinessRuleError("child not found")
	}

	// Remover el hijo
	parent.RemoveChildFromSlice(childID)
	parent.SetUpdatedAt(time.Now())
	return nil
}

// GetAllDescendants retorna todos los descendientes de una unidad de forma recursiva
func (s *AcademicUnitDomainService) GetAllDescendants(unit *entity.AcademicUnit) []*entity.AcademicUnit {
	descendants := make([]*entity.AcademicUnit, 0)

	// Agregar hijos directos
	for _, child := range unit.Children() {
		descendants = append(descendants, child)

		// Agregar descendientes de cada hijo (recursivo)
		childDescendants := s.GetAllDescendants(child)
		descendants = append(descendants, childDescendants...)
	}

	return descendants
}

// GetDepth retorna la profundidad máxima del árbol desde un nodo
func (s *AcademicUnitDomainService) GetDepth(unit *entity.AcademicUnit) int {
	if !s.HasChildren(unit) {
		return 0
	}

	maxChildDepth := 0
	for _, child := range unit.Children() {
		childDepth := s.GetDepth(child)
		if childDepth > maxChildDepth {
			maxChildDepth = childDepth
		}
	}

	return maxChildDepth + 1
}

// UpdateInfo actualiza la información de una unidad
func (s *AcademicUnitDomainService) UpdateInfo(
	unit *entity.AcademicUnit,
	displayName, description string,
) error {
	if displayName == "" && description == "" {
		return errors.NewValidationError("at least one field must be provided")
	}

	if displayName != "" {
		if len(displayName) < 3 {
			return errors.NewValidationError("display_name must be at least 3 characters")
		}
		unit.SetDisplayName(displayName)
	}

	if description != "" {
		unit.SetDescription(description)
	}

	unit.SetUpdatedAt(time.Now())
	return nil
}

// UpdateDisplayName actualiza el nombre de visualización
func (s *AcademicUnitDomainService) UpdateDisplayName(
	unit *entity.AcademicUnit,
	displayName string,
) error {
	if displayName == "" {
		return errors.NewValidationError("display_name is required")
	}

	if len(displayName) < 3 {
		return errors.NewValidationError("display_name must be at least 3 characters")
	}

	unit.SetDisplayName(displayName)
	unit.SetUpdatedAt(time.Now())
	return nil
}

// CanHaveChildren verifica si una unidad puede tener hijos
func (s *AcademicUnitDomainService) CanHaveChildren(unit *entity.AcademicUnit) bool {
	return unit.UnitType().CanHaveChildren()
}

// HasChildren verifica si una unidad tiene hijos
func (s *AcademicUnitDomainService) HasChildren(unit *entity.AcademicUnit) bool {
	return len(unit.Children()) > 0
}

// IsRoot verifica si una unidad es raíz (sin padre)
func (s *AcademicUnitDomainService) IsRoot(unit *entity.AcademicUnit) bool {
	return unit.ParentUnitID() == nil
}

// SoftDelete marca una unidad como eliminada
func (s *AcademicUnitDomainService) SoftDelete(unit *entity.AcademicUnit) error {
	if unit.IsDeleted() {
		return errors.NewBusinessRuleError("unit is already deleted")
	}

	now := time.Now()
	unit.SetDeletedAt(&now)
	unit.SetUpdatedAt(now)
	return nil
}

// Restore restaura una unidad eliminada
func (s *AcademicUnitDomainService) Restore(unit *entity.AcademicUnit) error {
	if !unit.IsDeleted() {
		return errors.NewBusinessRuleError("unit is not deleted")
	}

	unit.SetDeletedAt(nil)
	unit.SetUpdatedAt(time.Now())
	return nil
}

// Validate valida el estado completo de una unidad
func (s *AcademicUnitDomainService) Validate(unit *entity.AcademicUnit) error {
	if unit.SchoolID().IsZero() {
		return errors.NewValidationError("school_id is required")
	}

	if !unit.UnitType().IsValid() {
		return errors.NewValidationError("invalid unit type")
	}

	if unit.DisplayName() == "" {
		return errors.NewValidationError("display_name is required")
	}

	if len(unit.DisplayName()) < 3 {
		return errors.NewValidationError("display_name must be at least 3 characters")
	}

	return nil
}
