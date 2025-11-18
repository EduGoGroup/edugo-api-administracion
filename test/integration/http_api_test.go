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
	// 4. GET /api/v1/schools/:schoolId/units/tree
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
	// 2. GET /api/v1/schools/:schoolId/units (con filtro depth=1 si se implementa)
	// 3. Verificar que solo retorna nivel 1
	// 4. GET /api/v1/schools/:schoolId/units (con filtro depth=2)
	// 5. Verificar que solo retorna nivel 2
}

// TestAPI_ErrorHandling verifica manejo de errores
func TestAPI_ErrorHandling(t *testing.T) {
	t.Skip("STUB_FASE2: Requiere servidor HTTP - Completar en FASE 2")

	// TODO_FASE2: Descomentar y ejecutar
	// 1. POST con JSON inválido -> 400
	// 2. GET con ID inexistente -> 404
	// 3. POST con código duplicado -> 400 o 409
	// 4. PUT para crear ciclo -> 400
}

// TestUnitAPI_GetHierarchyPath verifica obtención de path jerárquico (ltree!)
func TestUnitAPI_GetHierarchyPath(t *testing.T) {
	t.Skip("STUB_FASE2: Requiere servidor HTTP - Completar en FASE 2")

	// TODO_FASE2: Descomentar y ejecutar
	// 1. Crear jerarquía: School -> Grade -> Section -> Club
	// 2. GET /api/v1/units/:club_id/hierarchy-path
	// 3. Verificar que retorna el path completo desde la raíz
	// 4. Validar que el orden es correcto (de raíz a hoja)
}

// TestSchoolAPI_ListAll verifica listado de escuelas
func TestSchoolAPI_ListAll(t *testing.T) {
	t.Skip("STUB_FASE2: Requiere servidor HTTP - Completar en FASE 2")

	// TODO_FASE2: Descomentar y ejecutar
	// 1. Crear múltiples escuelas
	// 2. GET /api/v1/schools
	// 3. Verificar que retorna todas las escuelas
	// 4. Validar formato de respuesta
}

// TestSchoolAPI_UpdateAndDelete verifica actualización y eliminación
func TestSchoolAPI_UpdateAndDelete(t *testing.T) {
	t.Skip("STUB_FASE2: Requiere servidor HTTP - Completar en FASE 2")

	// TODO_FASE2: Descomentar y ejecutar
	// 1. Crear escuela
	// 2. PUT /api/v1/schools/:id con cambios
	// 3. GET /api/v1/schools/:id y verificar cambios
	// 4. DELETE /api/v1/schools/:id
	// 5. GET /api/v1/schools/:id debe retornar 404
}

// TestUnitAPI_RestoreDeleted verifica restauración de unidades eliminadas
func TestUnitAPI_RestoreDeleted(t *testing.T) {
	t.Skip("STUB_FASE2: Requiere servidor HTTP - Completar en FASE 2")

	// TODO_FASE2: Descomentar y ejecutar
	// 1. Crear unidad
	// 2. DELETE /api/v1/units/:id (soft delete)
	// 3. Verificar que GET retorna 404
	// 4. POST /api/v1/units/:id/restore
	// 5. Verificar que GET retorna la unidad restaurada
}
