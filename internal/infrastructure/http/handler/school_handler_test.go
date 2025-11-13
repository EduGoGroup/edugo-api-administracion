package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EduGoGroup/edugo-shared/logger"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/EduGoGroup/edugo-api-administracion/internal/application/dto"
	commonErrors "github.com/EduGoGroup/edugo-shared/common/errors"
)

// MockSchoolService es un mock del SchoolService
type MockSchoolService struct {
	mock.Mock
}

func (m *MockSchoolService) CreateSchool(ctx context.Context, req dto.CreateSchoolRequest) (*dto.SchoolResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.SchoolResponse), args.Error(1)
}

func (m *MockSchoolService) GetSchool(ctx context.Context, id string) (*dto.SchoolResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.SchoolResponse), args.Error(1)
}

func (m *MockSchoolService) GetSchoolByCode(ctx context.Context, code string) (*dto.SchoolResponse, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.SchoolResponse), args.Error(1)
}

func (m *MockSchoolService) ListSchools(ctx context.Context) ([]dto.SchoolResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.SchoolResponse), args.Error(1)
}

func (m *MockSchoolService) UpdateSchool(ctx context.Context, id string, req dto.UpdateSchoolRequest) (*dto.SchoolResponse, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.SchoolResponse), args.Error(1)
}

func (m *MockSchoolService) DeleteSchool(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockLogger es un mock simple del logger
type MockLogger struct{}

func (m *MockLogger) Debug(msg string, fields ...interface{})  {}
func (m *MockLogger) Info(msg string, fields ...interface{})   {}
func (m *MockLogger) Warn(msg string, fields ...interface{})   {}
func (m *MockLogger) Error(msg string, fields ...interface{})  {}
func (m *MockLogger) Fatal(msg string, fields ...interface{})  {}
func (m *MockLogger) With(fields ...interface{}) logger.Logger { return m }
func (m *MockLogger) Sync() error                              { return nil }

func setupSchoolHandler() (*SchoolHandler, *MockSchoolService) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockSchoolService)
	mockLogger := &MockLogger{}
	handler := NewSchoolHandler(mockService, mockLogger)
	return handler, mockService
}

func TestSchoolHandler_CreateSchool_Success(t *testing.T) {
	handler, mockService := setupSchoolHandler()

	expectedResp := &dto.SchoolResponse{
		ID:   "school-123",
		Name: "Test School",
		Code: "TEST001",
	}

	mockService.On("CreateSchool", mock.Anything, mock.MatchedBy(func(req dto.CreateSchoolRequest) bool {
		return req.Name == "Test School" && req.Code == "TEST001"
	})).Return(expectedResp, nil)

	reqBody := dto.CreateSchoolRequest{
		Name: "Test School",
		Code: "TEST001",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/v1/schools", bytes.NewReader(bodyBytes))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.CreateSchool(c)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp dto.SchoolResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "school-123", resp.ID)
	assert.Equal(t, "Test School", resp.Name)

	mockService.AssertExpectations(t)
}

func TestSchoolHandler_CreateSchool_InvalidJSON(t *testing.T) {
	handler, _ := setupSchoolHandler()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/v1/schools", bytes.NewReader([]byte("{invalid json")))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.CreateSchool(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSchoolHandler_CreateSchool_AlreadyExists(t *testing.T) {
	handler, mockService := setupSchoolHandler()

	appErr := commonErrors.NewAlreadyExistsError("school").WithField("code", "TEST001")
	mockService.On("CreateSchool", mock.Anything, mock.Anything).Return(nil, appErr)

	reqBody := dto.CreateSchoolRequest{
		Name: "Test School",
		Code: "TEST001",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/v1/schools", bytes.NewReader(bodyBytes))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.CreateSchool(c)

	assert.Equal(t, http.StatusConflict, w.Code)
	mockService.AssertExpectations(t)
}

func TestSchoolHandler_GetSchool_Success(t *testing.T) {
	handler, mockService := setupSchoolHandler()

	expectedResp := &dto.SchoolResponse{
		ID:   "school-123",
		Name: "Test School",
	}

	mockService.On("GetSchool", mock.Anything, "school-123").Return(expectedResp, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/v1/schools/school-123", nil)
	c.Params = gin.Params{{Key: "id", Value: "school-123"}}

	handler.GetSchool(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.SchoolResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "school-123", resp.ID)

	mockService.AssertExpectations(t)
}

func TestSchoolHandler_GetSchool_NotFound(t *testing.T) {
	handler, mockService := setupSchoolHandler()

	appErr := commonErrors.NewNotFoundError("school")
	mockService.On("GetSchool", mock.Anything, "school-999").Return(nil, appErr)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/v1/schools/school-999", nil)
	c.Params = gin.Params{{Key: "id", Value: "school-999"}}

	handler.GetSchool(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockService.AssertExpectations(t)
}

func TestSchoolHandler_ListSchools_Success(t *testing.T) {
	handler, mockService := setupSchoolHandler()

	expectedResp := []dto.SchoolResponse{
		{ID: "school-1", Name: "School 1"},
		{ID: "school-2", Name: "School 2"},
	}

	mockService.On("ListSchools", mock.Anything).Return(expectedResp, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/v1/schools", nil)

	handler.ListSchools(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp []dto.SchoolResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 2)

	mockService.AssertExpectations(t)
}

func TestSchoolHandler_UpdateSchool_Success(t *testing.T) {
	handler, mockService := setupSchoolHandler()

	expectedResp := &dto.SchoolResponse{
		ID:   "school-123",
		Name: "Updated School",
	}

	mockService.On("UpdateSchool", mock.Anything, "school-123", mock.Anything).Return(expectedResp, nil)

	reqBody := dto.UpdateSchoolRequest{
		Name: ptrString("Updated School"),
	}
	bodyBytes, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("PUT", "/v1/schools/school-123", bytes.NewReader(bodyBytes))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "school-123"}}

	handler.UpdateSchool(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestSchoolHandler_DeleteSchool_Success(t *testing.T) {
	handler, mockService := setupSchoolHandler()

	mockService.On("DeleteSchool", mock.Anything, "school-123").Return(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("DELETE", "/v1/schools/school-123", nil)
	c.Params = gin.Params{{Key: "id", Value: "school-123"}}

	handler.DeleteSchool(c)

	assert.True(t, w.Code == http.StatusNoContent || w.Code == http.StatusOK, "Expected 204 or 200, got %d", w.Code)
	mockService.AssertExpectations(t)
}

func TestSchoolHandler_DeleteSchool_InternalError(t *testing.T) {
	handler, mockService := setupSchoolHandler()

	mockService.On("DeleteSchool", mock.Anything, "school-123").Return(errors.New("database error"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("DELETE", "/v1/schools/school-123", nil)
	c.Params = gin.Params{{Key: "id", Value: "school-123"}}

	handler.DeleteSchool(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

// Helper function
func ptrString(s string) *string {
	return &s
}
