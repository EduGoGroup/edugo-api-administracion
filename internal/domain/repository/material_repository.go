package repository

import (
	"context"

	"github.com/google/uuid"
)

// MaterialRepository define las operaciones de persistencia para Material
// Nota: Por ahora solo implementamos Delete, ya que los demás endpoints
// están en api-mobile
type MaterialRepository interface {
	// Delete elimina un material (soft delete)
	Delete(ctx context.Context, id uuid.UUID) error

	// Exists verifica si un material existe
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
}
