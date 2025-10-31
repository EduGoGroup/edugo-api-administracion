package valueobject

import (
	"strings"

	"github.com/EduGoGroup/edugo-shared/pkg/errors"
	"github.com/EduGoGroup/edugo-shared/pkg/validator"
)

// Email representa un email válido
type Email struct {
	value string
}

// NewEmail crea y valida un nuevo Email
func NewEmail(value string) (Email, error) {
	// Normalizar (trim + lowercase)
	normalized := validator.Normalize(value)

	// Validar formato
	if !validator.IsValidEmail(normalized) {
		return Email{}, errors.NewValidationError("invalid email format").
			WithField("email", value)
	}

	return Email{value: normalized}, nil
}

// String retorna la representación en string del email
func (e Email) String() string {
	return e.value
}

// IsZero verifica si es el valor cero
func (e Email) IsZero() bool {
	return e.value == ""
}

// Equals compara dos emails (case-insensitive)
func (e Email) Equals(other Email) bool {
	return strings.EqualFold(e.value, other.value)
}

// Domain retorna el dominio del email (después del @)
func (e Email) Domain() string {
	parts := strings.Split(e.value, "@")
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
}
