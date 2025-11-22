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
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/repository"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
)

type MockSchoolRepository struct {
	mock.Mock
}

func (m *MockSchoolRepository) Create(ctx context.Context, school *entities.School) error {
	args := m.Called(ctx, school)
	return args.Error(0)
}

func (m *MockSchoolRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.School, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.School), args.Error(1)
}

func (m *MockSchoolRepository) FindByCode(ctx context.Context, code string) (*entities.School, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.School), args.Error(1)
}

func (m *MockSchoolRepository) FindByName(ctx context.Context, name string) (*entities.School, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.School), args.Error(1)
}

func (m *MockSchoolRepository) Update(ctx context.Context, school *entities.School) error {
	args := m.Called(ctx, school)
	return args.Error(0)
}

func (m *MockSchoolRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockSchoolRepository) List(ctx context.Context, filters repository.ListFilters) ([]*entities.School, error) {
	args := m.Called(ctx, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.School), args.Error(1)
}

func (m *MockSchoolRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	args := m.Called(ctx, name)
	return args.Bool(0), args.Error(1)
}

func (m *MockSchoolRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	args := m.Called(ctx, code)
	return args.Bool(0), args.Error(1)
}

func TestCreateSchool_Success(t *testing.T) {
	mockRepo := new(MockSchoolRepository)
	mockLogger := newTestLogger()
	service := NewSchoolService(mockRepo, mockLogger)

	req := dto.CreateSchoolRequest{
		Name:    "Test School",
		Code:    "TS001",
		Address: "123 Main St",
	}

	mockRepo.On("ExistsByCode", mock.Anything, req.Code).Return(false, nil)
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*entities.School")).Return(nil)

	result, err := service.CreateSchool(context.Background(), req)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.Name, result.Name)
	assert.Equal(t, req.Code, result.Code)
	mockRepo.AssertExpectations(t)
}

func TestCreateSchool_CodeAlreadyExists(t *testing.T) {
	mockRepo := new(MockSchoolRepository)
	mockLogger := newTestLogger()
	service := NewSchoolService(mockRepo, mockLogger)

	req := dto.CreateSchoolRequest{
		Name:    "Test School",
		Code:    "EXISTING",
		Address: "123 Main St",
	}

	mockRepo.On("ExistsByCode", mock.Anything, req.Code).Return(true, nil)

	_, err := service.CreateSchool(context.Background(), req)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
	mockRepo.AssertNotCalled(t, "Create")
	mockRepo.AssertExpectations(t)
}

func TestUpdateSchool_Success(t *testing.T) {
	mockRepo := new(MockSchoolRepository)
	mockLogger := newTestLogger()
	service := NewSchoolService(mockRepo, mockLogger)

	schoolID := uuid.New()
	existingSchool := &entities.School{
		ID:        schoolID,
		Name:      "Old Name",
		Code:      "SC001",
		Address:   strPtr("Old Address"),
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	newName := "New Name"
	req := dto.UpdateSchoolRequest{
		Name: &newName,
	}

	mockRepo.On("FindByID", mock.Anything, schoolID).Return(existingSchool, nil)
	mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*entities.School")).Return(nil)

	result, err := service.UpdateSchool(context.Background(), schoolID.String(), req)

	require.NoError(t, err)
	assert.Equal(t, newName, result.Name)
	mockRepo.AssertExpectations(t)
}

func TestGetSchool_Success(t *testing.T) {
	mockRepo := new(MockSchoolRepository)
	mockLogger := newTestLogger()
	service := NewSchoolService(mockRepo, mockLogger)

	schoolID := uuid.New()
	school := &entities.School{
		ID:        schoolID,
		Name:      "Test School",
		Code:      "TS001",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("FindByID", mock.Anything, schoolID).Return(school, nil)

	result, err := service.GetSchool(context.Background(), schoolID.String())

	require.NoError(t, err)
	assert.Equal(t, school.Name, result.Name)
	mockRepo.AssertExpectations(t)
}

func strPtr(s string) *string {
	return &s
}
