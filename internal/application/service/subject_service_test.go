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

type MockSubjectRepository struct {
	mock.Mock
}

func (m *MockSubjectRepository) Create(ctx context.Context, subject *entities.Subject) error {
	args := m.Called(ctx, subject)
	return args.Error(0)
}

func (m *MockSubjectRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Subject, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Subject), args.Error(1)
}

func (m *MockSubjectRepository) Update(ctx context.Context, subject *entities.Subject) error {
	args := m.Called(ctx, subject)
	return args.Error(0)
}

func (m *MockSubjectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockSubjectRepository) List(ctx context.Context) ([]*entities.Subject, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Subject), args.Error(1)
}

func (m *MockSubjectRepository) FindBySchoolID(ctx context.Context, schoolID uuid.UUID) ([]*entities.Subject, error) {
	args := m.Called(ctx, schoolID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Subject), args.Error(1)
}

func (m *MockSubjectRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	args := m.Called(ctx, name)
	return args.Bool(0), args.Error(1)
}

func TestCreateSubject_Success(t *testing.T) {
	mockRepo := new(MockSubjectRepository)
	mockLogger := newTestLogger()
	service := NewSubjectService(mockRepo, mockLogger)

	req := dto.CreateSubjectRequest{
		Name:        "Mathematics",
		Description: "Basic math",
		Metadata:    "{}",
	}

	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*entities.Subject")).Return(nil)

	result, err := service.CreateSubject(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, req.Name, result.Name)
	mockRepo.AssertExpectations(t)
}

func TestUpdateSubject_Success(t *testing.T) {
	mockRepo := new(MockSubjectRepository)
	mockLogger := newTestLogger()
	service := NewSubjectService(mockRepo, mockLogger)

	subjectID := uuid.New()
	existing := &entities.Subject{
		ID:        subjectID,
		Name:      "Old Name",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	newName := "New Name"
	req := dto.UpdateSubjectRequest{
		Name: &newName,
	}

	mockRepo.On("FindByID", mock.Anything, subjectID).Return(existing, nil)
	mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*entities.Subject")).Return(nil)

	result, err := service.UpdateSubject(context.Background(), subjectID.String(), req)

	require.NoError(t, err)
	assert.Equal(t, newName, result.Name)
	mockRepo.AssertExpectations(t)
}

func TestGetSubject_Success(t *testing.T) {
	mockRepo := new(MockSubjectRepository)
	mockLogger := newTestLogger()
	service := NewSubjectService(mockRepo, mockLogger)

	subjectID := uuid.New()
	subject := &entities.Subject{
		ID:        subjectID,
		Name:      "Math",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("FindByID", mock.Anything, subjectID).Return(subject, nil)

	result, err := service.GetSubject(context.Background(), subjectID.String())

	require.NoError(t, err)
	assert.Equal(t, "Math", result.Name)
	mockRepo.AssertExpectations(t)
}
