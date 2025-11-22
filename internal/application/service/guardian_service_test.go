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

type MockGuardianRepository struct {
	mock.Mock
}

func (m *MockGuardianRepository) Create(ctx context.Context, relation *entities.GuardianRelation) error {
	args := m.Called(ctx, relation)
	return args.Error(0)
}

func (m *MockGuardianRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.GuardianRelation, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.GuardianRelation), args.Error(1)
}

func (m *MockGuardianRepository) FindByGuardian(ctx context.Context, guardianID uuid.UUID) ([]*entities.GuardianRelation, error) {
	args := m.Called(ctx, guardianID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.GuardianRelation), args.Error(1)
}

func (m *MockGuardianRepository) FindByStudent(ctx context.Context, studentID uuid.UUID) ([]*entities.GuardianRelation, error) {
	args := m.Called(ctx, studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.GuardianRelation), args.Error(1)
}

func (m *MockGuardianRepository) Update(ctx context.Context, relation *entities.GuardianRelation) error {
	args := m.Called(ctx, relation)
	return args.Error(0)
}

func (m *MockGuardianRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockGuardianRepository) ExistsActiveRelation(ctx context.Context, guardianID, studentID uuid.UUID) (bool, error) {
	args := m.Called(ctx, guardianID, studentID)
	return args.Bool(0), args.Error(1)
}

func (m *MockGuardianRepository) CreateRelation(ctx context.Context, relation *entities.GuardianRelation) error {
	return m.Create(ctx, relation)
}

func (m *MockGuardianRepository) FindRelationByID(ctx context.Context, id uuid.UUID) (*entities.GuardianRelation, error) {
	return m.FindByID(ctx, id)
}

func (m *MockGuardianRepository) FindRelationsByGuardian(ctx context.Context, guardianID uuid.UUID) ([]*entities.GuardianRelation, error) {
	return m.FindByGuardian(ctx, guardianID)
}

func (m *MockGuardianRepository) FindRelationsByStudent(ctx context.Context, studentID uuid.UUID) ([]*entities.GuardianRelation, error) {
	return m.FindByStudent(ctx, studentID)
}

func (m *MockGuardianRepository) UpdateRelation(ctx context.Context, relation *entities.GuardianRelation) error {
	return m.Update(ctx, relation)
}

func (m *MockGuardianRepository) DeleteRelation(ctx context.Context, id uuid.UUID) error {
	return m.Delete(ctx, id)
}

func TestCreateGuardianRelation_Success(t *testing.T) {
	mockRepo := new(MockGuardianRepository)
	mockLogger := newTestLogger()
	service := NewGuardianService(mockRepo, mockLogger)

	guardianID := uuid.New()
	studentID := uuid.New()

	req := dto.CreateGuardianRelationRequest{
		GuardianID:       guardianID.String(),
		StudentID:        studentID.String(),
		RelationshipType: "father",
	}

	mockRepo.On("ExistsActiveRelation", mock.Anything, guardianID, studentID).Return(false, nil)
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*entities.GuardianRelation")).Return(nil)

	result, err := service.CreateGuardianRelation(context.Background(), req, "admin@test.com")

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, guardianID.String(), result.GuardianID)
	assert.Equal(t, studentID.String(), result.StudentID)
	assert.Equal(t, "father", result.RelationshipType)
	assert.True(t, result.IsActive)
	mockRepo.AssertExpectations(t)
}

func TestCreateGuardianRelation_AlreadyExists(t *testing.T) {
	mockRepo := new(MockGuardianRepository)
	mockLogger := newTestLogger()
	service := NewGuardianService(mockRepo, mockLogger)

	guardianID := uuid.New()
	studentID := uuid.New()

	req := dto.CreateGuardianRelationRequest{
		GuardianID:       guardianID.String(),
		StudentID:        studentID.String(),
		RelationshipType: "mother",
	}

	mockRepo.On("ExistsActiveRelation", mock.Anything, guardianID, studentID).Return(true, nil)

	_, err := service.CreateGuardianRelation(context.Background(), req, "admin@test.com")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
	mockRepo.AssertNotCalled(t, "Create")
	mockRepo.AssertExpectations(t)
}

func TestCreateGuardianRelation_InvalidRelationshipType(t *testing.T) {
	mockRepo := new(MockGuardianRepository)
	mockLogger := newTestLogger()
	service := NewGuardianService(mockRepo, mockLogger)

	guardianID := uuid.New()
	studentID := uuid.New()

	req := dto.CreateGuardianRelationRequest{
		GuardianID:       guardianID.String(),
		StudentID:        studentID.String(),
		RelationshipType: "invalid_type",
	}

	mockRepo.On("ExistsActiveRelation", mock.Anything, guardianID, studentID).Return(false, nil)

	_, err := service.CreateGuardianRelation(context.Background(), req, "admin@test.com")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "relationship_type")
	mockRepo.AssertNotCalled(t, "Create")
}

func TestGetGuardianRelation_Success(t *testing.T) {
	mockRepo := new(MockGuardianRepository)
	mockLogger := newTestLogger()
	service := NewGuardianService(mockRepo, mockLogger)

	relationID := uuid.New()
	guardianID := uuid.New()
	studentID := uuid.New()

	relation := &entities.GuardianRelation{
		ID:               relationID,
		GuardianID:       guardianID,
		StudentID:        studentID,
		RelationshipType: "father",
		IsActive:         true,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		CreatedBy:        "admin@test.com",
	}

	mockRepo.On("FindByID", mock.Anything, relationID).Return(relation, nil)

	result, err := service.GetGuardianRelation(context.Background(), relationID.String())

	require.NoError(t, err)
	assert.Equal(t, relationID.String(), result.ID)
	assert.Equal(t, "father", result.RelationshipType)
	mockRepo.AssertExpectations(t)
}

func TestGetGuardianRelations_Success(t *testing.T) {
	mockRepo := new(MockGuardianRepository)
	mockLogger := newTestLogger()
	service := NewGuardianService(mockRepo, mockLogger)

	guardianID := uuid.New()
	relations := []*entities.GuardianRelation{
		{
			ID:               uuid.New(),
			GuardianID:       guardianID,
			StudentID:        uuid.New(),
			RelationshipType: "father",
			IsActive:         true,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		},
	}

	mockRepo.On("FindByGuardian", mock.Anything, guardianID).Return(relations, nil)

	results, err := service.GetGuardianRelations(context.Background(), guardianID.String())

	require.NoError(t, err)
	assert.Len(t, results, 1)
	mockRepo.AssertExpectations(t)
}
