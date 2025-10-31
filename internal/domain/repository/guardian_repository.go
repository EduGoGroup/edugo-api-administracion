package repository

import (
	"context"

	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/entity"
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/valueobject"
	"github.com/EduGoGroup/edugo-shared/pkg/types"
)

// GuardianRepository define las operaciones de persistencia para GuardianRelation
type GuardianRepository interface {
	// Create crea una nueva relación guardian-estudiante
	Create(ctx context.Context, relation *entity.GuardianRelation) error

	// FindByID busca una relación por ID
	FindByID(ctx context.Context, id types.UUID) (*entity.GuardianRelation, error)

	// FindByGuardianAndStudent busca una relación por guardian y estudiante
	FindByGuardianAndStudent(
		ctx context.Context,
		guardianID valueobject.GuardianID,
		studentID valueobject.StudentID,
	) (*entity.GuardianRelation, error)

	// FindByGuardian busca todas las relaciones de un guardian
	FindByGuardian(ctx context.Context, guardianID valueobject.GuardianID) ([]*entity.GuardianRelation, error)

	// FindByStudent busca todas las relaciones de un estudiante
	FindByStudent(ctx context.Context, studentID valueobject.StudentID) ([]*entity.GuardianRelation, error)

	// Update actualiza una relación existente
	Update(ctx context.Context, relation *entity.GuardianRelation) error

	// Delete elimina una relación (soft delete)
	Delete(ctx context.Context, id types.UUID) error

	// ExistsActiveRelation verifica si existe una relación activa entre guardian y estudiante
	ExistsActiveRelation(
		ctx context.Context,
		guardianID valueobject.GuardianID,
		studentID valueobject.StudentID,
	) (bool, error)
}
