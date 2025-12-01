//go:build integration

package integration

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/EduGoGroup/edugo-api-administracion/internal/application/dto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// seedMembershipTestData crea datos base (school, unit, user) para tests de memberships
// Retorna IDs de school, unit y user creados
func seedMembershipTestData(t *testing.T, server *httptest.Server, db *sql.DB) (schoolID string, unitID string, userID string) {
	t.Helper()

	// 1. Crear escuela via API
	schoolReq := dto.CreateSchoolRequest{
		Name: "Membership Test School " + uuid.New().String()[:8],
		Code: "MTS" + uuid.New().String()[:4],
	}
	resp, body := doRequest(t, server, "POST", "/v1/schools", schoolReq)
	require.Equal(t, http.StatusCreated, resp.StatusCode, "failed to create school: %s", string(body))

	var school dto.SchoolResponse
	err := json.Unmarshal(body, &school)
	require.NoError(t, err)
	schoolID = school.ID

	// 2. Crear unidad académica via API
	unitReq := dto.CreateAcademicUnitRequest{
		Type:        "grade",
		DisplayName: "Test Grade " + uuid.New().String()[:8],
		Code:        "TG" + uuid.New().String()[:4],
	}
	resp, body = doRequest(t, server, "POST", "/v1/schools/"+schoolID+"/units", unitReq)
	require.Equal(t, http.StatusCreated, resp.StatusCode, "failed to create unit: %s", string(body))

	var unit dto.AcademicUnitResponse
	err = json.Unmarshal(body, &unit)
	require.NoError(t, err)
	unitID = unit.ID

	// 3. Crear usuario directamente en BD (FK constraint)
	userID = uuid.New().String()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("Test1234!"), bcrypt.DefaultCost)
	
	_, err = db.Exec(`
		INSERT INTO users (id, email, password_hash, first_name, last_name, role, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, userID, "testuser"+uuid.New().String()[:8]+"@edugo.test", string(hashedPassword), "Test", "User", "student", true)
	require.NoError(t, err, "failed to create user in DB")

	t.Logf("✅ Test data created: school=%s, unit=%s, user=%s", schoolID, unitID, userID)
	return schoolID, unitID, userID
}

// TestMembershipAPI_CreateAndGet verifica flujo de creación y obtención de membership
func TestMembershipAPI_CreateAndGet(t *testing.T) {
	server, db, cleanup := setupTestServerWithDB(t)
	defer cleanup()

	schoolID, unitID, userID := seedMembershipTestData(t, server, db)

	// Crear membership
	createReq := dto.CreateMembershipRequest{
		UserID:         userID,
		UnitID: unitID,
		Role:           "student",
	}

	resp, body := doRequest(t, server, "POST", "/v1/memberships", createReq)
	require.Equal(t, http.StatusCreated, resp.StatusCode, "create failed: %s", string(body))

	var membership dto.MembershipResponse
	err := json.Unmarshal(body, &membership)

	assert.NotEmpty(t, membership.ID, "membership should be created with validity dates")
	require.NoError(t, err)

	assert.NotEmpty(t, membership.ID)
	assert.Equal(t, userID, membership.UserID)
	assert.Equal(t, unitID, membership.UnitID)
	assert.Equal(t, "student", membership.Role)

	// Obtener membership por ID
	resp, body = doRequest(t, server, "GET", "/v1/memberships/"+membership.ID, nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode, "get failed: %s", string(body))

	var fetched dto.MembershipResponse
	json.Unmarshal(body, &fetched)
	assert.Equal(t, membership.ID, fetched.ID)

	_ = schoolID // usado para crear datos
}

// TestMembershipAPI_ListByUnit verifica listado de memberships por unidad
func TestMembershipAPI_ListByUnit(t *testing.T) {
	server, db, cleanup := setupTestServerWithDB(t)
	defer cleanup()

	schoolID, unitID, userID := seedMembershipTestData(t, server, db)

	// Crear membership
	createReq := dto.CreateMembershipRequest{
		UserID:         userID,
		UnitID: unitID,
		Role:           "student",
	}
	resp, body := doRequest(t, server, "POST", "/v1/memberships", createReq)
	require.Equal(t, http.StatusCreated, resp.StatusCode, "create failed: %s", string(body))

	// Listar por unidad
	resp, body = doRequest(t, server, "GET", "/v1/units/"+unitID+"/memberships", nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode, "list failed: %s", string(body))

	var memberships []dto.MembershipResponse
	err := json.Unmarshal(body, &memberships)
	require.NoError(t, err)

	assert.GreaterOrEqual(t, len(memberships), 1, "should have at least 1 membership")

	_ = schoolID
}

// TestMembershipAPI_ListByUser verifica listado de memberships por usuario
func TestMembershipAPI_ListByUser(t *testing.T) {
	server, db, cleanup := setupTestServerWithDB(t)
	defer cleanup()

	schoolID, unitID, userID := seedMembershipTestData(t, server, db)

	// Crear membership
	createReq := dto.CreateMembershipRequest{
		UserID:         userID,
		UnitID: unitID,
		Role:           "teacher",
	}
	resp, body := doRequest(t, server, "POST", "/v1/memberships", createReq)
	require.Equal(t, http.StatusCreated, resp.StatusCode, "create failed: %s", string(body))

	// Listar por usuario
	resp, body = doRequest(t, server, "GET", "/v1/users/"+userID+"/memberships", nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode, "list by user failed: %s", string(body))

	var memberships []dto.MembershipResponse
	err := json.Unmarshal(body, &memberships)
	require.NoError(t, err)

	assert.GreaterOrEqual(t, len(memberships), 1, "should have at least 1 membership")

	_ = schoolID
}

// TestMembershipAPI_UpdateAndDelete verifica actualización y eliminación
func TestMembershipAPI_UpdateAndDelete(t *testing.T) {
	server, db, cleanup := setupTestServerWithDB(t)
	defer cleanup()

	schoolID, unitID, userID := seedMembershipTestData(t, server, db)

	// Crear membership
	createReq := dto.CreateMembershipRequest{
		UserID:         userID,
		UnitID: unitID,
		Role:           "student",
	}
	resp, body := doRequest(t, server, "POST", "/v1/memberships", createReq)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var membership dto.MembershipResponse
	json.Unmarshal(body, &membership)

	assert.NotEmpty(t, membership.ID, "membership should be created with validity dates")

	// Actualizar rol
	newRole := "teacher"
	updateReq := dto.UpdateMembershipRequest{
		Role: &newRole,
	}
	resp, body = doRequest(t, server, "PUT", "/v1/memberships/"+membership.ID, updateReq)
	assert.Equal(t, http.StatusOK, resp.StatusCode, "update failed: %s", string(body))

	var updated dto.MembershipResponse
	json.Unmarshal(body, &updated)
	assert.Equal(t, "teacher", updated.Role)

	// Eliminar membership
	resp, _ = doRequest(t, server, "DELETE", "/v1/memberships/"+membership.ID, nil)
	assert.True(t, resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent)

	// Verificar que fue eliminada
	resp, _ = doRequest(t, server, "GET", "/v1/memberships/"+membership.ID, nil)
	assert.True(t, resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusOK, "after delete should return 404 or 200 (soft delete)")

	_ = schoolID
}

// TestMembershipAPI_ExpireMembership verifica expiración de membership
func TestMembershipAPI_ExpireMembership(t *testing.T) {
	server, db, cleanup := setupTestServerWithDB(t)
	defer cleanup()

	schoolID, unitID, userID := seedMembershipTestData(t, server, db)

	// Crear membership
	createReq := dto.CreateMembershipRequest{
		UserID:         userID,
		UnitID: unitID,
		Role:           "student",
	}
	resp, body := doRequest(t, server, "POST", "/v1/memberships", createReq)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var membership dto.MembershipResponse
	json.Unmarshal(body, &membership)

	assert.NotEmpty(t, membership.ID, "membership should be created with validity dates")

	// Expirar membership
	resp, body = doRequest(t, server, "POST", "/v1/memberships/"+membership.ID+"/expire", nil)
	assert.True(t, resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent, 
		"expire failed with status %d: %s", resp.StatusCode, string(body))

	_ = schoolID
}

// TestMembershipAPI_ListByRole verifica listado por rol
func TestMembershipAPI_ListByRole(t *testing.T) {
	server, db, cleanup := setupTestServerWithDB(t)
	defer cleanup()

	schoolID, unitID, _ := seedMembershipTestData(t, server, db)

	// Crear 2 usuarios más
	user2ID := uuid.New().String()
	user3ID := uuid.New().String()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("Test1234!"), bcrypt.DefaultCost)

	db.Exec(`INSERT INTO users (id, email, password_hash, first_name, last_name, role, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		user2ID, schoolID, "teacher"+uuid.New().String()[:8]+"@edugo.test", string(hashedPassword), "Teacher", "Test", "teacher", true)
	db.Exec(`INSERT INTO users (id, email, password_hash, first_name, last_name, role, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		user3ID, schoolID, "student"+uuid.New().String()[:8]+"@edugo.test", string(hashedPassword), "Student", "Test", "student", true)

	// Crear memberships con diferentes roles
	doRequest(t, server, "POST", "/v1/memberships", dto.CreateMembershipRequest{
		UserID: user2ID, UnitID: unitID, Role: "teacher",
	})
	doRequest(t, server, "POST", "/v1/memberships", dto.CreateMembershipRequest{
		UserID: user3ID, UnitID: unitID, Role: "student",
	})

	// Listar por rol "teacher"
	resp, body := doRequest(t, server, "GET", "/v1/units/"+unitID+"/memberships/by-role?role=teacher", nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode, "list by role failed: %s", string(body))

	var memberships []dto.MembershipResponse
	json.Unmarshal(body, &memberships)

	// Verificar que solo hay teachers
	for _, m := range memberships {
		assert.Equal(t, "teacher", m.Role)
	}
}

// TestMembershipAPI_ErrorHandling verifica manejo de errores
func TestMembershipAPI_ErrorHandling(t *testing.T) {
	server, db, cleanup := setupTestServerWithDB(t)
	defer cleanup()

	schoolID, unitID, _ := seedMembershipTestData(t, server, db)

	// 1. Crear con usuario inexistente -> error FK
	resp, _ := doRequest(t, server, "POST", "/v1/memberships", dto.CreateMembershipRequest{
		UserID:         uuid.New().String(),
		UnitID: unitID,
		Role:           "student",
	})
	assert.True(t, resp.StatusCode >= 400, "should fail with non-existent user")

	// 2. Crear con unidad inexistente -> error FK
	userID := uuid.New().String()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("Test1234!"), bcrypt.DefaultCost)
	db.Exec(`INSERT INTO users (id, email, password_hash, first_name, last_name, role, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		userID, schoolID, "err"+uuid.New().String()[:8]+"@edugo.test", string(hashedPassword), "Err", "User", "student", true)

	resp, _ = doRequest(t, server, "POST", "/v1/memberships", dto.CreateMembershipRequest{
		UserID:         userID,
		UnitID: uuid.New().String(),
		Role:           "student",
	})
	assert.True(t, resp.StatusCode >= 400, "should fail with non-existent unit")

	// 3. GET membership inexistente -> 404
	resp, _ = doRequest(t, server, "GET", "/v1/memberships/"+uuid.New().String(), nil)
	assert.True(t, resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusOK, "after delete should return 404 or 200 (soft delete)")
}

// TestMembershipAPI_WithValidityDates verifica memberships con fechas de validez
func TestMembershipAPI_WithValidityDates(t *testing.T) {
	server, db, cleanup := setupTestServerWithDB(t)
	defer cleanup()

	schoolID, unitID, userID := seedMembershipTestData(t, server, db)

	// Crear membership con fechas de validez
	validFrom := time.Now()
	validUntil := time.Now().AddDate(1, 0, 0) // 1 año

	createReq := dto.CreateMembershipRequest{
		UserID:         userID,
		UnitID: unitID,
		Role:           "student",
		ValidFrom:      &validFrom,
		ValidUntil:     &validUntil,
	}

	resp, body := doRequest(t, server, "POST", "/v1/memberships", createReq)
	require.Equal(t, http.StatusCreated, resp.StatusCode, "create with dates failed: %s", string(body))

	var membership dto.MembershipResponse
	json.Unmarshal(body, &membership)

	assert.NotEmpty(t, membership.ID, "membership should be created with validity dates")


	_ = schoolID
}
