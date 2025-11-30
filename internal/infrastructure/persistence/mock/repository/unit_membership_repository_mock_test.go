package repository

import (
	"context"
	"testing"
	"time"

	mockData "github.com/EduGoGroup/edugo-api-administracion/internal/infrastructure/persistence/mock/data"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
)

func TestMockUnitMembershipRepository_Create(t *testing.T) {
	repo := NewMockUnitMembershipRepository()
	ctx := context.Background()

	unitID := uuid.New()
	newMembership := &entities.Membership{
		UserID:         uuid.New(),
		SchoolID:       uuid.New(),
		AcademicUnitID: &unitID,
		Role:           "student",
		IsActive:       true,
		EnrolledAt:     time.Now(),
		Metadata:       []byte(`{"test": "data"}`),
	}

	err := repo.Create(ctx, newMembership)
	if err != nil {
		t.Fatalf("Error al crear membership: %v", err)
	}

	if newMembership.ID == uuid.Nil {
		t.Error("El ID no fue generado")
	}
}

func TestMockUnitMembershipRepository_FindByID(t *testing.T) {
	repo := NewMockUnitMembershipRepository()
	ctx := context.Background()

	membership, err := repo.FindByID(ctx, mockData.MembershipTeacherMaria)
	if err != nil {
		t.Fatalf("Error al buscar membership: %v", err)
	}

	if membership.ID != mockData.MembershipTeacherMaria {
		t.Errorf("ID incorrecto: esperado %s, obtenido %s", mockData.MembershipTeacherMaria, membership.ID)
	}
}

func TestMockUnitMembershipRepository_FindByUser(t *testing.T) {
	repo := NewMockUnitMembershipRepository()
	ctx := context.Background()

	memberships, err := repo.FindByUser(ctx, mockData.TeacherMathID)
	if err != nil {
		t.Fatalf("Error al buscar memberships por usuario: %v", err)
	}

	if len(memberships) == 0 {
		t.Error("No se encontraron memberships para el usuario")
	}
}

func TestMockUnitMembershipRepository_ExistsByUnitAndUser(t *testing.T) {
	repo := NewMockUnitMembershipRepository()
	ctx := context.Background()

	// UUID de seccionPrimer1A (c4eebc99-9c0b-4ef8-bb6d-6bb9bd380a44)
	seccionPrimer1A := uuid.MustParse("c4eebc99-9c0b-4ef8-bb6d-6bb9bd380a44")

	exists, err := repo.ExistsByUnitAndUser(ctx, seccionPrimer1A, mockData.TeacherMathID)
	if err != nil {
		t.Fatalf("Error al verificar existencia: %v", err)
	}

	if !exists {
		t.Error("Debería existir una membership para este usuario y unidad")
	}
}

func TestMockUnitMembershipRepository_Reset(t *testing.T) {
	repo := NewMockUnitMembershipRepository().(*MockUnitMembershipRepository)
	ctx := context.Background()

	// Eliminar una membership
	err := repo.Delete(ctx, mockData.MembershipTeacherMaria)
	if err != nil {
		t.Fatalf("Error al eliminar membership: %v", err)
	}

	// Verificar que fue eliminada
	_, err = repo.FindByID(ctx, mockData.MembershipTeacherMaria)
	if err == nil {
		t.Error("La membership debería estar eliminada")
	}

	// Reset
	repo.Reset()

	// Verificar que fue restaurada
	_, err = repo.FindByID(ctx, mockData.MembershipTeacherMaria)
	if err != nil {
		t.Error("La membership debería estar restaurada después del reset")
	}
}
