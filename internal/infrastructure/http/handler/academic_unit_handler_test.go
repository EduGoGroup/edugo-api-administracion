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

// MockAcademicUnitService es un mock del AcademicUnitService
type MockAcademicUnitService struct {
	mock.Mock
}

func (m *MockAcademicUnitService) CreateUnit(ctx context.Context, schoolID string, req dto.CreateAcademicUnitRequest) (*dto.AcademicUnitResponse, error) {
	args := m.Called(ctx, schoolID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.AcademicUnitResponse), args.Error(1)
}

func (m *MockAcademicUnitService) GetUnit(ctx context.Context, id string) (*dto.AcademicUnitResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.AcademicUnitResponse), args.Error(1)
}

func (m *MockAcademicUnitService) GetUnitTree(ctx context.Context, schoolID string) ([]*dto.UnitTreeNode, error) {
	args := m.Called(ctx, schoolID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*dto.UnitTreeNode), args.Error(1)
}

func (m *MockAcademicUnitService) ListUnitsBySchool(ctx context.Context, schoolID string, includeDeleted bool) ([]dto.AcademicUnitResponse, error) {
	args := m.Called(ctx, schoolID, includeDeleted)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.AcademicUnitResponse), args.Error(1)
}

func (m *MockAcademicUnitService) ListUnitsByType(ctx context.Context, schoolID string, unitType string) ([]dto.AcademicUnitResponse, error) {
	args := m.Called(ctx, schoolID, unitType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.AcademicUnitResponse), args.Error(1)
}

func (m *MockAcademicUnitService) UpdateUnit(ctx context.Context, id string, req dto.UpdateAcademicUnitRequest) (*dto.AcademicUnitResponse, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.AcademicUnitResponse), args.Error(1)
}

func (m *MockAcademicUnitService) DeleteUnit(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAcademicUnitService) RestoreUnit(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAcademicUnitService) GetHierarchyPath(ctx context.Context, id string) ([]dto.AcademicUnitResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.AcademicUnitResponse), args.Error(1)
}

func setupAcademicUnitHandler() (*AcademicUnitHandler, *MockAcademicUnitService) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockAcademicUnitService)
	mockLogger := &MockLogger{}
	handler := NewAcademicUnitHandler(mockService, mockLogger)
	return handler, mockService
}

func TestAcademicUnitHandler_CreateUnit_Success(t *testing.T) {
	handler, mockService := setupAcademicUnitHandler()

	expectedResp := &dto.AcademicUnitResponse{
		ID:          "unit-123",
		SchoolID:    "school-123",
		Type:        "grade",
		DisplayName: "1st Grade",
	}

	mockService.On("CreateUnit", mock.Anything, "school-123", mock.Anything).Return(expectedResp, nil)

	reqBody := dto.CreateAcademicUnitRequest{
		Type:        "grade",
		DisplayName: "1st Grade",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/v1/schools/school-123/units", bytes.NewReader(bodyBytes))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "school-123"}} // Cambio: schoolId → id

	handler.CreateUnit(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}

func TestAcademicUnitHandler_GetUnit_Success(t *testing.T) {
	handler, mockService := setupAcademicUnitHandler()

	expectedResp := &dto.AcademicUnitResponse{
		ID:          "unit-123",
		DisplayName: "1st Grade",
	}

	mockService.On("GetUnit", mock.Anything, "unit-123").Return(expectedResp, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/v1/units/unit-123", nil)
	c.Params = gin.Params{{Key: "id", Value: "unit-123"}}

	handler.GetUnit(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestAcademicUnitHandler_GetUnitTree_Success(t *testing.T) {
	handler, mockService := setupAcademicUnitHandler()

	expectedResp := []*dto.UnitTreeNode{
		{
			ID:          "unit-1",
			DisplayName: "Grade 1",
			Type:        "grade",
			Children:    []*dto.UnitTreeNode{},
		},
	}

	mockService.On("GetUnitTree", mock.Anything, "school-123").Return(expectedResp, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/v1/schools/school-123/units/tree", nil)
	c.Params = gin.Params{{Key: "id", Value: "school-123"}} // Cambio: schoolId → id

	handler.GetUnitTree(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestAcademicUnitHandler_ListUnitsBySchool_Success(t *testing.T) {
	handler, mockService := setupAcademicUnitHandler()

	expectedResp := []dto.AcademicUnitResponse{
		{ID: "unit-1", DisplayName: "Unit 1"},
		{ID: "unit-2", DisplayName: "Unit 2"},
	}

	mockService.On("ListUnitsBySchool", mock.Anything, "school-123", false).Return(expectedResp, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/v1/schools/school-123/units", nil)
	c.Params = gin.Params{{Key: "id", Value: "school-123"}} // Cambio: schoolId → id

	handler.ListUnitsBySchool(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestAcademicUnitHandler_ListUnitsByType_Success(t *testing.T) {
	handler, mockService := setupAcademicUnitHandler()

	expectedResp := []dto.AcademicUnitResponse{
		{ID: "unit-1", Type: "grade"},
	}

	mockService.On("ListUnitsByType", mock.Anything, "school-123", "grade").Return(expectedResp, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/v1/schools/school-123/units/by-type?type=grade", nil)
	c.Params = gin.Params{{Key: "id", Value: "school-123"}} // Cambio: schoolId → id

	handler.ListUnitsByType(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestAcademicUnitHandler_UpdateUnit_Success(t *testing.T) {
	handler, mockService := setupAcademicUnitHandler()

	expectedResp := &dto.AcademicUnitResponse{
		ID:          "unit-123",
		DisplayName: "Updated Unit",
	}

	mockService.On("UpdateUnit", mock.Anything, "unit-123", mock.Anything).Return(expectedResp, nil)

	reqBody := dto.UpdateAcademicUnitRequest{
		DisplayName: ptrString("Updated Unit"),
	}
	bodyBytes, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("PUT", "/v1/units/unit-123", bytes.NewReader(bodyBytes))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "unit-123"}}

	handler.UpdateUnit(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestAcademicUnitHandler_DeleteUnit_Success(t *testing.T) {
	handler, mockService := setupAcademicUnitHandler()

	mockService.On("DeleteUnit", mock.Anything, "unit-123").Return(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("DELETE", "/v1/units/unit-123", nil)
	c.Params = gin.Params{{Key: "id", Value: "unit-123"}}

	handler.DeleteUnit(c)

	assert.True(t, w.Code == http.StatusNoContent || w.Code == http.StatusOK)
	mockService.AssertExpectations(t)
}

func TestAcademicUnitHandler_RestoreUnit_Success(t *testing.T) {
	handler, mockService := setupAcademicUnitHandler()

	mockService.On("RestoreUnit", mock.Anything, "unit-123").Return(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/v1/units/unit-123/restore", nil)
	c.Params = gin.Params{{Key: "id", Value: "unit-123"}}

	handler.RestoreUnit(c)

	assert.True(t, w.Code == http.StatusNoContent || w.Code == http.StatusOK)
	mockService.AssertExpectations(t)
}

func TestAcademicUnitHandler_GetHierarchyPath_Success(t *testing.T) {
	handler, mockService := setupAcademicUnitHandler()

	expectedResp := []dto.AcademicUnitResponse{
		{ID: "unit-1", DisplayName: "Root"},
		{ID: "unit-2", DisplayName: "Child"},
	}

	mockService.On("GetHierarchyPath", mock.Anything, "unit-123").Return(expectedResp, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/v1/units/unit-123/hierarchy-path", nil)
	c.Params = gin.Params{{Key: "id", Value: "unit-123"}}

	handler.GetHierarchyPath(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestAcademicUnitHandler_GetUnit_NotFound(t *testing.T) {
	handler, mockService := setupAcademicUnitHandler()

	appErr := commonErrors.NewNotFoundError("unit")
	mockService.On("GetUnit", mock.Anything, "unit-999").Return(nil, appErr)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/v1/units/unit-999", nil)
	c.Params = gin.Params{{Key: "id", Value: "unit-999"}}

	handler.GetUnit(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockService.AssertExpectations(t)
}
