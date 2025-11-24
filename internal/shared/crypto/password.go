// Package crypto proporciona utilidades criptográficas
package crypto

// PasswordHasher maneja el hashing de passwords con bcrypt
// Será implementado en FASE 2 del Sprint 1
type PasswordHasher interface {
	// Hash genera un hash bcrypt del password
	// Hash(password string) (string, error)

	// Compare compara un password con su hash
	// Compare(password, hash string) error

	// Validate verifica que el password cumple los requisitos
	// Validate(password string) error
}
