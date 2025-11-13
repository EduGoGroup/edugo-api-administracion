package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/EduGoGroup/edugo-api-administracion/internal/application/dto"
	commonErrors "github.com/EduGoGroup/edugo-shared/common/errors"
)

// MockUnitMembershipService es un mock del UnitMembershipService
type MockUnitMembershipService struct {
	mock.Mock
}

func (m *MockUnitMembershipService) CreateMembership(ctx context.Context, req dto.CreateMembershipRequest) (*dto.MembershipResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.MembershipResponse), args.Error(1)
}

func (m *MockUnitMembershipService) GetMembership(ctx context.Context, id string) (*dto.MembershipResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.MembershipResponse), args.Error(1)
}

func (m *MockUnitMembershipService) ListMembershipsByUnit(ctx context.Context, unitID string, activeOnly bool) ([]dto.MembershipResponse, error) {
	args := m.Called(ctx, unitID, activeOnly)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.MembershipResponse), args.Error(1)
}

func (m *MockUnitMembershipService) ListMembershipsByUser(ctx context.Context, userID string, activeOnly bool) ([]dto.MembershipResponse, error) {
	args := m.Called(ctx, userID, activeOnly)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.MembershipResponse), args.Error(1)
}

func (m *MockUnitMembershipService) ListMembershipsByRole(ctx context.Context, unitID string, role string, activeOnly bool) ([]dto.MembershipResponse, error) {
	args := m.Called(ctx, unitID, role, activeOnly)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.MembershipResponse), args.Error(1)
}

func (m *MockUnitMembershipService) UpdateMembership(ctx context.Context, id string, req dto.UpdateMembershipRequest) (*dto.MembershipResponse, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.MembershipResponse), args.Error(1)
}

func (m *MockUnitMembershipService) ExpireMembership(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUnitMembershipService) DeleteMembership(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func setupMembershipHandler() (*UnitMembershipHandler, *MockUnitMembershipService) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockUnitMembershipService)
	mockLogger := &MockLogger{}
	handler := NewUnitMembershipHandler(mockService, mockLogger)
	return handler, mockService
}

func TestMembershipHandler_CreateMembership_Success(t *testing.T) {
	handler, mockService := setupMembershipHandler()

	expectedResp := &dto.MembershipResponse{
		ID:     "membership-123",
		UnitID: "unit-123",
		UserID: "user-123",
		Role:   "student",
	}

	mockService.On("CreateMembership", mock.Anything, mock.Anything).Return(expectedResp, nil)

	reqBody := dto.CreateMembershipRequest{
		UnitID: "unit-123",
		UserID: "user-123",
		Role:   "student",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/v1/memberships", bytes.NewReader(bodyBytes))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.CreateMembership(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}

func TestMembershipHandler_GetMembership_Success(t *testing.T) {
	handler, mockService := setupMembershipHandler()

	expectedResp := &dto.MembershipResponse{
		ID: "membership-123",
	}

	mockService.On("GetMembership", mock.Anything, "membership-123").Return(expectedResp, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/v1/memberships/membership-123", nil)
	c.Params = gin.Params{{Key: "id", Value: "membership-123"}}

	handler.GetMembership(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestMembershipHandler_ListByUnit_Success(t *testing.T) {
	handler, mockService := setupMembershipHandler()

	expectedResp := []dto.MembershipResponse{
		{ID: "m-1", Role: "student"},
	}

	mockService.On("ListMembershipsByUnit", mock.Anything, "unit-123", true).Return(expectedResp, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/v1/units/unit-123/memberships", nil)
	c.Params = gin.Params{{Key: "unitId", Value: "unit-123"}}

	handler.ListMembershipsByUnit(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestMembershipHandler_ListByUser_Success(t *testing.T) {
	handler, mockService := setupMembershipHandler()

	expectedResp := []dto.MembershipResponse{
		{ID: "m-1"},
	}

	mockService.On("ListMembershipsByUser", mock.Anything, "user-123", true).Return(expectedResp, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/v1/users/user-123/memberships", nil)
	c.Params = gin.Params{{Key: "userId", Value: "user-123"}}

	handler.ListMembershipsByUser(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestMembershipHandler_UpdateMembership_Success(t *testing.T) {
	handler, mockService := setupMembershipHandler()

	expectedResp := &dto.MembershipResponse{
		ID:   "membership-123",
		Role: "teacher",
	}

	mockService.On("UpdateMembership", mock.Anything, "membership-123", mock.Anything).Return(expectedResp, nil)

	reqBody := dto.UpdateMembershipRequest{
		Role: ptrString("teacher"),
	}
	bodyBytes, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("PUT", "/v1/memberships/membership-123", bytes.NewReader(bodyBytes))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "membership-123"}}

	handler.UpdateMembership(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestMembershipHandler_ExpireMembership_Success(t *testing.T) {
	handler, mockService := setupMembershipHandler()

	mockService.On("ExpireMembership", mock.Anything, "membership-123").Return(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/v1/memberships/membership-123/expire", nil)
	c.Params = gin.Params{{Key: "id", Value: "membership-123"}}

	handler.ExpireMembership(c)

	assert.True(t, w.Code == http.StatusNoContent || w.Code == http.StatusOK)
	mockService.AssertExpectations(t)
}

func TestMembershipHandler_DeleteMembership_Success(t *testing.T) {
	handler, mockService := setupMembershipHandler()

	mockService.On("DeleteMembership", mock.Anything, "membership-123").Return(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("DELETE", "/v1/memberships/membership-123", nil)
	c.Params = gin.Params{{Key: "id", Value: "membership-123"}}

	handler.DeleteMembership(c)

	assert.True(t, w.Code == http.StatusNoContent || w.Code == http.StatusOK)
	mockService.AssertExpectations(t)
}

func TestMembershipHandler_GetMembership_NotFound(t *testing.T) {
	handler, mockService := setupMembershipHandler()

	appErr := commonErrors.NewNotFoundError("membership")
	mockService.On("GetMembership", mock.Anything, "membership-999").Return(nil, appErr)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/v1/memberships/membership-999", nil)
	c.Params = gin.Params{{Key: "id", Value: "membership-999"}}

	handler.GetMembership(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockService.AssertExpectations(t)
}
