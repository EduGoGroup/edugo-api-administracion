package dto

import (
	"time"

	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/entity"
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

// ToAcademicUnitResponse convierte una entidad AcademicUnit a response
func ToAcademicUnitResponse(unit *entity.AcademicUnit) AcademicUnitResponse {
	var parentID *string
	if unit.ParentUnitID() != nil {
		id := unit.ParentUnitID().String()
		parentID = &id
	}

	return AcademicUnitResponse{
		ID:           unit.ID().String(),
		ParentUnitID: parentID,
		SchoolID:     unit.SchoolID().String(),
		Type:         unit.UnitType().String(),
		DisplayName:  unit.DisplayName(),
		Code:         unit.Code(),
		Description:  unit.Description(),
		Metadata:     unit.Metadata(),
		CreatedAt:    unit.CreatedAt(),
		UpdatedAt:    unit.UpdatedAt(),
		DeletedAt:    unit.DeletedAt(),
	}
}

// ToAcademicUnitResponseList convierte una lista de entidades a lista de responses
func ToAcademicUnitResponseList(units []*entity.AcademicUnit) []AcademicUnitResponse {
	responses := make([]AcademicUnitResponse, len(units))
	for i, unit := range units {
		responses[i] = ToAcademicUnitResponse(unit)
	}
	return responses
}

// BuildUnitTree construye un árbol jerárquico desde una lista plana de unidades
func BuildUnitTree(units []*entity.AcademicUnit) []*UnitTreeNode {
	// Mapas para construcción eficiente
	nodeMap := make(map[string]*UnitTreeNode)
	roots := []*UnitTreeNode{}

	// Crear todos los nodos
	for _, unit := range units {
		node := &UnitTreeNode{
			ID:          unit.ID().String(),
			Type:        unit.UnitType().String(),
			DisplayName: unit.DisplayName(),
			Code:        unit.Code(),
			Children:    []*UnitTreeNode{},
		}
		nodeMap[unit.ID().String()] = node
	}

	// Construir relaciones padre-hijo
	for _, unit := range units {
		node := nodeMap[unit.ID().String()]

		if unit.ParentUnitID() == nil {
			// Es raíz
			roots = append(roots, node)
		} else {
			// Tiene padre, agregarlo como hijo
			parentID := unit.ParentUnitID().String()
			if parent, exists := nodeMap[parentID]; exists {
				parent.Children = append(parent.Children, node)
			}
		}
	}

	// Calcular profundidad recursivamente
	var calculateDepth func(node *UnitTreeNode, depth int)
	calculateDepth = func(node *UnitTreeNode, depth int) {
		node.Depth = depth
		for _, child := range node.Children {
			calculateDepth(child, depth+1)
		}
	}

	for _, root := range roots {
		calculateDepth(root, 1)
	}

	return roots
}
