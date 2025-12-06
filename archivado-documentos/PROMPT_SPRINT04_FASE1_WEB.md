# PROMPT SPRINT-04 FASE 1 - CLAUDE CODE WEB

**Proyecto:** edugo-api-administracion  
**Sprint:** Sprint-04 - Services/API  
**Ejecutor:** Claude Code Web  
**Duración estimada:** 2-3 horas  
**Branch:** `feature/sprint-04-services-api`  
**Base branch:** `dev` ⚠️ IMPORTANTE: Crear desde dev, NO desde main

---

## 🎯 TU OBJETIVO (Fase 1 - Sin servidor local)

Implementar la estructura completa de servicios de aplicación y endpoints REST, dejando con **STUBS** aquellas partes que requieran:
- Ejecución de servidor HTTP
- Tests E2E que necesiten servidor corriendo
- Validaciones que requieran base de datos real

**Lo que SÍ puedes hacer:**
- ✅ Crear toda la estructura de código
- ✅ Implementar DTOs, handlers, servicios
- ✅ Escribir tests unitarios
- ✅ Configurar router de Gin
- ✅ Implementar lógica de negocio en servicios

**Lo que debes dejar con STUB para Fase 2:**
- ⚠️ Tests de integración HTTP (requieren servidor)
- ⚠️ Validación E2E de flujos completos
- ⚠️ Pruebas con Postman/curl

---

## 📋 PREREQUISITOS

Antes de empezar, verifica:
1. Estás en el proyecto `edugo-api-administracion`
2. Branch `main` está actualizada (Sprint-03 mergeado)
3. Dependencia de Gin disponible en go.mod

---

## 📋 TAREAS FASE 1

### TASK-01: Setup Inicial (15min)

```bash
# 1. Crear branch desde dev actualizada (IMPORTANTE: usar dev, no main)
git checkout dev
git pull origin dev
git checkout -b feature/sprint-04-services-api

# 2. Verificar/agregar dependencias
go get github.com/gin-gonic/gin@v1.9.1
go get github.com/go-playground/validator/v10

# 3. Verificar que compila
go build ./...
```

**Commit:** `chore(deps): add Gin and validator dependencies for Sprint-04`

---

### TASK-02: DTOs (Data Transfer Objects) (30min)

**Archivo:** `internal/infrastructure/http/dto/school_dto.go`

```go
package dto

import (
	"time"
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/valueobject"
)

// CreateSchoolRequest DTO para crear escuela
type CreateSchoolRequest struct {
	Name    string `json:"name" binding:"required,min=3,max=100"`
	Code    string `json:"code" binding:"required,min=3,max=20"`
	Address string `json:"address" binding:"required"`
} // @name CreateSchoolRequest

// UpdateSchoolRequest DTO para actualizar escuela
type UpdateSchoolRequest struct {
	Name    *string `json:"name" binding:"omitempty,min=3,max=100"`
	Address *string `json:"address" binding:"omitempty"`
} // @name UpdateSchoolRequest

// SchoolResponse DTO de respuesta
type SchoolResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
} // @name SchoolResponse

// ToSchoolResponse convierte entity a DTO
func ToSchoolResponse(school *entity.School) *SchoolResponse {
	return &SchoolResponse{
		ID:        school.ID().String(),
		Name:      school.Name(),
		Code:      school.Code(),
		Address:   school.Address(),
		CreatedAt: school.CreatedAt(),
		UpdatedAt: school.UpdatedAt(),
	}
}
```

**Archivo:** `internal/infrastructure/http/dto/academic_unit_dto.go`

```go
package dto

// CreateUnitRequest DTO para crear unidad académica
type CreateUnitRequest struct {
	ParentUnitID *string `json:"parent_unit_id" binding:"omitempty,uuid"`
	SchoolID     string  `json:"school_id" binding:"required,uuid"`
	Type         string  `json:"type" binding:"required,oneof=grade section club department"`
	Name         string  `json:"name" binding:"required,min=3,max=100"`
	Code         string  `json:"code" binding:"required,min=2,max=50"`
	Description  *string `json:"description" binding:"omitempty,max=500"`
} // @name CreateUnitRequest

// UpdateUnitRequest DTO para actualizar unidad
type UpdateUnitRequest struct {
	ParentUnitID *string `json:"parent_unit_id" binding:"omitempty,uuid"`
	Name         *string `json:"name" binding:"omitempty,min=3,max=100"`
	Description  *string `json:"description" binding:"omitempty,max=500"`
} // @name UpdateUnitRequest

// UnitResponse DTO de respuesta simple
type UnitResponse struct {
	ID           string    `json:"id"`
	ParentUnitID *string   `json:"parent_unit_id,omitempty"`
	SchoolID     string    `json:"school_id"`
	Type         string    `json:"type"`
	Name         string    `json:"name"`
	Code         string    `json:"code"`
	Description  *string   `json:"description,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
} // @name UnitResponse

// UnitTreeNode DTO para árbol jerárquico (usa ltree!)
type UnitTreeNode struct {
	ID       string          `json:"id"`
	Name     string          `json:"name"`
	Code     string          `json:"code"`
	Type     string          `json:"type"`
	Depth    int             `json:"depth"`
	Children []*UnitTreeNode `json:"children,omitempty"`
} // @name UnitTreeNode

// ToUnitResponse convierte entity a DTO
func ToUnitResponse(unit *entity.AcademicUnit) *UnitResponse {
	var parentID *string
	if unit.ParentUnitID() != nil {
		pid := unit.ParentUnitID().String()
		parentID = &pid
	}

	var desc *string
	if unit.Description() != "" {
		d := unit.Description()
		desc = &d
	}

	return &UnitResponse{
		ID:           unit.ID().String(),
		ParentUnitID: parentID,
		SchoolID:     unit.SchoolID().String(),
		Type:         unit.UnitType().String(),
		Name:         unit.DisplayName(),
		Code:         unit.Code(),
		Description:  desc,
		CreatedAt:    unit.CreatedAt(),
		UpdatedAt:    unit.UpdatedAt(),
	}
}
```

**Archivo:** `internal/infrastructure/http/dto/common_dto.go`

```go
package dto

// ErrorResponse DTO para errores
type ErrorResponse struct {
	Error   string            `json:"error"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
} // @name ErrorResponse

// SuccessResponse DTO genérico de éxito
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
} // @name SuccessResponse

// PaginationMeta metadata de paginación
type PaginationMeta struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
} // @name PaginationMeta
```

**Commit:** `feat(dto): add DTOs for schools and academic units`

---

### TASK-03: Application Service - HierarchyService (45min)

**Archivo:** `internal/application/service/hierarchy_service.go`

```go
package service

import (
	"context"
	"fmt"

	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/entity"
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/repository"
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/service"
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/valueobject"
	"github.com/EduGoGroup/edugo-shared/common/errors"
)

// HierarchyService servicio de aplicación para operaciones jerárquicas
type HierarchyService struct {
	unitRepo        repository.AcademicUnitRepository
	schoolRepo      repository.SchoolRepository
	domainService   *service.AcademicUnitDomainService
}

// NewHierarchyService constructor
func NewHierarchyService(
	unitRepo repository.AcademicUnitRepository,
	schoolRepo repository.SchoolRepository,
	domainService *service.AcademicUnitDomainService,
) *HierarchyService {
	return &HierarchyService{
		unitRepo:      unitRepo,
		schoolRepo:    schoolRepo,
		domainService: domainService,
	}
}

// CreateUnit crea una nueva unidad académica
func (s *HierarchyService) CreateUnit(
	ctx context.Context,
	parentUnitID *valueobject.UnitID,
	schoolID valueobject.SchoolID,
	unitType valueobject.UnitType,
	name string,
	code string,
	description string,
) (*entity.AcademicUnit, error) {
	// 1. Validar que la escuela existe
	school, err := s.schoolRepo.FindByID(ctx, schoolID)
	if err != nil {
		return nil, errors.NewNotFoundError(fmt.Sprintf("school %s not found", schoolID))
	}

	// 2. Validar que el código no esté duplicado en la escuela
	exists, err := s.unitRepo.ExistsBySchoolIDAndCode(ctx, schoolID, code)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.NewValidationError(fmt.Sprintf("code %s already exists in school", code))
	}

	// 3. Si tiene padre, validar que existe y no crea ciclo
	var parent *entity.AcademicUnit
	if parentUnitID != nil {
		parent, err = s.unitRepo.FindByID(ctx, *parentUnitID, false)
		if err != nil {
			return nil, errors.NewNotFoundError(fmt.Sprintf("parent unit %s not found", parentUnitID))
		}

		// Validar que el padre pertenece a la misma escuela
		if parent.SchoolID() != schoolID {
			return nil, errors.NewValidationError("parent unit must belong to the same school")
		}
	}

	// 4. Crear la unidad
	unit, err := entity.NewAcademicUnit(schoolID, unitType, name, code)
	if err != nil {
		return nil, err
	}

	if description != "" {
		unit.SetDescription(description)
	}

	// 5. Establecer padre si existe
	if parent != nil {
		if err := s.domainService.SetParent(unit, parent.ID(), parent.UnitType()); err != nil {
			return nil, err
		}
	}

	// 6. Persistir
	if err := s.unitRepo.Create(ctx, unit); err != nil {
		return nil, err
	}

	return unit, nil
}

// GetUnitTree obtiene el árbol jerárquico completo de una unidad usando ltree
func (s *HierarchyService) GetUnitTree(ctx context.Context, unitID valueobject.UnitID) (*entity.AcademicUnit, []*entity.AcademicUnit, error) {
	// 1. Obtener la unidad raíz
	root, err := s.unitRepo.FindByID(ctx, unitID, false)
	if err != nil {
		return nil, nil, err
	}

	// 2. Obtener todos los descendientes usando ltree (¡Sprint-03!)
	descendants, err := s.unitRepo.FindDescendants(ctx, unitID)
	if err != nil {
		return nil, nil, err
	}

	return root, descendants, nil
}

// MoveUnit mueve una unidad a un nuevo padre (o a raíz si newParentID es nil)
func (s *HierarchyService) MoveUnit(
	ctx context.Context,
	unitID valueobject.UnitID,
	newParentID *valueobject.UnitID,
) error {
	// 1. Validar que la unidad existe
	unit, err := s.unitRepo.FindByID(ctx, unitID, false)
	if err != nil {
		return err
	}

	// 2. Si hay nuevo padre, validar
	if newParentID != nil {
		newParent, err := s.unitRepo.FindByID(ctx, *newParentID, false)
		if err != nil {
			return errors.NewNotFoundError(fmt.Sprintf("new parent unit %s not found", newParentID))
		}

		// Validar misma escuela
		if newParent.SchoolID() != unit.SchoolID() {
			return errors.NewValidationError("cannot move unit to different school")
		}

		// Validar que no crea ciclo (nuevo padre no puede ser descendiente)
		descendants, err := s.unitRepo.FindDescendants(ctx, unitID)
		if err != nil {
			return err
		}

		for _, desc := range descendants {
			if desc.ID() == *newParentID {
				return errors.NewValidationError("cannot move unit: would create a circular reference")
			}
		}
	}

	// 3. Mover usando ltree (Sprint-03!)
	return s.unitRepo.MoveSubtree(ctx, unitID, newParentID)
}

// ValidateNoCircularReference valida que mover una unidad no cree un ciclo
func (s *HierarchyService) ValidateNoCircularReference(
	ctx context.Context,
	unitID valueobject.UnitID,
	newParentID valueobject.UnitID,
) error {
	// Obtener todos los descendientes de la unidad
	descendants, err := s.unitRepo.FindDescendants(ctx, unitID)
	if err != nil {
		return err
	}

	// Verificar que el nuevo padre no sea uno de los descendientes
	for _, desc := range descendants {
		if desc.ID() == newParentID {
			return errors.NewValidationError("circular reference detected")
		}
	}

	// Verificar que no sea la misma unidad
	if unitID == newParentID {
		return errors.NewValidationError("unit cannot be its own parent")
	}

	return nil
}
```

**Archivo:** `internal/application/service/hierarchy_service_test.go`

```go
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

// TODO: Implementar resto de métodos del mock según se necesiten

// MockSchoolRepository mock del repositorio de escuelas
type MockSchoolRepository struct {
	mock.Mock
}

func (m *MockSchoolRepository) FindByID(ctx context.Context, id valueobject.SchoolID) (*entity.School, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.School), args.Error(1)
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
		parentID := valueobject.NewUnitID()

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
}
```

**Commit:** `feat(service): implement HierarchyService with ltree support`

---

### TASK-04: HTTP Handlers (1h 30min)

**Archivo:** `internal/infrastructure/http/handler/school_handler.go`

```go
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/repository"
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/valueobject"
	"github.com/EduGoGroup/edugo-api-administracion/internal/infrastructure/http/dto"
	"github.com/EduGoGroup/edugo-shared/common/errors"
)

// SchoolHandler handler para endpoints de escuelas
type SchoolHandler struct {
	schoolRepo repository.SchoolRepository
}

// NewSchoolHandler constructor
func NewSchoolHandler(schoolRepo repository.SchoolRepository) *SchoolHandler {
	return &SchoolHandler{
		schoolRepo: schoolRepo,
	}
}

// CreateSchool godoc
// @Summary Create a new school
// @Tags schools
// @Accept json
// @Produce json
// @Param school body dto.CreateSchoolRequest true "School data"
// @Success 201 {object} dto.SchoolResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /schools [post]
func (h *SchoolHandler) CreateSchool(c *gin.Context) {
	var req dto.CreateSchoolRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	// Crear entidad
	school, err := entity.NewSchool(req.Name, req.Code, req.Address)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "domain_error",
			Message: err.Error(),
		})
		return
	}

	// Persistir
	if err := h.schoolRepo.Create(c.Request.Context(), school); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "database_error",
			Message: "failed to create school",
		})
		return
	}

	c.JSON(http.StatusCreated, dto.ToSchoolResponse(school))
}

// GetSchool godoc
// @Summary Get school by ID
// @Tags schools
// @Produce json
// @Param id path string true "School ID"
// @Success 200 {object} dto.SchoolResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /schools/{id} [get]
func (h *SchoolHandler) GetSchool(c *gin.Context) {
	idStr := c.Param("id")
	
	schoolID, err := valueobject.ParseSchoolID(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_id",
			Message: "invalid school ID format",
		})
		return
	}

	school, err := h.schoolRepo.FindByID(c.Request.Context(), schoolID)
	if err != nil {
		if errors.IsNotFoundError(err) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "not_found",
				Message: "school not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "database_error",
			Message: "failed to get school",
		})
		return
	}

	c.JSON(http.StatusOK, dto.ToSchoolResponse(school))
}

// ListSchools godoc
// @Summary List all schools
// @Tags schools
// @Produce json
// @Success 200 {array} dto.SchoolResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /schools [get]
func (h *SchoolHandler) ListSchools(c *gin.Context) {
	schools, err := h.schoolRepo.FindAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "database_error",
			Message: "failed to list schools",
		})
		return
	}

	response := make([]*dto.SchoolResponse, len(schools))
	for i, school := range schools {
		response[i] = dto.ToSchoolResponse(school)
	}

	c.JSON(http.StatusOK, response)
}

// UpdateSchool godoc
// @Summary Update school
// @Tags schools
// @Accept json
// @Produce json
// @Param id path string true "School ID"
// @Param school body dto.UpdateSchoolRequest true "Update data"
// @Success 200 {object} dto.SchoolResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /schools/{id} [put]
func (h *SchoolHandler) UpdateSchool(c *gin.Context) {
	idStr := c.Param("id")
	
	schoolID, err := valueobject.ParseSchoolID(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_id",
			Message: "invalid school ID format",
		})
		return
	}

	var req dto.UpdateSchoolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	// Obtener escuela existente
	school, err := h.schoolRepo.FindByID(c.Request.Context(), schoolID)
	if err != nil {
		if errors.IsNotFoundError(err) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "not_found",
				Message: "school not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "database_error",
			Message: "failed to get school",
		})
		return
	}

	// Actualizar campos
	if req.Name != nil {
		school.SetName(*req.Name)
	}
	if req.Address != nil {
		school.SetAddress(*req.Address)
	}

	// Persistir
	if err := h.schoolRepo.Update(c.Request.Context(), school); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "database_error",
			Message: "failed to update school",
		})
		return
	}

	c.JSON(http.StatusOK, dto.ToSchoolResponse(school))
}

// DeleteSchool godoc
// @Summary Delete school
// @Tags schools
// @Param id path string true "School ID"
// @Success 204
// @Failure 404 {object} dto.ErrorResponse
// @Router /schools/{id} [delete]
func (h *SchoolHandler) DeleteSchool(c *gin.Context) {
	idStr := c.Param("id")
	
	schoolID, err := valueobject.ParseSchoolID(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_id",
			Message: "invalid school ID format",
		})
		return
	}

	if err := h.schoolRepo.Delete(c.Request.Context(), schoolID); err != nil {
		if errors.IsNotFoundError(err) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "not_found",
				Message: "school not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "database_error",
			Message: "failed to delete school",
		})
		return
	}

	c.Status(http.StatusNoContent)
}
```

**Archivo:** `internal/infrastructure/http/handler/academic_unit_handler.go`

```go
package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	
	"github.com/EduGoGroup/edugo-api-administracion/internal/application/service"
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/repository"
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/valueobject"
	"github.com/EduGoGroup/edugo-api-administracion/internal/infrastructure/http/dto"
	"github.com/EduGoGroup/edugo-shared/common/errors"
)

// AcademicUnitHandler handler para endpoints de unidades académicas
type AcademicUnitHandler struct {
	unitRepo         repository.AcademicUnitRepository
	hierarchyService *service.HierarchyService
}

// NewAcademicUnitHandler constructor
func NewAcademicUnitHandler(
	unitRepo repository.AcademicUnitRepository,
	hierarchyService *service.HierarchyService,
) *AcademicUnitHandler {
	return &AcademicUnitHandler{
		unitRepo:         unitRepo,
		hierarchyService: hierarchyService,
	}
}

// CreateUnit godoc
// @Summary Create academic unit
// @Tags units
// @Accept json
// @Produce json
// @Param unit body dto.CreateUnitRequest true "Unit data"
// @Success 201 {object} dto.UnitResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /units [post]
func (h *AcademicUnitHandler) CreateUnit(c *gin.Context) {
	var req dto.CreateUnitRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	// Parse IDs
	schoolID, err := valueobject.ParseSchoolID(req.SchoolID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_school_id",
			Message: "invalid school ID format",
		})
		return
	}

	var parentID *valueobject.UnitID
	if req.ParentUnitID != nil {
		pid, err := valueobject.ParseUnitID(*req.ParentUnitID)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "invalid_parent_id",
				Message: "invalid parent unit ID format",
			})
			return
		}
		parentID = &pid
	}

	unitType, err := valueobject.NewUnitType(req.Type)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_type",
			Message: err.Error(),
		})
		return
	}

	description := ""
	if req.Description != nil {
		description = *req.Description
	}

	// Usar servicio de aplicación
	unit, err := h.hierarchyService.CreateUnit(
		c.Request.Context(),
		parentID,
		schoolID,
		unitType,
		req.Name,
		req.Code,
		description,
	)
	if err != nil {
		if errors.IsValidationError(err) {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "validation_error",
				Message: err.Error(),
			})
			return
		}
		if errors.IsNotFoundError(err) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "not_found",
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "server_error",
			Message: "failed to create unit",
		})
		return
	}

	c.JSON(http.StatusCreated, dto.ToUnitResponse(unit))
}

// GetUnit godoc
// @Summary Get unit by ID
// @Tags units
// @Produce json
// @Param id path string true "Unit ID"
// @Success 200 {object} dto.UnitResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /units/{id} [get]
func (h *AcademicUnitHandler) GetUnit(c *gin.Context) {
	idStr := c.Param("id")
	
	unitID, err := valueobject.ParseUnitID(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_id",
			Message: "invalid unit ID format",
		})
		return
	}

	unit, err := h.unitRepo.FindByID(c.Request.Context(), unitID, false)
	if err != nil {
		if errors.IsNotFoundError(err) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "not_found",
				Message: "unit not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "database_error",
			Message: "failed to get unit",
		})
		return
	}

	c.JSON(http.StatusOK, dto.ToUnitResponse(unit))
}

// GetUnitTree godoc
// @Summary Get unit hierarchy tree (uses ltree!)
// @Tags units
// @Produce json
// @Param id path string true "Unit ID"
// @Success 200 {object} dto.UnitTreeNode
// @Failure 404 {object} dto.ErrorResponse
// @Router /units/{id}/tree [get]
func (h *AcademicUnitHandler) GetUnitTree(c *gin.Context) {
	idStr := c.Param("id")
	
	unitID, err := valueobject.ParseUnitID(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_id",
			Message: "invalid unit ID format",
		})
		return
	}

	// Usar servicio que usa ltree! (Sprint-03)
	root, descendants, err := h.hierarchyService.GetUnitTree(c.Request.Context(), unitID)
	if err != nil {
		if errors.IsNotFoundError(err) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "not_found",
				Message: "unit not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "server_error",
			Message: "failed to get unit tree",
		})
		return
	}

	// Construir árbol jerárquico
	tree := buildTree(root, descendants)

	c.JSON(http.StatusOK, tree)
}

// buildTree construye el árbol jerárquico recursivamente
func buildTree(root *entity.AcademicUnit, allDescendants []*entity.AcademicUnit) *dto.UnitTreeNode {
	node := &dto.UnitTreeNode{
		ID:       root.ID().String(),
		Name:     root.DisplayName(),
		Code:     root.Code(),
		Type:     root.UnitType().String(),
		Depth:    0,
		Children: []*dto.UnitTreeNode{},
	}

	// Encontrar hijos directos
	for _, desc := range allDescendants {
		if desc.ParentUnitID() != nil && *desc.ParentUnitID() == root.ID() {
			childNode := buildTreeRecursive(desc, allDescendants, 1)
			node.Children = append(node.Children, childNode)
		}
	}

	return node
}

func buildTreeRecursive(current *entity.AcademicUnit, allDescendants []*entity.AcademicUnit, depth int) *dto.UnitTreeNode {
	node := &dto.UnitTreeNode{
		ID:       current.ID().String(),
		Name:     current.DisplayName(),
		Code:     current.Code(),
		Type:     current.UnitType().String(),
		Depth:    depth,
		Children: []*dto.UnitTreeNode{},
	}

	// Encontrar hijos directos
	for _, desc := range allDescendants {
		if desc.ParentUnitID() != nil && *desc.ParentUnitID() == current.ID() {
			childNode := buildTreeRecursive(desc, allDescendants, depth+1)
			node.Children = append(node.Children, childNode)
		}
	}

	return node
}

// ListUnits godoc
// @Summary List units by school
// @Tags units
// @Produce json
// @Param school_id query string true "School ID"
// @Param depth query int false "Filter by depth (uses ltree nlevel)"
// @Success 200 {array} dto.UnitResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /units [get]
func (h *AcademicUnitHandler) ListUnits(c *gin.Context) {
	schoolIDStr := c.Query("school_id")
	if schoolIDStr == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "missing_parameter",
			Message: "school_id is required",
		})
		return
	}

	schoolID, err := valueobject.ParseSchoolID(schoolIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_school_id",
			Message: "invalid school ID format",
		})
		return
	}

	var units []*entity.AcademicUnit

	// Si hay filtro por profundidad, usar ltree!
	depthStr := c.Query("depth")
	if depthStr != "" {
		depth, err := strconv.Atoi(depthStr)
		if err != nil || depth < 1 {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "invalid_depth",
				Message: "depth must be a positive integer",
			})
			return
		}

		units, err = h.unitRepo.FindBySchoolIDAndDepth(c.Request.Context(), schoolID, depth)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "database_error",
				Message: "failed to list units",
			})
			return
		}
	} else {
		units, err = h.unitRepo.FindBySchoolID(c.Request.Context(), schoolID, false)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "database_error",
				Message: "failed to list units",
			})
			return
		}
	}

	response := make([]*dto.UnitResponse, len(units))
	for i, unit := range units {
		response[i] = dto.ToUnitResponse(unit)
	}

	c.JSON(http.StatusOK, response)
}

// UpdateUnit godoc
// @Summary Update unit
// @Tags units
// @Accept json
// @Produce json
// @Param id path string true "Unit ID"
// @Param unit body dto.UpdateUnitRequest true "Update data"
// @Success 200 {object} dto.UnitResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /units/{id} [put]
func (h *AcademicUnitHandler) UpdateUnit(c *gin.Context) {
	idStr := c.Param("id")
	
	unitID, err := valueobject.ParseUnitID(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_id",
			Message: "invalid unit ID format",
		})
		return
	}

	var req dto.UpdateUnitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	// Obtener unidad existente
	unit, err := h.unitRepo.FindByID(c.Request.Context(), unitID, false)
	if err != nil {
		if errors.IsNotFoundError(err) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "not_found",
				Message: "unit not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "database_error",
			Message: "failed to get unit",
		})
		return
	}

	// Actualizar nombre si se proporciona
	if req.Name != nil {
		unit.SetDisplayName(*req.Name)
	}

	// Actualizar descripción si se proporciona
	if req.Description != nil {
		unit.SetDescription(*req.Description)
	}

	// Si cambia el padre, usar MoveUnit (usa ltree!)
	if req.ParentUnitID != nil {
		var newParentID *valueobject.UnitID
		if *req.ParentUnitID != "" {
			pid, err := valueobject.ParseUnitID(*req.ParentUnitID)
			if err != nil {
				c.JSON(http.StatusBadRequest, dto.ErrorResponse{
					Error:   "invalid_parent_id",
					Message: "invalid parent unit ID format",
				})
				return
			}
			newParentID = &pid
		}

		if err := h.hierarchyService.MoveUnit(c.Request.Context(), unitID, newParentID); err != nil {
			if errors.IsValidationError(err) {
				c.JSON(http.StatusBadRequest, dto.ErrorResponse{
					Error:   "validation_error",
					Message: err.Error(),
				})
				return
			}
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "server_error",
				Message: "failed to move unit",
			})
			return
		}

		// Recargar la unidad después de moverla
		unit, err = h.unitRepo.FindByID(c.Request.Context(), unitID, false)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "database_error",
				Message: "failed to reload unit",
			})
			return
		}
	} else {
		// Solo actualizar sin mover
		if err := h.unitRepo.Update(c.Request.Context(), unit); err != nil {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "database_error",
				Message: "failed to update unit",
			})
			return
		}
	}

	c.JSON(http.StatusOK, dto.ToUnitResponse(unit))
}

// DeleteUnit godoc
// @Summary Delete unit (soft delete)
// @Tags units
// @Param id path string true "Unit ID"
// @Success 204
// @Failure 404 {object} dto.ErrorResponse
// @Router /units/{id} [delete]
func (h *AcademicUnitHandler) DeleteUnit(c *gin.Context) {
	idStr := c.Param("id")
	
	unitID, err := valueobject.ParseUnitID(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_id",
			Message: "invalid unit ID format",
		})
		return
	}

	if err := h.unitRepo.SoftDelete(c.Request.Context(), unitID); err != nil {
		if errors.IsNotFoundError(err) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "not_found",
				Message: "unit not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "database_error",
			Message: "failed to delete unit",
		})
		return
	}

	c.Status(http.StatusNoContent)
}
```

**Commit:** `feat(handlers): implement school and academic unit HTTP handlers`

---

### TASK-05: Router Configuration (20min)

**Archivo:** `internal/infrastructure/http/router/router.go`

```go
package router

import (
	"github.com/gin-gonic/gin"
	
	"github.com/EduGoGroup/edugo-api-administracion/internal/application/service"
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/repository"
	domainService "github.com/EduGoGroup/edugo-api-administracion/internal/domain/service"
	"github.com/EduGoGroup/edugo-api-administracion/internal/infrastructure/http/handler"
)

// Config configuración para el router
type Config struct {
	SchoolRepo repository.SchoolRepository
	UnitRepo   repository.AcademicUnitRepository
}

// SetupRouter configura todas las rutas de la API
func SetupRouter(cfg *Config) *gin.Engine {
	router := gin.Default()

	// Middleware global
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1
	v1 := router.Group("/api/v1")
	{
		// Inicializar servicios
		domainSvc := domainService.NewAcademicUnitDomainService()
		hierarchySvc := service.NewHierarchyService(cfg.UnitRepo, cfg.SchoolRepo, domainSvc)

		// Handlers
		schoolHandler := handler.NewSchoolHandler(cfg.SchoolRepo)
		unitHandler := handler.NewAcademicUnitHandler(cfg.UnitRepo, hierarchySvc)

		// School routes
		schools := v1.Group("/schools")
		{
			schools.POST("", schoolHandler.CreateSchool)
			schools.GET("", schoolHandler.ListSchools)
			schools.GET("/:id", schoolHandler.GetSchool)
			schools.PUT("/:id", schoolHandler.UpdateSchool)
			schools.DELETE("/:id", schoolHandler.DeleteSchool)
		}

		// Academic Unit routes
		units := v1.Group("/units")
		{
			units.POST("", unitHandler.CreateUnit)
			units.GET("", unitHandler.ListUnits)
			units.GET("/:id", unitHandler.GetUnit)
			units.GET("/:id/tree", unitHandler.GetUnitTree) // ltree endpoint!
			units.PUT("/:id", unitHandler.UpdateUnit)
			units.DELETE("/:id", unitHandler.DeleteUnit)
		}
	}

	return router
}

// corsMiddleware configuración básica de CORS
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
```

**Commit:** `feat(router): configure API routes with Gin`

---

### TASK-06: Integration Tests con STUBS (30min)

**Archivo:** `test/integration/http_api_test.go`

```go
//go:build integration

package integration

import (
	"testing"
)

// STUB_FASE2: Estos tests requieren servidor HTTP corriendo
// Completar en FASE 2 con Claude Code Local

// TestSchoolAPI_CreateAndGet verifica flujo de creación y obtención
func TestSchoolAPI_CreateAndGet(t *testing.T) {
	t.Skip("STUB_FASE2: Requiere servidor HTTP - Completar en FASE 2")

	// TODO_FASE2: Descomentar y ejecutar
	// 1. Levantar servidor Gin en puerto de test
	// 2. POST /api/v1/schools
	// 3. Verificar response 201
	// 4. GET /api/v1/schools/:id
	// 5. Verificar que devuelve la escuela creada
}

// TestUnitAPI_CreateTree verifica creación de jerarquía
func TestUnitAPI_CreateTree(t *testing.T) {
	t.Skip("STUB_FASE2: Requiere servidor HTTP - Completar en FASE 2")

	// TODO_FASE2: Descomentar y ejecutar
	// 1. Crear escuela via API
	// 2. Crear grado (raíz) via API
	// 3. Crear sección (hijo) via API
	// 4. GET /api/v1/units/:id/tree
	// 5. Verificar que el árbol usa ltree correctamente
}

// TestUnitAPI_MoveSubtree verifica mover jerarquía
func TestUnitAPI_MoveSubtree(t *testing.T) {
	t.Skip("STUB_FASE2: Requiere servidor HTTP - Completar en FASE 2")

	// TODO_FASE2: Descomentar y ejecutar
	// 1. Crear jerarquía: Grade1 -> Section -> Club
	//                      Grade2
	// 2. PUT /api/v1/units/:section_id (mover a Grade2)
	// 3. Verificar que Section y Club se movieron
	// 4. Usar endpoint /tree para validar
}

// TestUnitAPI_ListByDepth verifica filtro por profundidad (ltree!)
func TestUnitAPI_ListByDepth(t *testing.T) {
	t.Skip("STUB_FASE2: Requiere servidor HTTP - Completar en FASE 2")

	// TODO_FASE2: Descomentar y ejecutar
	// 1. Crear jerarquía de 3 niveles
	// 2. GET /api/v1/units?school_id=X&depth=1
	// 3. Verificar que solo retorna nivel 1
	// 4. GET /api/v1/units?school_id=X&depth=2
	// 5. Verificar que solo retorna nivel 2
}

// TestAPI_ErrorHandling verifica manejo de errores
func TestAPI_ErrorHandling(t *testing.T) {
	t.Skip("STUB_FASE2: Requiere servidor HTTP - Completar en FASE 2")

	// TODO_FASE2: Descomentar y ejecutar
	// 1. POST con JSON inválido -> 400
	// 2. GET con ID inexistente -> 404
	// 3. POST con código duplicado -> 400
	// 4. PUT para crear ciclo -> 400
}
```

**Commit:** `test(integration): add HTTP API test stubs for Phase 2`

---

### TASK-07: Crear main.go de ejemplo (15min)

**Archivo:** `cmd/api/main.go`

```go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	
	"github.com/EduGoGroup/edugo-api-administracion/internal/infrastructure/http/router"
	"github.com/EduGoGroup/edugo-api-administracion/internal/infrastructure/persistence/postgres/repository"
)

func main() {
	// STUB_FASE2: Configuración real de DB en Fase 2
	// Por ahora, solo ejemplo de estructura

	// Configuración de base de datos (desde env vars)
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "edugo_user")
	dbPass := getEnv("DB_PASSWORD", "edugo_pass")
	dbName := getEnv("DB_NAME", "edugo_admin")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("✅ Connected to database")

	// Inicializar repositorios
	schoolRepo := repository.NewPostgresSchoolRepository(db)
	unitRepo := repository.NewPostgresAcademicUnitRepository(db)

	// Configurar router
	routerCfg := &router.Config{
		SchoolRepo: schoolRepo,
		UnitRepo:   unitRepo,
	}

	r := router.SetupRouter(routerCfg)

	// Iniciar servidor
	port := getEnv("PORT", "8080")
	log.Printf("🚀 Server starting on port %s", port)
	
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
```

**Commit:** `feat(main): add API server entry point`

---

### TASK-08: Documentación y Handoff (30min)

**Archivo:** `HANDOFF_SPRINT04_FASE1_TO_FASE2.md`

```markdown
# Handoff: Sprint-04 Fase 1 (Web) → Fase 2 (Local)

**Sprint:** Sprint-04 - Services/API
**Ejecutor Fase 1:** Claude Code Web
**Fecha:** [AUTO-GENERATED]
**Branch:** `feature/sprint-04-services-api`

---

## ✅ COMPLETADO EN FASE 1 (Claude Web)

### 1. DTOs Implementados

**Archivos creados:**
- `internal/infrastructure/http/dto/school_dto.go`
- `internal/infrastructure/http/dto/academic_unit_dto.go`
- `internal/infrastructure/http/dto/common_dto.go`

**DTOs disponibles:**
- CreateSchoolRequest, UpdateSchoolRequest, SchoolResponse
- CreateUnitRequest, UpdateUnitRequest, UnitResponse
- UnitTreeNode (para árbol jerárquico con ltree)
- ErrorResponse, SuccessResponse, PaginationMeta

**Validaciones:** Usando `binding` tags de Gin/validator

---

### 2. Application Service - HierarchyService

**Archivo:** `internal/application/service/hierarchy_service.go`

**Métodos implementados:**
- `CreateUnit()` - Crea unidad con validaciones (escuela existe, código único, etc.)
- `GetUnitTree()` - Obtiene árbol usando ltree (Sprint-03!)
- `MoveUnit()` - Mueve unidad usando MoveSubtree ltree
- `ValidateNoCircularReference()` - Previene ciclos usando FindDescendants ltree

**Tests unitarios:** `hierarchy_service_test.go` con mocks

**Aprovecha Sprint-03:**
- ✅ FindDescendants para obtener árbol completo
- ✅ MoveSubtree para reorganizar jerarquías
- ✅ FindBySchoolIDAndDepth para filtrado por nivel

---

### 3. HTTP Handlers

**Archivos creados:**
- `internal/infrastructure/http/handler/school_handler.go`
- `internal/infrastructure/http/handler/academic_unit_handler.go`

**Endpoints implementados (10):**

**Schools (5):**
- POST   /api/v1/schools
- GET    /api/v1/schools
- GET    /api/v1/schools/:id
- PUT    /api/v1/schools/:id
- DELETE /api/v1/schools/:id

**Academic Units (5):**
- POST   /api/v1/units
- GET    /api/v1/units (con filtro ?depth= usando ltree!)
- GET    /api/v1/units/:id
- GET    /api/v1/units/:id/tree (árbol completo con ltree!)
- PUT    /api/v1/units/:id (incluye mover unidad)
- DELETE /api/v1/units/:id

**Características:**
- Validación de DTOs con Gin binding
- Manejo de errores con códigos HTTP apropiados
- Conversión entity ↔ DTO
- Documentación Swagger con anotaciones

---

### 4. Router Configuration

**Archivo:** `internal/infrastructure/http/router/router.go`

**Configurado:**
- Gin router con middleware de recovery
- CORS básico
- Health check en /health
- Agrupación de rutas en /api/v1
- Inyección de dependencias (repositorios)

---

### 5. Main Entry Point

**Archivo:** `cmd/api/main.go`

**Funcionalidad:**
- Conexión a PostgreSQL desde env vars
- Inicialización de repositorios
- Configuración de router
- Servidor HTTP en puerto configurable

---

### 6. Tests con STUBS

**Archivo:** `test/integration/http_api_test.go`

**Tests estructurados (todos con t.Skip):**
1. TestSchoolAPI_CreateAndGet
2. TestUnitAPI_CreateTree
3. TestUnitAPI_MoveSubtree
4. TestUnitAPI_ListByDepth
5. TestAPI_ErrorHandling

**Cada test tiene:**
- ⚠️ `t.Skip("STUB_FASE2: Requiere servidor HTTP")`
- Comentarios `TODO_FASE2` con pasos detallados

---

## ⏸️ PENDIENTE PARA FASE 2 (Claude Local)

### 1. Ejecutar Servidor HTTP ⚠️ CRÍTICO

**Razón:** Requiere levantar Gin server en local

**Tareas Fase 2:**
1. Configurar variables de entorno (DB_HOST, DB_PORT, etc.)
2. Ejecutar `go run cmd/api/main.go`
3. Verificar que servidor levanta en puerto 8080
4. Probar health check: `curl http://localhost:8080/health`

---

### 2. Descomentar y Ejecutar Tests E2E ⚠️ CRÍTICO

**Archivo:** `test/integration/http_api_test.go`

**Para cada test:**
1. Quitar `t.Skip()`
2. Descomentar código
3. Implementar helper para levantar servidor de test
4. Ejecutar requests HTTP (usar httptest o testcontainers)

**Ejemplo de helper necesario:**

```go
func setupTestServer(t *testing.T) (*gin.Engine, *sql.DB, func()) {
	db, cleanup := setupTestDB(t)
	
	cfg := &router.Config{
		SchoolRepo: repository.NewPostgresSchoolRepository(db),
		UnitRepo:   repository.NewPostgresAcademicUnitRepository(db),
	}
	
	router := router.SetupRouter(cfg)
	
	return router, db, cleanup
}
```

---

### 3. Validaciones Específicas

**Test:** `TestSchoolAPI_CreateAndGet`
- Crear escuela via POST
- Verificar response 201 con SchoolResponse válido
- GET por ID debe retornar misma escuela
- Validar que timestamps se generan

**Test:** `TestUnitAPI_CreateTree`
- Crear grado (raíz)
- Crear 2 secciones bajo el grado
- Crear club bajo sección
- GET /units/:id/tree del grado
- Verificar que árbol tiene estructura correcta
- **Validar que usa ltree** (verificar que hijos están ordenados por path)

**Test:** `TestUnitAPI_MoveSubtree`
- Crear Grade1 -> Section -> Club
- Crear Grade2 (vacío)
- PUT /units/:section_id con parent_unit_id = Grade2
- Verificar que Section se movió
- Verificar que Club sigue siendo hijo de Section
- GET /units/:grade2_id/tree debe mostrar Section y Club

**Test:** `TestUnitAPI_ListByDepth`
- Crear jerarquía de 3 niveles
- GET /units?school_id=X&depth=1 debe retornar solo grados
- GET /units?school_id=X&depth=2 debe retornar solo secciones
- Verificar que cuenta de resultados es correcta

**Test:** `TestAPI_ErrorHandling`
- POST con JSON inválido → 400
- POST con field faltante → 400 con detalles
- GET con UUID inválido → 400
- GET con ID inexistente → 404
- POST con código duplicado → 400
- PUT para crear ciclo → 400 con mensaje claro

---

### 4. Tests Manuales con Postman/curl

**Collection de Postman recomendada:**

```bash
# 1. Health check
curl http://localhost:8080/health

# 2. Crear escuela
curl -X POST http://localhost:8080/api/v1/schools \
  -H "Content-Type: application/json" \
  -d '{"name": "Test School", "code": "TS001", "address": "123 Main St"}'

# 3. Crear grado
curl -X POST http://localhost:8080/api/v1/units \
  -H "Content-Type: application/json" \
  -d '{"school_id": "SCHOOL_ID", "type": "grade", "name": "Grade 1", "code": "G1"}'

# 4. Obtener árbol
curl http://localhost:8080/api/v1/units/UNIT_ID/tree

# 5. Filtrar por profundidad
curl http://localhost:8080/api/v1/units?school_id=SCHOOL_ID&depth=1
```

---

### 5. Validar Integración con ltree (Sprint-03)

**Endpoints que DEBEN usar ltree:**

| Endpoint | Método ltree usado | Validación |
|----------|-------------------|------------|
| GET /units/:id/tree | FindDescendants | Árbol completo en 1 query |
| GET /units?depth=N | FindBySchoolIDAndDepth | Filtra por nlevel() |
| PUT /units/:id (move) | MoveSubtree | Actualiza paths en cascada |

**Cómo validar:**
1. Crear jerarquía de 100+ unidades
2. Medir tiempo de GET /units/:id/tree
3. Verificar que es rápido (< 100ms)
4. Confirmar en logs de PostgreSQL que usa índice GIST

---

### 6. Documentación Swagger (Opcional)

Si hay tiempo, generar documentación automática:

```bash
go get -u github.com/swaggo/swag/cmd/swag
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files

swag init -g cmd/api/main.go
```

Endpoint: `http://localhost:8080/swagger/index.html`

---

## 📊 COBERTURA ESPERADA POST-FASE 2

### Código
- **HierarchyService:** >= 80% cobertura
- **Handlers:** >= 60% cobertura (por validaciones de entrada)
- **Router:** 100% (es simple)

### Funcionalidad
- ✅ CRUD completo de escuelas
- ✅ CRUD completo de unidades
- ✅ Árbol jerárquico con ltree
- ✅ Mover unidades con validación de ciclos
- ✅ Filtrado por profundidad
- ✅ Manejo de errores HTTP

---

## 🚀 COMANDOS PARA FASE 2

```bash
# Checkout
git checkout feature/sprint-04-services-api
git pull origin feature/sprint-04-services-api

# Levantar PostgreSQL (si no está corriendo)
docker-compose up -d postgres

# Ejecutar migraciones (incluye 013 de Sprint-03)
migrate -path migrations -database "postgresql://edugo_user:edugo_pass@localhost:5432/edugo_admin?sslmode=disable" up

# Levantar servidor
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=edugo_user
export DB_PASSWORD=edugo_pass
export DB_NAME=edugo_admin
export PORT=8080

go run cmd/api/main.go

# En otra terminal, ejecutar tests E2E
go test -tags=integration ./test/integration/... -v -run TestAPI

# Tests manuales
curl http://localhost:8080/health
```

---

## 📝 NOTAS

- Sprint-04 **aprovecha completamente** el trabajo de Sprint-03 (ltree)
- El endpoint `/tree` sería muy lento sin ltree
- El filtro por profundidad no sería posible sin `nlevel()`
- MoveSubtree garantiza consistencia de paths

**¡Éxito en Fase 2!** 🚀

---

**Fin del documento de handoff**
```

**Commit:** `docs: add Sprint-04 Phase 1 handoff document`

---

### TASK-09: Validación Final Fase 1

```bash
# Compilación
go build ./...

# Tests unitarios
go test ./internal/application/service/... -v
go test ./internal/infrastructure/http/... -v

# Verificar estructura
tree internal/
```

**Checklist:**
- [ ] DTOs compilando
- [ ] HierarchyService con tests unitarios
- [ ] Handlers implementados
- [ ] Router configurado
- [ ] main.go creado
- [ ] Tests E2E con stubs
- [ ] Handoff document creado
- [ ] Código sin TODOs (excepto TODO_FASE2)

**Commit final:** `chore: Sprint-04 Phase 1 complete - ready for Phase 2`

---

## 📊 CHECKLIST FINAL FASE 1

Antes de crear PR, verifica:

### Código
- [ ] DTOs con validaciones Gin
- [ ] HierarchyService implementado
- [ ] Tests unitarios de servicio
- [ ] 10 endpoints HTTP implementados
- [ ] Router configurado
- [ ] main.go funcional (sin ejecutar)
- [ ] Compilación OK

### Tests
- [ ] Tests unitarios pasando
- [ ] Tests E2E con t.Skip()
- [ ] Comentarios TODO_FASE2 claros

### Documentación
- [ ] HANDOFF completo
- [ ] Comentarios Swagger en handlers
- [ ] README actualizado (opcional)

---

## 🎯 ENTREGABLES FASE 1

Al finalizar, debes tener:
1. ✅ Branch `feature/sprint-04-services-api`
2. ✅ 7-8 commits organizados
3. ✅ Código completo que compila
4. ✅ Tests unitarios pasando
5. ✅ HANDOFF document para Fase 2
6. ✅ Sin código ejecutable de servidor (eso es Fase 2)

---

**¿Listo para comenzar Sprint-04 Fase 1?** 🚀
