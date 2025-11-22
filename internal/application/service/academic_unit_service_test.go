package service

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/EduGoGroup/edugo-api-administracion/internal/application/dto"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
)

type MockAcademicUnitRepository struct {
	mock.Mock
}

func (m *MockAcademicUnitRepository) Create(ctx context.Context, unit *entities.AcademicUnit) error {
	args := m.Called(ctx, unit)
	return args.Error(0)
}

func (m *MockAcademicUnitRepository) FindByID(ctx context.Context, id uuid.UUID, includeDeleted bool) (*entities.AcademicUnit, error) {
	args := m.Called(ctx, id, includeDeleted)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.AcademicUnit), args.Error(1)
}

func (m *MockAcademicUnitRepository) FindBySchoolIDAndCode(ctx context.Context, schoolID uuid.UUID, code string) (*entities.AcademicUnit, error) {
	args := m.Called(ctx, schoolID, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.AcademicUnit), args.Error(1)
}

func (m *MockAcademicUnitRepository) FindBySchoolID(ctx context.Context, schoolID uuid.UUID, includeDeleted bool) ([]*entities.AcademicUnit, error) {
	args := m.Called(ctx, schoolID, includeDeleted)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.AcademicUnit), args.Error(1)
}

func (m *MockAcademicUnitRepository) FindByParentID(ctx context.Context, parentID uuid.UUID, includeDeleted bool) ([]*entities.AcademicUnit, error) {
	args := m.Called(ctx, parentID, includeDeleted)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.AcademicUnit), args.Error(1)
}

func (m *MockAcademicUnitRepository) FindRootUnits(ctx context.Context, schoolID uuid.UUID) ([]*entities.AcademicUnit, error) {
	args := m.Called(ctx, schoolID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.AcademicUnit), args.Error(1)
}

func (m *MockAcademicUnitRepository) FindByType(ctx context.Context, schoolID uuid.UUID, unitType string, includeDeleted bool) ([]*entities.AcademicUnit, error) {
	args := m.Called(ctx, schoolID, unitType, includeDeleted)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.AcademicUnit), args.Error(1)
}

func (m *MockAcademicUnitRepository) Update(ctx context.Context, unit *entities.AcademicUnit) error {
	args := m.Called(ctx, unit)
	return args.Error(0)
}

func (m *MockAcademicUnitRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAcademicUnitRepository) Restore(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAcademicUnitRepository) HardDelete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAcademicUnitRepository) GetHierarchyPath(ctx context.Context, id uuid.UUID) ([]*entities.AcademicUnit, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.AcademicUnit), args.Error(1)
}

func (m *MockAcademicUnitRepository) ExistsBySchoolIDAndCode(ctx context.Context, schoolID uuid.UUID, code string) (bool, error) {
	args := m.Called(ctx, schoolID, code)
	return args.Bool(0), args.Error(1)
}

func TestGetUnit_Success(t *testing.T) {
	mockUnitRepo := new(MockAcademicUnitRepository)
	mockSchoolRepo := new(MockSchoolRepository)
	mockLogger := newTestLogger()
	service := NewAcademicUnitService(mockUnitRepo, mockSchoolRepo, mockLogger)

	unitID := uuid.New()
	unit := &entities.AcademicUnit{
		ID:        unitID,
		SchoolID:  uuid.New(),
		Name:      "Grade 1",
		Code:      "G1",
		Type:      "grade",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockUnitRepo.On("FindByID", mock.Anything, unitID, false).Return(unit, nil)

	result, err := service.GetUnit(context.Background(), unitID.String())

	require.NoError(t, err)
	assert.Equal(t, "Grade 1", result.DisplayName)
	mockUnitRepo.AssertExpectations(t)
}

func TestGetUnitTree_Success(t *testing.T) {
	mockUnitRepo := new(MockAcademicUnitRepository)
	mockSchoolRepo := new(MockSchoolRepository)
	mockLogger := newTestLogger()
	service := NewAcademicUnitService(mockUnitRepo, mockSchoolRepo, mockLogger)

	schoolID := uuid.New()
	units := []*entities.AcademicUnit{
		{
			ID:        uuid.New(),
			SchoolID:  schoolID,
			Name:      "Grade 1",
			Code:      "G1",
			Type:      "grade",
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	mockUnitRepo.On("FindBySchoolID", mock.Anything, schoolID, false).Return(units, nil)

	result, err := service.GetUnitTree(context.Background(), schoolID.String())

	require.NoError(t, err)
	assert.Len(t, result, 1)
	mockUnitRepo.AssertExpectations(t)
}

func TestDeleteUnit_Success(t *testing.T) {
	mockUnitRepo := new(MockAcademicUnitRepository)
	mockSchoolRepo := new(MockSchoolRepository)
	mockLogger := newTestLogger()
	service := NewAcademicUnitService(mockUnitRepo, mockSchoolRepo, mockLogger)

	unitID := uuid.New()
	unit := &entities.AcademicUnit{
		ID:        unitID,
		SchoolID:  uuid.New(),
		Name:      "Grade 1",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockUnitRepo.On("FindByID", mock.Anything, unitID, false).Return(unit, nil)
	mockUnitRepo.On("SoftDelete", mock.Anything, unitID).Return(nil)

	err := service.DeleteUnit(context.Background(), unitID.String())

	require.NoError(t, err)
	mockUnitRepo.AssertExpectations(t)
}
