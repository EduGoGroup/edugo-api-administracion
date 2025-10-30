package dto

// GlobalStatsResponse representa las estadísticas globales del sistema
type GlobalStatsResponse struct {
	TotalUsers            int `json:"total_users"`
	TotalActiveUsers      int `json:"total_active_users"`
	TotalSchools          int `json:"total_schools"`
	TotalSubjects         int `json:"total_subjects"`
	TotalGuardianRelations int `json:"total_guardian_relations"`
}
