package entity

import (
	"time"

	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/valueobject"
	"github.com/EduGoGroup/edugo-shared/common/errors"
)

// AcademicUnit representa una unidad académica en la jerarquía (grado, sección, club, etc.)
type AcademicUnit struct {
	id           valueobject.UnitID
	parentUnitID *valueobject.UnitID
	schoolID     valueobject.SchoolID
	unitType     valueobject.UnitType
	displayName  string
	code         string
	description  string
	metadata     map[string]interface{}
	createdAt    time.Time
	updatedAt    time.Time
	deletedAt    *time.Time
}

// NewAcademicUnit crea una nueva unidad académica
func NewAcademicUnit(
	schoolID valueobject.SchoolID,
	unitType valueobject.UnitType,
	displayName string,
	code string,
) (*AcademicUnit, error) {
	// Validaciones de negocio
	if schoolID.IsZero() {
		return nil, errors.NewValidationError("school_id is required")
	}

	if !unitType.IsValid() {
		return nil, errors.NewValidationError("invalid unit type")
	}

	if displayName == "" {
		return nil, errors.NewValidationError("display_name is required")
	}

	if len(displayName) < 3 {
		return nil, errors.NewValidationError("display_name must be at least 3 characters")
	}

	now := time.Now()

	return &AcademicUnit{
		id:          valueobject.NewUnitID(),
		schoolID:    schoolID,
		unitType:    unitType,
		displayName: displayName,
		code:        code,
		metadata:    make(map[string]interface{}),
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

// ReconstructAcademicUnit reconstruye una AcademicUnit desde la base de datos
func ReconstructAcademicUnit(
	id valueobject.UnitID,
	parentUnitID *valueobject.UnitID,
	schoolID valueobject.SchoolID,
	unitType valueobject.UnitType,
	displayName string,
	code string,
	description string,
	metadata map[string]interface{},
	createdAt time.Time,
	updatedAt time.Time,
	deletedAt *time.Time,
) *AcademicUnit {
	if metadata == nil {
		metadata = make(map[string]interface{})
	}

	return &AcademicUnit{
		id:           id,
		parentUnitID: parentUnitID,
		schoolID:     schoolID,
		unitType:     unitType,
		displayName:  displayName,
		code:         code,
		description:  description,
		metadata:     metadata,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
		deletedAt:    deletedAt,
	}
}

// Getters

func (au *AcademicUnit) ID() valueobject.UnitID {
	return au.id
}

func (au *AcademicUnit) ParentUnitID() *valueobject.UnitID {
	return au.parentUnitID
}

func (au *AcademicUnit) SchoolID() valueobject.SchoolID {
	return au.schoolID
}

func (au *AcademicUnit) UnitType() valueobject.UnitType {
	return au.unitType
}

func (au *AcademicUnit) DisplayName() string {
	return au.displayName
}

func (au *AcademicUnit) Code() string {
	return au.code
}

func (au *AcademicUnit) Description() string {
	return au.description
}

func (au *AcademicUnit) Metadata() map[string]interface{} {
	copy := make(map[string]interface{})
	for k, v := range au.metadata {
		copy[k] = v
	}
	return copy
}

func (au *AcademicUnit) CreatedAt() time.Time {
	return au.createdAt
}

func (au *AcademicUnit) UpdatedAt() time.Time {
	return au.updatedAt
}

func (au *AcademicUnit) DeletedAt() *time.Time {
	return au.deletedAt
}

// Business Logic Methods

// SetParent establece la unidad padre en la jerarquía
func (au *AcademicUnit) SetParent(parentID valueobject.UnitID, parentType valueobject.UnitType) error {
	// Validar que el tipo padre puede tener hijos
	if !parentType.CanHaveChildren() {
		return errors.NewBusinessRuleError("parent unit type cannot have children: " + parentType.String())
	}

	// Validar que el tipo de hijo está permitido
	allowedTypes := parentType.AllowedChildTypes()
	isAllowed := false
	for _, allowed := range allowedTypes {
		if au.unitType == allowed {
			isAllowed = true
			break
		}
	}

	if !isAllowed {
		return errors.NewBusinessRuleError(
			"unit type " + au.unitType.String() + " cannot be child of " + parentType.String(),
		)
	}

	// No puede ser su propio padre
	if au.id.Equals(parentID) {
		return errors.NewBusinessRuleError("unit cannot be its own parent")
	}

	au.parentUnitID = &parentID
	au.updatedAt = time.Now()
	return nil
}

// RemoveParent remueve la unidad padre (convierte en unidad raíz)
func (au *AcademicUnit) RemoveParent() {
	au.parentUnitID = nil
	au.updatedAt = time.Now()
}

// UpdateInfo actualiza la información de la unidad
func (au *AcademicUnit) UpdateInfo(displayName, description string) error {
	if displayName == "" && description == "" {
		return errors.NewValidationError("at least one field must be provided")
	}

	if displayName != "" {
		if len(displayName) < 3 {
			return errors.NewValidationError("display_name must be at least 3 characters")
		}
		au.displayName = displayName
	}

	if description != "" {
		au.description = description
	}

	au.updatedAt = time.Now()
	return nil
}

// CanHaveChildren verifica si esta unidad puede tener hijos
func (au *AcademicUnit) CanHaveChildren() bool {
	return au.unitType.CanHaveChildren()
}

// IsRoot verifica si es una unidad raíz (sin padre)
func (au *AcademicUnit) IsRoot() bool {
	return au.parentUnitID == nil
}

// IsDeleted verifica si la unidad está eliminada (soft delete)
func (au *AcademicUnit) IsDeleted() bool {
	return au.deletedAt != nil
}

// SoftDelete marca la unidad como eliminada
func (au *AcademicUnit) SoftDelete() error {
	if au.IsDeleted() {
		return errors.NewBusinessRuleError("unit is already deleted")
	}

	now := time.Now()
	au.deletedAt = &now
	au.updatedAt = now
	return nil
}

// Restore restaura una unidad eliminada
func (au *AcademicUnit) Restore() error {
	if !au.IsDeleted() {
		return errors.NewBusinessRuleError("unit is not deleted")
	}

	au.deletedAt = nil
	au.updatedAt = time.Now()
	return nil
}

// SetMetadata establece un valor en el metadata
func (au *AcademicUnit) SetMetadata(key string, value interface{}) {
	if au.metadata == nil {
		au.metadata = make(map[string]interface{})
	}
	au.metadata[key] = value
	au.updatedAt = time.Now()
}

// GetMetadata obtiene un valor del metadata
func (au *AcademicUnit) GetMetadata(key string) (interface{}, bool) {
	if au.metadata == nil {
		return nil, false
	}
	val, exists := au.metadata[key]
	return val, exists
}

// Validate valida el estado completo de la entidad
func (au *AcademicUnit) Validate() error {
	if au.schoolID.IsZero() {
		return errors.NewValidationError("school_id is required")
	}

	if !au.unitType.IsValid() {
		return errors.NewValidationError("invalid unit type")
	}

	if au.displayName == "" {
		return errors.NewValidationError("display_name is required")
	}

	if len(au.displayName) < 3 {
		return errors.NewValidationError("display_name must be at least 3 characters")
	}

	return nil
}
