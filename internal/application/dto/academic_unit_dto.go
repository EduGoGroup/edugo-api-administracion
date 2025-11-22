package dto

import (
	"encoding/json"
	"time"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
)

// CreateAcademicUnitRequest representa la solicitud para crear una unidad académica
type CreateAcademicUnitRequest struct {
	ParentUnitID *string                `json:"parent_unit_id" validate:"omitempty,uuid"`
	Type         string                 `json:"type" validate:"required,oneof=school grade section club department"`
	DisplayName  string                 `json:"display_name" validate:"required,min=3,max=255"`
	Code         string                 `json:"code" validate:"omitempty,min=2,max=50"`
	Description  string                 `json:"description"`
	Metadata     map[string]interface{} `json:"metadata"`
}

// UpdateAcademicUnitRequest representa la solicitud para actualizar una unidad
type UpdateAcademicUnitRequest struct {
	ParentUnitID *string                `json:"parent_unit_id" validate:"omitempty,uuid"`
	DisplayName  *string                `json:"display_name" validate:"omitempty,min=3,max=255"`
	Description  *string                `json:"description"`
	Metadata     map[string]interface{} `json:"metadata"`
}

// AcademicUnitResponse representa la respuesta con datos de una unidad académica
type AcademicUnitResponse struct {
	ID           string                 `json:"id"`
	ParentUnitID *string                `json:"parent_unit_id,omitempty"`
	SchoolID     string                 `json:"school_id"`
	Type         string                 `json:"type"`
	DisplayName  string                 `json:"display_name"`
	Code         string                 `json:"code,omitempty"`
	Description  string                 `json:"description,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
	DeletedAt    *time.Time             `json:"deleted_at,omitempty"`
}

// UnitTreeNode representa un nodo en el árbol jerárquico
type UnitTreeNode struct {
	ID          string          `json:"id"`
	Type        string          `json:"type"`
	DisplayName string          `json:"display_name"`
	Code        string          `json:"code,omitempty"`
	Depth       int             `json:"depth"`
	Children    []*UnitTreeNode `json:"children,omitempty"`
}

// ToAcademicUnitResponse convierte una entidad AcademicUnit de infrastructure a response
func ToAcademicUnitResponse(unit *entities.AcademicUnit) AcademicUnitResponse {
	var parentID *string
	if unit.ParentUnitID != nil {
		id := unit.ParentUnitID.String()
		parentID = &id
	}

	// Deserializar metadata de []byte a map
	var metadata map[string]interface{}
	if len(unit.Metadata) > 0 {
		_ = json.Unmarshal(unit.Metadata, &metadata)
	}

	desc := ""
	if unit.Description != nil {
		desc = *unit.Description
	}

	return AcademicUnitResponse{
		ID:           unit.ID.String(),
		ParentUnitID: parentID,
		SchoolID:     unit.SchoolID.String(),
		Type:         unit.Type,
		DisplayName:  unit.Name,
		Code:         unit.Code,
		Description:  desc,
		Metadata:     metadata,
		CreatedAt:    unit.CreatedAt,
		UpdatedAt:    unit.UpdatedAt,
		DeletedAt:    unit.DeletedAt,
	}
}

// ToAcademicUnitResponseList convierte una lista de entidades a lista de responses
func ToAcademicUnitResponseList(units []*entities.AcademicUnit) []AcademicUnitResponse {
	responses := make([]AcademicUnitResponse, len(units))
	for i, unit := range units {
		responses[i] = ToAcademicUnitResponse(unit)
	}
	return responses
}

// BuildUnitTree construye árbol jerárquico desde lista plana
func BuildUnitTree(units []*entities.AcademicUnit) []*UnitTreeNode {
	if len(units) == 0 {
		return []*UnitTreeNode{}
	}

	// Mapear units por ID
	unitMap := make(map[string]*UnitTreeNode)
	var roots []*UnitTreeNode

	// Crear nodos
	for _, unit := range units {
		node := &UnitTreeNode{
			ID:          unit.ID.String(),
			Type:        unit.Type,
			DisplayName: unit.Name,
			Code:        unit.Code,
			Depth:       0,
			Children:    []*UnitTreeNode{},
		}
		unitMap[unit.ID.String()] = node
	}

	// Construir árbol
	for _, unit := range units {
		node := unitMap[unit.ID.String()]
		if unit.ParentUnitID == nil {
			roots = append(roots, node)
		} else {
			parentID := unit.ParentUnitID.String()
			if parent, exists := unitMap[parentID]; exists {
				node.Depth = parent.Depth + 1
				parent.Children = append(parent.Children, node)
			} else {
				roots = append(roots, node)
			}
		}
	}

	return roots
}
