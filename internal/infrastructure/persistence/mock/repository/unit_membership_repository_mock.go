package repository

import (
	"context"
	"sync"
	"time"

	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/repository"
	mockData "github.com/EduGoGroup/edugo-api-administracion/internal/infrastructure/persistence/mock/data"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-shared/common/errors"
	"github.com/google/uuid"
)

// MockUnitMembershipRepository es una implementación en memoria del UnitMembershipRepository para testing
type MockUnitMembershipRepository struct {
	mu          sync.RWMutex
	memberships map[uuid.UUID]*entities.Membership
}

// NewMockUnitMembershipRepository crea una nueva instancia de MockUnitMembershipRepository
// Pre-carga las memberships desde mockData.GetMemberships()
func NewMockUnitMembershipRepository() repository.UnitMembershipRepository {
	memberships := make(map[uuid.UUID]*entities.Membership)

	// Pre-cargar datos desde mockData
	mockMemberships := mockData.GetMemberships()
	for id, membership := range mockMemberships {
		// Hacer una copia de la membership para evitar modificaciones externas
		membershipCopy := *membership
		// Copiar también el puntero de AcademicUnitID si existe
		if membership.AcademicUnitID != nil {
			unitID := *membership.AcademicUnitID
			membershipCopy.AcademicUnitID = &unitID
		}
		// Copiar también el puntero de WithdrawnAt si existe
		if membership.WithdrawnAt != nil {
			withdrawnAt := *membership.WithdrawnAt
			membershipCopy.WithdrawnAt = &withdrawnAt
		}
		// Copiar metadata si existe
		if membership.Metadata != nil {
			metadataCopy := make([]byte, len(membership.Metadata))
			copy(metadataCopy, membership.Metadata)
			membershipCopy.Metadata = metadataCopy
		}
		memberships[id] = &membershipCopy
	}

	return &MockUnitMembershipRepository{
		memberships: memberships,
	}
}

// Create crea una nueva membership en el repositorio
func (r *MockUnitMembershipRepository) Create(ctx context.Context, membership *entities.Membership) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Validar que no exista una membership activa con el mismo usuario y unidad
	for _, existingMembership := range r.memberships {
		if existingMembership.UserID == membership.UserID &&
			existingMembership.SchoolID == membership.SchoolID &&
			existingMembership.IsActive &&
			existingMembership.WithdrawnAt == nil {
			// Comparar AcademicUnitID considerando que puede ser nil
			if (existingMembership.AcademicUnitID == nil && membership.AcademicUnitID == nil) ||
				(existingMembership.AcademicUnitID != nil && membership.AcademicUnitID != nil &&
					*existingMembership.AcademicUnitID == *membership.AcademicUnitID) {
				return errors.NewConflictError("membership already exists for this user and unit")
			}
		}
	}

	// Generar ID si no existe
	if membership.ID == uuid.Nil {
		membership.ID = uuid.New()
	}

	// Establecer timestamps
	now := time.Now()
	membership.CreatedAt = now
	membership.UpdatedAt = now

	// Guardar una copia de la membership
	membershipCopy := *membership
	// Copiar también el puntero de AcademicUnitID si existe
	if membership.AcademicUnitID != nil {
		unitID := *membership.AcademicUnitID
		membershipCopy.AcademicUnitID = &unitID
	}
	// Copiar también el puntero de WithdrawnAt si existe
	if membership.WithdrawnAt != nil {
		withdrawnAt := *membership.WithdrawnAt
		membershipCopy.WithdrawnAt = &withdrawnAt
	}
	// Copiar metadata si existe
	if membership.Metadata != nil {
		metadataCopy := make([]byte, len(membership.Metadata))
		copy(metadataCopy, membership.Metadata)
		membershipCopy.Metadata = metadataCopy
	}

	r.memberships[membership.ID] = &membershipCopy

	return nil
}

// FindByID busca una membership por ID
func (r *MockUnitMembershipRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Membership, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	membership, exists := r.memberships[id]
	if !exists {
		return nil, errors.NewNotFoundError("membership not found")
	}

	// Retornar una copia para evitar modificaciones externas
	return r.copyMembership(membership), nil
}

// FindByUserAndUnit busca una membership por usuario y unidad académica
func (r *MockUnitMembershipRepository) FindByUserAndUnit(ctx context.Context, userID, unitID uuid.UUID) (*entities.Membership, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, membership := range r.memberships {
		if membership.UserID == userID &&
			membership.AcademicUnitID != nil &&
			*membership.AcademicUnitID == unitID {
			// Retornar una copia para evitar modificaciones externas
			return r.copyMembership(membership), nil
		}
	}

	return nil, errors.NewNotFoundError("membership not found")
}

// FindByUser busca todas las memberships de un usuario
func (r *MockUnitMembershipRepository) FindByUser(ctx context.Context, userID uuid.UUID) ([]*entities.Membership, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*entities.Membership

	for _, membership := range r.memberships {
		if membership.UserID == userID {
			// Agregar copia de la membership
			result = append(result, r.copyMembership(membership))
		}
	}

	return result, nil
}

// FindByUnit busca todas las memberships de una unidad académica
func (r *MockUnitMembershipRepository) FindByUnit(ctx context.Context, unitID uuid.UUID) ([]*entities.Membership, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*entities.Membership

	for _, membership := range r.memberships {
		if membership.AcademicUnitID != nil && *membership.AcademicUnitID == unitID {
			// Agregar copia de la membership
			result = append(result, r.copyMembership(membership))
		}
	}

	return result, nil
}

// Update actualiza una membership existente
func (r *MockUnitMembershipRepository) Update(ctx context.Context, membership *entities.Membership) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Verificar que la membership existe
	existingMembership, exists := r.memberships[membership.ID]
	if !exists {
		return errors.NewNotFoundError("membership not found")
	}

	// Validar que no haya conflicto con otra membership (mismo user + unit)
	for id, m := range r.memberships {
		if id != membership.ID && m.UserID == membership.UserID && m.SchoolID == membership.SchoolID {
			// Comparar AcademicUnitID considerando que puede ser nil
			if (m.AcademicUnitID == nil && membership.AcademicUnitID == nil) ||
				(m.AcademicUnitID != nil && membership.AcademicUnitID != nil &&
					*m.AcademicUnitID == *membership.AcademicUnitID) {
				if m.IsActive && m.WithdrawnAt == nil {
					return errors.NewConflictError("membership already exists for this user and unit")
				}
			}
		}
	}

	// Actualizar timestamp
	membership.UpdatedAt = time.Now()

	// Preservar CreatedAt original
	membership.CreatedAt = existingMembership.CreatedAt

	// Guardar una copia de la membership actualizada
	membershipCopy := *membership
	// Copiar también el puntero de AcademicUnitID si existe
	if membership.AcademicUnitID != nil {
		unitID := *membership.AcademicUnitID
		membershipCopy.AcademicUnitID = &unitID
	}
	// Copiar también el puntero de WithdrawnAt si existe
	if membership.WithdrawnAt != nil {
		withdrawnAt := *membership.WithdrawnAt
		membershipCopy.WithdrawnAt = &withdrawnAt
	}
	// Copiar metadata si existe
	if membership.Metadata != nil {
		metadataCopy := make([]byte, len(membership.Metadata))
		copy(metadataCopy, membership.Metadata)
		membershipCopy.Metadata = metadataCopy
	}

	r.memberships[membership.ID] = &membershipCopy

	return nil
}

// Delete elimina una membership
func (r *MockUnitMembershipRepository) Delete(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.memberships[id]
	if !exists {
		return errors.NewNotFoundError("membership not found")
	}

	// Eliminar la membership del mapa
	delete(r.memberships, id)

	return nil
}

// ExistsByUnitAndUser verifica si existe una membership para una unidad y usuario
func (r *MockUnitMembershipRepository) ExistsByUnitAndUser(ctx context.Context, unitID, userID uuid.UUID) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, membership := range r.memberships {
		if membership.UserID == userID &&
			membership.AcademicUnitID != nil &&
			*membership.AcademicUnitID == unitID {
			return true, nil
		}
	}

	return false, nil
}

// Reset reinicia el repositorio a su estado inicial (útil para testing)
func (r *MockUnitMembershipRepository) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Recargar datos desde mockData
	memberships := make(map[uuid.UUID]*entities.Membership)
	mockMemberships := mockData.GetMemberships()
	for id, membership := range mockMemberships {
		membershipCopy := *membership
		// Copiar también el puntero de AcademicUnitID si existe
		if membership.AcademicUnitID != nil {
			unitID := *membership.AcademicUnitID
			membershipCopy.AcademicUnitID = &unitID
		}
		// Copiar también el puntero de WithdrawnAt si existe
		if membership.WithdrawnAt != nil {
			withdrawnAt := *membership.WithdrawnAt
			membershipCopy.WithdrawnAt = &withdrawnAt
		}
		// Copiar metadata si existe
		if membership.Metadata != nil {
			metadataCopy := make([]byte, len(membership.Metadata))
			copy(metadataCopy, membership.Metadata)
			membershipCopy.Metadata = metadataCopy
		}
		memberships[id] = &membershipCopy
	}

	r.memberships = memberships
}

// copyMembership es una función auxiliar que crea una copia profunda de una membership
func (r *MockUnitMembershipRepository) copyMembership(membership *entities.Membership) *entities.Membership {
	membershipCopy := *membership

	// Copiar también el puntero de AcademicUnitID si existe
	if membership.AcademicUnitID != nil {
		unitID := *membership.AcademicUnitID
		membershipCopy.AcademicUnitID = &unitID
	}

	// Copiar también el puntero de WithdrawnAt si existe
	if membership.WithdrawnAt != nil {
		withdrawnAt := *membership.WithdrawnAt
		membershipCopy.WithdrawnAt = &withdrawnAt
	}

	// Copiar metadata si existe
	if membership.Metadata != nil {
		metadataCopy := make([]byte, len(membership.Metadata))
		copy(metadataCopy, membership.Metadata)
		membershipCopy.Metadata = metadataCopy
	}

	return &membershipCopy
}
