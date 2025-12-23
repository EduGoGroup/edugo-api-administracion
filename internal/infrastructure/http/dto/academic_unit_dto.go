package dto

import (
	"time"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
)

// CreateUnitRequest DTO para crear unidad académica
type CreateUnitRequest struct {
	ParentUnitID *string `json:"parent_unit_id" binding:"omitempty,uuid"`
	SchoolID     string  `json:"school_id" binding:"required,uuid"`
	Type         string  `json:"type" binding:"required"` // Validación realizada por valueobject.ParseUnitType
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

// UnitTreeNode DTO para árbol jerárquico
type UnitTreeNode struct {
	ID       string          `json:"id"`
	Name     string          `json:"name"`
	Code     string          `json:"code"`
	Type     string          `json:"type"`
	Depth    int             `json:"depth"`
	Children []*UnitTreeNode `json:"children,omitempty"`
} // @name UnitTreeNode

// ToUnitResponse convierte entity de infrastructure a DTO
func ToUnitResponse(unit *entities.AcademicUnit) *UnitResponse {
	var parentID *string
	if unit.ParentUnitID != nil {
		pid := unit.ParentUnitID.String()
		parentID = &pid
	}

	var desc *string
	if unit.Description != nil && *unit.Description != "" {
		desc = unit.Description
	}

	return &UnitResponse{
		ID:           unit.ID.String(),
		ParentUnitID: parentID,
		SchoolID:     unit.SchoolID.String(),
		Type:         unit.Type,
		Name:         unit.Name,
		Code:         unit.Code,
		Description:  desc,
		CreatedAt:    unit.CreatedAt,
		UpdatedAt:    unit.UpdatedAt,
	}
}

// ToAcademicUnitResponse alias para compatibilidad
func ToAcademicUnitResponse(unit *entities.AcademicUnit) *UnitResponse {
	return ToUnitResponse(unit)
}

// BuildUnitTree construye árbol jerárquico desde lista plana
func BuildUnitTree(units []*entities.AcademicUnit) []*UnitTreeNode {
	if len(units) == 0 {
		return []*UnitTreeNode{}
	}

	// Mapear units por ID para acceso rápido
	unitMap := make(map[string]*UnitTreeNode)
	var roots []*UnitTreeNode

	// Crear nodos
	for _, unit := range units {
		node := &UnitTreeNode{
			ID:       unit.ID.String(),
			Name:     unit.Name,
			Code:     unit.Code,
			Type:     unit.Type,
			Depth:    0,
			Children: []*UnitTreeNode{},
		}
		unitMap[unit.ID.String()] = node
	}

	// Construir árbol
	for _, unit := range units {
		node := unitMap[unit.ID.String()]
		if unit.ParentUnitID == nil {
			// Es raíz
			roots = append(roots, node)
		} else {
			// Tiene padre
			parentID := unit.ParentUnitID.String()
			if parent, exists := unitMap[parentID]; exists {
				node.Depth = parent.Depth + 1
				parent.Children = append(parent.Children, node)
			} else {
				// Padre no encontrado, tratar como raíz
				roots = append(roots, node)
			}
		}
	}

	return roots
}
