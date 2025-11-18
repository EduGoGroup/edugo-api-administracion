package dto

import (
	"time"

	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/entity"
)

// CreateSchoolRequest DTO para crear escuela
type CreateSchoolRequest struct {
	Name    string `json:"name" binding:"required,min=3,max=100"`
	Code    string `json:"code" binding:"required,min=3,max=20"`
	Address string `json:"address" binding:"required"`
} // @name CreateSchoolRequest

// UpdateSchoolRequest DTO para actualizar escuela
type UpdateSchoolRequest struct {
	Name    *string `json:"name" binding:"omitempty,min=3,max=100"`
	Address *string `json:"address" binding:"omitempty"`
} // @name UpdateSchoolRequest

// SchoolResponse DTO de respuesta
type SchoolResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
} // @name SchoolResponse

// ToSchoolResponse convierte entity a DTO
func ToSchoolResponse(school *entity.School) *SchoolResponse {
	return &SchoolResponse{
		ID:        school.ID().String(),
		Name:      school.Name(),
		Code:      school.Code(),
		Address:   school.Address(),
		CreatedAt: school.CreatedAt(),
		UpdatedAt: school.UpdatedAt(),
	}
}
