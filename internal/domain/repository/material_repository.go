package repository

import (
	"context"

	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/valueobject"
)

// MaterialRepository define las operaciones de persistencia para Material
// Nota: Por ahora solo implementamos Delete, ya que los demás endpoints
// están en api-mobile
type MaterialRepository interface {
	// Delete elimina un material (soft delete)
	Delete(ctx context.Context, id valueobject.MaterialID) error

	// Exists verifica si un material existe
	Exists(ctx context.Context, id valueobject.MaterialID) (bool, error)
}
