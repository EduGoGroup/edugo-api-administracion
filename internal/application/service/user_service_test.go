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
	"github.com/EduGoGroup/edugo-shared/common/types/enum"
)

// MockUserRepository mock implementation
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *entities.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *entities.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) List(ctx context.Context, filters repository.ListFilters) ([]*entities.User, error) {
	args := m.Called(ctx, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.User), args.Error(1)
}

func (m *MockUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

// Tests

func TestCreateUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockLogger := newTestLogger()
	service := NewUserService(mockRepo, mockLogger)

	req := dto.CreateUserRequest{
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Role:      "teacher",
	}

	mockRepo.On("ExistsByEmail", mock.Anything, req.Email).Return(false, nil)
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*entities.User")).Return(nil)

	result, err := service.CreateUser(context.Background(), req)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.Email, result.Email)
	assert.Equal(t, req.FirstName, result.FirstName)
	assert.Equal(t, req.LastName, result.LastName)
	assert.Equal(t, "teacher", result.Role)
	assert.True(t, result.IsActive)
	mockRepo.AssertExpectations(t)
}

func TestCreateUser_InvalidEmail(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockLogger := newTestLogger()
	service := NewUserService(mockRepo, mockLogger)

	req := dto.CreateUserRequest{
		Email:     "", // Email vacío
		FirstName: "John",
		LastName:  "Doe",
		Role:      "teacher",
	}

	_, err := service.CreateUser(context.Background(), req)

	require.Error(t, err)
	mockRepo.AssertNotCalled(t, "Create")
}

func TestCreateUser_EmailAlreadyExists(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockLogger := newTestLogger()
	service := NewUserService(mockRepo, mockLogger)

	req := dto.CreateUserRequest{
		Email:     "existing@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Role:      "teacher",
	}

	mockRepo.On("ExistsByEmail", mock.Anything, req.Email).Return(true, nil)

	_, err := service.CreateUser(context.Background(), req)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
	mockRepo.AssertNotCalled(t, "Create")
	mockRepo.AssertExpectations(t)
}

func TestCreateUser_CannotCreateAdmin(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockLogger := newTestLogger()
	service := NewUserService(mockRepo, mockLogger)

	req := dto.CreateUserRequest{
		Email:     "admin@example.com",
		FirstName: "Admin",
		LastName:  "User",
		Role:      string(enum.SystemRoleAdmin),
	}

	mockRepo.On("ExistsByEmail", mock.Anything, req.Email).Return(false, nil)

	_, err := service.CreateUser(context.Background(), req)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot create admin users")
	mockRepo.AssertNotCalled(t, "Create")
}

func TestUpdateUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockLogger := newTestLogger()
	service := NewUserService(mockRepo, mockLogger)

	userID := uuid.New()
	existingUser := &entities.User{
		ID:        userID,
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Role:      "teacher",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	firstName := "Jane"
	lastName := "Smith"
	req := dto.UpdateUserRequest{
		FirstName: &firstName,
		LastName:  &lastName,
	}

	mockRepo.On("FindByID", mock.Anything, userID).Return(existingUser, nil)
	mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*entities.User")).Return(nil)

	result, err := service.UpdateUser(context.Background(), userID.String(), req)

	require.NoError(t, err)
	assert.Equal(t, firstName, result.FirstName)
	assert.Equal(t, lastName, result.LastName)
	mockRepo.AssertExpectations(t)
}

func TestUpdateUser_CannotPromoteToAdmin(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockLogger := newTestLogger()
	service := NewUserService(mockRepo, mockLogger)

	userID := uuid.New()
	existingUser := &entities.User{
		ID:        userID,
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Role:      "teacher",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	adminRole := string(enum.SystemRoleAdmin)
	req := dto.UpdateUserRequest{
		Role: &adminRole,
	}

	mockRepo.On("FindByID", mock.Anything, userID).Return(existingUser, nil)

	_, err := service.UpdateUser(context.Background(), userID.String(), req)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot promote to admin")
	mockRepo.AssertNotCalled(t, "Update")
}

func TestUpdateUser_ActivateInactiveUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockLogger := newTestLogger()
	service := NewUserService(mockRepo, mockLogger)

	userID := uuid.New()
	existingUser := &entities.User{
		ID:        userID,
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Role:      "teacher",
		IsActive:  false, // Usuario inactivo
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	isActive := true
	req := dto.UpdateUserRequest{
		IsActive: &isActive,
	}

	mockRepo.On("FindByID", mock.Anything, userID).Return(existingUser, nil)
	mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*entities.User")).Return(nil)

	result, err := service.UpdateUser(context.Background(), userID.String(), req)

	require.NoError(t, err)
	assert.True(t, result.IsActive)
	mockRepo.AssertExpectations(t)
}

func TestGetUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockLogger := newTestLogger()
	service := NewUserService(mockRepo, mockLogger)

	userID := uuid.New()
	existingUser := &entities.User{
		ID:        userID,
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Role:      "teacher",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("FindByID", mock.Anything, userID).Return(existingUser, nil)

	result, err := service.GetUser(context.Background(), userID.String())

	require.NoError(t, err)
	assert.Equal(t, userID.String(), result.ID)
	assert.Equal(t, "test@example.com", result.Email)
	mockRepo.AssertExpectations(t)
}

func TestGetUser_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockLogger := newTestLogger()
	service := NewUserService(mockRepo, mockLogger)

	userID := uuid.New()

	mockRepo.On("FindByID", mock.Anything, userID).Return(nil, nil)

	_, err := service.GetUser(context.Background(), userID.String())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
	mockRepo.AssertExpectations(t)
}

func TestDeleteUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockLogger := newTestLogger()
	service := NewUserService(mockRepo, mockLogger)

	userID := uuid.New()
	existingUser := &entities.User{
		ID:        userID,
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Role:      "teacher",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("FindByID", mock.Anything, userID).Return(existingUser, nil)
	mockRepo.On("Delete", mock.Anything, userID).Return(nil)

	err := service.DeleteUser(context.Background(), userID.String())

	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
