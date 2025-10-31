package dto

import (
	"time"

	"github.com/edugo/api-administracion/internal/domain/entity"
	"github.com/edugo/shared/pkg/validator"
)

// CreateUnitRequest solicitud para crear unidad
type CreateUnitRequest struct {
	SchoolID     string  `json:"school_id"`
	ParentUnitID *string `json:"parent_unit_id,omitempty"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
}

func (r *CreateUnitRequest) Validate() error {
	v := validator.New()

	v.Required(r.SchoolID, "school_id")
	v.UUID(r.SchoolID, "school_id")

	if r.ParentUnitID != nil && *r.ParentUnitID != "" {
		v.UUID(*r.ParentUnitID, "parent_unit_id")
	}

	v.Required(r.Name, "name")
	v.MinLength(r.Name, 2, "name")
	v.MaxLength(r.Name, 100, "name")

	v.MaxLength(r.Description, 500, "description")

	return v.GetError()
}

// UpdateUnitRequest solicitud para actualizar unidad
type UpdateUnitRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

func (r *UpdateUnitRequest) Validate() error {
	v := validator.New()

	if r.Name != nil {
		v.MinLength(*r.Name, 2, "name")
		v.MaxLength(*r.Name, 100, "name")
	}

	if r.Description != nil {
		v.MaxLength(*r.Description, 500, "description")
	}

	return v.GetError()
}

// AssignMemberRequest solicitud para asignar miembro
type AssignMemberRequest struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}

func (r *AssignMemberRequest) Validate() error {
	v := validator.New()

	v.Required(r.UserID, "user_id")
	v.UUID(r.UserID, "user_id")

	v.Required(r.Role, "role")
	v.InSlice(r.Role, []string{"teacher", "student"}, "role")

	return v.GetError()
}

// UnitResponse respuesta de unidad
type UnitResponse struct {
	ID           string    `json:"id"`
	SchoolID     string    `json:"school_id"`
	ParentUnitID *string   `json:"parent_unit_id,omitempty"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func ToUnitResponse(unit *entity.Unit) *UnitResponse {
	response := &UnitResponse{
		ID:          unit.ID().String(),
		SchoolID:    unit.SchoolID().String(),
		Name:        unit.Name(),
		Description: unit.Description(),
		IsActive:    unit.IsActive(),
		CreatedAt:   unit.CreatedAt(),
		UpdatedAt:   unit.UpdatedAt(),
	}

	if !unit.IsRootUnit() {
		parentID := unit.ParentUnitID().String()
		response.ParentUnitID = &parentID
	}

	return response
}
