package dto

import (
	"time"

	"github.com/edugo/api-administracion/internal/domain/entity"
	"github.com/edugo/shared/pkg/types/enum"
	"github.com/edugo/shared/pkg/validator"
)

// CreateUserRequest representa la solicitud para crear un usuario
type CreateUserRequest struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`
}

// Validate valida el request usando shared/validator
func (r *CreateUserRequest) Validate() error {
	v := validator.New()

	v.Required(r.Email, "email")
	v.Email(r.Email, "email")
	v.MaxLength(r.Email, 100, "email")

	v.Required(r.FirstName, "first_name")
	v.MinLength(r.FirstName, 2, "first_name")
	v.MaxLength(r.FirstName, 50, "first_name")
	v.Name(r.FirstName, "first_name")

	v.Required(r.LastName, "last_name")
	v.MinLength(r.LastName, 2, "last_name")
	v.MaxLength(r.LastName, 50, "last_name")
	v.Name(r.LastName, "last_name")

	v.Required(r.Role, "role")
	v.InSlice(r.Role, enum.AllSystemRolesStrings(), "role")

	return v.GetError()
}

// UpdateUserRequest representa la solicitud para actualizar un usuario
type UpdateUserRequest struct {
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	Role      *string `json:"role,omitempty"`
	IsActive  *bool   `json:"is_active,omitempty"`
}

// Validate valida el request
func (r *UpdateUserRequest) Validate() error {
	v := validator.New()

	if r.FirstName != nil {
		v.MinLength(*r.FirstName, 2, "first_name")
		v.MaxLength(*r.FirstName, 50, "first_name")
		v.Name(*r.FirstName, "first_name")
	}

	if r.LastName != nil {
		v.MinLength(*r.LastName, 2, "last_name")
		v.MaxLength(*r.LastName, 50, "last_name")
		v.Name(*r.LastName, "last_name")
	}

	if r.Role != nil {
		v.InSlice(*r.Role, enum.AllSystemRolesStrings(), "role")
	}

	return v.GetError()
}

// UserResponse representa la respuesta de un usuario
type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	FullName  string    `json:"full_name"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToUserResponse convierte una entidad User a DTO de respuesta
func ToUserResponse(user *entity.User) *UserResponse {
	return &UserResponse{
		ID:        user.ID().String(),
		Email:     user.Email().String(),
		FirstName: user.FirstName(),
		LastName:  user.LastName(),
		FullName:  user.FullName(),
		Role:      user.Role().String(),
		IsActive:  user.IsActive(),
		CreatedAt: user.CreatedAt(),
		UpdatedAt: user.UpdatedAt(),
	}
}

// ToUserResponses convierte una lista de entidades a DTOs
func ToUserResponses(users []*entity.User) []*UserResponse {
	responses := make([]*UserResponse, len(users))
	for i, user := range users {
		responses[i] = ToUserResponse(user)
	}
	return responses
}
