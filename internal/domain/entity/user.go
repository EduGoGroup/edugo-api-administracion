package entity

import (
	"time"

	"github.com/edugo/api-administracion/internal/domain/valueobject"
	"github.com/edugo/shared/pkg/errors"
	"github.com/edugo/shared/pkg/types/enum"
)

// User representa un usuario del sistema
type User struct {
	id        valueobject.UserID
	email     valueobject.Email
	firstName string
	lastName  string
	role      enum.SystemRole
	isActive  bool
	createdAt time.Time
	updatedAt time.Time
}

// NewUser crea un nuevo usuario con validaciones de negocio
func NewUser(
	email valueobject.Email,
	firstName string,
	lastName string,
	role enum.SystemRole,
) (*User, error) {
	// Validaciones de negocio
	if email.IsZero() {
		return nil, errors.NewValidationError("email is required")
	}

	if firstName == "" {
		return nil, errors.NewValidationError("first_name is required")
	}

	if lastName == "" {
		return nil, errors.NewValidationError("last_name is required")
	}

	if !role.IsValid() {
		return nil, errors.NewValidationError("invalid role")
	}

	// No permitir crear admin users (regla de negocio)
	if role == enum.SystemRoleAdmin {
		return nil, errors.NewBusinessRuleError("cannot create admin users through this endpoint")
	}

	now := time.Now()

	return &User{
		id:        valueobject.NewUserID(),
		email:     email,
		firstName: firstName,
		lastName:  lastName,
		role:      role,
		isActive:  true,
		createdAt: now,
		updatedAt: now,
	}, nil
}

// ReconstructUser reconstruye un User desde la base de datos
func ReconstructUser(
	id valueobject.UserID,
	email valueobject.Email,
	firstName string,
	lastName string,
	role enum.SystemRole,
	isActive bool,
	createdAt time.Time,
	updatedAt time.Time,
) *User {
	return &User{
		id:        id,
		email:     email,
		firstName: firstName,
		lastName:  lastName,
		role:      role,
		isActive:  isActive,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

// Getters

func (u *User) ID() valueobject.UserID {
	return u.id
}

func (u *User) Email() valueobject.Email {
	return u.email
}

func (u *User) FirstName() string {
	return u.firstName
}

func (u *User) LastName() string {
	return u.lastName
}

func (u *User) FullName() string {
	return u.firstName + " " + u.lastName
}

func (u *User) Role() enum.SystemRole {
	return u.role
}

func (u *User) IsActive() bool {
	return u.isActive
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}

// Business Logic Methods

// Deactivate desactiva el usuario
func (u *User) Deactivate() error {
	if !u.isActive {
		return errors.NewBusinessRuleError("user is already inactive")
	}

	u.isActive = false
	u.updatedAt = time.Now()
	return nil
}

// Activate activa el usuario
func (u *User) Activate() error {
	if u.isActive {
		return errors.NewBusinessRuleError("user is already active")
	}

	u.isActive = true
	u.updatedAt = time.Now()
	return nil
}

// UpdateName actualiza el nombre del usuario
func (u *User) UpdateName(firstName, lastName string) error {
	if firstName == "" || lastName == "" {
		return errors.NewValidationError("first_name and last_name are required")
	}

	u.firstName = firstName
	u.lastName = lastName
	u.updatedAt = time.Now()
	return nil
}

// ChangeRole cambia el rol del usuario
func (u *User) ChangeRole(newRole enum.SystemRole) error {
	if !newRole.IsValid() {
		return errors.NewValidationError("invalid role")
	}

	if u.role == newRole {
		return errors.NewBusinessRuleError("new role is the same as current role")
	}

	// No permitir promover a admin
	if newRole == enum.SystemRoleAdmin {
		return errors.NewBusinessRuleError("cannot promote to admin role")
	}

	u.role = newRole
	u.updatedAt = time.Now()
	return nil
}

// IsTeacher verifica si el usuario es teacher
func (u *User) IsTeacher() bool {
	return u.role == enum.SystemRoleTeacher
}

// IsStudent verifica si el usuario es student
func (u *User) IsStudent() bool {
	return u.role == enum.SystemRoleStudent
}

// IsGuardian verifica si el usuario es guardian
func (u *User) IsGuardian() bool {
	return u.role == enum.SystemRoleGuardian
}
