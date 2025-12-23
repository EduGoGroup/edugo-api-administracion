package dto

// ErrorResponse representa una respuesta de error estándar
// Error: mensaje de error principal para el usuario
// Code: código de error para identificación programática (ej: "INVALID_REQUEST", "NOT_FOUND")
// Details: detalles adicionales opcionales (validaciones, campos específicos, etc)
type ErrorResponse struct {
	Error   string            `json:"error"`
	Code    string            `json:"code"`
	Details map[string]string `json:"details,omitempty"`
} // @name ErrorResponse

// SuccessResponse representa una respuesta exitosa genérica
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
} // @name SuccessResponse

// PaginatedResponse representa una respuesta paginada
// Data: array de elementos de la página actual
// Pagination: metadata de paginación (página actual, total, etc)
type PaginatedResponse struct {
	Data       interface{}    `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
} // @name PaginatedResponse

// PaginationMeta contiene metadata de paginación
// Page: número de página actual (base 1)
// PerPage: cantidad de elementos por página
// Total: cantidad total de elementos
// TotalPages: cantidad total de páginas
type PaginationMeta struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
} // @name PaginationMeta

// IDResponse representa una respuesta con solo un ID
type IDResponse struct {
	ID string `json:"id"`
} // @name IDResponse
