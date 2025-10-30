package response

import "time"

// SuccessResponse representa una respuesta genérica de éxito
type SuccessResponse struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message,omitempty" example:"Operation completed successfully"`
} // @name SuccessResponse

// ErrorResponse representa una respuesta de error
type ErrorResponse struct {
	Error   string `json:"error" example:"Invalid request"`
	Message string `json:"message,omitempty" example:"Field 'guardian_id' is required"`
} // @name ErrorResponse

// GuardianRelationResponse representa la respuesta al crear un vínculo tutor-estudiante
type GuardianRelationResponse struct {
	ID               string    `json:"id" example:"uuid-relation-789"`
	GuardianID       string    `json:"guardian_id" example:"uuid-guardian-123"`
	StudentID        string    `json:"student_id" example:"uuid-student-456"`
	RelationshipType string    `json:"relationship_type" example:"father"`
	CreatedAt        time.Time `json:"created_at" example:"2025-01-15T10:00:00Z"`
} // @name GuardianRelationResponse

// SubjectResponse representa la respuesta con información de una materia
type SubjectResponse struct {
	ID          string                 `json:"id" example:"uuid-subject-101"`
	Name        string                 `json:"name" example:"Matemáticas"`
	Description string                 `json:"description,omitempty" example:"Curso de matemáticas"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt   time.Time              `json:"created_at" example:"2025-01-01T00:00:00Z"`
	UpdatedAt   *time.Time             `json:"updated_at,omitempty" example:"2025-01-15T10:00:00Z"`
} // @name SubjectResponse
