package repository

import (
	"context"

	"github.com/edugo/api-administracion/internal/domain/entity"
	"github.com/edugo/api-administracion/internal/domain/valueobject"
)

// SubjectRepository define las operaciones de persistencia para Subject
type SubjectRepository interface {
	Create(ctx context.Context, subject *entity.Subject) error
	FindByID(ctx context.Context, id valueobject.SubjectID) (*entity.Subject, error)
	Update(ctx context.Context, subject *entity.Subject) error
	Delete(ctx context.Context, id valueobject.SubjectID) error
}
