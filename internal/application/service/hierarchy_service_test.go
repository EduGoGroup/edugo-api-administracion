package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/entity"
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/service"
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/valueobject"
)

// MockUnitRepository mock del repositorio
type MockUnitRepository struct {
	mock.Mock
}

func (m *MockUnitRepository) Create(ctx context.Context, unit *entity.AcademicUnit) error {
	args := m.Called(ctx, unit)
	return args.Error(0)
}

func (m *MockUnitRepository) FindByID(ctx context.Context, id valueobject.UnitID, includeDeleted bool) (*entity.AcademicUnit, error) {
	args := m.Called(ctx, id, includeDeleted)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.AcademicUnit), args.Error(1)
}

func (m *MockUnitRepository) ExistsBySchoolIDAndCode(ctx context.Context, schoolID valueobject.SchoolID, code string) (bool, error) {
	args := m.Called(ctx, schoolID, code)
	return args.Bool(0), args.Error(1)
}

func (m *MockUnitRepository) FindDescendants(ctx context.Context, unitID valueobject.UnitID) ([]*entity.AcademicUnit, error) {
	args := m.Called(ctx, unitID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.AcademicUnit), args.Error(1)
}

func (m *MockUnitRepository) MoveSubtree(ctx context.Context, unitID valueobject.UnitID, newParentID *valueobject.UnitID) error {
	args := m.Called(ctx, unitID, newParentID)
	return args.Error(0)
}

func (m *MockUnitRepository) Update(ctx context.Context, unit *entity.AcademicUnit) error {
	args := m.Called(ctx, unit)
	return args.Error(0)
}

func (m *MockUnitRepository) FindBySchoolID(ctx context.Context, schoolID valueobject.SchoolID, includeDeleted bool) ([]*entity.AcademicUnit, error) {
	args := m.Called(ctx, schoolID, includeDeleted)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.AcademicUnit), args.Error(1)
}

func (m *MockUnitRepository) FindBySchoolIDAndDepth(ctx context.Context, schoolID valueobject.SchoolID, depth int) ([]*entity.AcademicUnit, error) {
	args := m.Called(ctx, schoolID, depth)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.AcademicUnit), args.Error(1)
}

func (m *MockUnitRepository) FindByType(ctx context.Context, schoolID valueobject.SchoolID, unitType valueobject.UnitType, includeDeleted bool) ([]*entity.AcademicUnit, error) {
	args := m.Called(ctx, schoolID, unitType, includeDeleted)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.AcademicUnit), args.Error(1)
}

func (m *MockUnitRepository) SoftDelete(ctx context.Context, id valueobject.UnitID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUnitRepository) Restore(ctx context.Context, id valueobject.UnitID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUnitRepository) HasChildren(ctx context.Context, id valueobject.UnitID) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockUnitRepository) GetHierarchyPath(ctx context.Context, id valueobject.UnitID) ([]*entity.AcademicUnit, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.AcademicUnit), args.Error(1)
}

// MockSchoolRepository mock del repositorio de escuelas
type MockSchoolRepository struct {
	mock.Mock
}

func (m *MockSchoolRepository) Create(ctx context.Context, school *entity.School) error {
	args := m.Called(ctx, school)
	return args.Error(0)
}

func (m *MockSchoolRepository) FindByID(ctx context.Context, id valueobject.SchoolID) (*entity.School, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.School), args.Error(1)
}

func (m *MockSchoolRepository) FindAll(ctx context.Context) ([]*entity.School, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.School), args.Error(1)
}

func (m *MockSchoolRepository) Update(ctx context.Context, school *entity.School) error {
	args := m.Called(ctx, school)
	return args.Error(0)
}

func (m *MockSchoolRepository) Delete(ctx context.Context, id valueobject.SchoolID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Tests
func TestHierarchyService_CreateUnit(t *testing.T) {
	ctx := context.Background()

	t.Run("should create unit without parent", func(t *testing.T) {
		// Setup
		mockUnitRepo := new(MockUnitRepository)
		mockSchoolRepo := new(MockSchoolRepository)
		domainService := service.NewAcademicUnitDomainService()
		hierarchyService := NewHierarchyService(mockUnitRepo, mockSchoolRepo, domainService)

		schoolID := valueobject.NewSchoolID()
		school, _ := entity.NewSchool("Test School", "TS001", "Address")

		mockSchoolRepo.On("FindByID", ctx, schoolID).Return(school, nil)
		mockUnitRepo.On("ExistsBySchoolIDAndCode", ctx, schoolID, "G1").Return(false, nil)
		mockUnitRepo.On("Create", ctx, mock.AnythingOfType("*entity.AcademicUnit")).Return(nil)

		// Execute
		unitType, _ := valueobject.NewUnitType("grade")
		unit, err := hierarchyService.CreateUnit(ctx, nil, schoolID, unitType, "Grade 1", "G1", "")

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, unit)
		assert.Equal(t, "Grade 1", unit.DisplayName())
		mockSchoolRepo.AssertExpectations(t)
		mockUnitRepo.AssertExpectations(t)
	})

	// TODO_FASE2: Agregar más tests cuando tengamos DB real
	// - should create unit with parent
	// - should fail when school not found
	// - should fail when code already exists
	// - should fail when parent not in same school
}

func TestHierarchyService_ValidateNoCircularReference(t *testing.T) {
	ctx := context.Background()

	t.Run("should detect circular reference", func(t *testing.T) {
		// Setup
		mockUnitRepo := new(MockUnitRepository)
		mockSchoolRepo := new(MockSchoolRepository)
		domainService := service.NewAcademicUnitDomainService()
		hierarchyService := NewHierarchyService(mockUnitRepo, mockSchoolRepo, domainService)

		unitID := valueobject.NewUnitID()

		schoolID := valueobject.NewSchoolID()
		unitType, _ := valueobject.NewUnitType("grade")

		// El descendiente es el que queremos establecer como padre (ciclo!)
		descendant, _ := entity.NewAcademicUnit(schoolID, unitType, "Descendant", "DESC")

		mockUnitRepo.On("FindDescendants", ctx, unitID).Return([]*entity.AcademicUnit{descendant}, nil)

		// Execute
		err := hierarchyService.ValidateNoCircularReference(ctx, unitID, descendant.ID())

		// Assert
		require.Error(t, err)
		assert.Contains(t, err.Error(), "circular reference")
		mockUnitRepo.AssertExpectations(t)
	})

	t.Run("should detect self reference", func(t *testing.T) {
		// Setup
		mockUnitRepo := new(MockUnitRepository)
		mockSchoolRepo := new(MockSchoolRepository)
		domainService := service.NewAcademicUnitDomainService()
		hierarchyService := NewHierarchyService(mockUnitRepo, mockSchoolRepo, domainService)

		unitID := valueobject.NewUnitID()

		mockUnitRepo.On("FindDescendants", ctx, unitID).Return([]*entity.AcademicUnit{}, nil)

		// Execute
		err := hierarchyService.ValidateNoCircularReference(ctx, unitID, unitID)

		// Assert
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be its own parent")
		mockUnitRepo.AssertExpectations(t)
	})
}
