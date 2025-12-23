package dto

// ErrorResponse representa una respuesta de error estándar
type ErrorResponse struct {
	Error   string            `json:"error"`
	Message string            `json:"message,omitempty"`
	Code    string            `json:"code"`
	Details map[string]string `json:"details,omitempty"`
} // @name ErrorResponse

// SuccessResponse representa una respuesta exitosa genérica
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
} // @name SuccessResponse

// PaginatedResponse representa una respuesta paginada
type PaginatedResponse struct {
	Data       interface{}    `json:"data"`
	TotalCount int64          `json:"total_count"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	Meta       PaginationMeta `json:"meta,omitempty"`
} // @name PaginatedResponse

// PaginationMeta metadata de paginación
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
