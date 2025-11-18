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
