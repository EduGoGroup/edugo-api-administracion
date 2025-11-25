// Package repository define las interfaces de persistencia para autenticación
package repository

// UserRepository define las operaciones de persistencia para usuarios
// en el contexto de autenticación
// Será implementado en FASE 2 del Sprint 1
type UserRepository interface {
	// FindByEmail busca un usuario por email para autenticación
	// FindByEmail(ctx context.Context, email string) (*entity.User, error)

	// FindByID busca un usuario por ID
	// FindByID(ctx context.Context, id string) (*entity.User, error)

	// UpdateLastLogin actualiza la fecha de último login
	// UpdateLastLogin(ctx context.Context, userID string) error
}
