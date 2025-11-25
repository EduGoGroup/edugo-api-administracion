// Package crypto proporciona utilidades criptográficas
package crypto

import (
	"errors"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

// Errores de validación de password
var (
	ErrPasswordTooShort    = errors.New("password debe tener al menos 8 caracteres")
	ErrPasswordNoUpper     = errors.New("password debe tener al menos una mayúscula")
	ErrPasswordNoLower     = errors.New("password debe tener al menos una minúscula")
	ErrPasswordNoNumber    = errors.New("password debe tener al menos un número")
	ErrPasswordMismatch    = errors.New("password incorrecto")
)

// PasswordHasher maneja el hashing de passwords con bcrypt
type PasswordHasher struct {
	cost int
}

// NewPasswordHasher crea un nuevo hasher con el costo especificado
// El costo por defecto de bcrypt es 10, recomendado 12 para producción
func NewPasswordHasher(cost int) *PasswordHasher {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		cost = bcrypt.DefaultCost
	}
	return &PasswordHasher{cost: cost}
}

// Hash genera un hash bcrypt del password
func (h *PasswordHasher) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// Compare compara un password con su hash
// Retorna nil si coinciden, ErrPasswordMismatch si no
func (h *PasswordHasher) Compare(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrPasswordMismatch
		}
		return err
	}
	return nil
}

// Validate verifica que el password cumple los requisitos mínimos
func (h *PasswordHasher) Validate(password string) error {
	if len(password) < 8 {
		return ErrPasswordTooShort
	}

	var hasUpper, hasLower, hasNumber bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		}
	}

	if !hasUpper {
		return ErrPasswordNoUpper
	}
	if !hasLower {
		return ErrPasswordNoLower
	}
	if !hasNumber {
		return ErrPasswordNoNumber
	}

	return nil
}
