package dto

import (
	"time"

	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/entity"
)

// CreateUnitRequest DTO para crear unidad académica
type CreateUnitRequest struct {
	ParentUnitID *string `json:"parent_unit_id" binding:"omitempty,uuid"`
	SchoolID     string  `json:"school_id" binding:"required,uuid"`
	Type         string  `json:"type" binding:"required,oneof=grade section club department"`
	Name         string  `json:"name" binding:"required,min=3,max=100"`
	Code         string  `json:"code" binding:"required,min=2,max=50"`
	Description  *string `json:"description" binding:"omitempty,max=500"`
} // @name CreateUnitRequest

// UpdateUnitRequest DTO para actualizar unidad
type UpdateUnitRequest struct {
	ParentUnitID *string `json:"parent_unit_id" binding:"omitempty,uuid"`
	Name         *string `json:"name" binding:"omitempty,min=3,max=100"`
	Description  *string `json:"description" binding:"omitempty,max=500"`
} // @name UpdateUnitRequest

// UnitResponse DTO de respuesta simple
type UnitResponse struct {
	ID           string    `json:"id"`
	ParentUnitID *string   `json:"parent_unit_id,omitempty"`
	SchoolID     string    `json:"school_id"`
	Type         string    `json:"type"`
	Name         string    `json:"name"`
	Code         string    `json:"code"`
	Description  *string   `json:"description,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
} // @name UnitResponse

// UnitTreeNode DTO para árbol jerárquico (usa ltree!)
type UnitTreeNode struct {
	ID       string          `json:"id"`
	Name     string          `json:"name"`
	Code     string          `json:"code"`
	Type     string          `json:"type"`
	Depth    int             `json:"depth"`
	Children []*UnitTreeNode `json:"children,omitempty"`
} // @name UnitTreeNode

// ToUnitResponse convierte entity a DTO
func ToUnitResponse(unit *entity.AcademicUnit) *UnitResponse {
	var parentID *string
	if unit.ParentUnitID() != nil {
		pid := unit.ParentUnitID().String()
		parentID = &pid
	}

	var desc *string
	if unit.Description() != "" {
		d := unit.Description()
		desc = &d
	}

	return &UnitResponse{
		ID:           unit.ID().String(),
		ParentUnitID: parentID,
		SchoolID:     unit.SchoolID().String(),
		Type:         unit.UnitType().String(),
		Name:         unit.DisplayName(),
		Code:         unit.Code(),
		Description:  desc,
		CreatedAt:    unit.CreatedAt(),
		UpdatedAt:    unit.UpdatedAt(),
	}
}
