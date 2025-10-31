package entity

import (
	"time"

	"github.com/edugo/api-administracion/internal/domain/valueobject"
	"github.com/edugo/shared/pkg/errors"
)

// Unit representa una unidad organizacional (departamento, grado, grupo, etc.)
type Unit struct {
	id           valueobject.UnitID
	schoolID     valueobject.SchoolID
	parentUnitID *valueobject.UnitID // nil si es unidad raíz
	name         string
	description  string
	isActive     bool
	createdAt    time.Time
	updatedAt    time.Time
}

// NewUnit crea una nueva unidad
func NewUnit(
	schoolID valueobject.SchoolID,
	parentUnitID *valueobject.UnitID,
	name string,
	description string,
) (*Unit, error) {
	// Validaciones de negocio
	if schoolID.IsZero() {
		return nil, errors.NewValidationError("school_id is required")
	}

	if name == "" {
		return nil, errors.NewValidationError("name is required")
	}

	if len(name) < 2 {
		return nil, errors.NewValidationError("name must be at least 2 characters")
	}

	// No puede ser su propio padre
	unitID := valueobject.NewUnitID()
	if parentUnitID != nil && unitID.Equals(*parentUnitID) {
		return nil, errors.NewBusinessRuleError("unit cannot be its own parent")
	}

	now := time.Now()

	return &Unit{
		id:           unitID,
		schoolID:     schoolID,
		parentUnitID: parentUnitID,
		name:         name,
		description:  description,
		isActive:     true,
		createdAt:    now,
		updatedAt:    now,
	}, nil
}

// ReconstructUnit reconstruye una Unit desde la base de datos
func ReconstructUnit(
	id valueobject.UnitID,
	schoolID valueobject.SchoolID,
	parentUnitID *valueobject.UnitID,
	name string,
	description string,
	isActive bool,
	createdAt time.Time,
	updatedAt time.Time,
) *Unit {
	return &Unit{
		id:           id,
		schoolID:     schoolID,
		parentUnitID: parentUnitID,
		name:         name,
		description:  description,
		isActive:     isActive,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}
}

// Getters

func (u *Unit) ID() valueobject.UnitID {
	return u.id
}

func (u *Unit) SchoolID() valueobject.SchoolID {
	return u.schoolID
}

func (u *Unit) ParentUnitID() *valueobject.UnitID {
	return u.parentUnitID
}

func (u *Unit) Name() string {
	return u.name
}

func (u *Unit) Description() string {
	return u.description
}

func (u *Unit) IsActive() bool {
	return u.isActive
}

func (u *Unit) CreatedAt() time.Time {
	return u.createdAt
}

func (u *Unit) UpdatedAt() time.Time {
	return u.updatedAt
}

// Business Logic Methods

// IsRootUnit verifica si es una unidad raíz (sin padre)
func (u *Unit) IsRootUnit() bool {
	return u.parentUnitID == nil
}

// UpdateInfo actualiza la información de la unidad
func (u *Unit) UpdateInfo(name, description *string) error {
	if name != nil && *name != "" {
		if len(*name) < 2 {
			return errors.NewValidationError("name must be at least 2 characters")
		}
		u.name = *name
	}

	if description != nil {
		u.description = *description
	}

	u.updatedAt = time.Now()
	return nil
}

// ChangeParent cambia el padre de la unidad
func (u *Unit) ChangeParent(newParentID *valueobject.UnitID) error {
	// No puede ser su propio padre
	if newParentID != nil && u.id.Equals(*newParentID) {
		return errors.NewBusinessRuleError("unit cannot be its own parent")
	}

	u.parentUnitID = newParentID
	u.updatedAt = time.Now()
	return nil
}

// Deactivate desactiva la unidad
func (u *Unit) Deactivate() error {
	if !u.isActive {
		return errors.NewBusinessRuleError("unit is already inactive")
	}

	u.isActive = false
	u.updatedAt = time.Now()
	return nil
}
