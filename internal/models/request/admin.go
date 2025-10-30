package request

// CreateGuardianRelationRequest representa la solicitud para crear un vínculo tutor-estudiante
type CreateGuardianRelationRequest struct {
	GuardianID       string `json:"guardian_id" binding:"required" example:"uuid-guardian-123"`
	StudentID        string `json:"student_id" binding:"required" example:"uuid-student-456"`
	RelationshipType string `json:"relationship_type" binding:"required" example:"father" enums:"father,mother,guardian,other"`
} // @name CreateGuardianRelationRequest

// UpdateSubjectRequest representa la solicitud para actualizar una materia
type UpdateSubjectRequest struct {
	Name        string                 `json:"name,omitempty" example:"Matemáticas Avanzadas"`
	Description string                 `json:"description,omitempty" example:"Curso avanzado de matemáticas"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
} // @name UpdateSubjectRequest
