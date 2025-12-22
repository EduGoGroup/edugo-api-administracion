# 🟠 FASE 2: Eliminación de Código Deprecado

**Prioridad**: Alta  
**Estimación**: 2 horas  
**Rama**: `chore/fase-2-eliminar-deprecado`

---

## Preparación Git

```bash
git checkout dev
git pull origin dev
git checkout -b chore/fase-2-eliminar-deprecado
```

---

## 2.1 Eliminar `legacy_handlers.go`

### Ubicación
```
cmd/legacy_handlers.go
```

### Descripción
Archivo con 174 líneas de código muerto que contiene handlers HTTP que ya no se usan. Retornan HTTP 410 (Gone) con datos mock.

### Contenido a Eliminar
- `CreateUser` - Handler deprecado
- `UpdateUser` - Handler deprecado
- `DeleteUser` - Handler deprecado
- `CreateSubject` - Handler deprecado
- `DeleteMaterial` - Handler deprecado
- `GetGlobalStats` - Handler deprecado
- `DeprecatedResponse` - Tipo deprecado
- `SuccessResponse` - Tipo deprecado (duplicado)
- `CreateUserResponse` - Tipo deprecado
- `CreateSubjectResponse` - Tipo deprecado
- `GlobalStatsResponse` - Tipo deprecado

### Tareas
1. Eliminar el archivo completo
2. Verificar que no hay imports rotos
3. Regenerar documentación Swagger

### Comandos
```bash
rm cmd/legacy_handlers.go
make swagger
go mod tidy
```

### Esfuerzo
30 minutos

---

## 2.2 Centralizar Response Types

### Problema
Tipos `ErrorResponse`, `SuccessResponse` están duplicados en múltiples handlers.

### Tareas
1. Crear archivo centralizado de response types
2. Refactorizar handlers para usar los tipos centralizados
3. Eliminar definiciones duplicadas

### Crear Archivo
```
internal/infrastructure/http/dto/response.go
```

### Código
```go
package dto

// ErrorResponse representa una respuesta de error estándar
type ErrorResponse struct {
    Error   string `json:"error"`
    Message string `json:"message,omitempty"`
    Code    string `json:"code"`
    Details any    `json:"details,omitempty"`
}

// SuccessResponse representa una respuesta exitosa genérica
type SuccessResponse struct {
    Message string `json:"message"`
    Data    any    `json:"data,omitempty"`
}

// PaginatedResponse representa una respuesta paginada
type PaginatedResponse struct {
    Data       any   `json:"data"`
    TotalCount int64 `json:"total_count"`
    Page       int   `json:"page"`
    PageSize   int   `json:"page_size"`
}

// IDResponse representa una respuesta con solo un ID
type IDResponse struct {
    ID string `json:"id"`
}
```

### Refactorizar Handlers
Actualizar cada handler para importar y usar:
```go
import httpdto "github.com/EduGoGroup/edugo-api-administracion/internal/infrastructure/http/dto"

// Uso:
c.JSON(http.StatusBadRequest, httpdto.ErrorResponse{
    Error: "validation error",
    Code:  "VALIDATION_ERROR",
})
```

### Esfuerzo
1.5 horas

---

## Documentación a Actualizar

Al completar esta fase, actualizar:

- `documents/improvements/DEPRECATED_CODE.md` - Eliminar secciones 1, 2 y 3 (código eliminado)
- `documents/improvements/README.md` - Actualizar estado
- `documents/ARCHITECTURE.md` - Documentar ubicación de response types en la sección de Infrastructure Layer

---

## Finalización

```bash
git add .
git commit -m "chore: eliminar código deprecado y centralizar response types"
git push origin chore/fase-2-eliminar-deprecado
```

### Crear PR a dev con:
- Título: `chore: eliminar código deprecado y centralizar response types`
- Descripción: Fase 2 del plan de mejoras - Limpieza de código

---

## Checklist

- [ ] `cmd/legacy_handlers.go` eliminado
- [ ] Swagger regenerado
- [ ] `internal/infrastructure/http/dto/response.go` creado
- [ ] Handlers refactorizados para usar tipos centralizados
- [ ] No hay imports rotos
- [ ] `go build` exitoso
- [ ] Documentación actualizada
- [ ] PR creado a dev
