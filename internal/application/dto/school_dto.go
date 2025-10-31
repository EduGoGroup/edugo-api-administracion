package dto

import (
	"time"

	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/entity"
	"github.com/EduGoGroup/edugo-shared/pkg/validator"
)

// CreateSchoolRequest representa la solicitud para crear una escuela
type CreateSchoolRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

// Validate valida el request
func (r *CreateSchoolRequest) Validate() error {
	v := validator.New()

	v.Required(r.Name, "name")
	v.MinLength(r.Name, 3, "name")
	v.MaxLength(r.Name, 100, "name")

	v.Required(r.Address, "address")
	v.MinLength(r.Address, 5, "address")
	v.MaxLength(r.Address, 200, "address")

	return v.GetError()
}

// SchoolResponse representa la respuesta de una escuela
type SchoolResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToSchoolResponse convierte una entidad School a DTO
func ToSchoolResponse(school *entity.School) *SchoolResponse {
	return &SchoolResponse{
		ID:        school.ID().String(),
		Name:      school.Name(),
		Address:   school.Address(),
		IsActive:  school.IsActive(),
		CreatedAt: school.CreatedAt(),
		UpdatedAt: school.UpdatedAt(),
	}
}
