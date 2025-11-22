package dto

import (
	"encoding/json"
	"time"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
)

// CreateSchoolRequest representa la solicitud para crear una escuela
type CreateSchoolRequest struct {
	Name         string                 `json:"name" validate:"required,min=3"`
	Code         string                 `json:"code" validate:"required,min=3"`
	Address      string                 `json:"address"`
	ContactEmail string                 `json:"contact_email" validate:"omitempty,email"`
	ContactPhone string                 `json:"contact_phone"`
	Metadata     map[string]interface{} `json:"metadata"`
}

// UpdateSchoolRequest representa la solicitud para actualizar una escuela
type UpdateSchoolRequest struct {
	Name         *string                `json:"name" validate:"omitempty,min=3"`
	Address      *string                `json:"address"`
	ContactEmail *string                `json:"contact_email" validate:"omitempty,email"`
	ContactPhone *string                `json:"contact_phone"`
	Metadata     map[string]interface{} `json:"metadata"`
}

// SchoolResponse representa la respuesta con datos de una escuela
type SchoolResponse struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Code         string                 `json:"code"`
	Address      string                 `json:"address"`
	ContactEmail string                 `json:"contact_email,omitempty"`
	ContactPhone string                 `json:"contact_phone,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

// ToSchoolResponse convierte una entidad School de infrastructure a SchoolResponse
func ToSchoolResponse(school *entities.School) SchoolResponse {
	var email string
	if school.Email != nil {
		email = *school.Email
	}

	var phone string
	if school.Phone != nil {
		phone = *school.Phone
	}

	var address string
	if school.Address != nil {
		address = *school.Address
	}

	// Deserializar metadata
	var metadata map[string]interface{}
	if len(school.Metadata) > 0 {
		_ = json.Unmarshal(school.Metadata, &metadata)
	}

	return SchoolResponse{
		ID:           school.ID.String(),
		Name:         school.Name,
		Code:         school.Code,
		Address:      address,
		ContactEmail: email,
		ContactPhone: phone,
		Metadata:     metadata,
		CreatedAt:    school.CreatedAt,
		UpdatedAt:    school.UpdatedAt,
	}
}

// ToSchoolResponseList convierte una lista de entidades a lista de responses
func ToSchoolResponseList(schools []*entities.School) []SchoolResponse {
	responses := make([]SchoolResponse, len(schools))
	for i, school := range schools {
		responses[i] = ToSchoolResponse(school)
	}
	return responses
}
