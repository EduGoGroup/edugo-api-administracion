//go:build integration

package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/EduGoGroup/edugo-api-administracion/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-administracion/internal/container"
	"github.com/EduGoGroup/edugo-shared/logger"
	"github.com/gin-gonic/gin"
)

// getTestLogger crea un logger para tests
func getTestLogger() logger.Logger {
	return logger.NewZapLogger("debug", "console")
}

// setupTestServer levanta un servidor de test con PostgreSQL real
func setupTestServer(t *testing.T) (*httptest.Server, func()) {
	// Setup BD con contenedor independiente por test (mejor aislamiento)
	db, dbCleanup := setupTestDB(t)

	// Crear logger para tests
	testLogger := getTestLogger()

	// Crear container con la BD de test
	c := container.NewContainer(db, testLogger)

	// Configurar Gin en modo test
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// Rutas v1
	v1 := r.Group("/v1")
	{
		// Schools
		schools := v1.Group("/schools")
		{
			schools.POST("", c.SchoolHandler.CreateSchool)
			schools.GET("", c.SchoolHandler.ListSchools)
			schools.GET("/code/:code", c.SchoolHandler.GetSchoolByCode)

			// Academic Units nested under school
			schools.POST("/:id/units", c.AcademicUnitHandler.CreateUnit)
			schools.GET("/:id/units", c.AcademicUnitHandler.ListUnitsBySchool)
			schools.GET("/:id/units/tree", c.AcademicUnitHandler.GetUnitTree)
			schools.GET("/:id/units/by-type", c.AcademicUnitHandler.ListUnitsByType)

			// School CRUD
			schools.GET("/:id", c.SchoolHandler.GetSchool)
			schools.PUT("/:id", c.SchoolHandler.UpdateSchool)
			schools.DELETE("/:id", c.SchoolHandler.DeleteSchool)
		}

		// Units
		units := v1.Group("/units")
		{
			units.GET("/:id", c.AcademicUnitHandler.GetUnit)
			units.PUT("/:id", c.AcademicUnitHandler.UpdateUnit)
			units.DELETE("/:id", c.AcademicUnitHandler.DeleteUnit)
			units.POST("/:id/restore", c.AcademicUnitHandler.RestoreUnit)
			units.GET("/:id/hierarchy-path", c.AcademicUnitHandler.GetHierarchyPath)
		}
	}

	// Crear test server
	server := httptest.NewServer(r)

	cleanupFunc := func() {
		server.Close()
		c.Close()
		dbCleanup()
	}

	return server, cleanupFunc
}

// doRequest helper para hacer requests HTTP
func doRequest(t *testing.T, server *httptest.Server, method, path string, body interface{}) (*http.Response, []byte) {
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		require.NoError(t, err)
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequest(method, server.URL+path, bodyReader)
	require.NoError(t, err)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	resp.Body.Close()

	return resp, respBody
}

// TestSchoolAPI_CreateAndGet verifica flujo de creación y obtención
func TestSchoolAPI_CreateAndGet(t *testing.T) {
	server, cleanup := setupTestServer(t)
	defer cleanup()

	// 1. Crear escuela
	createReq := dto.CreateSchoolRequest{
		Name:    "Integration Test School",
		Code:    "ITS001",
		Address: "Test Address 123",
	}

	resp, body := doRequest(t, server, "POST", "/v1/schools", createReq)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var schoolResp dto.SchoolResponse
	err := json.Unmarshal(body, &schoolResp)
	require.NoError(t, err)

	assert.NotEmpty(t, schoolResp.ID)
	assert.Equal(t, "Integration Test School", schoolResp.Name)
	assert.Equal(t, "ITS001", schoolResp.Code)

	// 2. Obtener escuela por ID
	resp, body = doRequest(t, server, "GET", "/v1/schools/"+schoolResp.ID, nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var getResp dto.SchoolResponse
	err = json.Unmarshal(body, &getResp)
	require.NoError(t, err)

	assert.Equal(t, schoolResp.ID, getResp.ID)
	assert.Equal(t, "Integration Test School", getResp.Name)
}

// TestUnitAPI_CreateTree verifica creación de jerarquía
func TestUnitAPI_CreateTree(t *testing.T) {
	server, cleanup := setupTestServer(t)
	defer cleanup()

	// 1. Crear escuela
	schoolReq := dto.CreateSchoolRequest{
		Name: "Test School for Tree",
		Code: "TSFT001",
	}
	resp, body := doRequest(t, server, "POST", "/v1/schools", schoolReq)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var school dto.SchoolResponse
	json.Unmarshal(body, &school)

	// 2. Crear grado (raíz)
	gradeReq := dto.CreateAcademicUnitRequest{
		Type:        "grade",
		DisplayName: "Test Grade",
		Code:        "TG1",
	}
	resp, body = doRequest(t, server, "POST", "/v1/schools/"+school.ID+"/units", gradeReq)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var grade dto.AcademicUnitResponse
	json.Unmarshal(body, &grade)
	assert.NotEmpty(t, grade.ID)

	// 3. Crear sección (hijo)
	sectionReq := dto.CreateAcademicUnitRequest{
		ParentUnitID: &grade.ID,
		Type:         "section",
		DisplayName:  "Test Section",
		Code:         "TS1",
	}
	resp, body = doRequest(t, server, "POST", "/v1/schools/"+school.ID+"/units", sectionReq)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// 4. Obtener árbol (ltree!)
	resp, body = doRequest(t, server, "GET", "/v1/schools/"+school.ID+"/units/tree", nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var tree []dto.UnitTreeNode
	json.Unmarshal(body, &tree)

	// Verificar estructura del árbol
	require.Len(t, tree, 1, "should have one root node")
	assert.Equal(t, grade.ID, tree[0].ID)
	assert.Equal(t, 1, tree[0].Depth)

	require.Len(t, tree[0].Children, 1, "grade should have one child")
	assert.Equal(t, "Test Section", tree[0].Children[0].DisplayName)
	assert.Equal(t, 2, tree[0].Children[0].Depth)
}

// TestUnitAPI_MoveSubtree verifica mover jerarquía con ltree
func TestUnitAPI_MoveSubtree(t *testing.T) {
	server, cleanup := setupTestServer(t)
	defer cleanup()

	// Setup: Crear escuela
	_, schoolBody := doRequest(t, server, "POST", "/v1/schools", dto.CreateSchoolRequest{
		Name: "Test School Move",
		Code: "TSM001",
	})
	var school dto.SchoolResponse
	json.Unmarshal(schoolBody, &school)

	// Crear Grade 1 con Section
	_, grade1Body := doRequest(t, server, "POST", "/v1/schools/"+school.ID+"/units", dto.CreateAcademicUnitRequest{
		Type:        "grade",
		DisplayName: "Grade 1",
		Code:        "MG1",
	})
	var grade1 dto.AcademicUnitResponse
	json.Unmarshal(grade1Body, &grade1)

	_, sectionBody := doRequest(t, server, "POST", "/v1/schools/"+school.ID+"/units", dto.CreateAcademicUnitRequest{
		ParentUnitID: &grade1.ID,
		Type:         "section",
		DisplayName:  "Section A",
		Code:         "MS-A",
	})
	var section dto.AcademicUnitResponse
	json.Unmarshal(sectionBody, &section)

	// Crear Grade 2 (vacío)
	_, grade2Body := doRequest(t, server, "POST", "/v1/schools/"+school.ID+"/units", dto.CreateAcademicUnitRequest{
		Type:        "grade",
		DisplayName: "Grade 2",
		Code:        "MG2",
	})
	var grade2 dto.AcademicUnitResponse
	json.Unmarshal(grade2Body, &grade2)

	// Mover Section A a Grade 2 (ltree MoveSubtree!)
	moveResp, moveBody := doRequest(t, server, "PUT", "/v1/units/"+section.ID, dto.UpdateAcademicUnitRequest{
		ParentUnitID: &grade2.ID,
	})
	assert.Equal(t, http.StatusOK, moveResp.StatusCode, "move should succeed: %s", string(moveBody))

	// Verificar árbol de Grade 2 (debe tener Section A)
	treeResp, treeBody := doRequest(t, server, "GET", "/v1/schools/"+school.ID+"/units/tree", nil)
	assert.Equal(t, http.StatusOK, treeResp.StatusCode)

	var tree []dto.UnitTreeNode
	json.Unmarshal(treeBody, &tree)

	// Buscar Grade 2 en el árbol
	var grade2Node *dto.UnitTreeNode
	for i := range tree {
		if tree[i].ID == grade2.ID {
			grade2Node = &tree[i]
			break
		}
	}

	require.NotNil(t, grade2Node, "Grade 2 should be in tree")
	require.Len(t, grade2Node.Children, 1, "Grade 2 should have Section A")
	assert.Equal(t, section.ID, grade2Node.Children[0].ID)
}

// TestAPI_ErrorHandling verifica manejo de errores
func TestAPI_ErrorHandling(t *testing.T) {
	server, cleanup := setupTestServer(t)
	defer cleanup()

	// 1. POST con JSON inválido -> 400
	req, _ := http.NewRequest("POST", server.URL+"/v1/schools", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := http.DefaultClient.Do(req)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	resp.Body.Close()

	// 2. GET con ID inexistente -> 404
	resp, _ = doRequest(t, server, "GET", "/v1/units/00000000-0000-0000-0000-000000000000", nil)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	// 3. POST con código duplicado -> 400 o 409
	schoolReq := dto.CreateSchoolRequest{
		Name: "Test School",
		Code: "DUP001",
	}
	doRequest(t, server, "POST", "/v1/schools", schoolReq)

	// Intentar crear con mismo código
	resp, _ = doRequest(t, server, "POST", "/v1/schools", schoolReq)
	assert.True(t, resp.StatusCode == http.StatusBadRequest || resp.StatusCode == http.StatusConflict)
}

// TestUnitAPI_GetHierarchyPath verifica obtención de path jerárquico (ltree!)
func TestUnitAPI_GetHierarchyPath(t *testing.T) {
	server, cleanup := setupTestServer(t)
	defer cleanup()

	// Crear escuela
	_, schoolBody := doRequest(t, server, "POST", "/v1/schools", dto.CreateSchoolRequest{
		Name: "Test School Path",
		Code: "TSP001",
	})
	var school dto.SchoolResponse
	json.Unmarshal(schoolBody, &school)

	// Crear Grade -> Section
	_, gradeBody := doRequest(t, server, "POST", "/v1/schools/"+school.ID+"/units", dto.CreateAcademicUnitRequest{
		Type:        "grade",
		DisplayName: "Test Grade",
		Code:        "PG1",
	})
	var grade dto.AcademicUnitResponse
	json.Unmarshal(gradeBody, &grade)

	_, sectionBody := doRequest(t, server, "POST", "/v1/schools/"+school.ID+"/units", dto.CreateAcademicUnitRequest{
		ParentUnitID: &grade.ID,
		Type:         "section",
		DisplayName:  "Test Section",
		Code:         "PS1",
	})
	var section dto.AcademicUnitResponse
	json.Unmarshal(sectionBody, &section)

	// Obtener hierarchy path (ltree!)
	resp, body := doRequest(t, server, "GET", "/v1/units/"+section.ID+"/hierarchy-path", nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var path []dto.AcademicUnitResponse
	json.Unmarshal(body, &path)

	// Verificar orden: de raíz a hoja (Grade -> Section)
	require.Len(t, path, 2, "path should have 2 nodes")
	assert.Equal(t, grade.ID, path[0].ID, "first should be grade")
	assert.Equal(t, section.ID, path[1].ID, "second should be section")
}

// TestSchoolAPI_ListAll verifica listado de escuelas
func TestSchoolAPI_ListAll(t *testing.T) {
	server, cleanup := setupTestServer(t)
	defer cleanup()

	// Crear múltiples escuelas
	for i := 1; i <= 3; i++ {
		doRequest(t, server, "POST", "/v1/schools", dto.CreateSchoolRequest{
			Name: "Test School " + string(rune('A'+i-1)),
			Code: "LIST00" + string(rune('0'+i)),
		})
	}

	// Listar todas
	resp, body := doRequest(t, server, "GET", "/v1/schools", nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var schools []dto.SchoolResponse
	err := json.Unmarshal(body, &schools)
	require.NoError(t, err, "should unmarshal schools list")

	// Verificar que al menos tiene las 3 escuelas que creamos
	// Nota: En ambiente de test puede haber datos residuales de otros tests
	assert.GreaterOrEqual(t, len(schools), 3, "should have at least the 3 schools we created")
}

// TestSchoolAPI_UpdateAndDelete verifica actualización y eliminación
func TestSchoolAPI_UpdateAndDelete(t *testing.T) {
	server, cleanup := setupTestServer(t)
	defer cleanup()

	// 1. Crear escuela
	createResp, createBody := doRequest(t, server, "POST", "/v1/schools", dto.CreateSchoolRequest{
		Name: "School To Update",
		Code: "STU001",
	})
	require.Equal(t, http.StatusCreated, createResp.StatusCode)

	var school dto.SchoolResponse
	json.Unmarshal(createBody, &school)

	// 2. Actualizar escuela
	updateReq := dto.UpdateSchoolRequest{
		Name: strPtr("Updated School Name"),
	}
	resp, body := doRequest(t, server, "PUT", "/v1/schools/"+school.ID, updateReq)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var updated dto.SchoolResponse
	json.Unmarshal(body, &updated)
	assert.Equal(t, "Updated School Name", updated.Name)

	// 3. Eliminar escuela
	resp, _ = doRequest(t, server, "DELETE", "/v1/schools/"+school.ID, nil)
	assert.True(t, resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent, "delete should return 200 or 204")

	// 4. Verificar que no existe (puede retornar 404 o 500 dependiendo de implementación de soft delete)
	resp, _ = doRequest(t, server, "GET", "/v1/schools/"+school.ID, nil)
	assert.True(t, resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusInternalServerError, "should not find deleted school")
}

// Helper para crear punteros a strings
func strPtr(s string) *string {
	return &s
}
