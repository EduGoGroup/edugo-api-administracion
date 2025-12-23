# 🗑️ Código Deprecado

> Código marcado para eliminación en próximas versiones

---

## 1. Legacy Handlers (CRÍTICO)

### Ubicación
```
cmd/legacy_handlers.go
```

### Descripción
Archivo completo con handlers HTTP que ya no se usan. Fueron marcados como deprecados en v0.5.0 y deben eliminarse en v0.6.0.

### Código Afectado

```go
// ==================== LEGACY HANDLERS ====================
//
// DEPRECATED: Estos endpoints están deprecated y serán removidos en v0.6.0
// NO implementan lógica real, solo retornan datos mock para compatibilidad.

func CreateUser(c *gin.Context)    // Deprecated - No se usa
func UpdateUser(c *gin.Context)    // Deprecated - No se usa  
func DeleteUser(c *gin.Context)    // Deprecated - No se usa
func CreateSubject(c *gin.Context) // Deprecated - No se usa
func DeleteMaterial(c *gin.Context) // Deprecated - No se usa
func GetGlobalStats(c *gin.Context) // Deprecated - No se usa

// Tipos deprecated
type DeprecatedResponse struct {...}
type SuccessResponse struct {...}
type CreateUserResponse struct {...}
type CreateSubjectResponse struct {...}
type GlobalStatsResponse struct {...}
```

### Problema
- **174 líneas de código muerto**
- Retornan HTTP 410 (Gone) con datos mock
- No están conectados a ninguna ruta en main.go
- Confunden a nuevos desarrolladores
- Swagger los documenta como existentes

### Acción Requerida
```bash
# Eliminar el archivo completo
rm cmd/legacy_handlers.go

# Regenerar Swagger
make swagger
```

### Impacto de Eliminación
- ✅ Ningún impacto funcional (no están en uso)
- ✅ Reduce tamaño del binario
- ✅ Limpia documentación Swagger
- ⚠️ Verificar que no haya clientes usando estos endpoints

### Fecha Límite
**v0.6.0** (según comentarios en el código)

---

## 2. Response Types Duplicados

### Ubicación
```
cmd/legacy_handlers.go (líneas 157-173)
internal/infrastructure/http/handler/ (varios archivos)
```

### Descripción
Tipos de respuesta duplicados en diferentes lugares.

### Código Afectado

```go
// En legacy_handlers.go (DEPRECADO)
type SuccessResponse struct {
    Message string `json:"message"`
}

// En handler/school_handler.go (ACTIVO)
type ErrorResponse struct {
    Error string `json:"error"`
    Code  string `json:"code"`
}
```

### Problema
- Tipos similares definidos en múltiples lugares
- No hay un paquete centralizado de DTOs HTTP
- Inconsistencia en estructura de respuestas

### Acción Requerida
1. Crear `internal/infrastructure/http/dto/response.go`
2. Centralizar tipos de respuesta comunes
3. Eliminar duplicados

### Código Sugerido
```go
// internal/infrastructure/http/dto/response.go
package dto

type ErrorResponse struct {
    Error   string `json:"error"`
    Message string `json:"message,omitempty"`
    Code    string `json:"code"`
    Details any    `json:"details,omitempty"`
}

type SuccessResponse struct {
    Message string `json:"message"`
    Data    any    `json:"data,omitempty"`
}

type PaginatedResponse struct {
    Data       any   `json:"data"`
    TotalCount int64 `json:"total_count"`
    Page       int   `json:"page"`
    PageSize   int   `json:"page_size"`
}
```

---

## 3. Comentarios DEPRECATED en Código

### Ubicación
Varios archivos

### Descripción
Comentarios `// DEPRECATED` o `// TODO: deprecated` que indican código a eliminar.

### Archivos Afectados

| Archivo | Línea | Comentario |
|---------|-------|------------|
| `cmd/legacy_handlers.go` | 9 | `DEPRECATED: Estos endpoints están deprecated` |
| `cmd/legacy_handlers.go` | 156 | `DEPRECATED: Legacy response types` |

### Acción Requerida
- Buscar todos los comentarios DEPRECATED
- Evaluar si el código puede eliminarse
- Crear tickets para cada uno

```bash
# Buscar deprecados
grep -rn "DEPRECATED\|deprecated" --include="*.go" .
```

---

## 4. Imports No Utilizados (Potencial)

### Descripción
Después de eliminar legacy_handlers.go, verificar que no queden imports sin usar.

### Verificación
```bash
# Verificar imports no usados
go mod tidy
goimports -w .
```

---

## 📊 Resumen de Eliminación

| Item | Líneas | Esfuerzo | Riesgo |
|------|--------|----------|--------|
| legacy_handlers.go | 174 | 15 min | Bajo |
| Response types | 20 | 1 hora | Bajo |
| Comentarios deprecated | N/A | 30 min | Ninguno |

**Total estimado**: 2 horas de trabajo

---

## ✅ Checklist de Eliminación

```
[ ] Verificar que endpoints legacy no estén en uso
[ ] Eliminar cmd/legacy_handlers.go
[ ] Regenerar documentación Swagger
[ ] Ejecutar tests completos
[ ] Centralizar tipos de respuesta
[ ] Actualizar CHANGELOG.md
[ ] Crear release v0.6.0
```
